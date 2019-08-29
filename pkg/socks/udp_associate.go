package socks

import (
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

// udp server for socks5.
type udpServer struct {
	net.PacketConn
	dst *net.UDPAddr
}

func (s *udpServer) Serve(cmdFunc func(pc net.PacketConn, remote net.Addr, dst net.Addr, b []byte) error) {
	b := make([]byte, 1500)
	err := s.SetReadDeadline(time.Now().Add(time.Minute))
	if err != nil {
		logrus.WithField("server", "udp").Debugf("set read dead line failed")
		return
	}
	defer func() {
		if err != nil {
			s.Close()
		}
	}()
	n, remoteAddr, err := s.ReadFrom(b)
	if err != nil {
		logrus.WithError(err).WithField("server", "udp").WithField("addr", s.LocalAddr().String()).
			Debugf("read packet failed")
		return
	}
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
	ip := s.LocalAddr().(*net.UDPAddr).IP.String()
	if ip != "0.0.0.0" && ip != "::" && dstAddr.String() != s.LocalAddr().String() {
		return
	}
	err = cmdFunc(s, remoteAddr, s.dst, b[off+l:n])
	if err != nil {
		logrus.WithError(err).WithField("server", "udp").Debugf("cmd func failed")
		return
	}
}

// NewUDPServer returns a new udpServer.
func NewUDPServer(addr string, dst *net.UDPAddr) (*udpServer, error) {
	var err error
	s := new(udpServer)
	s.dst = dst
	s.PacketConn, err = net.ListenPacket("udp", addr)
	if err != nil {
		return nil, err
	}
	logrus.WithField("addr", s.LocalAddr()).Debugf("start udp read on: %s", s.LocalAddr().String())
	return s, nil
}
