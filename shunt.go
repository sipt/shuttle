package shuttle

import (
	"io"
	"github.com/sipt/shuttle/pool"
)

// 分流
func NewShunt(w1, w2 io.Writer) *Shunt {
	return &Shunt{
		w1: w1,
		w2: w2,
	}
}

type Shunt struct {
	w1, w2 io.Writer
}

func (s *Shunt) Write(p []byte) (n int, err error) {
	var buf []byte
	if pool.BufferSize >= len(p) {
		buf = pool.GetBuf()
	} else {
		buf = make([]byte, len(p))
	}
	copy(buf[0:], p)
	go func() {
		if s.w2 != nil {
			_, err := s.w2.Write(buf[:n])
			if err != nil {
				Logger.Errorf("[Shunt] [Sub2] Write data failed: %s", err.Error())
			}
		}
	}()
	if s.w1 != nil {
		n, err = s.w1.Write(p)
	}
	return
}

func ToWriter(w func([]byte) (int, error)) io.Writer {
	return &writer{
		w: w,
	}
}

type writer struct {
	w func([]byte) (int, error)
}

func (w *writer) Write(p []byte) (n int, err error) {
	return w.w(p)
}
