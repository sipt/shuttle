package protocol

import (
	"net"
	"github.com/sipt/shuttle"
	"golang.org/x/net/proxy"
	"github.com/sipt/shuttle/log"
	"crypto/tls"
	"fmt"
)

const (
	ConfigSocksTLSSkipVerify = "skip-verify"
)

func init() {
	shuttle.RegisterProxyProtocolCreator("socks-tls", NewSocks5TLSProtocol)
}

func NewSocks5TLSProtocol(params []string) (shuttle.IProtocol, error) {
	//[]string{"addr", "port", "skip-verify","username", "password"}
	if len(params) != 5 && len(params) != 3 {
		log.Logger.Errorf(`[SOCKS5 over TLS Server] init socks5 server failed params count must be 5 or 3, but: %v`, params)
		return nil, fmt.Errorf(`[SOCKS5 over TLS Server] init socks5 server failed params count must be 5 or 3, but: %v`, params)
	}
	ser := &socksTLSProtocol{
		Addr:               params[0],
		Port:               params[1],
		InsecureSkipVerify: params[2] == ConfigSocksTLSSkipVerify,
	}
	if len(params) == 5 {
		ser.UserName = params[3]
		ser.Password = params[4]
	}
	return ser, nil
}

//implement protocol.IServer
//type IServer interface {
//	//获取服务器连接
//	Conn(request shuttle.Request) (shuttle.IConn, error)
//}
type socksTLSProtocol struct {
	Addr               string
	Port               string
	UserName           string
	Password           string
	InsecureSkipVerify bool
}

func (s *socksTLSProtocol) Conn(request *shuttle.Request) (shuttle.IConn, error) {
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
	dialer, err := proxy.SOCKS5(request.Network(), net.JoinHostPort(addr, s.Port), auth, s)
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

func (s *socksTLSProtocol) Dial(network, addr string) (c net.Conn, err error) {
	return tls.Dial(network, addr, &tls.Config{
		InsecureSkipVerify: s.InsecureSkipVerify,
	})
}
