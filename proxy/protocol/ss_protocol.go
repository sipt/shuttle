package protocol

import (
	"bytes"
	"fmt"
	"github.com/sipt/shuttle"
	"github.com/sipt/shuttle/ciphers"
	connect "github.com/sipt/shuttle/conn"
	"github.com/sipt/shuttle/dns"
	"github.com/sipt/shuttle/log"
	sproxy "github.com/sipt/shuttle/proxy"
	"net"
	"strconv"
)

func init() {
	sproxy.RegisterProxyProtocolCreator("ss", NewSsProtocol)
}

func NewSsProtocol(params []string) (sproxy.IProtocol, error) {
	//[]string{"addr", "port", "method", "password"}
	if len(params) != 4 {
		log.Logger.Errorf(`[SOCKS5 Server] init socks5 server failed params must be ["addr", "port", "method", "password"], but: %v`, params)
		return nil, fmt.Errorf(`[SOCKS5 Server] init socks5 server failed params must be ["addr", "port", "method", "password"], but: %v`, params)
	}
	ser := &ssProtocol{
		Addr:     params[0],
		Port:     params[1],
		Method:   params[2],
		Password: params[3],
	}
	return ser, nil
}

type ssProtocol struct {
	Addr     string
	Port     string
	Method   string
	Password string
}

func (s *ssProtocol) Conn(req sproxy.IRequest) (connect.IConn, error) {
	network := req.Network()

	var addr = s.Addr
	answer, err := dns.ResolveDomainByCache(s.Addr)
	if err != nil {
		log.Logger.Errorf("[SsProtocol] [Conn] Resolve domain failed [%s]: %v", s.Addr, err)
	} else if answer != nil {
		addr = answer.GetIP()
	}
	conn, err := net.DialTimeout(network, net.JoinHostPort(addr, s.Port), connect.DefaultTimeOut)
	if err != nil {
		return nil, err
	}
	c, err := connect.DefaultDecorate(conn, network)
	if err != nil {
		return nil, err
	}
	c, err = connect.TrafficDecorate(c)
	if err != nil {
		return nil, err
	}
	if network == connect.UDP {
		c, err = connect.BufferDecorate(c)
		if err != nil {
			return nil, err
		}
	}
	rc, err := ciphers.CipherDecorate(s.Password, s.Method, c)
	if err != nil {
		return nil, err
	}
	rawAddr, err := AddressEncoding(req)
	if err != nil {
		return nil, err
	}
	_, err = rc.Write(rawAddr)
	if err != nil {
		return nil, err
	}
	return rc, nil
}

func AddressEncoding(req sproxy.IRequest) ([]byte, error) {
	var atyp uint8 = shuttle.AddrTypeDomain
	ip := net.ParseIP(req.Domain())
	addr := []byte(req.Domain())
	if ip != nil {
		if len(ip) == net.IPv4len {
			atyp = shuttle.AddrTypeIPv4
		} else {
			atyp = shuttle.AddrTypeIPv6
		}
		addr = []byte(ip)
	}
	if len(addr) == 0 {
		ip = net.ParseIP(req.IP())
		if len(ip) == net.IPv4len {
			atyp = shuttle.AddrTypeIPv4
		} else {
			atyp = shuttle.AddrTypeIPv6
		}
		addr = []byte(ip)
	}
	if len(addr) == 0 {
		return nil, fmt.Errorf("addr error [%s]", req.Host())
	}
	port, err := strconv.ParseUint(req.Port(), 10, 16)
	if err != nil {
		return nil, err
	}
	portBytes := []byte{byte(port >> 8), byte(port & 0xff)}
	buffer := bytes.NewBuffer([]byte{})
	switch atyp {
	case shuttle.AddrTypeIPv4, shuttle.AddrTypeIPv6:
		buffer.WriteByte(atyp)
		buffer.Write(addr)
		buffer.Write(portBytes)
	case shuttle.AddrTypeDomain:
		buffer.WriteByte(atyp)
		buffer.WriteByte(byte(len(addr)))
		buffer.Write(addr)
		buffer.Write(portBytes)
	default:
	}
	return buffer.Bytes(), nil
}
