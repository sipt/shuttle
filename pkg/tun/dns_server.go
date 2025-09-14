package tun

import (
	"net"
	"sync"
	"time"

	"github.com/sipt/shuttle/pkg/dns"
	"github.com/sirupsen/logrus"
)

// DNSServerManager manages DNS server functionality for TUN interface
type DNSServerManager struct {
	dnsServer   *dns.DNSServer
	logger      *logrus.Logger
	running     bool
	mu          sync.RWMutex
	connections map[string]*dnsConnection
}

// dnsConnection represents a DNS connection
type dnsConnection struct {
	conn       UdpConn
	packetConn *dnsPacketConn
	createdAt  time.Time
}

// NewDNSServerManager creates a new DNS server manager
func NewDNSServerManager(upstreamServers []string) *DNSServerManager {
	handler := dns.NewDefaultDNSHandler(upstreamServers)
	dnsServer := dns.NewDNSServer(handler)

	return &DNSServerManager{
		dnsServer:   dnsServer,
		logger:      logrus.New(),
		connections: make(map[string]*dnsConnection),
	}
}

// Start starts the DNS server manager
func (m *DNSServerManager) Start() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return
	}

	m.running = true
	m.logger.Info("DNS server manager started")

	// Start cleanup goroutine for stale connections
	go m.cleanupConnections()
}

// Stop stops the DNS server manager
func (m *DNSServerManager) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		return
	}

	m.running = false

	// Close all connections
	for _, conn := range m.connections {
		conn.conn.Close()
	}
	m.connections = make(map[string]*dnsConnection)

	m.logger.Info("DNS server manager stopped")
}

// HandleDNSConnection handles a DNS UDP connection
func (m *DNSServerManager) HandleDNSConnection(conn UdpConn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		conn.Close()
		return
	}

	connKey := conn.RemoteAddr().String()

	// Check if connection already exists
	if existing, exists := m.connections[connKey]; exists {
		existing.conn.Close()
	}

	// Create packet connection wrapper
	packetConn := &dnsPacketConn{
		conn:       conn,
		localAddr:  conn.LocalAddr(),
		remoteAddr: conn.RemoteAddr(),
	}

	dnsConn := &dnsConnection{
		conn:       conn,
		packetConn: packetConn,
		createdAt:  time.Now(),
	}

	m.connections[connKey] = dnsConn

	m.logger.WithFields(logrus.Fields{
		"local":  conn.LocalAddr(),
		"remote": conn.RemoteAddr(),
	}).Info("Handling DNS connection")

	// Handle the DNS connection in a separate goroutine
	go func() {
		defer func() {
			m.mu.Lock()
			delete(m.connections, connKey)
			m.mu.Unlock()
			conn.Close()
		}()

		if err := m.dnsServer.ServeDNSOverUDP(packetConn); err != nil {
			m.logger.WithError(err).Error("DNS server error")
		}
	}()
}

// cleanupConnections periodically cleans up stale connections
func (m *DNSServerManager) cleanupConnections() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.cleanupStaleConnections()
		}

		m.mu.RLock()
		running := m.running
		m.mu.RUnlock()

		if !running {
			return
		}
	}
}

// cleanupStaleConnections removes connections older than 5 minutes
func (m *DNSServerManager) cleanupStaleConnections() {
	m.mu.Lock()
	defer m.mu.Unlock()

	cutoff := time.Now().Add(-5 * time.Minute)
	var toDelete []string

	for key, conn := range m.connections {
		if conn.createdAt.Before(cutoff) {
			toDelete = append(toDelete, key)
			conn.conn.Close()
		}
	}

	for _, key := range toDelete {
		delete(m.connections, key)
	}

	if len(toDelete) > 0 {
		m.logger.WithField("count", len(toDelete)).Info("Cleaned up stale DNS connections")
	}
}

// dnsPacketConn implements net.PacketConn for DNS server
type dnsPacketConn struct {
	conn       UdpConn
	localAddr  net.Addr
	remoteAddr net.Addr
}

func (c *dnsPacketConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	n, err = c.conn.Read(p)
	if err != nil {
		return 0, nil, err
	}
	return n, c.remoteAddr, nil
}

func (c *dnsPacketConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	return c.conn.Write(p)
}

func (c *dnsPacketConn) Close() error {
	return c.conn.Close()
}

func (c *dnsPacketConn) LocalAddr() net.Addr {
	return c.localAddr
}

func (c *dnsPacketConn) SetDeadline(t time.Time) error {
	// Not implemented for this use case
	return nil
}

func (c *dnsPacketConn) SetReadDeadline(t time.Time) error {
	// Not implemented for this use case
	return nil
}

func (c *dnsPacketConn) SetWriteDeadline(t time.Time) error {
	// Not implemented for this use case
	return nil
}
