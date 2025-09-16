package enhance

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	connpkg "github.com/sipt/shuttle/conn"
	"github.com/sipt/shuttle/constant"
	"github.com/sipt/shuttle/handle"
	"github.com/sipt/shuttle/inbound"
	"github.com/sipt/shuttle/pkg/tun"
	"github.com/sirupsen/logrus"
)

type tcpRequest struct {
	id          int64
	network     string
	domain      string
	uri         string
	ip          net.IP
	port        int
	countryCode string
}

func (r *tcpRequest) ID() int64 {
	if r.id == 0 {
		r.id = inbound.GetRequestID()
	}
	return r.id
}
func (r *tcpRequest) Network() string {
	return r.network
}
func (r *tcpRequest) Domain() string {
	return r.domain
}
func (r *tcpRequest) URI() string {
	if r.uri != "" {
		return r.uri
	}
	uri := ""
	if r.domain != "" {
		uri = r.domain
	} else if r.ip != nil {
		uri = r.ip.String()
	}
	if r.port > 0 {
		return fmt.Sprintf("%s:%d", uri, r.port)
	}
	return uri
}
func (r *tcpRequest) IP() net.IP {
	return r.ip
}
func (r *tcpRequest) CountryCode() string {
	return r.countryCode
}
func (r *tcpRequest) Port() int {
	return r.port
}
func (r *tcpRequest) SetIP(in net.IP) {
	r.ip = in
}
func (r *tcpRequest) SetPort(in int) {
	r.port = in
}
func (r *tcpRequest) SetCountryCode(in string) {
	r.countryCode = in
}

// ============================
// TCP
// ============================
func (e *EnhanceMode) handleTcp() {
	for {
		conn, err := e.tunDevice.Listener.TcpListener.Accept()
		if err != nil {
			return
		}
		go e.handleTcpConn(conn)
	}
}

func (e *EnhanceMode) handleTcpConn(conn tun.TcpConn) {
	defer conn.Close()
	logger := logrus.WithField("mode", "enhance").WithField("method", "tcp")
	logger.Debugf("handleTcpConn: %v", conn.LocalAddr())
	req := &tcpRequest{
		network: conn.RemoteAddr().Network(),
	}
	var localIP net.IP
	if tcpAddr, ok := conn.LocalAddr().(*net.TCPAddr); ok {
		localIP = tcpAddr.IP
		req.port = tcpAddr.Port
	} else if items := strings.Split(conn.LocalAddr().String(), ":"); len(items) == 2 {
		localIP = net.ParseIP(items[0])
		req.port, _ = strconv.Atoi(items[1])
	} else {
		logger.Errorf("invalid local address: %v", conn.LocalAddr())
		return
	}

	e.dnsHandler.ReverseLookup(localIP)
	if domain, ok := e.dnsHandler.ReverseLookup(localIP); ok {
		req.domain = domain
	}

	ctx := context.WithValue(e.ctx, constant.KeyRequestInfo, req)
	ctx = context.WithValue(ctx, constant.KeyProtocol, inbound.ProtocolTCP)
	handle := handle.Handle()
	handle(connpkg.NewConn(conn, ctx))
}
