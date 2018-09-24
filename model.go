package shuttle

import (
	"net"
	"strconv"
)

const (
	CmdTCP = 0x01
	CmdUDP = 0x03

	ProtocolSocks = "SOCKS"
	ProtocolHttp  = "HTTP"
	ProtocolHttps = "HTTPS"

	AddrTypeIPv4   = 0x01 //    0x01：IPv4
	AddrTypeDomain = 0x03 //    0x03：域名
	AddrTypeIPv6   = 0x04 //    0x04：IPv6
)

//|VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
type Request struct {
	Ver        uint8
	Cmd        uint8
	Rsv        uint8
	Atyp       uint8
	Addr       string
	IP         net.IP
	Port       uint16
	DomainHost DomainHost
	Data       []byte //for udp
	Protocol   string
	Target     string
	ConnID     int64
}

func (r *Request) Host() string {
	if len(r.IP) > 0 {
		return net.JoinHostPort(r.IP.String(), strconv.Itoa(int(r.Port)))
	} else {
		return net.JoinHostPort(r.Addr, strconv.Itoa(int(r.Port)))
	}
}

func (r *Request) Host2() string {
	if len(r.Addr) > 0 {
		return net.JoinHostPort(r.Addr, strconv.Itoa(int(r.Port)))
	} else {
		return net.JoinHostPort(r.IP.String(), strconv.Itoa(int(r.Port)))
	}
}

func (r *Request) Network() string {
	switch r.Cmd {
	case CmdTCP:
		return TCP
	case CmdUDP:
		return UDP
	}
	return ""
}

func (r *Request) GetIP() net.IP {
	if len(r.IP) > 0 {
		return r.IP
	}
	if r.IP = net.ParseIP(r.Addr); len(r.IP) > 0 {
		return r.IP
	}
	err := ResolveDomain(r)
	log.Logger.Errorf("[Request] GetIP error: %v", err)
	return nil
}
