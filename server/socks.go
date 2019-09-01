package server

import (
	"context"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/conn"
	"github.com/sipt/shuttle/dns"
	"github.com/sipt/shuttle/pkg/socks"
)

const (
	Socks5 = "socks5"

	ParamsKeyAuthType = "auth_type"
	ParamsKeyUser     = "user"
	ParamsKeyPassword = "password"

	AuthTypeBasic = "basic"
)

func init() {
	Register(Socks5, func(name, host string, port int, params map[string]string, dnsHandle dns.Handle) (server IServer, err error) {
		s := &Socks5Server{
			RWMutex:   &sync.RWMutex{},
			rtt:       make(map[string]time.Duration),
			port:      port,
			name:      name,
			dnsHandle: dnsHandle,
		}
		if ip := net.ParseIP(host); len(ip) > 0 {
			s.ip = ip
		} else {
			s.domain = host
		}
		switch params[ParamsKeyAuthType] {
		case AuthTypeBasic:
			if u, p := params[ParamsKeyUser], params[ParamsKeyPassword]; len(u) == 0 || len(p) == 0 {
				return nil, errors.Errorf("server[%s].user[%s] or password[%s] invalid", name, u, p)
			} else {
				s.auth = &socks.UsernamePassword{
					Username: u,
					Password: p,
				}
			}
		case "":
		default:
			return nil, errors.Errorf("server[%s].auth_type[%s] invalid", name, params[ParamsKeyAuthType])
		}
		return s, nil
	})
}

type Socks5Server struct {
	IServer   // just for not implement: TestRtt
	name      string
	rtt       map[string]time.Duration
	auth      *socks.UsernamePassword
	dnsHandle dns.Handle
	domain    string
	ip        net.IP
	port      int
	*sync.RWMutex
}

func (s *Socks5Server) Typ() string {
	return Socks5
}
func (s *Socks5Server) Name() string {
	return s.name
}
func (s *Socks5Server) SetRtt(key string, rtt time.Duration) {
	s.Lock()
	defer s.Unlock()
	s.rtt[key] = rtt
}
func (s *Socks5Server) Rtt(key string) time.Duration {
	s.RLock()
	defer s.RUnlock()
	return s.rtt[key]
}

func (s *Socks5Server) Dial(ctx context.Context, network string, info Info, dial conn.DialFunc) (conn.ICtxConn, error) {
	var port = strconv.Itoa(s.port)
	var host = s.domain
	if len(s.ip) > 0 {
		host = s.ip.String()
	} else {
		timeOutCtx, _ := context.WithTimeout(ctx, time.Second)
		answer := s.dnsHandle(timeOutCtx, host)
		if answer == nil {
			return nil, errors.Errorf("resolve domain[%s] failed", host)
		}
		host = answer.CurrentIP.String()
	}
	dialer := socks.NewDialer(network, net.JoinHostPort(host, port))
	if s.auth != nil {
		dialer.Authenticate = s.auth.Authenticate
	}
	dialer.ProxyDial = func(ctx context.Context, network string, addr string) (net.Conn, error) {
		return dial(ctx, network, host, port)
	}
	var targetHost string
	if len(info.Domain()) == 0 {
		targetHost = info.IP().String()
	} else {
		targetHost = info.Domain()
	}
	var targetPort = strconv.Itoa(info.Port())
	sc, err := dialer.DialContext(ctx, network, net.JoinHostPort(targetHost, targetPort))
	if err != nil {
		return nil, err
	}
	return sc.(*socks.Conn).Conn.(conn.ICtxConn), nil
}
