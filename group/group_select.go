package group

import (
	"context"
	"net/url"
	"sync"

	"github.com/sipt/shuttle/constant/typ"
	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/dns"
	"github.com/sipt/shuttle/server"
)

const TypSelect = "select"

func init() {
	Register(TypSelect, func(ctx context.Context, runtime typ.Runtime, name string, params map[string]string, _ dns.Handle) (group IGroup, e error) {
		s := &SelectGroup{
			name:     name,
			RWMutex:  &sync.RWMutex{},
			testUrl:  params[ParamsKeyTestURI],
			udpRelay: params[ParamsKeyUdpRelay] == "true",
			runtime:  runtime,
		}
		if s.testUrl == "" {
			s.testUrl = DefaultTestURL
		} else if testUrl, err := url.Parse(s.testUrl); err != nil || len(testUrl.Scheme) == 0 || len(testUrl.Hostname()) == 0 {
			err = errors.Errorf("[group: %s] [%s: %s] is invalid", name, ParamsKeyTestURI, s.testUrl)
			return nil, err
		}
		return s, nil
	})
}

type SelectGroup struct {
	name     string
	servers  []IServerX
	current  IServerX
	testUrl  string
	udpRelay bool
	runtime  typ.Runtime
	*sync.RWMutex
}

func (s *SelectGroup) Append(servers []IServerX) {
	if len(servers) == 0 {
		return
	}
	s.Lock()
	defer s.Unlock()
	if len(s.servers) == 0 {
		s.servers = servers
	} else {
		s.servers = append(s.servers, servers...)
	}
	selected, ok := s.runtime.Get("selected").(string)
	if !ok {
		s.current = s.servers[0]
	} else {
		for i, v := range s.servers {
			if v.Name() == selected {
				s.current = s.servers[i]
				break
			}
		}
		if s.current == nil {
			s.current = s.servers[0]
		}
	}
	err := s.runtime.Set("selected", s.current.Name())
	if err != nil {
		logrus.WithField("select_group", s.name).WithError(err).Error("save runtime failed")
	}
}
func (s *SelectGroup) Items() []IServerX {
	return s.servers
}
func (s *SelectGroup) Reset() {
	s.testAllRTT()
}
func (s *SelectGroup) Typ() string {
	return TypSelect
}
func (s *SelectGroup) Name() string {
	return s.name
}
func (s *SelectGroup) Trace() []string {
	trace := make([]string, 0, len(s.current.Trace())+1)
	return append(append(trace, s.name), s.current.Trace()...)
}
func (s *SelectGroup) UdpRelay() bool {
	return s.udpRelay
}
func (s *SelectGroup) Server() server.IServer {
	s.RLock()
	defer s.RUnlock()
	if s.current != nil {
		return s.current.Server()
	} else if len(s.servers) > 0 {
		return s.servers[0].Server()
	} else {
		return nil
	}
}
func (s *SelectGroup) Selected() IServerX {
	return s.current
}
func (s *SelectGroup) Select(name string) error {
	s.Lock()
	defer s.Unlock()
	if s.current != nil && s.current.Name() == name {
		return nil
	}
	for i, v := range s.servers {
		if v.Name() == name {
			s.current = s.servers[i]
			err := s.runtime.Set(s.name, s.current.Name())
			if err != nil {
				logrus.WithField("select_group", s.name).WithError(err).Error("save runtime failed")
			}
			return nil
		}
	}
	return errors.Errorf("server[%s] not exist in group[%s]", name, s.name)
}
func (s *SelectGroup) Clear() {
	if len(s.servers) == 0 {
		return
	}
	s.Lock()
	defer s.Unlock()
	s.servers = nil // clear all
}

func (s *SelectGroup) testAllRTT() {
	if len(s.servers) == 0 {
		return
	}
	var wg = &sync.WaitGroup{}
	for _, v := range s.servers {
		wg.Add(1)
		go func(sx IServerX) {
			if g, ok := v.(IGroup); ok {
				g.Reset()
			} else {
				v.Server().TestRtt(s.name, s.testUrl)
			}
			wg.Done()
		}(v)
	}
	wg.Wait()
}
