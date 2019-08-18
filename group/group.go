package group

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/sipt/shuttle/server"

	"github.com/sipt/shuttle/conf/model"
)

func ApplyConfig(ctx context.Context, config *model.Config, servers []server.IServer) ([]IGroup, error) {
	serverMap := make(map[string]IServerX)
	for _, v := range servers {
		serverMap[v.Name()] = &serverx{IServer: v}
	}
	groups := make([]IGroup, 0, len(config.ServerGroup))
	var (
		g   IGroup
		err error
		ok  bool
	)
	for name, v := range config.ServerGroup {
		g, err = Get(ctx, v.Typ, name, v.Params)
		if err != nil {
			return nil, err
		}
		if _, ok = serverMap[name]; ok {
			return nil, errors.Errorf("group name duplicate: %s", name)
		}
		serverMap[name] = g
		groups = append(groups, g)
	}
	for name, g := range config.ServerGroup {
		ss := make([]IServerX, 0, len(g.Servers))
		for _, s := range g.Servers {
			ss = append(ss, serverMap[s])
		}
		serverMap[name].(IGroup).Append(ss)
	}
	return groups, nil
}

type NewFunc func(ctx context.Context, name string, params map[string]string) (IGroup, error)

var creator = make(map[string]NewFunc)

// Register: register {key: NewFunc}
func Register(key string, f NewFunc) {
	creator[key] = f
}

// Get: get group by key
func Get(ctx context.Context, typ string, name string, params map[string]string) (IGroup, error) {
	f, ok := creator[typ]
	if !ok {
		return nil, fmt.Errorf("server not support: %s", typ)
	}
	return f(ctx, name, params)
}

type IGroup interface {
	Append(servers []IServerX)
	Select(name string) error
	IServerX
}

type IServerX interface {
	Typ() string
	Name() string
	// connect to server
	Server() server.IServer
	Trace() []string
}

type serverx struct {
	server.IServer
}

func (s *serverx) Server() server.IServer {
	return s.IServer
}

func (s *serverx) Trace() []string {
	return []string{s.Name()}
}
