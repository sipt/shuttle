package proxy

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

func Duration2Str(d time.Duration) string {
	if d == -1 {
		return "failed"
	}
	if d == 0 {
		return "???"
	} else if d > time.Second {
		return fmt.Sprintf("%ds", d/time.Second)
	} else if d > time.Millisecond {
		return fmt.Sprintf("%dms", d/time.Millisecond)
	} else if d > time.Microsecond {
		return fmt.Sprintf("%dus", d/time.Microsecond)
	} else if d > time.Microsecond {
		return fmt.Sprintf("%dns", d/time.Nanosecond)
	}
	return "0s"
}

func TestRTT(s IServer, testURL string) (rtt time.Duration, err error) {
	var server *Server
	server, err = s.GetServer()
	if err != nil {
		return 0, nil
	}
	var sc net.Conn
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, addr string) (i net.Conn, e error) {
				rttRequest := &_HttpRequest{
					network: "tcp",
				}
				for i := len(addr) - 1; i >= 0; i-- {
					if addr[i] == ':' {
						rttRequest.port = addr[i+1:]
						rttRequest.domain = addr[:i]
						if ip := net.ParseIP(rttRequest.domain); len(ip) > 0 {
							rttRequest.port = rttRequest.domain
							rttRequest.domain = ""
						}
					}
				}
				var err error
				sc, err = server.Conn(rttRequest)
				return sc, err
			},
		},
	}
	start := time.Now()
	resp, err := client.Get(testURL)
	defer func() {
		if sc != nil {
			_ = sc.Close()
		}
	}()
	if err != nil {
		return 0, err
	}
	rtt = time.Now().Sub(start)
	if !resp.Close && resp.Body != nil {
		_ = resp.Body.Close()
	}
	return
}

//HTTP Request
type _HttpRequest struct {
	network string
	domain  string
	ip      string
	port    string
}

func (r *_HttpRequest) Network() string {
	return r.network
}
func (r *_HttpRequest) Domain() string {
	return r.domain
}
func (r *_HttpRequest) IP() string {
	return r.ip
}
func (r *_HttpRequest) Port() string {
	return r.port
}

//return domain!=""?domain:ip
func (r *_HttpRequest) Addr() string {
	if len(r.domain) > 0 {
		return r.domain
	}
	return r.ip
}

//return [domain/ip]:[port]
func (r *_HttpRequest) Host() string {
	if len(r.IP()) > 0 {
		return net.JoinHostPort(r.IP(), r.Port())
	}
	return net.JoinHostPort(r.Addr(), r.Port())
}
