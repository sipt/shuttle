package group

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/sipt/shuttle/conn"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/server"
)

const (
	TypRTT          = "rtt"
	DefaultTestURL  = "http://www.gstatic.com/generate_204"
	DefaultInterval = 10 * time.Minute

	ParamsKeyTestURL  = "test_url"
	ParamsKeyInterval = "interval"
)

func init() {
	Register(TypRTT, newRttGroup)
}

func newRttGroup(ctx context.Context, name string, params map[string]string) (group IGroup, err error) {
	rtt := &rttGroup{
		ctx:     ctx,
		name:    name,
		RWMutex: &sync.RWMutex{},
	}
	rtt.testUrl = params[ParamsKeyTestURL]
	if rtt.testUrl == "" {
		rtt.testUrl = DefaultTestURL
	} else if testUrl, err := url.Parse(rtt.testUrl); err != nil || len(testUrl.Scheme) == 0 || len(testUrl.Hostname()) == 0 {
		err = errors.Errorf("[group: %s] [%s: %s] is invalid", name, ParamsKeyTestURL, rtt.testUrl)
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
	return rtt, nil
}

type rttGroup struct {
	ctx      context.Context
	name     string
	servers  []IServerX
	current  IServerX
	testUrl  string
	interval time.Duration
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
	go r.autoSelectByRTT()
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
func (r *rttGroup) Select(name string) error {
	return nil
}

func (r *rttGroup) autoSelectByRTT() {
	r.testAllRTT()
	ticker := time.NewTicker(r.interval)
	for {
		select {
		case <-ticker.C:
			r.testAllRTT()
		case <-r.ctx.Done():
			return
		}
	}
	return
}

func (r *rttGroup) testAllRTT() {
	if len(r.servers) == 0 {
		return
	}
	ctx, _ := context.WithTimeout(r.ctx, time.Second*10)
	var (
		reply   = make(chan IServerX, len(r.servers))
		current IServerX
	)
	for _, v := range r.servers {
		go r.testServerRTT(ctx, v, reply)
	}
	for range r.servers {
		if current = <-reply; current != nil {
			r.current = current
			close(reply)
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
				return s.Server().DialTCP(ctx, host, fmt.Sprint(port), conn.DefaultDialTCP)
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
