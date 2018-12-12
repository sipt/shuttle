package shuttle

import (
	"github.com/sipt/shuttle/conn"
	"github.com/sipt/shuttle/dns"
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
type SocksRequest struct {
	ver      uint8
	cmd      uint8
	rsv      uint8
	atyp     uint8
	addr     string
	ip       net.IP
	port     uint16
	data     []byte //for udp
	protocol string
	target   string
	connID   int64
	answer   *dns.Answer
}

func (r *SocksRequest) Network() string {
	switch r.cmd {
	case CmdTCP:
		return conn.TCP
	case CmdUDP:
		return conn.UDP
	}
	return ""
}
func (r *SocksRequest) Domain() string {
	return r.addr
}
func (r *SocksRequest) IP() string {
	if len(r.ip) > 0 {
		return r.ip.String()
	} else if r.answer != nil {
		r.ip = net.ParseIP(r.answer.GetIP())
		return r.ip.String()
	}
	return ""
}
func (r *SocksRequest) Port() string {
	if r.port != 0 {
		return strconv.FormatInt(int64(r.port), 10)
	}
	return ""
}
func (r *SocksRequest) Answer() *dns.Answer {
	return r.answer
}
func (r *SocksRequest) SetAnswer(answer *dns.Answer) {
	r.answer = answer
}

//return request id
func (r *SocksRequest) ID() int64 {
	return r.connID
}

//return domain!=""?domain:ip
func (r *SocksRequest) Addr() string {
	if len(r.addr) > 0 {
		return r.addr
	}
	return r.IP()
}

//return [domain/ip]:[port]
func (r *SocksRequest) Host() string {
	if len(r.IP()) > 0 {
		return net.JoinHostPort(r.IP(), strconv.FormatInt(int64(r.port), 10))
	}
	return net.JoinHostPort(r.Addr(), strconv.FormatInt(int64(r.port), 10))
}

func NewHttpRequest(network string, domain string, ip string, port string, protocol string,
	target string, connID int64, answer *dns.Answer) *HttpRequest {
	return &HttpRequest{
		network:  network,
		domain:   domain,
		ip:       ip,
		port:     port,
		protocol: protocol,
		target:   target,
		connID:   connID,
		answer:   answer,
	}
}

//HTTP Request
type HttpRequest struct {
	network  string
	domain   string
	ip       string
	port     string
	protocol string
	target   string
	connID   int64
	answer   *dns.Answer
}

func (r *HttpRequest) Network() string {
	return r.network
}
func (r *HttpRequest) Domain() string {
	return r.domain
}
func (r *HttpRequest) IP() string {
	if len(r.ip) == 0 && r.answer != nil {
		r.ip = r.answer.GetIP()
	}
	return r.ip
}
func (r *HttpRequest) Port() string {
	if len(r.port) == 0 {
		if r.answer != nil && len(r.answer.Port) > 0 {
			r.port = r.answer.Port
		} else {
			switch r.protocol {
			case HTTP:
				r.port = "80"
			case HTTPS:
				r.port = "443"
			}
		}
	}
	return r.port
}
func (r *HttpRequest) Answer() *dns.Answer {
	return r.answer
}
func (r *HttpRequest) SetAnswer(answer *dns.Answer) {
	r.answer = answer
}

//return request id
func (r *HttpRequest) ID() int64 {
	return r.connID
}

//return domain!=""?domain:ip
func (r *HttpRequest) Addr() string {
	if len(r.domain) > 0 {
		return r.domain
	}
	return r.ip
}

//return [domain/ip]:[port]
func (r *HttpRequest) Host() string {
	if len(r.IP()) > 0 {
		return net.JoinHostPort(r.IP(), r.Port())
	}
	return net.JoinHostPort(r.Addr(), r.Port())
}
