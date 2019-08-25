package span

import (
	"context"
)

var tagKey = "tag"
var spanKey = "span"

func ExtractTag(ctx context.Context) (context.Context, ISpan) {
	return Extract(ctx, tagKey)
}

func ExtractSpan(ctx context.Context) (context.Context, ISpan) {
	return Extract(ctx, spanKey)
}

func Extract(ctx context.Context, key interface{}) (context.Context, ISpan) {
	span := ctx.Value(key)
	if span == nil {
		var m ISpan = &mapSpan{values: make(map[string]interface{})}
		ctx = context.WithValue(ctx, key, m)
		return ctx, m
	}
	return ctx, span.(ISpan)
}

type ISpan interface {
	Set(key string, value interface{}) ISpan
	Get(key string) interface{}
	Has(key string) bool
	Values() map[string]interface{}
}

type mapSpan struct {
	values map[string]interface{}
}

func (m *mapSpan) Set(key string, value interface{}) ISpan {
	m.values[key] = value
	return m
}

func (m *mapSpan) Has(key string) bool {
	_, ok := m.values[key]
	return ok
}

func (m *mapSpan) Values() map[string]interface{} {
	return m.values
}

func (m *mapSpan) Get(key string) interface{} {
	return m.values[key]
}
