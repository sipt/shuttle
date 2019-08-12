package include

import (
	"context"
	"fmt"
)

type NewFunc func(map[string]string) (IInclude, error)

var creator = make(map[string]NewFunc)

// Register: register {key: Include}
func Register(key string, f NewFunc) {
	creator[key] = f
}

// Get: get Include by key
func Get(key string, params map[string]string) (IInclude, error) {
	f, ok := creator[key]
	if !ok {
		return nil, fmt.Errorf("marshal not support: %s", key)
	}
	return f(params)
}

type IInclude interface {
	Load() ([]byte, error)
	RegisterNotify(ctx context.Context, notify func()) error
	Save([]byte) error
}
