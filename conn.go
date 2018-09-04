package shuttle

import (
	"net"
	"errors"
)

const (
	TCP = "tcp"
	UDP = "udp"
)

type IConn interface {
	net.Conn
	GetID() int64
	GetRecordID() int64
	SetRecordID(id int64)
	GetNetwork() string
	Flush() (int, error)
}

func NewDefaultConn(conn net.Conn, network string) (IConn, error) {
	c, err := DefaultDecorate(conn, network)
	return c, err
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
	conn, err := net.DialTimeout(req.Network(), req.Host(), DefaultTimeOut)
	if err != nil {
		return nil, err
	}
	c, err := NewDefaultConn(conn, req.Network())
	if err == nil {
		c, err = TrafficDecorate(c)
	}
	return c, err
}
