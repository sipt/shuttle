package dns

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/sirupsen/logrus"
)

// DNSServer represents a custom DNS server
type DNSServer struct {
	handler DNSHandler
	logger  *logrus.Logger
}

// DNSHandler defines the interface for handling DNS queries
type DNSHandler interface {
	HandleQuery(ctx context.Context, domain string, qtype uint16) ([]net.IP, error)
}

// NewDNSServer creates a new DNS server instance
func NewDNSServer(handler DNSHandler) *DNSServer {
	return &DNSServer{
		handler: handler,
		logger:  logrus.New(),
	}
}

// HandleDNSRequest processes incoming DNS requests
func (s *DNSServer) HandleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Authoritative = true

	for _, q := range r.Question {
		s.logger.WithFields(logrus.Fields{
			"domain": q.Name,
			"type":   dns.TypeToString[q.Qtype],
		}).Info("Processing DNS query")

		switch q.Qtype {
		case dns.TypeA:
			s.handleARecord(ctx, msg, q)
		case dns.TypeAAAA:
			s.handleAAAARecord(ctx, msg, q)
		default:
			s.logger.WithField("type", dns.TypeToString[q.Qtype]).Warn("Unsupported query type")
			msg.Rcode = dns.RcodeNotImplemented
		}
	}

	if err := w.WriteMsg(msg); err != nil {
		s.logger.WithError(err).Error("Failed to write DNS response")
	}
}

// handleARecord handles A record queries (IPv4)
func (s *DNSServer) handleARecord(ctx context.Context, msg *dns.Msg, q dns.Question) {
	domain := strings.TrimSuffix(q.Name, ".")

	ips, err := s.handler.HandleQuery(ctx, domain, dns.TypeA)
	if err != nil {
		s.logger.WithError(err).WithField("domain", domain).Error("Failed to resolve domain")
		msg.Rcode = dns.RcodeServerFailure
		return
	}

	if len(ips) == 0 {
		s.logger.WithField("domain", domain).Warn("No A records found")
		msg.Rcode = dns.RcodeNameError
		return
	}

	for _, ip := range ips {
		if ip.To4() != nil { // IPv4 address
			rr := &dns.A{
				Hdr: dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    300, // 5 minutes TTL
				},
				A: ip,
			}
			msg.Answer = append(msg.Answer, rr)
		}
	}

	if len(msg.Answer) == 0 {
		msg.Rcode = dns.RcodeNameError
	}
}

// handleAAAARecord handles AAAA record queries (IPv6)
func (s *DNSServer) handleAAAARecord(ctx context.Context, msg *dns.Msg, q dns.Question) {
	domain := strings.TrimSuffix(q.Name, ".")

	ips, err := s.handler.HandleQuery(ctx, domain, dns.TypeAAAA)
	if err != nil {
		s.logger.WithError(err).WithField("domain", domain).Error("Failed to resolve domain")
		msg.Rcode = dns.RcodeServerFailure
		return
	}

	if len(ips) == 0 {
		s.logger.WithField("domain", domain).Warn("No AAAA records found")
		msg.Rcode = dns.RcodeNameError
		return
	}

	for _, ip := range ips {
		if ip.To4() == nil && ip.To16() != nil { // IPv6 address
			rr := &dns.AAAA{
				Hdr: dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeAAAA,
					Class:  dns.ClassINET,
					Ttl:    300, // 5 minutes TTL
				},
				AAAA: ip,
			}
			msg.Answer = append(msg.Answer, rr)
		}
	}

	if len(msg.Answer) == 0 {
		msg.Rcode = dns.RcodeNameError
	}
}

// DefaultDNSHandler provides a default implementation using existing DNS resolution logic
type DefaultDNSHandler struct {
	upstreamServers []string
}

// NewDefaultDNSHandler creates a new default DNS handler
func NewDefaultDNSHandler(upstreamServers []string) *DefaultDNSHandler {
	if len(upstreamServers) == 0 {
		upstreamServers = []string{"8.8.8.8:53", "1.1.1.1:53"}
	}
	return &DefaultDNSHandler{
		upstreamServers: upstreamServers,
	}
}

