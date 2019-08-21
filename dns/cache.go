package dns

import (
	"context"
	"sync"
	"time"
)

var (
	clearCacheExpiredInterval = 10 * time.Minute
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
	ClearExpired()
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
	d, ok := c.m[key]
	if ok {
		return *d
	}
	return DNS{}
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

func (c *cache) ClearExpired() {
	c.RLock()
	queue := make([]int, 0, 8)
	for i, v := range c.l {
		if v.ExpireAt.Sub(time.Now()).Minutes() > 30 {
			queue = append(queue, i)
		}
	}
	c.RUnlock()
	c.Lock()
	for i := len(queue) - 1; i >= 0; i-- {
		if c.l[i].ExpireAt.Sub(time.Now()).Minutes() > 30 {
			delete(c.m, c.l[i].Domain)
			if i+1 == len(c.l) {
				c.l = c.l[:i]
			} else {
				c.l = append(c.l[:i], c.l[i+1:]...)
			}
		}
	}
	c.Unlock()
}

func newCacheHandle(next Handle) (Handle, error) {
	cache := NewCache()
	go func() {
		timer := time.NewTimer(clearCacheExpiredInterval)
		for {
			<-timer.C
			cache.ClearExpired()
			timer.Reset(clearCacheExpiredInterval)
		}
	}()
	return func(ctx context.Context, domain string) *DNS {
		dns := cache.Get(domain)
		dnsPtr := &dns
		if dnsPtr.IsNil() {
			dnsPtr = next(ctx, domain)
			cache.Set(domain, *dnsPtr)
		} else if dnsPtr.ExpireAt.Before(time.Now()) {
			go func() {
				dnsPtr = next(ctx, domain)
				cache.Set(domain, *dnsPtr)
			}()
		}
		return dnsPtr
	}, nil
}
