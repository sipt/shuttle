package conf

import "fmt"

type NewMarshalFunc func(map[string]string) (IMarshal, error)

var marshalCreator = make(map[string]NewMarshalFunc)

// Register: register {key: marshal}
func RegisterMarshal(key string, f NewMarshalFunc) {
	marshalCreator[key] = f
}

// GetMarshal: get Marshal by key
func GetMarshal(key string, params map[string]string) (IMarshal, error) {
	f, ok := marshalCreator[key]
	if !ok {
		return nil, fmt.Errorf("marshal not support: %s", key)
	}
	return f(params)
}

type IMarshal interface {
	Marshal(*Config) ([]byte, error)
	UnMarshal([]byte) (*Config, error)
}
