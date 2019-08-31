package listener

import (
	"context"
	"fmt"

	"github.com/sipt/shuttle/constant/typ"
)

type NewFunc func(addr string) (func(context.Context, typ.HandleFunc), error)

var creator = make(map[string]NewFunc)

// Register: register {key: NewFunc}
func Register(key string, f NewFunc) {
	creator[key] = f
}

// Get: get listener by key
func Get(typ, addr string) (func(context.Context, typ.HandleFunc), error) {
	f, ok := creator[typ]
	if !ok {
		return nil, fmt.Errorf("inbound not support: %s", typ)
	}
	return f(addr)
}
