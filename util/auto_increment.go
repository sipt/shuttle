package util

import "sync/atomic"

var id int64

func getNextID() int64 {
	reply := atomic.AddInt64(&id, 1)
	return reply
}

func NextID() int64 {
	return getNextID()
}
