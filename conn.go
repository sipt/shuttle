package shuttle

import (
	"net"
	"bytes"
	"strconv"
	"errors"
	"time"
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
	return AddressEncoding(addrTypeDomain, []byte(domain), uint16(p))
}

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

func ConnectToServer(req *Request) (rconn IConn, err error) {
	//DNS
	if len(req.IP) == 0 {
		req.IP = net.ParseIP(req.Addr)
		if len(req.IP) == 0 {
			err := ResolveDomain(req)
			if err != nil {
				return nil, err
			}
		}
	}
	//Rules filter
	rule, err := Filter(req)
	if err != nil {
		return nil, err
	}
	var s *Server
	if rule == nil {
		Logger.Debugf("[%s] rule: [%v]", req.Host(), PolicyDirect)
		s, err = GetServer(PolicyDirect) // 没有匹配规则，直连
	} else {
		Logger.Debugf("[RULE] [%s, %s, %s] rule: [%s,%s,%s]", req.Host(), req.Addr, req.DomainHost.Country, rule.Type, rule.Value, rule.Policy)
		//Select proxy server
		s, err = GetServer(rule.Policy)
		if err != nil {
			return nil, errors.New(err.Error() + ":" + rule.Policy)
		}
		Logger.Debugf("get server by policy [%s] => %v", rule.Policy, s)
	}
	switch s.Name {
	case PolicyDirect:
		rconn, err = DirectConn(req)
	case PolicyReject:
		return nil, ErrorReject
	default:
		sc, err := s.Conn(req.Network())
		if err != nil {
			return nil, err
		}
		addr, err := AddressEncoding(req.Atyp, []byte(req.Addr), req.Port)
		if err != nil {
			return nil, err
		}
		_, err = sc.Write(addr)
		if err != nil {
			return nil, err
		}
		rconn, err = sc, nil
	}
	recordChan <- &Record{
		ID:       rconn.GetID(),
		Protocol: req.Protocol,
		Created:  time.Now(),
		Proxy:    s,
		Status:   RecordStatusActive,
		URL:      req.Target,
		Rule:     rule,
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
