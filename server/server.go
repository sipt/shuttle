package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/conf/model"
	"github.com/sipt/shuttle/conn"
)

const (
	Direct = "DIRECT"
	Reject = "REJECT"
)

var (
	ErrRejected = errors.New("connect was rejected")
	defaults    = []string{Direct, Reject}
)

func ApplyConfig(config *model.Config) (map[string]IServer, error) {
	servers := make(map[string]IServer, len(config.Server)+len(defaults))
	var (
		s   IServer
		err error
	)
	for name, v := range config.Server {
		s, err = Get(v.Typ, name, v.Addr, v.Port, v.Params)
		if err != nil {
			return nil, err
		}
		servers[s.Name()] = s
	}
	for _, v := range defaults {
		s, err = Get(v, Direct, "", "", nil)
		if err != nil {
			return nil, err
		}
		servers[s.Name()] = s
	}
	return servers, nil
}

type NewFunc func(name, addr, port string, params map[string]string) (IServer, error)

var creator = make(map[string]NewFunc)

// Register: register {key: NewFunc}
func Register(key string, f NewFunc) {
	creator[key] = f
}

// Get: get server by key
func Get(typ, name, addr, port string, params map[string]string) (IServer, error) {
	f, ok := creator[typ]
	if !ok {
		return nil, fmt.Errorf("server not support: %s", typ)
	}
	return f(name, addr, port, params)
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
