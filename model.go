package shuttle

import (
	"net"
	"strconv"
)

const (
	cmdTCP = 0x01
	cmdUDP = 0x03

	ProtocolSocks = "SOCKS"
	ProtocolHttp  = "HTTP"
	ProtocolHttps = "HTTPS"
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
	case cmdTCP:
		return TCP
	case cmdUDP:
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
	Logger.Errorf("[Request] GetIP error: %v", err)
	return nil
}
