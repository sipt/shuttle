package selector

import (
	"github.com/sipt/shuttle"
	"net/http"
	"net"
	"time"
	"context"
	"sync/atomic"
)

const (
	timerDulation = 10 * time.Minute
)

func init() {
	shuttle.RegisterSelector("rtt", func(group *shuttle.ServerGroup) (shuttle.ISelector, error) {
		s := &rttSelector{
			group: group,
			timer: time.NewTimer(timerDulation),
		}
		go func() {
			for {
				select {
				case <-s.timer.C:
				}
				s.autoTest()
				s.timer.Reset(timerDulation)
			}
		}()
		go s.Refresh()
		return s, nil
	})
}

type rttSelector struct {
	group    *shuttle.ServerGroup
	selected shuttle.IServer
	status   uint32
	timer    *time.Timer
}

func (m *rttSelector) Get() (*shuttle.Server, error) {
	return m.selected.GetServer()
}
func (m *rttSelector) Select(name string) error {
	return nil
}
func (m *rttSelector) Refresh() error {
	m.autoTest()
	return nil
}
func (m *rttSelector) Reset(group *shuttle.ServerGroup) error {
	m.group = group
	m.selected = m.group.Servers[0].(shuttle.IServer)
	go m.autoTest()
	return nil
}
func (m *rttSelector) autoTest() {
	if m.status == 0 {
		if ok := atomic.CompareAndSwapUint32(&m.status, 0, 1); !ok {
			return
		}
	}
	var is shuttle.IServer
	var s *shuttle.Server
	var err error
	c := make(chan *shuttle.Server, 1)
	start := time.Now()
	for _, v := range m.group.Servers {
		is = v.(shuttle.IServer)
		s, err = is.GetServer()
		if err != nil {
			continue
		}
		go urlTest(s, c)
	}
	s = <-c
	shuttle.Logger.Infof("[RTT Select] %s  RTT: %dms", s.Name, time.Now().Sub(start).Nanoseconds()/1000)
	close(c)
	m.selected = s
	m.timer.Reset(timerDulation)
	atomic.CompareAndSwapUint32(&m.status, 1, 0)
}

const url = "http://www.gstatic.com/generate_204"

func urlTest(s *shuttle.Server, c chan *shuttle.Server) {
	tr := &http.Transport{
		DialContext: func(_ context.Context, _, addr string) (net.Conn, error) {
			return s.Conn(addr)
		},
	}
	client := &http.Client{Transport: tr, Timeout: 2 * time.Second}
	resp, err := client.Get(url)
	if err == nil && resp.StatusCode != 204 {
		select {
		case c <- s:
		default:
		}
	}
}
