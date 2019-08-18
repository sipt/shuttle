package group

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/server"
)

const TypSelect = "select"

func init() {
	Register(TypSelect, func(ctx context.Context, name string, params map[string]string) (group IGroup, e error) {
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
