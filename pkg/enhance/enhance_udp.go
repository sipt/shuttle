package enhance

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	connpkg "github.com/sipt/shuttle/conn"
	"github.com/sipt/shuttle/constant"
	"github.com/sipt/shuttle/inbound"
	"github.com/sipt/shuttle/pkg/tun"
	"github.com/sirupsen/logrus"
)

type udpRequest struct {
	id          int64
	network     string
	domain      string
	uri         string
	ip          net.IP
	port        int
	countryCode string
}

func (r *udpRequest) ID() int64 {
	if r.id == 0 {
		r.id = inbound.GetRequestID()
	}
	return r.id
}
func (r *udpRequest) Network() string {
	return r.network
}
func (r *udpRequest) Domain() string {
	return r.domain
}
func (r *udpRequest) URI() string {
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
func (r *udpRequest) IP() net.IP {
	return r.ip
}
func (r *udpRequest) CountryCode() string {
	return r.countryCode
}
func (r *udpRequest) Port() int {
	return r.port
}
func (r *udpRequest) SetIP(in net.IP) {
	r.ip = in
}
func (r *udpRequest) SetPort(in int) {
	r.port = in
}
func (r *udpRequest) SetCountryCode(in string) {
	r.countryCode = in
}

// ============================
// UDP
// ============================
func (e *EnhanceMode) handleUdp() {
	for {
		conn, err := e.tunDevice.Listener.UdpListener.Accept()
		if err != nil {
			return
		}
		go e.handleUdpConn(conn)
	}
}

func (e *EnhanceMode) handleUdpConn(conn tun.UdpConn) {
	logger := logrus.WithField("mode", "enhance").WithField("method", "udp")
	logger.Debugf("handleUdpConn: %v", conn.LocalAddr())
	defer conn.Close()
	if items := strings.Split(conn.LocalAddr().String(), ":"); len(items) != 2 || items[0] == "198.18.0.2" {
		if items[1] != "53" {
			return
		}
		e.handleDnsPacketConn(conn)
		return
	}

	req := &udpRequest{
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
	if e.handle != nil {
		e.handle(connpkg.NewConn(conn, ctx))
	}
}

func (e *EnhanceMode) handleDnsPacketConn(conn tun.UdpConn) {
	logger := logrus.WithField("method", "mock_dns")
	if e.dnsServer == nil {
		return
	}

	// 设置超时
	conn.SetDeadline(time.Now().Add(time.Second))
	err := e.dnsServer.ServeDNSOverUDP(conn)
	if err != nil {
		logger.Errorf("DNS server error: %v", err)
	}
}
