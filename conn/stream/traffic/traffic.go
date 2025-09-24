package stream

import (
	"context"
	"io"
	"sync/atomic"
	"time"

	"github.com/sipt/shuttle/conn"
	"github.com/sipt/shuttle/conn/stream"
	"github.com/sipt/shuttle/constant/typ"
)

var down, up int64 = 0, 0
var oldDown, oldUp int64 = 0, 0

func init() {
	stream.RegisterStream("traffic", newTrafficMetrics)
}

func newTrafficMetrics(ctx context.Context, _ typ.Runtime, _ map[string]string) (typ.DecorateFunc, error) {
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				oldDown = atomic.SwapInt64(&down, 0)
				oldUp = atomic.SwapInt64(&up, 0)
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
	return func(c conn.ICtxConn) conn.ICtxConn {
		writerTo, ok1 := c.(io.WriterTo)
		readerFrom, ok2 := c.(io.ReaderFrom)
		if ok1 && ok2 {
			return &trafficConn2{
				WriterTo:   writerTo,
				ReaderFrom: readerFrom,
				trafficConn: &trafficConn{
					ICtxConn: c,
				},
			}
		}
		return &trafficConn{
			ICtxConn: c,
		}
	}, nil
}

type trafficConn struct {
	conn.ICtxConn
}

func (t *trafficConn) Read(b []byte) (n int, err error) {
	n, err = t.ICtxConn.Read(b)
	atomic.AddInt64(&up, int64(n))
	return
}

func (t *trafficConn) Write(b []byte) (n int, err error) {
	n, err = t.ICtxConn.Write(b)
	atomic.AddInt64(&down, int64(n))
	return
}

type trafficWrite struct {
	io.Writer
}

func (t *trafficWrite) Write(b []byte) (n int, err error) {
	n, err = t.Writer.Write(b)
	atomic.AddInt64(&up, int64(n))
	return n, err
}

type trafficRead struct {
	io.Reader
}

func (t *trafficRead) Read(b []byte) (n int, err error) {
	n, err = t.Reader.Read(b)
	atomic.AddInt64(&down, int64(n))
	return n, err
}

type trafficConn2 struct {
	io.WriterTo
	io.ReaderFrom
	*trafficConn
}

func (t *trafficConn2) WriteTo(w io.Writer) (n int64, err error) {
	trafficWrite := &trafficWrite{Writer: w}
	n, err = t.WriterTo.WriteTo(trafficWrite)
	atomic.AddInt64(&up, int64(n))
	return n, err
}

func (t *trafficConn2) ReadFrom(r io.Reader) (n int64, err error) {
	trafficRead := &trafficRead{Reader: r}
	n, err = t.ReaderFrom.ReadFrom(trafficRead)
	atomic.AddInt64(&down, int64(n))
	return n, err
}
