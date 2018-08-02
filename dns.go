package shuttle

import "net"

type DomainHost struct {
	IP            net.IP
	Country       string
	DNS           string
	RemoteResolve bool
}

type DNS struct {
	Domain        string
	IP            net.IP
	DNS           []string
	RemoteResolve bool
}

var (
	_DNS      []string
	_LocalDNS []DNS
)

func ResolveDomain(req *Request) error {
	//DomainHost
	return nil
}
