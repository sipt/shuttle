package conn

import (
	"bytes"
	"context"
	"net"
	"time"

	"github.com/sipt/shuttle/log"
	"github.com/sipt/shuttle/pool"
	"github.com/sipt/shuttle/util"
)

var DefaultTimeOut = 10 * time.Second

//
func DefaultDecorate(c net.Conn, network string) (IConn, error) {
	id := util.GetLongID()
	return &DefaultConn{
		Conn:    c,
		ID:      id,
		Network: network,
		context: context.Background(),
	}, nil
}

func DefaultDecorateForTls(c net.Conn, network string, id int64) (IConn, error) {
	return &DefaultConn{
		Conn:    c,
		ID:      id,
		Network: network,
		context: context.Background(),
	}, nil
}

type DefaultConn struct {
	net.Conn
	ID       int64
	RecordID int64
	Network  string
	context  context.Context
}

func (c *DefaultConn) GetID() int64 {
	return c.ID
}

func (c *DefaultConn) GetRecordID() int64 {
	return c.RecordID
}

func (c *DefaultConn) SetRecordID(id int64) {
	c.RecordID = id
}

func (c *DefaultConn) Flush() (int, error) {
	return 0, nil
}

func (c *DefaultConn) GetNetwork() string {
	return c.Network
}
func (c *DefaultConn) Context() context.Context {
	return c.context
}
func (c *DefaultConn) SetContext(ctx context.Context) {
	c.context = ctx
}

func (c *DefaultConn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	log.Logger.Tracef("[ID:%d] Read Data: %v", c.GetID(), b[:n])
	return
}

func (c *DefaultConn) Write(b []byte) (n int, err error) {
	log.Logger.Tracef("[ID:%d] Write Data: %v", c.GetID(), b)
	return c.Conn.Write(b)
}
func (c *DefaultConn) Close() error {
	log.Logger.Debugf("[ID:%d] close connection", c.GetID())
	return c.Conn.Close()
}

//超时装饰
func TimerDecorate(c IConn, rto, wto time.Duration) (IConn, error) {
	if rto == 0 {
		rto = DefaultTimeOut
	}
	if wto == 0 {
		wto = DefaultTimeOut
	}
	return &TimerConn{
		IConn:        c,
		ReadTimeOut:  rto,
		WriteTimeOut: wto,
	}, nil
}

type TimerConn struct {
	IConn
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
}

func (c *TimerConn) resetReadDeadline() {
	if c.ReadTimeOut > -1 {
		c.SetReadDeadline(time.Now().Add(c.ReadTimeOut))
	}
}

func (c *TimerConn) resetWriteDeadline() {
	if c.WriteTimeOut > -1 {
		c.SetWriteDeadline(time.Now().Add(c.WriteTimeOut))
	}
}

func (c *TimerConn) Read(b []byte) (n int, err error) {
	c.resetReadDeadline()
	n, err = c.IConn.Read(b)
	c.resetWriteDeadline()
	return
}

func (c *TimerConn) Write(b []byte) (n int, err error) {
	c.resetWriteDeadline()
	n, err = c.IConn.Write(b)
	c.resetReadDeadline()
	return
}

//缓冲装饰
func BufferDecorate(c IConn) (IConn, error) {
	return &BufferConn{
		IConn:  c,
		buffer: bytes.NewBuffer(pool.GetBuf()[:0]),
	}, nil
}

type BufferConn struct {
	IConn
	buffer *bytes.Buffer
}

func (c *BufferConn) Write(b []byte) (n int, err error) {
	return c.buffer.Write(b)
}

func (c *BufferConn) Flush() (n int, err error) {
	n, err = c.IConn.Write(c.buffer.Bytes())
	if err != nil {
		return
	}
	c.buffer.Reset()
	n, err = c.IConn.Flush()
	return
}

//实时写出
func RealTimeDecorate(c IConn) (IConn, error) {
	return &RealTimeFlush{
		IConn: c,
	}, nil
}

type RealTimeFlush struct {
	IConn
}

func (r *RealTimeFlush) Write(b []byte) (n int, err error) {
	n, err = r.IConn.Write(b)
	if err != nil {
		return
	}
	_, err = r.IConn.Flush()
	return
}

//导出装饰器
var upload, download func(int64, int) = nil, nil

func InitTrafficChannel(up, down func(int64, int)) {
	upload, download = up, down
}
func TrafficDecorate(c IConn, ) (IConn, error) {
	return &Traffic{
		IConn: c,
	}, nil
}

type Traffic struct {
	IConn
}

func (t *Traffic) Read(b []byte) (n int, err error) {
	n, err = t.IConn.Read(b)
	if download != nil && t.GetRecordID() > 0 && n > 0 {
		download(t.GetRecordID(), n)
	}
	return
}

func (t *Traffic) Write(b []byte) (n int, err error) {
	n, err = t.IConn.Write(b)
	if upload != nil && t.GetRecordID() > 0 && n > 0 {
		upload(t.GetRecordID(), n)
	}
	return
}
