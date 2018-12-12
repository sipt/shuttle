package selector

import (
	"github.com/sipt/shuttle/log"
	"github.com/sipt/shuttle/proxy"
	"sync/atomic"
	"time"
)

const (
	timerDulation = 10 * time.Minute
)

func init() {
	proxy.RegisterSelector("rtt", func(group *proxy.ServerGroup) (proxy.ISelector, error) {
		s := &rttSelector{
			group:    group,
			timer:    time.NewTimer(timerDulation),
			selected: group.Servers[0].(proxy.IServer),
			cancel:   make(chan bool, 1),
		}
		go func() {
			for {
				select {
				case <-s.timer.C:
				case <-s.cancel:
					return
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
	group    *proxy.ServerGroup
	selected proxy.IServer
	status   uint32
	timer    *time.Timer
	cancel   chan bool
}

func (r *rttSelector) Get() (*proxy.Server, error) {
	return r.selected.GetServer()
}
func (r *rttSelector) Select(name string) error {
	return nil
}
func (r *rttSelector) Refresh() error {
	r.autoTest()
	return nil
}
func (r *rttSelector) Reset(group *proxy.ServerGroup) error {
	r.group = group
	r.selected = r.group.Servers[0].(proxy.IServer)
	go r.autoTest()
	return nil
}
func (r *rttSelector) Destroy() {
	r.cancel <- true
}
func (r *rttSelector) autoTest() {
	if r.status == 0 {
		if ok := atomic.CompareAndSwapUint32(&r.status, 0, 1); !ok {
			return
		}
	}
	r.timer.Stop()
	log.Logger.Debug("[Rtt-Selector] start testing ...")
	var is proxy.IServer
	var s *proxy.Server
	var err error
	c := make(chan *proxy.Server, 1)
	for _, v := range r.group.Servers {
		is = v.(proxy.IServer)
		s, err = is.GetServer()
		if err != nil {
			continue
		}
		go urlTest(s, r.group.GetRttRrl(), c)
	}
	s = <-c
	log.Logger.Infof("[Rtt-Select] rtt select server: [%s]", s.Name)
	r.selected = s
	r.timer.Reset(timerDulation)
	atomic.CompareAndSwapUint32(&r.status, 1, 0)
}

func urlTest(s *proxy.Server, rttUrl string, c chan *proxy.Server) {
	rtt, err := proxy.TestRTT(s, rttUrl)
	if err != nil {
		s.Rtt = -1
		log.Logger.Debugf("[Rtt-Select] [%s]  url test result: <failed> %v", s.Name, err)
		return
	}
	s.Rtt = rtt
	select {
	case c <- s:
	default:
	}
	log.Logger.Debugf("[Rtt-Select] [%s]  Rtt:[%dms]", s.Name, s.Rtt.Nanoseconds()/1000000)
}
func (r *rttSelector) Current() proxy.IServer {
	return r.selected
}
