package shuttle

import (
	"net"
	"bytes"
	"fmt"
)

const (
	TCP = "tcp"
	UDP = "udp"
)

func AddressEncoding(atyp uint8, addr []byte, port uint16) ([]byte, error) {
	portBytes := []byte{byte(port >> 8), byte(port & 0xff)}
	//binary.LittleEndian.PutUint16(portBytes, port)
	buffer := bytes.NewBuffer([]byte{})
	switch atyp {
	case addrTypeIPv4, addrTypeIPv6:
		buffer.WriteByte(atyp)
		buffer.Write(addr)
		buffer.Write(portBytes)
	case addrTypeDomain:
		buffer.WriteByte(atyp)
		buffer.WriteByte(byte(len(addr)))
		buffer.Write(addr)
		buffer.Write(portBytes)
	default:
	}
	fmt.Println(buffer.Bytes())
	return buffer.Bytes(), nil
}

type IConn interface {
	net.Conn
	GetID() int64
	GetNetwork() string
	Flush() (int, error)
}

func NewDefaultConn(conn net.Conn, network string) (IConn, error) {
	c, err := DefaultDecorate(conn, network)
	if err != nil {
		return nil, err
	}
	return TimerDecorate(c, 0)
}

func ConnectToServer(req *Request) (IConn, error) {
	//DNS
	err := ResolveDomain(req)
	if err != nil {
		return nil, err
	}
	//Rules filter
	rule, err := Filter(req)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		Logger.Debugf("[%s] rule: [%v]", req.Host(), PolicyDirect)
		return DirectConn(req) // 没有匹配规则，直连
	}
	Logger.Debugf("[%s] rule: [%v]", req.Host(), rule)
	//Select proxy server
	s, err := GetServer(rule.Policy)
	if err != nil {
		return nil, err
	}
	Logger.Debugf("get server by policy [%s] => %v", rule.Policy, s)
	switch s.Name {
	case PolicyDirect:
		return DirectConn(req)
	case PolicyReject:
		return nil, ErrorReject
	default:
		sc, err := s.Conn(req.Network())
		if err != nil {
			return nil, err
		}
		var host []byte
		if len(req.IP) == 0 {
			host = []byte(req.Addr)
		} else {
			host = req.IP
		}
		addr, err := AddressEncoding(req.Atyp, host, req.Port)
		if err != nil {
			return nil, err
		}
		_, err = sc.Write(addr)
		if err != nil {
			return nil, err
		}
		return sc, nil
	}
}

func DirectConn(req *Request) (IConn, error) {
	c, err := net.DialTimeout(req.Network(), req.Host(), defaultTimeOut)
	if err != nil {
		return nil, err
	}
	return NewDefaultConn(c, req.Network())
}
