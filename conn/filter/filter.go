package filter

import (
	"context"
	"fmt"

	"github.com/sipt/shuttle/conf/model"
	"github.com/sipt/shuttle/constant/typ"
)

func ApplyConfig(ctx context.Context, config *model.Config) (filter FilterFunc, err error) {
	filter = Nop
	for _, v := range config.Filter {
		filter, err = Get(ctx, v.Typ, v.Params, filter)
		if err != nil {
			return
		}
	}
	return
}

type FilterFunc func(typ.HandleFunc) typ.HandleFunc
type NewFunc func(context.Context, map[string]string, FilterFunc) (FilterFunc, error)

func Nop(h typ.HandleFunc) typ.HandleFunc {
	return h
}

var filterMap = make(map[string]NewFunc)

// Register: register {key: filterFunc}
func Register(key string, f NewFunc) {
	filterMap[key] = f
}

// Get: get filter by key
func Get(ctx context.Context, typ string, params map[string]string, filter FilterFunc) (FilterFunc, error) {
	f, ok := filterMap[typ]
	if !ok {
		return nil, fmt.Errorf("filter not support: %s", typ)
	}
	return f(ctx, params, filter)
}
