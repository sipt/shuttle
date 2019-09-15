package rule

import (
	"context"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/conf/model"
	"github.com/sipt/shuttle/dns"
)

func ApplyConfig(ctx context.Context, config *model.Config, isUDP bool, proxyName map[string]bool, fallback Handle, dnsHandle dns.Handle) (handle Handle, err error) {
	handle = fallback
	rules := config.Rule
	if isUDP {
		rules = config.UDPRule
	}
	for i := len(rules) - 1; i >= 0; i-- {
		v := rules[i]
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
		handle, err = Get(ctx, rule.Typ, rule, handle, dnsHandle)
		if err != nil {
			logrus.WithError(err).WithField("name", rule.Typ).WithField("value", rule.Value).
				WithField("proxy", rule.Proxy).Error("init rule failed")
			return
		}
	}
	return
}

// simple of RequestInfo
type RequestInfo interface {
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

func (r *Rule) String() string {
	return fmt.Sprintf("%s:%s [%s]", r.Typ, r.Value, r.Proxy)
}

type Handle func(ctx context.Context, info RequestInfo) *Rule
type NewFunc func(ctx context.Context, rule *Rule, handle Handle, dnsHandle dns.Handle) (Handle, error)

var creator = make(map[string]NewFunc)

// Register: register {key: NewFunc}
func Register(key string, f NewFunc) {
	creator[key] = f
}

// Get: get rule by key
func Get(ctx context.Context, typ string, rule *Rule, handle Handle, dnsHandle dns.Handle) (Handle, error) {
	f, ok := creator[typ]
	if !ok {
		return nil, fmt.Errorf("rule not support: %s", typ)
	}
	return f(ctx, rule, handle, dnsHandle)
}
