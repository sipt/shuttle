package stream

import (
	"context"
	"fmt"

	"github.com/sipt/shuttle/constant/typ"

	"github.com/sipt/shuttle/conf/model"
	"github.com/sipt/shuttle/conn"
	"github.com/sirupsen/logrus"
)

func ApplyConfig(ctx context.Context, config *model.Config) (before, after DecorateFunc, err error) {
	before, after = NopConn, NopConn
	var bs = make([]DecorateFunc, 0, len(config.Stream.Before))
	var as = make([]DecorateFunc, 0, len(config.Stream.After))
	var decorate DecorateFunc
	for _, v := range config.Stream.Before {
		decorate, err = GetDecorateFunc(ctx, v.Typ, v.Params)
		if err != nil {
			return
		}
		bs = append(bs, decorate)
		logrus.WithField("stream-before", v.Typ).Info("init stream-before success")
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
		decorate, err = GetDecorateFunc(ctx, v.Typ, v.Params)
		if err != nil {
			return
		}
		as = append(as, decorate)
		logrus.WithField("stream-after", v.Typ).Info("init stream-after success")
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

func ApplyRuntime(_ context.Context, _ typ.Runtime) error {
	return nil
}

type DecorateFunc func(conn.ICtxConn) conn.ICtxConn
type NewFunc func(context.Context, map[string]string) (DecorateFunc, error)

func NopConn(conn conn.ICtxConn) conn.ICtxConn {
	return conn
}

var streamMap = make(map[string]NewFunc)

// RegisterStream: register {key: filterFunc}
func RegisterStream(key string, f NewFunc) {
	streamMap[key] = f
}

// Get: get filter by key
func GetDecorateFunc(ctx context.Context, typ string, params map[string]string) (DecorateFunc, error) {
	f, ok := streamMap[typ]
	if !ok {
		return nil, fmt.Errorf("stream not support: %s", typ)
	}
	return f(ctx, params)
}
