package rule

import (
	"context"
	"fmt"
	"net"

	"github.com/sipt/shuttle/dns"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/conf/model"
)

func ApplyConfig(config *model.Config, proxyName map[string]bool, fallback Handle, dnsHandle dns.Handle) (handle Handle, err error) {
	handle = fallback
	for i := len(config.Rule) - 1; i >= 0; i-- {
		v := config.Rule[i]
		rule := &Rule{
			Typ:     v.Typ,
			Value:   v.Value,
			Proxy:   v.Proxy,
			Params:  v.Params,
			Profile: config.Info.Name,
		}
		if !proxyName[rule.Proxy] {
			err = errors.Errorf("rule:[%s, %s, %s, %v], proxy:[%s] not found",
				rule.Typ, rule.Value, rule.Proxy, rule.Params, rule.Proxy)
			return
		}
		handle, err = Get(rule.Typ, rule, handle, dnsHandle)
		if err != nil {
			return
		}
	}
	return
}

type Info interface {
	Domain() string
	URI() string
	IP() net.IP
	CountryCode() string
	Port() int
	SetIP(net.IP)
	SetPort(int)
	SetCountryCode(string)
}

type Rule struct {
	Parent  *Rule
	Profile string
	Typ     string
	Value   string
	Proxy   string
	Params  map[string]string
}

type Handle func(ctx context.Context, info Info) *Rule
type NewFunc func(rule *Rule, handle Handle, dnsHandle dns.Handle) (Handle, error)

var creator = make(map[string]NewFunc)

// Register: register {key: NewFunc}
func Register(key string, f NewFunc) {
	creator[key] = f
}

// Get: get rule by key
func Get(typ string, rule *Rule, handle Handle, dnsHandle dns.Handle) (Handle, error) {
	f, ok := creator[typ]
	if !ok {
		return nil, fmt.Errorf("rule not support: %s", typ)
	}
	return f(rule, handle, dnsHandle)
}
