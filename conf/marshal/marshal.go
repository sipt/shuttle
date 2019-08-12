package marshal

import (
	"fmt"

	"github.com/sipt/shuttle/conf/model"
)

type NewFunc func(map[string]string) (IMarshal, error)

var creator = make(map[string]NewFunc)

// Register: register {key: marshal}
func Register(key string, f NewFunc) {
	creator[key] = f
}

// GetMarshal: get Marshal by key
func Get(key string, params map[string]string) (IMarshal, error) {
	f, ok := creator[key]
	if !ok {
		return nil, fmt.Errorf("marshal not support: %s", key)
	}
	return f(params)
}

type IMarshal interface {
	Marshal(*model.Config) ([]byte, error)
	UnMarshal([]byte) (*model.Config, error)
}
