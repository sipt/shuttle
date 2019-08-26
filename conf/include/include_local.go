package include

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/conf/storage"
)

func init() {
	Register("local", NewLocalInclude)
}

func NewLocalInclude(params map[string]string) (IInclude, error) {
	s, err := storage.Get(storage.KeyFile, params)
	if err != nil {
		return nil, errors.Wrap(err, "NewLocalInclude failed")
	}
	return &localInclude{
		IStorage: s,
	}, nil
}

type localInclude struct {
	storage.IStorage
}

func (l *localInclude) Load() ([]byte, error) {
	return l.IStorage.Load()
}
func (l *localInclude) RegisterNotify(ctx context.Context, notify func()) error {
	return l.IStorage.RegisterNotify(ctx, notify)
}
func (l *localInclude) Save(data []byte) error {
	return l.IStorage.Save(data)
}
