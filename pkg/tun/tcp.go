package stack

import (
	"gvisor.dev/gvisor/pkg/tcpip"
	"gvisor.dev/gvisor/pkg/tcpip/header"
	"gvisor.dev/gvisor/pkg/tcpip/transport/tcp"
	"gvisor.dev/gvisor/pkg/waiter"
	"io"

	"gvisor.dev/gvisor/pkg/tcpip/stack"

	"gvisor.dev/gvisor/pkg/tcpip/adapters/gonet"
)

type TcpListener interface {
	Accept() (TcpConn, error)
}

type TcpConn interface {
	io.ReadWriteCloser
}

func newTcpListener(bufferSize int) *tcpListener {
	return &tcpListener{
		connPool: make(chan TcpConn, bufferSize),
	}
}

type tcpListener struct {
	connPool chan TcpConn
}

func (t *tcpListener) Accept() (TcpConn, error) {
	conn, ok := <-t.connPool
	if !ok {
		return nil, io.EOF
	}
	return conn, nil
}

type tcpConn struct {
	*gonet.TCPConn
	id stack.TransportEndpointID
}

func TcpForward(s *stack.Stack, bufferSize int) (TcpListener, error) {
	var tcpListener = newTcpListener(bufferSize)
	forwarder := tcp.NewForwarder(s, 0, 1<<16, func(r *tcp.ForwarderRequest) {
		wq := &waiter.Queue{}
		ep, err := r.CreateEndpoint(wq)
		if err != nil {
			r.Complete(true)
			return
		}
		defer r.Complete(false)
		// recv/send buffer
		var ss tcpip.TCPSendBufferSizeRangeOption
		if err := s.TransportProtocolOption(header.TCPProtocolNumber, &ss); err == nil {
			ep.SocketOptions().SetReceiveBufferSize(int64(ss.Default), false)
		}

		var rs tcpip.TCPReceiveBufferSizeRangeOption
		if err := s.TransportProtocolOption(header.TCPProtocolNumber, &rs); err == nil {
			ep.SocketOptions().SetReceiveBufferSize(int64(rs.Default), false)
		}

		// handle conn
		conn := &tcpConn{
			TCPConn: gonet.NewTCPConn(wq, ep),
			id:      r.ID(),
		}
		tcpListener.connPool <- conn
	})
	s.SetTransportProtocolHandler(tcp.ProtocolNumber, forwarder.HandlePacket)
	return tcpListener, nil
}
