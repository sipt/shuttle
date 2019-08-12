package conf

import (
	"bytes"
	"context"

	"github.com/sipt/shuttle/conf/include"
	"github.com/sipt/shuttle/conf/marshal"
	"github.com/sipt/shuttle/conf/model"
	"github.com/sipt/shuttle/conf/storage"
)

// LoadConfig
// typ:
func LoadConfig(ctx context.Context, typ, encode string, params map[string]string, notify func()) (*model.Config, error) {
	s, err := storage.GetStorage(typ, params)
	if err != nil {
		return nil, err
	}
	data, err := s.Load()
	if err != nil {
		return nil, err
	}
	m, err := marshal.Get(encode, params)
	if err != nil {
		return nil, err
	}
	config, err := m.UnMarshal(data)
	if err != nil {
		return nil, err
	}
	err = s.RegisterNotify(ctx, notify)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(data)
	for _, v := range config.Include {
		c, err := include.Get(v.Typ, v.Params)
		if err != nil {
			return nil, err
		}
		data, err = c.Load()
		if err != nil {
			return nil, err
		}
		buffer.Write(data)
		err = c.RegisterNotify(ctx, notify)
		if err != nil {
			return nil, err
		}
	}
	return m.UnMarshal(buffer.Bytes())
}
