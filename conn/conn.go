package conn

import (
	"context"
	"net"
	"sync/atomic"
)

const (
	KeyConnID = "conn-id"
)

var (
	connID      int64 = 0
	ctx, cancel       = context.WithCancel(context.Background())
)

func GetConnID() int64 {
	return atomic.AddInt64(&connID, 1)
}

type DialFunc func(ctx context.Context, network string, addr, port string) (ICtxConn, error)

func DefaultDial(ctx context.Context, network string, addr, port string) (ICtxConn, error) {
	conn, err := net.Dial(network, net.JoinHostPort(addr, port))
	if err != nil {
		return nil, err
	}
	return WrapConn(conn), nil
}

type ICtxConn interface {
	net.Conn
	context.Context

	GetConnID() int64
	WithContext(ctx context.Context)
	GetContext() context.Context
	WithValue(string, interface{})
}

type ctxConn struct {
	net.Conn
	context.Context
}

func (c *ctxConn) GetContext() context.Context {
	return c.Context
}

func (c *ctxConn) WithContext(ctx context.Context) {
	c.Context = ctx
}

func (c *ctxConn) WithValue(k string, v interface{}) {
	c.Context = context.WithValue(c.Context, k, v)
}

func (c *ctxConn) GetConnID() int64 {
	id, _ := c.Value(KeyConnID).(int64)
	return id
}

func WrapConn(conn net.Conn) ICtxConn {
	return &ctxConn{
		Conn:    conn,
		Context: context.WithValue(ctx, KeyConnID, GetConnID()),
	}
}

func NewConn(conn net.Conn, ctx context.Context) ICtxConn {
	if ctx.Value(KeyConnID) == nil {
		ctx = context.WithValue(ctx, KeyConnID, GetConnID())
	}
	return &ctxConn{
		Conn:    conn,
		Context: ctx,
	}
}
