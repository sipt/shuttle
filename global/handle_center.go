package global

import (
	"context"

	"github.com/sipt/shuttle/dns"
)

func init() {
	namespace = make(map[string]*handleCenter)
}

const defaultName = "default"

var namespace map[string]*handleCenter

type handleCenter struct {
	dnsHandle dns.Handle
	ctx       context.Context
	cancel    context.CancelFunc
}

func GetDnsHandle(name ...string) dns.Handle {
	if len(name) == 0 || len(name[0]) == 0 {
		return namespace[defaultName].dnsHandle
	}
	return namespace[name[0]].dnsHandle
}

func SetDnsHandle(handle dns.Handle, name ...string) {
	var tempName string
	if len(name) == 0 || len(name[0]) == 0 {
		tempName = defaultName
	} else {
		tempName = name[0]
	}
	if hc, ok := namespace[tempName]; ok {
		hc.dnsHandle = handle
	} else {
		namespace[tempName] = &handleCenter{
			dnsHandle: handle,
		}
	}
}

func SetContext(ctx context.Context, name ...string) {
	var tempName string
	if len(name) == 0 || len(name[0]) == 0 {
		tempName = defaultName
	} else {
		tempName = name[0]
	}
	if hc, ok := namespace[tempName]; ok {
		hc.ctx, hc.cancel = context.WithCancel(ctx)
	} else {
		hc := &handleCenter{}
		hc.ctx, hc.cancel = context.WithCancel(ctx)
		namespace[tempName] = hc
	}
}
