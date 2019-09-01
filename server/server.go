package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/conf/model"
	"github.com/sipt/shuttle/conn"
	"github.com/sipt/shuttle/dns"
)

const (
	Direct = "DIRECT"
	Reject = "REJECT"
)

var (
	ErrRejected = errors.New("connect was rejected")
	defaults    = []string{Direct, Reject}
)

func ApplyConfig(config *model.Config, dnsHandle dns.Handle) (map[string]IServer, error) {
	servers := make(map[string]IServer, len(config.Server)+len(defaults))
	var (
		s   IServer
		err error
	)
	for name, v := range config.Server {
		s, err = Get(v.Typ, name, v.Host, v.Port, v.Params, dnsHandle)
		if err != nil {
			return nil, err
		}
		servers[s.Name()] = s
	}
	for _, v := range defaults {
		s, err = Get(v, Direct, "", 0, nil, dnsHandle)
		if err != nil {
			return nil, err
		}
		servers[s.Name()] = s
	}
	return servers, nil
}

type NewFunc func(name, addr string, port int, params map[string]string, dnsHandle dns.Handle) (IServer, error)

var creator = make(map[string]NewFunc)

// Register: register {key: NewFunc}
func Register(key string, f NewFunc) {
	creator[key] = f
}

// Get: get server by key
func Get(typ, name, addr string, port int, params map[string]string, dnsHandle dns.Handle) (IServer, error) {
	f, ok := creator[typ]
	if !ok {
		return nil, fmt.Errorf("server not support: %s", typ)
	}
	return f(name, addr, port, params, dnsHandle)
}

type IServer interface {
	Typ() string
	Name() string
	SetRtt(key string, duration time.Duration)
	Rtt(key string) time.Duration
	// connect to server
	Dial(ctx context.Context, network string, info Info, dial conn.DialFunc) (conn.ICtxConn, error)
}

type Info interface {
	Domain() string
	IP() net.IP
	Port() int
}
