package shuttle

import (
	"github.com/sipt/shuttle/log"
	"github.com/sipt/shuttle/pool"
	"io"
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
	var l = len(p)
	if s.w2 != nil {
		var buf []byte
		if pool.BufferSize >= l {
			buf = pool.GetBuf()
		} else {
			buf = make([]byte, l)
		}
		copy(buf, p)
		_, err := s.w2.Write(buf[:l])
		if err != nil {
			log.Logger.Errorf("[Shunt] [Sub2] Write data failed: %s", err.Error())
		}
	}
	if s.w1 != nil {
		n, err = s.w1.Write(p)
	}
	if n == 0 {
		n = l
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
