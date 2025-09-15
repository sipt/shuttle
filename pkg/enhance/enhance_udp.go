package enhance

import (
	"strings"
	"time"

	"github.com/sipt/shuttle/pkg/tun"
	"github.com/sirupsen/logrus"
)

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
