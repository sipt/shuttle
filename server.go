package shuttle

import (
	"net"
	"fmt"
	"time"
)

var groups []*ServerGroup
var servers []*Server

func InitServers(gs []*ServerGroup, ss []*Server) error {
	groups = gs
	servers = ss
	var err error
	for i, v := range gs {
		v.Selector, err = seletors[gs[i].SelectType](v)
		if err != nil {
			return err
		}
	}
	return nil
}
func DestroyServers() {
	for _, v := range groups {
		v.Selector.Destroy()
	}
}
func GetGroups() []*ServerGroup {
	return groups
}
func SelectServer(groupName, serverName string) error {
	for _, g := range groups {
		if g.Name == groupName {
			return g.Selector.Select(serverName)
		}
	}
	return fmt.Errorf("group[%s] is not exist", groupName)
}
func SelectRefresh(groupName string) error {
	for _, g := range groups {
		if g.Name == groupName {
			return g.Selector.Refresh()
		}
	}
	return fmt.Errorf("group[%s] is not exist", groupName)
}

type IServer interface {
	GetName() string
	GetServer() (*Server, error)
}

type ServerGroup struct {
	Servers    []interface{}
	Name       string
	SelectType string
	Selector   ISelector
}

func (s *ServerGroup) GetName() string {
	return s.Name
}

func (s *ServerGroup) GetServer() (*Server, error) {
	return s.Selector.Get()
}

type Server struct {
	Name     string
	Host     string
	Port     string `json:"-"`
	Method   string `json:"-"`
	Password string `json:"-"`
	Rtt      time.Duration
}

func (s *Server) GetName() string {
	return s.Name
}
func (s *Server) GetServer() (*Server, error) {
	return s, nil
}

func (s *Server) Conn(req *Request) (IConn, error) {
	switch s.Name {
	case PolicyDirect:
		return DirectConn(req)
	case PolicyReject:
		return nil, ErrorReject
	}
	var network string
	switch req.Cmd {
	case CmdTCP:
		network = "tcp"
	case CmdUDP:
		network = "udp"
	}
	ssReq := &Request{
		Addr: s.Host,
	}
	addr := s.Host
	err := ResolveDomain(ssReq)
	if err != nil {
		Logger.Errorf("Resolve domain failed [%s]: %v", s.Host, err)
	} else {
		addr = ssReq.IP.String()
	}
	conn, err := net.DialTimeout(network, net.JoinHostPort(addr, s.Port), defaultTimeOut)
	if err != nil {
		return nil, err
	}
	c, err := NewDefaultConn(conn, network)
	if err != nil {
		return nil, err
	}
	if network == UDP {
		c, err = BufferDecorate(c)
		if err != nil {
			return nil, err
		}
	}
	rc, err := CipherDecorate(s.Password, s.Method, c)
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

func GetServer(name string) (*Server, error) {
	for _, v := range groups {
		if v.Name == name {
			return v.Selector.Get()
		}
	}
	for i := range servers {
		if servers[i].Name == name {
			return servers[i], nil
		}
	}
	return nil, ErrorServerNotFound
}
