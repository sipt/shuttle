package stream

import (
	"fmt"

	"github.com/sipt/shuttle/conf/model"
	"github.com/sipt/shuttle/conn"
)

func ApplyConfig(config *model.Config) (before, after DecorateFunc, err error) {
	before, after = NopConn, NopConn
	var bs = make([]DecorateFunc, 0, len(config.Stream.Before))
	var as = make([]DecorateFunc, 0, len(config.Stream.After))
	var decorate DecorateFunc
	for _, v := range config.Stream.Before {
		decorate, err = GetDecorateFunc(v.Typ, v.Params)
		if err != nil {
			return
		}
		bs = append(bs, decorate)
	}
	if len(bs) > 0 {
		before = func(conn conn.ICtxConn) conn.ICtxConn {
			for _, f := range bs {
				conn = f(conn)
			}
			return conn
		}
	}
	for _, v := range config.Stream.After {
		decorate, err = GetDecorateFunc(v.Typ, v.Params)
		if err != nil {
			return
		}
		as = append(as, decorate)
	}
	if len(as) > 0 {
		after = func(conn conn.ICtxConn) conn.ICtxConn {
			for _, f := range as {
				conn = f(conn)
			}
			return conn
		}
	}
	return
}

type DecorateFunc func(conn.ICtxConn) conn.ICtxConn
type NewFunc func(map[string]string) (DecorateFunc, error)

func NopConn(conn conn.ICtxConn) conn.ICtxConn {
	return conn
}

var streamMap = make(map[string]NewFunc)

// RegisterStream: register {key: filterFunc}
func RegisterStream(key string, f NewFunc) {
	streamMap[key] = f
}

// Get: get filter by key
func GetDecorateFunc(typ string, params map[string]string) (DecorateFunc, error) {
	f, ok := streamMap[typ]
	if !ok {
		return nil, fmt.Errorf("stream not support: %s", typ)
	}
	return f(params)
}
