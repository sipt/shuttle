package enhance

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/sipt/shuttle/pkg/dns"
	mockdns "github.com/sipt/shuttle/pkg/dns/mock"
	"github.com/sipt/shuttle/pkg/tun"
	"github.com/sirupsen/logrus"
)

type EnhanceModeState int

const (
	EnhanceModeStateInit EnhanceModeState = iota
	EnhanceModeStateRunning
	EnhanceModeStateStopped
)

type EnhanceMode struct {
	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	state     EnhanceModeState
	tunDevice *tun.TunDevice
	dnsServer *dns.DNSServer
}

func NewEnhanceMode() *EnhanceMode {
	return &EnhanceMode{state: EnhanceModeStateInit}
}

func (e *EnhanceMode) Start() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	ctx, cancel := context.WithCancel(context.Background())
	e.ctx = ctx
	e.cancel = cancel
	e.state = EnhanceModeStateRunning

	tunDevice, err := tun.OpenTun(ctx)
	if err != nil {
		return err
	}
	e.tunDevice = tunDevice

	// dns mock server
	dnsHandler := mockdns.NewMockHandle()
	e.dnsServer = dns.NewDNSServer(dnsHandler)

	go e.handleUdp()
	go e.handleTcp()

	return nil
}

func (e *EnhanceMode) Stop() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.state = EnhanceModeStateStopped
	if e.cancel != nil {
		e.cancel()
	}
	if e.tunDevice != nil {
		e.tunDevice.Close()
	}
	return nil
}

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

func (e *EnhanceMode) handleTcp() {
	for {
		conn, err := e.tunDevice.Listener.TcpListener.Accept()
		if err != nil {
			return
		}
		conn.Close()
	}
}
