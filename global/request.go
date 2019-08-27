package global

import "net"

type RequestInfo interface {
	Network() string
	Domain() string
	URI() string
	IP() net.IP
	CountryCode() string
	Port() int
	SetIP(net.IP)
	SetPort(int)
	SetCountryCode(string)
}
