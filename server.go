package shuttle

import (
	"net"
)

var groups []*ServerGroup
var servers []*Server

func InitServers(gs []*ServerGroup, ss []*Server) error {
	groups = gs
	servers = ss
	servers = make([]*Server, 0, 10)
	var err error
	for i, v := range gs {
		v.Selector, err = seletors[gs[i].SelectType](v)
		if err != nil {
			return err
		}
	}
	return nil
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
	Port     string
	Method   string
	Password string
}

func (s *Server) GetName() string {
	return s.Name
}
func (s *Server) GetServer() (*Server, error) {
	return s, nil
}

func (s *Server) Conn(network string) (IConn, error) {
	conn, err := net.DialTimeout(network, net.JoinHostPort(s.Host, s.Port), defaultTimeOut)
	if err != nil {
		return nil, err
	}
	c, err := NewDefaultConn(conn, network)
	if err != nil {
		return nil, err
	}
	c, err = CipherDecorate(s.Password, s.Method, c)
	if err != nil {
		return nil, err
	}
	if network == UDP {
		c, err = BufferDecorate(c)
	}
	return c, err
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
