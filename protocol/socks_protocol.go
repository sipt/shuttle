package protocol

import (
	"net"
	"github.com/sipt/shuttle"
	"golang.org/x/net/proxy"
	"fmt"
)

func init() {
	shuttle.RegisterProxyProtocolCreator("socks", NewSocks5Protocol)
}

func NewSocks5Protocol(params []string) (shuttle.IProtocol, error) {
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
//	Conn(request shuttle.Request) (shuttle.IConn, error)
//}
type socksProtocol struct {
	Addr     string
	Port     string
	UserName string
	Password string
}

func (s *socksProtocol) Conn(request *shuttle.Request) (shuttle.IConn, error) {
	var auth *proxy.Auth
	if len(s.UserName) > 0 {
		auth = &proxy.Auth{User: s.UserName, Password: s.Password}
	}
	addr := s.Addr
	ssReq := &shuttle.Request{
		Addr: s.Addr,
	}
	err := shuttle.ResolveDomain(ssReq)
	if err != nil {
		log.Logger.Errorf("[SocksProtocol] [Conn] Resolve domain failed [%s]: %v", s.Addr, err)
	} else {
		addr = ssReq.IP.String()
	}
	dialer, err := proxy.SOCKS5(request.Network(), net.JoinHostPort(addr, s.Port), auth, nil)
	if err != nil {
		return nil, err
	}
	conn, err := dialer.Dial(request.Network(), request.Host2())
	if err != nil {
		return nil, err
	}
	c, err := shuttle.DefaultDecorate(conn, request.Network())
	if err != nil {
		return nil, err
	}
	return shuttle.TrafficDecorate(c)
}
