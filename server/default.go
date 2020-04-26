package server

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/conn"
	"github.com/sipt/shuttle/dns"
	"github.com/sirupsen/logrus"
)

func init() {
	Register(Direct, func(_, _ string, _ int, _ map[string]string, _ dns.Handle) (server IServer, e error) {
		return &DirectServer{
			RWMutex: &sync.RWMutex{},
			rtt:     make(map[string]time.Duration),
		}, nil
	})
	Register(Reject, func(_, _ string, _ int, _ map[string]string, _ dns.Handle) (server IServer, e error) {
		return &RejectServer{}, nil
	})
}

type DirectServer struct {
	IServer // just for not implement: TestRtt
	rtt     map[string]time.Duration
	*sync.RWMutex
}

func (d *DirectServer) Typ() string {
	return Direct
}
func (d *DirectServer) Name() string {
	return Direct
}
func (d *DirectServer) SetRtt(key string, rtt time.Duration) {
	d.Lock()
	defer d.Unlock()
	d.rtt[key] = rtt
}
func (d *DirectServer) Rtt(key string) time.Duration {
	d.RLock()
	defer d.RUnlock()
	return d.rtt[key]
}
func (d *DirectServer) UdpRelay() bool {
	return true
}

func (d *DirectServer) Dial(ctx context.Context, network string, info Info, dial conn.DialFunc) (conn.ICtxConn, error) {
	var host string
	if len(info.IP()) == 0 {
		host = info.Domain()
	} else {
		host = info.IP().String()
	}
	return dial(ctx, network, host, strconv.Itoa(info.Port()))
}

type RejectServer struct{}

func (r *RejectServer) Typ() string {
	return Reject
}
func (r *RejectServer) Name() string {
	return Reject
}
func (r *RejectServer) SetRtt(_ string, _ time.Duration) {
}
func (r *RejectServer) Rtt(_ string) time.Duration {
	return time.Duration(-1)
}
func (r *RejectServer) TestRtt(_, _ string) time.Duration {
	return time.Duration(-1)
}
func (r *RejectServer) Dial(ctx context.Context, network string, info Info, dial conn.DialFunc) (conn.ICtxConn, error) {
	return nil, ErrRejected
}
func (r *RejectServer) UdpRelay() bool {
	return false
}

func NewRttServer(server IServer, params map[string]string) IServer {
	rtt := &RttServer{
		IServer: server,
		testUri: params[ParamsKeyTestURI],
	}
	if rtt.testUri == "" {
		rtt.testUri = DefaultTestURL
	} else if testUrl, err := url.Parse(rtt.testUri); err != nil || len(testUrl.Scheme) == 0 || len(testUrl.Hostname()) == 0 {
		err = errors.Errorf("[server: %s] [%s: %s] is invalid", server.Name(), ParamsKeyTestURI, rtt.testUri)
		rtt.testUri = DefaultTestURL
	}
	rtt.client = &http.Client{
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
				ctx, _ = context.WithTimeout(ctx, time.Second*3)
				conn, err := rtt.Dial(ctx, "tcp", req, conn.DefaultDial)
				return conn, err
			},
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	return rtt
}

type RttServer struct {
	IServer
	testUri string
	client  *http.Client
}

func (r *RttServer) TestRtt(key, uri string) time.Duration {
	log := logrus.WithField("method", "rtt-test").WithField("server", key)
	start := time.Now()
	if len(uri) == 0 {
		uri = r.testUri
	}
	resp, err := r.client.Get(uri)
	if err != nil {
		r.SetRtt(key, time.Duration(-1))
		log.WithError(err).WithField("uri", uri).WithField("rtt", "failed").
			Debug("rtt test failed")
		return r.Rtt(key)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		r.SetRtt(key, time.Now().Sub(start))
		log.WithField("rtt", r.Rtt(key).Round(time.Millisecond).String()).Debug("rtt test success")
	} else {
		r.SetRtt(key, time.Duration(-1))
		log.WithField("rtt", "failed").WithField("uri", uri).
			WithField("status_code", resp.StatusCode).Debug("rtt test failed")
	}
	return r.Rtt(key)
}
