package group

import (
	"context"
	"net/url"
	"sync"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/dns"
	"github.com/sipt/shuttle/server"
)

const TypSelect = "select"

func init() {
	Register(TypSelect, func(ctx context.Context, name string, params map[string]string, _ dns.Handle) (group IGroup, e error) {
		s := &SelectGroup{
			name:    name,
			RWMutex: &sync.RWMutex{},
			testUrl: params[ParamsKeyTestURI],
		}
		if s.testUrl == "" {
			s.testUrl = DefaultTestURL
		} else if testUrl, err := url.Parse(s.testUrl); err != nil || len(testUrl.Scheme) == 0 || len(testUrl.Hostname()) == 0 {
			err = errors.Errorf("[group: %s] [%s: %s] is invalid", name, ParamsKeyTestURI, s.testUrl)
			return nil, err
		}
		return &SelectGroup{
			name:    name,
			RWMutex: &sync.RWMutex{},
		}, nil
	})
}

type SelectGroup struct {
	name    string
	servers []IServerX
	current IServerX
	testUrl string
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
	s.current = s.servers[0]
}
func (s *SelectGroup) Items() []IServerX {
	return s.servers
}
func (s *SelectGroup) Reset() {
	for _, v := range s.servers {
		if g, ok := v.(IGroup); ok {
			g.Reset()
		} else {
			v.Server().TestRtt(s.name, s.testUrl)
		}
	}
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
func (s *SelectGroup) Select(name string) error {
	s.Lock()
	defer s.Unlock()
	if s.current != nil && s.current.Name() == name {
		return nil
	}
	for _, v := range s.servers {
		if v.Name() == name {
			s.current = v
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
