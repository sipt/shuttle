package group

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/conn"
	"github.com/sipt/shuttle/dns"
	"github.com/sipt/shuttle/server"
	"github.com/sirupsen/logrus"
)

const (
	TypRTT          = "rtt"
	DefaultTestURL  = "http://www.gstatic.com/generate_204"
	DefaultInterval = 10 * time.Minute

	ParamsKeyTestURI  = "test-url"
	ParamsKeyInterval = "interval"
	ParamsKeyUdpRelay = "udp-relay"
)

func init() {
	Register(TypRTT, newRttGroup)
}

func newRttGroup(ctx context.Context, name string, params map[string]string, _ dns.Handle) (group IGroup, err error) {
	rtt := &rttGroup{
		ctx:     ctx,
		name:    name,
		RWMutex: &sync.RWMutex{},
	}
	rtt.testUrl = params[ParamsKeyTestURI]
	if rtt.testUrl == "" {
		rtt.testUrl = DefaultTestURL
	} else if testUrl, err := url.Parse(rtt.testUrl); err != nil || len(testUrl.Scheme) == 0 || len(testUrl.Hostname()) == 0 {
		err = errors.Errorf("[group: %s] [%s: %s] is invalid", name, ParamsKeyTestURI, rtt.testUrl)
		return nil, err
	}
	interval := params[ParamsKeyInterval]
	if len(interval) == 0 {
		rtt.interval = DefaultInterval
	} else {
		rtt.interval, err = time.ParseDuration(interval)
		if err != nil {
			err = errors.Wrapf(err, "[group: %s] [%s: %s] is invalid", name, ParamsKeyInterval, interval)
			return
		}
	}
	rtt.reset = make(chan bool)
	rtt.udpRelay = params[ParamsKeyUdpRelay] == "true"
	return rtt, nil
}

type rttGroup struct {
	ctx      context.Context
	cancel   func()
	name     string
	servers  []IServerX
	current  IServerX
	testUrl  string
	interval time.Duration
	reset    chan bool
	udpRelay bool
	*sync.RWMutex
}

func (r *rttGroup) Append(servers []IServerX) {
	if len(servers) == 0 {
		return
	}
	r.Lock()
	defer r.Unlock()
	if len(r.servers) == 0 {
		r.servers = servers
	} else {
		r.servers = append(r.servers, servers...)
	}
	r.current = r.servers[0]
	var ctx context.Context
	ctx, r.cancel = context.WithCancel(r.ctx)
	go r.autoSelectByRTT(ctx)
}
func (r *rttGroup) Items() []IServerX {
	return r.servers
}
func (r *rttGroup) Typ() string {
	return TypRTT
}
func (r *rttGroup) Name() string {
	return r.name
}
func (r *rttGroup) Trace() []string {
	trace := make([]string, 0, len(r.current.Trace())+1)
	return append(append(trace, r.name), r.current.Trace()...)
}
func (r *rttGroup) UdpRelay() bool {
	return r.udpRelay
}
func (r *rttGroup) Server() server.IServer {
	r.RLock()
	defer r.RUnlock()
	if r.current != nil {
		return r.current.Server()
	} else if len(r.servers) > 0 {
		return r.servers[0].Server()
	} else {
		return nil
	}
}
func (r *rttGroup) Clear() {
	if len(r.servers) == 0 {
		return
	}
	r.Lock()
	defer r.Unlock()
	r.cancel()      // stop timer
	r.servers = nil // clear all
}
func (r *rttGroup) Selected() IServerX {
	return r.current
}
func (r *rttGroup) Select(name string) error {
	return nil
}
func (r *rttGroup) Reset() {
	r.reset <- true
}

func (r *rttGroup) autoSelectByRTT(ctx context.Context) {
	r.testAllRTT()
	timer := time.NewTimer(r.interval)
	for {
		select {
		case <-timer.C:
			r.testAllRTT()
		case <-r.reset:
			timer.Stop()
			r.testAllRTT()
			timer.Reset(r.interval)
		case <-ctx.Done():
			timer.Stop()
			return
		}
	}
	return
}

func (r *rttGroup) testAllRTT() {
	if len(r.servers) == 0 {
		return
	}
	var (
		reply   = make(chan IServerX, len(r.servers))
		current IServerX
	)
	for _, v := range r.servers {
		go func(s IServerX) {
			if s.Server().TestRtt(r.name, r.testUrl) > 0 {
				reply <- s
			} else {
				reply <- nil
			}
		}(v)
	}
	var skip = true
	for range r.servers {
		if current = <-reply; skip && current != nil {
			skip = false
			r.current = current
			return
		}
	}
}

func (r *rttGroup) testServerRTT(ctx context.Context, s IServerX, reply chan IServerX) {
	log := logrus.WithField("method", "rtt-test").WithField("server", strings.Join(r.Trace(), ":"))
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				host, port, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, err
				}
				req := &reqInfo{}
				req.port, _ = strconv.Atoi(port)
				req.ip = net.ParseIP(host)
				if len(req.ip) == 0 {
					req.domain = host
				}
				conn, err := s.Server().Dial(ctx, "tcp", req, conn.DefaultDial)
				return conn, err
			},
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	start := time.Now()
	resp, err := client.Get(r.testUrl)
	if err != nil {
		s.Server().SetRtt(r.name, time.Duration(-1))
		reply <- nil
		log.WithField("rtt", "failed").Debug("rtt test failed")
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		s.Server().SetRtt(r.name, time.Now().Sub(start))
		reply <- s
		log.WithField("rtt", s.Server().Rtt(r.name).Round(time.Millisecond).String()).Debug("rtt test success")
	} else {
		s.Server().SetRtt(r.name, time.Duration(-1))
		reply <- nil
		log.WithField("rtt", "failed").Debug("rtt test failed")
	}
}

type reqInfo struct {
	domain string
	ip     net.IP
	port   int
}

func (r *reqInfo) Domain() string {
	return r.domain
}
func (r *reqInfo) IP() net.IP {
	return r.ip
}
func (r *reqInfo) Port() int {
	return r.port
}
