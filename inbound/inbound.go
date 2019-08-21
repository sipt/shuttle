package inbound

import (
	"context"
	"fmt"
	"net"

	"github.com/sipt/shuttle/conf/model"
)

const (
	ProtocolTCP_HTTP  = "http"
	ProtocolTCP_HTTPS = "https"
	ProtocolUDP_DNS   = "dns"
)

func ApplyConfig(config *model.Config) ([]Inbound, error) {
	for _, v := range config.Listener {
	}
}

type NewFunc func(name, addr, port string, params map[string]string) (Inbound, error)

var creator = make(map[string]NewFunc)

// Register: register {key: NewFunc}
func Register(key string, f NewFunc) {
	creator[key] = f
}

// Get: get inbound by key
func Get(typ, name, addr, port string, params map[string]string) (Inbound, error) {
	f, ok := creator[typ]
	if !ok {
		return nil, fmt.Errorf("inbound not support: %s", typ)
	}
	return f(name, addr, port, params)
}

type Inbound interface {
	Listen(ctx context.Context, callback func(ctx context.Context, conn net.Conn)) (close func(), err error)
	ListConn() []net.Conn
	Protocol() string
}
