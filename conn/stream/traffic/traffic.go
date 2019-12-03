package stream

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/sipt/shuttle/conn"
	"github.com/sipt/shuttle/conn/stream"
)

var down, up int64 = 0, 0

func init() {
	stream.RegisterStream("traffic", newTrafficMetrics)
}

func newTrafficMetrics(ctx context.Context, _ map[string]string) (stream.DecorateFunc, error) {
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				atomic.SwapInt64(&down, 0)
				atomic.SwapInt64(&up, 0)
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
	return func(c conn.ICtxConn) conn.ICtxConn {
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
	atomic.AddInt64(&down, int64(n))
	return
}

func (t *trafficConn) Write(b []byte) (n int, err error) {
	n, err = t.ICtxConn.Write(b)
	atomic.AddInt64(&up, int64(n))
	return
}
