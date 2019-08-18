package conn

import (
	"context"
	"net"
)

type DialTCPFunc func(ctx context.Context, addr, port string) (*net.TCPConn, error)
type DialUDPFunc func(ctx context.Context, addr, port string) (*net.UDPConn, error)

func DefaultDialTCP(ctx context.Context, addr, port string) (*net.TCPConn, error) {
	conn, err := net.Dial("tcp", net.JoinHostPort(addr, port))
	if err != nil {
		return nil, err
	}
	return conn.(*net.TCPConn), nil
}

func DefaultDialUDP(ctx context.Context, addr, port string) (*net.UDPConn, error) {
	conn, err := net.Dial("udp", net.JoinHostPort(addr, port))
	if err != nil {
		return nil, err
	}
	return conn.(*net.UDPConn), nil
}
