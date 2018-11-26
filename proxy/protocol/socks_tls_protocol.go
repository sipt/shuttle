package protocol

import (
	"crypto/tls"
	"fmt"
	connect "github.com/sipt/shuttle/conn"
	"github.com/sipt/shuttle/dns"
	"github.com/sipt/shuttle/log"
	sproxy "github.com/sipt/shuttle/proxy"
	"golang.org/x/net/proxy"
	"net"
)

const (
	ConfigSocksTLSSkipVerify = "skip-verify"
)

func init() {
	sproxy.RegisterProxyProtocolCreator("socks-tls", NewSocks5TLSProtocol)
}

func NewSocks5TLSProtocol(params []string) (sproxy.IProtocol, error) {
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
//	Conn(request shuttle.SocksRequest) (shuttle.IConn, error)
//}
type socksTLSProtocol struct {
	Addr               string
	Port               string
	UserName           string
	Password           string
	InsecureSkipVerify bool
}

func (s *socksTLSProtocol) Conn(req sproxy.IRequest) (connect.IConn, error) {
	var auth *proxy.Auth
	if len(s.UserName) > 0 {
		auth = &proxy.Auth{User: s.UserName, Password: s.Password}
	}

	var addr = s.Addr
	answer, err := dns.ResolveDomainByCache(s.Addr)
	if err != nil {
		log.Logger.Errorf("[SocksOverTlsProtocol] [Conn] Resolve domain failed [%s]: %v", s.Addr, err)
	} else if answer != nil {
		addr = answer.GetIP()
	}
	dialer, err := proxy.SOCKS5(req.Network(), net.JoinHostPort(addr, s.Port), auth, s)
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

func (s *socksTLSProtocol) Dial(network, addr string) (c net.Conn, err error) {
	return tls.Dial(network, addr, &tls.Config{
		InsecureSkipVerify: s.InsecureSkipVerify,
		ServerName:         s.Addr,
	})
}
