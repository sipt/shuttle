package protocol

import (
	"net"
	"github.com/sipt/shuttle"
	"strconv"
	"bytes"
	"fmt"
	"github.com/sipt/shuttle/ciphers"
)

func init() {
	shuttle.RegisterProxyProtocolCreator("ss", NewSsProtocol)
}

func NewSsProtocol(params []string) (shuttle.IProtocol, error) {
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

func (s *ssProtocol) Conn(req *shuttle.Request) (shuttle.IConn, error) {
	network := req.Network()
	addr := s.Addr
	ssReq := &shuttle.Request{
		Addr: s.Addr,
	}
	err := shuttle.ResolveDomain(ssReq)
	if err != nil {
		log.Logger.Errorf("[SsProtocol] [Conn] Resolve domain failed [%s]: %v", s.Addr, err)
	} else {
		addr = ssReq.IP.String()
	}
	conn, err := net.DialTimeout(network, net.JoinHostPort(addr, s.Port), shuttle.DefaultTimeOut)
	if err != nil {
		return nil, err
	}
	c, err := shuttle.DefaultDecorate(conn, network)
	if err != nil {
		return nil, err
	}
	c, err = shuttle.TrafficDecorate(c)
	if err != nil {
		return nil, err
	}
	if network == shuttle.UDP {
		c, err = shuttle.BufferDecorate(c)
		if err != nil {
			return nil, err
		}
	}
	rc, err := ciphers.CipherDecorate(s.Password, s.Method, c)
	if err != nil {
		return nil, err
	}
	var addrBytes []byte
	if len(req.Addr) > 0 {
		addrBytes = []byte(req.Addr)
	} else {
		addrBytes = req.IP
	}
	rawAddr, err := AddressEncoding(req.Atyp, addrBytes, req.Port)
	if err != nil {
		return nil, err
	}
	_, err = rc.Write(rawAddr)
	if err != nil {
		return nil, err
	}
	return rc, nil
}

func DomainEncodeing(host string) ([]byte, error) {
	domain, port, err := net.SplitHostPort(host)
	if err != nil {
		return nil, err
	}
	p, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return nil, err
	}
	return AddressEncoding(shuttle.AddrTypeDomain, []byte(domain), uint16(p))
}

func AddressEncoding(atyp uint8, addr []byte, port uint16) ([]byte, error) {
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
