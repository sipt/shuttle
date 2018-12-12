package proxy

import (
	"errors"
	"fmt"
	"github.com/sipt/shuttle/conn"
	"github.com/sipt/shuttle/log"
	"time"
)

const (
	ProxyDirect = "DIRECT"
	ProxyReject = "REJECT"
	ProxyGlobal = "GLOBAL"
)

var (
	groups  []*ServerGroup
	servers []*Server

	ErrorReject         = errors.New("connection reject")
	ErrorServerNotFound = errors.New("server or server group not found")

	MockServer   = &Server{Name: "MOCK"}
	FailedServer = &Server{Name: "FAILED"}
	RejectServer = &Server{Name: "REJECT"}
)

type IProxyConfig interface {
	GetProxy() map[string][]string
	SetProxy(map[string][]string)
	GetProxyGroup() map[string][]string
	SetProxyGroup(map[string][]string)
}

type IRequest interface {
	Network() string
	Domain() string
	IP() string
	Port() string
	Host() string
}

func ApplyConfig(config IProxyConfig) (err error) {
	proxy := config.GetProxy()
	//Servers
	ss := make([]*Server, len(proxy)+2)
	index := 0
	ss[index] = &Server{Name: ProxyDirect} // 直连
	index ++
	ss[index] = &Server{Name: ProxyReject} // 拒绝
	for k, v := range proxy {
		index ++
		if len(v) < 2 {
			return fmt.Errorf("resolve config file [proxy] [%s] failed", k)
		}
		ss[index], err = NewServer(k, v)
		if err != nil {
			return
		}
	}

	proxyGroup := config.GetProxyGroup()
	gs := make([]*ServerGroup, len(proxyGroup))
	index = 0
	for k := range proxyGroup {
		gs[index] = &ServerGroup{Name: k}
		index ++
	}
	getServer := func(name string) interface{} {
		for i := range ss {
			if ss[i].Name == name {
				return ss[i]
			}
		}
		for i := range gs {
			if gs[i].Name == name {
				return gs[i]
			}
		}
		return nil
	}
	var cs []string
	for _, v := range gs {
		cs = proxyGroup[v.Name]
		if len(cs) < 2 {
			return fmt.Errorf("resolve config file [proxy_group] [%s] failed", v.Name)
		}
		v.SelectType = cs[0]
		v.Servers = make([]interface{}, len(cs)-1)
		for i := range v.Servers {
			v.Servers[i] = getServer(cs[i+1])
			if v.Servers[i] == nil {
				return fmt.Errorf("resolve config file [proxy_group] [%s] [%s] not found", v.Name, cs[i+1])
			}
		}
	}
	err = InitServers(gs, ss)
	if err != nil {
		return fmt.Errorf("init server failed: %v", err)
	}
	return nil
}

func InitServers(gs []*ServerGroup, ss []*Server) error {
	g := &ServerGroup{
		Name:       ProxyGlobal,
		SelectType: "select",
		Servers:    make([]interface{}, len(gs)+len(ss)),
	}
	index := 0
	for i := range ss {
		if ss[i].Name != ProxyDirect && ss[i].Name != ProxyReject {
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
		v.Selector, err = GetSelector(gs[i].SelectType, v)
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
	log.Logger.Infof("Support Proxy Protocol: [%s]", name)
}

type IServer interface {
	GetName() string
	GetServer() (*Server, error)
}

type NewProtocol func([]string) (IProtocol, error)

type IProtocol interface {
	//获取服务器连接
	Conn(request IRequest) (conn.IConn, error)
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
	IProtocol     `json:"-"`
}

func (s *Server) GetName() string {
	return s.Name
}
func (s *Server) GetServer() (*Server, error) {
	return s, nil
}

func (s *Server) Conn(req IRequest) (conn.IConn, error) {
	switch s.Name {
	case ProxyDirect:
		return conn.DirectConn(req.Network(), req.Host())
	case ProxyReject:
		return nil, ErrorReject
	}
	return s.IProtocol.Conn(req)
}

func GetServer(name string) (*Server, error) {
	if name == "REJECT" {
		return RejectServer, nil
	}
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
	return FailedServer, ErrorServerNotFound
}
