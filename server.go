package shuttle

import (
	"fmt"
	"time"
)

var groups []*ServerGroup
var servers []*Server

func InitServers(gs []*ServerGroup, ss []*Server) error {
	g := &ServerGroup{
		Name:       PolicyGlobal,
		SelectType: "select",
		Servers:    make([]interface{}, len(gs)+len(ss)),
	}
	index := 0
	for i := range ss {
		if ss[i].Name != PolicyDirect && ss[i].Name != PolicyReject {
			g.Servers[index] = ss[i]
			index ++
		}
	}
	for i := range gs {
		g.Servers[index] = gs[i]
		index++
	}
	g.Servers = g.Servers[:index]
	gs = append(gs, g)
	var err error
	for i, v := range gs {
		v.Selector, err = seletors[gs[i].SelectType](v)
		if err != nil {
			return err
		}
	}
	groups = gs
	servers = ss
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

var proxyProtocolCreator = make(map[string]NewProtocol)

func RegisterProxyProtocolCreator(name string, p NewProtocol) {
	proxyProtocolCreator[name] = p
	Logger.Infof("Support Proxy Protocol: [%s]", name)
}

type IServer interface {
	GetName() string
	GetServer() (*Server, error)
}

type NewProtocol func([]string) (IProtocol, error)

type IProtocol interface {
	//获取服务器连接
	Conn(request *Request) (IConn, error)
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

//创建Server
func NewServer(name string, params []string) (*Server, error) {
	if len(params) < 1 {
		return nil, fmt.Errorf("[Config] [InitServer] Invalid format: %v", params)
	}
	ser := &Server{
		Name:          name,
		ProxyProtocol: params[0],
	}
	n := proxyProtocolCreator[ser.ProxyProtocol]
	if n == nil {
		return nil, fmt.Errorf("[Config] [InitServer] Not support protocol: %s", ser.ProxyProtocol)
	}
	var err error
	ser.IProtocol, err = n(params[1:])
	return ser, err
}

type Server struct {
	Name          string
	Rtt           time.Duration
	ProxyProtocol string
	IProtocol `json:"-"`
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
	return s.IProtocol.Conn(req)
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