// HandleQuery implements DNSHandler interface
func (h *DefaultDNSHandler) HandleQuery(ctx context.Context, domain string, qtype uint16) ([]net.IP, error) {
	client := &dns.Client{
		Timeout: 3 * time.Second,
	}

	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), qtype)
	msg.RecursionDesired = true

	// Try each upstream server
	for _, server := range h.upstreamServers {
		resp, _, err := client.ExchangeContext(ctx, msg, server)
		if err != nil {
			logrus.WithError(err).WithField("server", server).Warn("Failed to query upstream DNS server")
			continue
		}

		if resp.Rcode != dns.RcodeSuccess {
			logrus.WithFields(logrus.Fields{
				"domain": domain,
				"server": server,
				"rcode":  dns.RcodeToString[resp.Rcode],
			}).Warn("DNS query returned error code")
			continue
		}

		var ips []net.IP
		for _, rr := range resp.Answer {
			switch record := rr.(type) {
			case *dns.A:
				ips = append(ips, record.A)
			case *dns.AAAA:
				ips = append(ips, record.AAAA)
			}
		}

		if len(ips) > 0 {
			return ips, nil
		}
	}

	return nil, fmt.Errorf("failed to resolve domain %s from any upstream server", domain)
}

// ServeDNSOverUDP serves DNS requests over UDP using a custom connection
func (s *DNSServer) ServeDNSOverUDP(conn net.PacketConn) error {
	s.logger.Info("Starting DNS server over UDP")

	for {
		buf := make([]byte, 512) // Standard DNS message size
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			s.logger.WithError(err).Error("Failed to read from UDP connection")
			continue
		}

		go s.handleUDPRequest(conn, addr, buf[:n])
	}
}

// handleUDPRequest processes a single UDP DNS request
func (s *DNSServer) handleUDPRequest(conn net.PacketConn, addr net.Addr, data []byte) {
	msg := new(dns.Msg)
	if err := msg.Unpack(data); err != nil {
		s.logger.WithError(err).Error("Failed to unpack DNS message")
		return
	}

	// Create a response writer that writes back to the UDP connection
	w := &udpResponseWriter{
		conn: conn,
		addr: addr,
	}

	s.HandleDNSRequest(w, msg)
}

// udpResponseWriter implements dns.ResponseWriter for UDP connections
type udpResponseWriter struct {
	conn net.PacketConn
	addr net.Addr
}

func (w *udpResponseWriter) LocalAddr() net.Addr {
	return w.conn.LocalAddr()
}

func (w *udpResponseWriter) RemoteAddr() net.Addr {
	return w.addr
}

func (w *udpResponseWriter) WriteMsg(msg *dns.Msg) error {
	data, err := msg.Pack()
	if err != nil {
		return err
	}

	_, err = w.conn.WriteTo(data, w.addr)
	return err
}

func (w *udpResponseWriter) Write(data []byte) (int, error) {
	return w.conn.WriteTo(data, w.addr)
}

func (w *udpResponseWriter) Close() error {
	return nil // UDP connections don't need to be closed
}

func (w *udpResponseWriter) TsigStatus() error {
	return nil
}

func (w *udpResponseWriter) TsigTimersOnly(bool) {}

func (w *udpResponseWriter) Hijack() {}

// FixedIPDNSHandler returns a fixed IP for all DNS queries
type FixedIPDNSHandler struct {
	fixedIP net.IP
	logger  *logrus.Logger
}

// NewFixedIPDNSHandler creates a new fixed IP DNS handler
func NewFixedIPDNSHandler(ip string) *FixedIPDNSHandler {
	fixedIP := net.ParseIP(ip)
	if fixedIP == nil {
		panic(fmt.Sprintf("Invalid IP address: %s", ip))
	}
	
	return &FixedIPDNSHandler{
		fixedIP: fixedIP,
		logger:  logrus.New(),
	}
}

// HandleQuery implements DNSHandler interface, always returns the fixed IP
func (h *FixedIPDNSHandler) HandleQuery(ctx context.Context, domain string, qtype uint16) ([]net.IP, error) {
	h.logger.WithFields(logrus.Fields{
		"domain": domain,
		"type":   dns.TypeToString[qtype],
		"ip":     h.fixedIP.String(),
	}).Info("Returning fixed IP for DNS query")
	
	// Only return IP for A records (IPv4) if our fixed IP is IPv4
	if qtype == dns.TypeA && h.fixedIP.To4() != nil {
		return []net.IP{h.fixedIP}, nil
	}
	
	// Only return IP for AAAA records (IPv6) if our fixed IP is IPv6
	if qtype == dns.TypeAAAA && h.fixedIP.To4() == nil && h.fixedIP.To16() != nil {
		return []net.IP{h.fixedIP}, nil
	}
	
	// For other query types or mismatched IP versions, return empty
	return []net.IP{}, nil
}
