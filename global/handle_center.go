package global

import "github.com/sipt/shuttle/dns"

func init() {
	namespace = make(map[string]*handleCenter)
}

const defaultName = "default"

var namespace map[string]*handleCenter

type handleCenter struct {
	dnsHandle dns.Handle
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
