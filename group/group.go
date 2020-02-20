package group

import (
	"context"
	"fmt"

	"github.com/sipt/shuttle/constant/typ"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/conf/model"
	"github.com/sipt/shuttle/dns"
	"github.com/sipt/shuttle/server"
)

var Global = "GLOBAL"

func ApplyConfig(ctx context.Context, config *model.Config, runtime typ.Runtime, servers map[string]server.IServer, dnsHandle dns.Handle) (map[interface{}]IGroup, error) {
	serverMap := make(map[string]IServerX)
	for _, v := range servers {
		serverMap[v.Name()] = WrapServer(v)
	}
	groups := make(map[interface{}]IGroup)
	var (
		g   IGroup
		err error
		ok  bool
	)
	getRuntime, err := ApplyRuntime(ctx, runtime)
	if err != nil {
		return nil, err
	}
	for name, v := range config.ServerGroup {
		if name == Global {
			return nil, errors.Errorf("group name [%s] is reserved", name)
		}
		if v.Params == nil {
			v.Params = map[string]string{ParamsKeyTestURI: config.General.DefaultTestURI}
		} else if _, ok := v.Params[ParamsKeyTestURI]; !ok {
			v.Params[ParamsKeyTestURI] = config.General.DefaultTestURI
		}
		g, err = Get(ctx, v.Typ, getRuntime(name), name, v.Params, dnsHandle)
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
	gl, err := Get(ctx, TypSelect, getRuntime(Global), Global, map[string]string{ParamsKeyTestURI: config.General.DefaultTestURI}, dnsHandle)
	if err != nil {
		return nil, err
	}
	gs := make([]IServerX, 0, len(config.ServerGroup))
	for gname, g := range config.ServerGroup {
		ss := make([]IServerX, 0, len(g.Servers))
		gEntity := serverMap[gname].(IGroup)
		for _, sname := range g.Servers {
			s := serverMap[sname]
			if s == nil {
				return nil, errors.Errorf("[group:%s] [server: %s] not exist in group/server", gname, sname)
			}
			if gEntity.UdpRelay() && !s.UdpRelay() {
				return nil, errors.Errorf("[group:%s] [server: %s] server not support udp-relay", gname, sname)
			}
			ss = append(ss, s)
		}
		gEntity.Append(ss)
		gs = append(gs, serverMap[gname])
	}
	gl.Append(gs)
	groups[Global] = gl
	return groups, nil
}

func ApplyRuntime(ctx context.Context, rt typ.Runtime) (func(string) typ.Runtime, error) {
	rt = newRuntime("group", rt)
	return func(name string) typ.Runtime {
		return newRuntime(name, rt)
	}, nil
}

func newRuntime(name string, rt typ.Runtime) typ.Runtime {
	r := &runtime{
		name:    name,
		Runtime: rt,
	}
	var ok bool
	r.current, ok = rt.Get(name).(map[string]interface{})
	if !ok {
		r.current = make(map[string]interface{})
	}
	return r
}

type runtime struct {
	name    string
	current map[string]interface{}
	typ.Runtime
}

func (r *runtime) Get(key string) interface{} {
	return r.current[key]
}
func (r *runtime) Set(key string, value interface{}) error {
	r.current[key] = value
	return r.Runtime.Set(r.name, r.current)
}

type NewFunc func(ctx context.Context, runtime typ.Runtime, name string, params map[string]string, dnsHandle dns.Handle) (IGroup, error)

var creator = make(map[string]NewFunc)

// Register: register {key: NewFunc}
func Register(key string, f NewFunc) {
	creator[key] = f
}

// Get: get group by key
func Get(ctx context.Context, typ string, runtime typ.Runtime, name string, params map[string]string, dnsHandle dns.Handle) (IGroup, error) {
	f, ok := creator[typ]
	if !ok {
		return nil, fmt.Errorf("server not support: %s", typ)
	}
	return f(ctx, runtime, name, params, dnsHandle)
}

type IGroup interface {
	Append(servers []IServerX)
	Select(name string) error
	Selected() IServerX
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
	UdpRelay() bool
}

func WrapServer(s server.IServer) IServerX {
	return &serverx{s}
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
