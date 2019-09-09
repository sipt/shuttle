package group

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/conf/model"
	"github.com/sipt/shuttle/server"
)

var Global = "GLOBAL"

func ApplyConfig(ctx context.Context, config *model.Config, servers map[string]server.IServer) (map[interface{}]IGroup, error) {
	serverMap := make(map[string]IServerX)
	for _, v := range servers {
		serverMap[v.Name()] = &serverx{IServer: v}
	}
	groups := make(map[interface{}]IGroup)
	var (
		g   IGroup
		err error
		ok  bool
	)
	for name, v := range config.ServerGroup {
		if name == Global {
			return nil, errors.Errorf("group name [%s] is reserved", name)
		}
		if v.Params == nil {
			v.Params = map[string]string{ParamsKeyTestURI: config.General.DefaultTestURI}
		} else if _, ok := v.Params[ParamsKeyTestURI]; !ok {
			v.Params[ParamsKeyTestURI] = config.General.DefaultTestURI
		}
		g, err = Get(ctx, v.Typ, name, v.Params)
		if err != nil {
			return nil, err
		}
		if _, ok = serverMap[name]; ok {
			return nil, errors.Errorf("group name duplicate: %s", name)
		}
		serverMap[name] = g
		groups[name] = g
	}
	// global group, when in GLOBAL_MODE
	gl, err := Get(ctx, TypSelect, Global, map[string]string{ParamsKeyTestURI: config.General.DefaultTestURI})
	if err != nil {
		return nil, err
	}
	gs := make([]IServerX, 0, len(config.ServerGroup))
	for name, g := range config.ServerGroup {
		ss := make([]IServerX, 0, len(g.Servers))
		for _, s := range g.Servers {
			ss = append(ss, serverMap[s])
			gs = append(gs, serverMap[s])
		}
		serverMap[name].(IGroup).Append(ss)
		gs = append(gs, serverMap[name])
	}
	gl.Append(gs)
	groups[Global] = gl
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
	Items() []IServerX
	Reset()
	Clear()
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
