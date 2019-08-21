package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	key := []string{"www.baidu.com", "www.google.com", "github.com"}
	cache := NewCache()
	list := make([]DNS, len(key))
	for i, v := range key {
		list[i] = DNS{Domain: v}
		cache.Set(v, list[i])
	}
	for _, v := range key {
		assert.Equal(t, cache.Get(v).Domain, v)
	}

	assert.ElementsMatch(t, cache.List(), list)
	cache.Clear()
	assert.Equal(t, len(cache.List()), 0)
}
