package inbound

import (
	"context"
	"fmt"

	"github.com/sipt/shuttle/conf/model"
	"github.com/sipt/shuttle/constant/typ"
	"github.com/sirupsen/logrus"
)

const (
	ProtocolTCP_HTTP  = "http"
	ProtocolTCP_HTTPS = "https"
	ProtocolUDP_DNS   = "dns"
)

var ctx = context.Background()
var inboundContext = make(map[string]context.CancelFunc)

func Cancel(addr string) {
	cancel := inboundContext[addr]
	if cancel != nil {
		cancel()
	}
}

func ApplyConfig(config *model.Config, handle typ.HandleFunc) error {
	for _, v := range config.Listener {
		f, err := Get(v.Typ, v.Addr, v.Params)
		if err != nil {
			return err
		}
		subCtx, cancel := context.WithCancel(context.WithValue(ctx, "addr", v.Addr))
		inboundContext[v.Addr] = cancel
		go f(subCtx, handle)
	}
	logrus.Debug("inbound init success")
	return nil
}

type NewFunc func(addr string, params map[string]string) (listen func(context.Context, typ.HandleFunc), err error)

var creator = make(map[string]NewFunc)

// Register: register {key: NewFunc}
func Register(key string, f NewFunc) {
	creator[key] = f
}

// Get: get inbound by key
func Get(typ, addr string, params map[string]string) (func(context.Context, typ.HandleFunc), error) {
	f, ok := creator[typ]
	if !ok {
		return nil, fmt.Errorf("inbound not support: %s", typ)
	}
	return f(addr, params)
}
