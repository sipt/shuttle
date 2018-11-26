package conn

import (
	"context"
	"net"
)

const (
	TCP = "tcp"
	UDP = "udp"
)

type IConn interface {
	net.Conn
	GetID() int64
	GetRecordID() int64
	SetRecordID(id int64)
	GetNetwork() string
	Flush() (int, error)
	Context() context.Context
	SetContext(context.Context)
}

func NewDefaultConn(conn net.Conn, network string) (IConn, error) {
	c, err := DefaultDecorate(conn, network)
	return c, err
}

func DirectConn(network, host string) (IConn, error) {
	conn, err := net.DialTimeout(network, host, DefaultTimeOut)
	if err != nil {
		return nil, err
	}
	c, err := NewDefaultConn(conn, network)
	if err == nil {
		c, err = TrafficDecorate(c)
	}
	return c, err
}
