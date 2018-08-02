package shuttle

import (
	"net"
	"strconv"
)

const (
	cmdTCP = 0x01
	cmdUDP = 0x03
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
	DomainHost *DomainHost
	Data       []byte //for udp
}

func (r *Request) Host() string {
	if r.IP != nil && len(r.IP) > 0 {
		return net.JoinHostPort(r.IP.String(), strconv.Itoa(int(r.Port)))
	} else {
		return net.JoinHostPort(r.Addr, strconv.Itoa(int(r.Port)))
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

type HttpRequest struct {
	*Request
	Scheme string
}
