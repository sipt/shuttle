package inbound

import (
	"context"
	"fmt"
	"net"

	"github.com/sipt/shuttle/listener"

	"github.com/sipt/shuttle/conf/model"
)

const (
	ProtocolTCP_HTTP  = "http"
	ProtocolTCP_HTTPS = "https"
	ProtocolUDP_DNS   = "dns"
)

func ApplyConfig(config *model.Config) ([]Inbound, error) {
	return nil, nil
}

type NewFunc func(addr string, params map[string]string) (listen func(listener.HandleFunc) error, err error)

var creator = make(map[string]NewFunc)

// Register: register {key: NewFunc}
func Register(key string, f NewFunc) {
	creator[key] = f
}

// Get: get inbound by key
func Get(typ, addr string, params map[string]string) (func(listener.HandleFunc) error, error) {
	f, ok := creator[typ]
	if !ok {
		return nil, fmt.Errorf("inbound not support: %s", typ)
	}
	return f(addr, params)
}

type Inbound interface {
	Listen(ctx context.Context, callback func(ctx context.Context, conn net.Conn)) (close func(), err error)
	ListConn() []net.Conn
	Protocol() string
}
