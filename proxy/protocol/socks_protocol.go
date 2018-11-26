package protocol

import (
	"fmt"
	"net"

	connect "github.com/sipt/shuttle/conn"
	"github.com/sipt/shuttle/dns"
	"github.com/sipt/shuttle/log"
	sproxy "github.com/sipt/shuttle/proxy"

	"golang.org/x/net/proxy"
)

func init() {
	sproxy.RegisterProxyProtocolCreator("socks", NewSocks5Protocol)
}

func NewSocks5Protocol(params []string) (sproxy.IProtocol, error) {
	//[]string{"addr", "port", "username", "password"}
	if len(params) != 4 && len(params) != 2 {
		log.Logger.Errorf(`[SOCKS5 Server] init socks5 server failed params must be ["addr", "port"] or ["addr", "port", "username", "password"], but: %v`, params)
		return nil, fmt.Errorf(`[SOCKS5 Server] init socks5 server failed params must be ["addr", "port"] or ["addr", "port", "username", "password"], but: %v`, params)
	}
	ser := &socksProtocol{
		Addr: params[0],
		Port: params[1],
	}
	if len(params) == 4 {
		ser.UserName = params[2]
		ser.Password = params[3]
	}
	return ser, nil
}

//implement protocol.IServer
//type IServer interface {
//	//获取服务器连接
//	Conn(request shuttle.SocksRequest) (shuttle.IConn, error)
//}
type socksProtocol struct {
	Addr     string
	Port     string
	UserName string
	Password string
}

func (s *socksProtocol) Conn(req sproxy.IRequest) (connect.IConn, error) {
	var auth *proxy.Auth
	if len(s.UserName) > 0 {
		auth = &proxy.Auth{User: s.UserName, Password: s.Password}
	}
	var addr = s.Addr
	answer, err := dns.ResolveDomainByCache(s.Addr)
	if err != nil {
		log.Logger.Errorf("[SocksProtocol] [Conn] Resolve domain failed [%s]: %v", s.Addr, err)
	} else if answer != nil {
		addr = answer.GetIP()
	}
	dialer, err := proxy.SOCKS5(req.Network(), net.JoinHostPort(addr, s.Port), auth, nil)
	if err != nil {
		return nil, err
	}
	addr = req.IP()
	if addr == "" {
		addr = req.Domain()
	}
	addr = net.JoinHostPort(addr, req.Port())
	conn, err := dialer.Dial(req.Network(), addr)
	if err != nil {
		return nil, err
	}
	c, err := connect.DefaultDecorate(conn, req.Network())
	if err != nil {
		return nil, err
	}
	return connect.TrafficDecorate(c)

}
