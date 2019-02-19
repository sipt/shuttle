package storage

import (
	"sync"
	"sync/atomic"
	"time"
)

var (
	traffic = &Traffic{
		Cancel: make(chan bool, 1),
	}
)

type Traffic struct {
	sync.RWMutex
	UpSpeed        int
	DownSpeed      int
	UpBytes        int
	TotalUpBytes   int
	DownBytes      int
	TotalDownBytes int
	Cancel         chan bool
	status         int32
}

func StartTrafficStatistics() {
	if atomic.CompareAndSwapInt32(&traffic.status, 0, 1) {
		go func() {
			ticker := time.NewTicker(time.Second)
			for {
				select {
				case <-ticker.C:
					traffic.Lock()
					traffic.UpSpeed, traffic.UpBytes = traffic.UpBytes, 0
					traffic.DownSpeed, traffic.DownBytes = traffic.DownBytes, 0
					traffic.Unlock()
				case <-traffic.Cancel:
					atomic.CompareAndSwapInt32(&traffic.status, 1, 0)
					return
				}
			}
		}()
	}
}

func Stop() {
	traffic.Cancel <- true
}

func GetSpeed() (up, down int) {
	traffic.RLock()
	defer traffic.RUnlock()
	return traffic.UpSpeed, traffic.DownSpeed
}

func TrafficUp(upBytes int) {
	traffic.Lock()
	defer traffic.Unlock()
	traffic.UpBytes += upBytes
	traffic.TotalUpBytes += upBytes
}

func TrafficDown(downBytes int) {
	traffic.Lock()
	defer traffic.Unlock()
	traffic.DownBytes += downBytes
	traffic.TotalDownBytes += downBytes
}

func GetTrafficBytes() (up, down int) {
	traffic.RLock()
	defer traffic.RUnlock()
	return traffic.TotalUpBytes, traffic.TotalDownBytes
}
