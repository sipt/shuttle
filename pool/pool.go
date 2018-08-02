package pool

import "sync"

const BufferSize = 4108

func init() {
	pool = &sync.Pool{
		New: func() interface{} {
			return make([]byte, BufferSize)
		},
	}
}

var pool *sync.Pool

func GetBuf() []byte {
	buf := pool.Get().([]byte)
	buf = buf[:cap(buf)]
	return buf
}

func PutBuf(buf []byte) {
	pool.Put(buf)
}
