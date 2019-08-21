package dns

import (
	"context"
	"sync"
)

func NewCache() ICache {
	return &cache{
		m:       make(map[string]*DNS),
		l:       make([]*DNS, 0, 16),
		RWMutex: &sync.RWMutex{},
	}
}

type ICache interface {
	Get(key string) DNS
	Set(key string, value DNS)
	List() (reply []DNS)
	Clear()
}

type cache struct {
	m map[string]*DNS
	l []*DNS
	*sync.RWMutex
}

func (c *cache) Clear() {
	c.Lock()
	defer c.Unlock()
	c.m = make(map[string]*DNS)
	c.l = make([]*DNS, 0, 16)
}

func (c *cache) Get(key string) DNS {
	c.RLock()
	defer c.RUnlock()
	return *c.m[key]
}

func (c *cache) Set(key string, value DNS) {
	c.Lock()
	defer c.Unlock()
	_, ok := c.m[key]
	c.m[key] = &value
	if ok {
		for i, v := range c.l {
			if v.Domain == key {
				c.l[i] = &value
				return
			}
		}
	} else {
		c.l = append(c.l, &value)
	}
}

func (c *cache) List() (reply []DNS) {
	reply = make([]DNS, len(c.l))
	c.RLock()
	defer c.RUnlock()
	for i, v := range c.l {
		reply[i] = *v
	}
	return
}

func newCacheHandle(next Handle) (Handle, error) {
	cache := NewCache()
	return func(ctx context.Context, domain string) *DNS {
		dns := cache.Get(domain)
		dnsPtr := &dns
		if dnsPtr.IsNil() {
			dnsPtr = next(ctx, domain)
			cache.Set(domain, *dnsPtr)
		}
		return dnsPtr
	}, nil
}
