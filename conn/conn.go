package conn

import (
	"bytes"
	"context"
	"io"
	"net"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
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
	dialer := new(net.Dialer)
	conn, err := dialer.DialContext(ctx, network, net.JoinHostPort(addr, port))
	if err != nil {
		return nil, err
	}
	logrus.WithField("network", network).WithField("addr", addr+":"+port).Debug("connect to server")
	c := WrapConn(conn)
	PushOutputConn(c)
	return c, nil
}

type ICtxConn interface {
	net.Conn
	context.Context

	GetConnID() int64
	WithContext(ctx context.Context)
	GetContext() context.Context
	WithValue(interface{}, interface{})
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

func (c *ctxConn) WithValue(k interface{}, v interface{}) {
	c.Context = context.WithValue(c.Context, k, v)
}

func (c *ctxConn) GetConnID() int64 {
	id, _ := c.Value(KeyConnID).(int64)
	return id
}

func (c *ctxConn) Close() error {
	Remove(c)
	logrus.WithField("conn-id", c.GetConnID()).Debug("close the connection")
	return c.Conn.Close()
}

func (c *ctxConn) Read(b []byte) (int, error) {
	n, err := c.Conn.Read(b)
	return n, err
}
func (c *ctxConn) Write(b []byte) (int, error) {
	n, err := c.Conn.Write(b)
	return n, err
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

type udpConn struct {
	context.Context
	buf    *bytes.Buffer
	local  net.Addr
	remote net.Addr
	c      net.PacketConn
}

func (u *udpConn) Read(b []byte) (n int, err error) {
	if u.buf != nil && u.buf.Len() > 0 {
		n, err = u.buf.Read(b)
		logrus.WithField("data", b[:n]).Debug("read udp data")
		return
	}
	return 0, io.EOF
}
func (u *udpConn) Write(b []byte) (n int, err error) {
	if u.c != nil {
		n, err = u.c.WriteTo(b, u.remote)
		logrus.WithField("data", b).WithField("remote", u.remote.String()).Debug("read udp data")
		return
	}
	return 0, nil
}
func (u *udpConn) Close() error {
	u.buf.Reset()
	return nil
}
func (u *udpConn) LocalAddr() net.Addr {
	return u.local
}
func (u *udpConn) RemoteAddr() net.Addr {
	return u.remote
}
func (u *udpConn) SetDeadline(t time.Time) error {
	return nil
}
func (u *udpConn) SetReadDeadline(t time.Time) error {
	return nil
}
func (u *udpConn) SetWriteDeadline(t time.Time) error {
	return u.c.SetWriteDeadline(t)
}
func (u *udpConn) GetContext() context.Context {
	return u.Context
}
func (u *udpConn) WithContext(ctx context.Context) {
	u.Context = ctx
}

func (u *udpConn) WithValue(k interface{}, v interface{}) {
	u.Context = context.WithValue(u.Context, k, v)
}
func (u *udpConn) GetConnID() int64 {
	id, _ := u.Value(KeyConnID).(int64)
	return id
}
func NewUDPConn(pc net.PacketConn, ctx context.Context, remoteAddr net.Addr, data []byte) ICtxConn {
	if ctx.Value(KeyConnID) == nil {
		ctx = context.WithValue(ctx, KeyConnID, GetConnID())
	}
	return &udpConn{
		Context: ctx,
		c:       pc,
		buf:     bytes.NewBuffer(data),
		local:   pc.LocalAddr(),
		remote:  remoteAddr,
	}
}

func Transfer(from, to ICtxConn) {
	end := make(chan bool, 1)
	go func() {
		_, err := io.Copy(to, from)
		if err != nil {
			end <- true
			logrus.WithError(err).WithField("from", from.GetConnID()).WithField("to", to.GetConnID()).
				Debug("io.copy failed")
		}
	}()
	go func() {
		_, err := io.Copy(from, to)
		if err != nil {
			logrus.WithError(err).WithField("from", from.GetConnID()).WithField("to", to.GetConnID()).
				Debug("io.copy failed")
		}
		end <- true
	}()
	<-end
	_ = from.Close()
	_ = to.Close()
}
