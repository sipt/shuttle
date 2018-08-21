package shuttle

import (
	"net"
	"bytes"
	"strconv"
	"errors"
)

const (
	TCP = "tcp"
	UDP = "udp"
)

func DomainEncodeing(host string) ([]byte, error) {
	domain, port, err := net.SplitHostPort(host)
	if err != nil {
		return nil, err
	}
	p, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return nil, err
	}
	return AddressEncoding(AddrTypeDomain, []byte(domain), uint16(p))
}

func AddressEncoding(atyp uint8, addr []byte, port uint16) ([]byte, error) {
	portBytes := []byte{byte(port >> 8), byte(port & 0xff)}
	//binary.LittleEndian.PutUint16(portBytes, port)
	buffer := bytes.NewBuffer([]byte{})
	switch atyp {
	case AddrTypeIPv4, AddrTypeIPv6:
		buffer.WriteByte(atyp)
		buffer.Write(addr)
		buffer.Write(portBytes)
	case AddrTypeDomain:
		buffer.WriteByte(atyp)
		buffer.WriteByte(byte(len(addr)))
		buffer.Write(addr)
		buffer.Write(portBytes)
	default:
	}
	return buffer.Bytes(), nil
}

type IConn interface {
	net.Conn
	GetID() int64
	GetNetwork() string
	Flush() (int, error)
}

func NewDefaultConn(conn net.Conn, network string) (IConn, error) {
	return DefaultDecorate(conn, network)
}

func FilterByReq(req *Request) (rule *Rule, s *Server, err error) {
	//DNS
	if len(req.IP) == 0 {
		req.IP = net.ParseIP(req.Addr)
		if len(req.IP) == 0 {
			err = ResolveDomain(req)
			if err != nil {
				return
			}
		} else if len(req.DomainHost.Country) == 0 {
			req.DomainHost.Country = GeoLookUp(req.IP)
		}
	} else if len(req.DomainHost.Country) == 0 {
		req.DomainHost.Country = GeoLookUp(req.IP)
	}
	//Rules filter
	rule, err = filter(req)
	if err != nil {
		return
	}
	if rule == nil {
		Logger.Debugf("[%s] rule: [%v]", req.Host(), PolicyDirect)
		s, err = GetServer(PolicyDirect) // 没有匹配规则，直连
	} else {
		Logger.Debugf("[RULE] [%s, %s, %s] rule: [%s,%s,%s]", req.Host(), req.Addr, req.DomainHost.Country, rule.Type, rule.Value, rule.Policy)
		//Select proxy server
		s, err = GetServer(rule.Policy)
		if err != nil {
			err = errors.New(err.Error() + ":" + rule.Policy)
			return
		}
		Logger.Debugf("get server by policy [%s] => %v", rule.Policy, s)
	}
	return
}

func DirectConn(req *Request) (IConn, error) {
	c, err := net.DialTimeout(req.Network(), req.Host(), defaultTimeOut)
	if err != nil {
		return nil, err
	}
	return NewDefaultConn(c, req.Network())
}
