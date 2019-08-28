package socks

import (
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

// udp server for socks5.
type udpServer struct {
	ln net.PacketConn
}

// Addr rerurns a server address.
func (s *udpServer) Addr() net.Addr {
	return s.ln.LocalAddr()
}

// Close closes the server.
func (s *udpServer) Close() error {
	return s.ln.Close()
}

func (s *udpServer) Serve(cmdFunc func(wt WriteTo, remote net.Addr, dst *Addr, b []byte) error) {
	b := make([]byte, 1500)
	n, remoteAddr, err := s.ln.ReadFrom(b)
	go s.Serve(cmdFunc)
	// handshake
	if b[0]|b[1] != 0 {
		logrus.WithField("server", "udp").Debugf("RSV [%x%x] not eq 0x0000", b[0], b[1])
		return
	}
	if b[2] != 0 {
		logrus.WithField("server", "udp").Debugf("FRAG [%x] not eq 0x00", b[2])
		return
	}
	l := 2
	off := 4
	dstAddr := &Addr{}
	switch b[3] {
	case AddrTypeIPv4:
		l += net.IPv4len
		dstAddr.IP = make(net.IP, net.IPv4len)
	case AddrTypeIPv6:
		l += net.IPv6len
		dstAddr.IP = make(net.IP, net.IPv6len)
	case AddrTypeFQDN:
		l += int(b[4])
		off = 5
	default:
		logrus.WithField("server", "udp").Debugf("ATYP [%x] unknown address type", b[3])
		return
	}
	if len(b[off:]) < l {
		logrus.WithField("server", "udp").Debugf("short cmd request")
		return
	}
	if dstAddr.IP != nil {
		copy(dstAddr.IP, b[off:])
	} else {
		dstAddr.Name = string(b[off : off+l-2])
	}
	dstAddr.Port = int(b[off+l-2])<<8 | int(b[off+l-1])
	if off+l >= n {
		logrus.WithField("server", "udp").Debugf("short cmd request")
		return
	}
	err = cmdFunc(s.ln, remoteAddr, dstAddr, b[off+l:n])
	if err != nil {
		logrus.WithError(err).WithField("server", "udp").Debugf("cmd func failed")
		return
	}
}

type WriteTo interface {
	WriteTo(p []byte, addr net.Addr) (n int, err error)
	LocalAddr() net.Addr
	SetWriteDeadline(t time.Time) error
}

// NewUDPServer returns a new udpServer.
func NewUDPServer(addr string) (*udpServer, error) {
	var err error
	s := new(udpServer)
	s.ln, err = net.ListenPacket("udp", addr)
	if err != nil {
		return nil, err
	}
	return s, nil
}
