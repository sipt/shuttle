package tun

import (
	"context"
	"io"
	"net"

	"gvisor.dev/gvisor/pkg/tcpip/adapters/gonet"
	"gvisor.dev/gvisor/pkg/tcpip/stack"
	"gvisor.dev/gvisor/pkg/tcpip/transport/udp"
	"gvisor.dev/gvisor/pkg/waiter"
)

type UdpListener interface {
	Accept() (UdpConn, error)
}

type UdpConn interface {
	net.PacketConn
	io.ReadWriteCloser
	RemoteAddr() net.Addr
}

func newUdpListener(ctx context.Context, bufferSize int) *udpListener {
	return &udpListener{
		connPool: make(chan UdpConn, bufferSize),
		ctx:      ctx,
	}
}

type udpListener struct {
	connPool chan UdpConn
	ctx      context.Context
}

func (t *udpListener) Accept() (UdpConn, error) {
	select {
	case <-t.ctx.Done():
		return nil, t.ctx.Err()
	case conn, ok := <-t.connPool:
		if !ok {
			return nil, io.EOF
		}
		return conn, nil
	}
}

type udpConn struct {
	*gonet.UDPConn
	id stack.TransportEndpointID
}

// LocalAddr returns the local network address
func (c *udpConn) LocalAddr() net.Addr {
	return c.UDPConn.LocalAddr()
}

// RemoteAddr returns the remote network address
func (c *udpConn) RemoteAddr() net.Addr {
	return c.UDPConn.RemoteAddr()
}

// IsDNSQuery checks if this UDP connection is targeting DNS port (53)
func (c *udpConn) IsDNSQuery() bool {
	if addr := c.LocalAddr(); addr != nil {
		if udpAddr, ok := addr.(*net.UDPAddr); ok {
			return udpAddr.Port == 53
		}
	}
	return false
}

func UdpForward(ctx context.Context, s *stack.Stack, bufferSize int) (UdpListener, error) {
	listener := newUdpListener(ctx, bufferSize)
	forwarder := udp.NewForwarder(s, func(r *udp.ForwarderRequest) (handled bool) {
		wq := &waiter.Queue{}
		ep, err := r.CreateEndpoint(wq)
		if err != nil {
			return false
		}
		conn := gonet.NewUDPConn(wq, ep)
		listener.connPool <- &udpConn{
			UDPConn: conn,
			id:      r.ID(),
		}
		return true
	})
	s.SetTransportProtocolHandler(udp.ProtocolNumber, forwarder.HandlePacket)
	return listener, nil
}
