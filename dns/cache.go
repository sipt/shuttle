package dns

import (
	"container/heap"
	"github.com/sipt/shuttle/log"
	"net"
	"strings"
	"sync"
	"time"
)

const CacheTTL = 10 * time.Minute

var dnsCacheManager *CacheManager

func InitDNSCache() {
	if dnsCacheManager == nil {
		dnsCacheManager = NewCacheManager()
	} else {
		dnsCacheManager.Clear()
	}
}

// resolve domain and through the DNS-Cache
func ResolveDomainByCache(domain string) (*Answer, error) {
	if net.ParseIP(domain) != nil {
		return nil, nil
	}

	matched := dnsCacheManager.Range(func(data interface{}) bool {
		answer := data.(*Answer)
		if answer.Domain == domain {
			return true
		}
		return false
	})
	var answer *Answer
	if matched != nil {
		answer = matched.(*Answer)
		log.Logger.Infof("[DNS] [Cache] resolve [%s] -> [%s] [%s]", domain, strings.Join(answer.IPs, ","), answer.Country)
		return answer, nil
	}
	//cache miss
	answer, err := ResolveDomain(domain)
	if err != nil {
		return nil, err
	}
	if answer != nil {
		dnsCacheManager.Push(answer, CacheTTL)
		log.Logger.Infof("[DNS] [Cache] resolve [%s] -> [%s] [%s]", domain, strings.Join(answer.IPs, ","), answer.Country)
	}
	return answer, nil
}

func ClearDNSCache() {
	dnsCacheManager.Clear()
}

func DNSCacheList() []*Answer {
	list := make([]*Answer, 0, 64)
	dnsCacheManager.Range(func(data interface{}) bool {
		list = append(list, data.(*Answer))
		return false
	})
	return list
}

type CacheEntity struct {
	data    interface{}
	expires time.Time
}

func NewCachePool() *CachePool {
	return &CachePool{
		list: make([]*CacheEntity, 0, 64),
	}
}

type CachePool struct {
	list []*CacheEntity
	sync.RWMutex
}

func (c *CachePool) Head() *CacheEntity {
	c.RLock()
	defer c.RUnlock()
	if len(c.list) > 0 {
		return c.list[0]
	}
	return nil
}

func (c *CachePool) Len() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.list)
}
func (c *CachePool) Less(i, j int) bool {
	c.RLock()
	defer c.RUnlock()
	return c.list[i].expires.Before(c.list[j].expires)
}

func (c *CachePool) Swap(i, j int) {
	l := len(c.list)
	if l == 0 || i < 0 || j < 0 || i >= l || j >= l {
		return
	}
	c.Lock()
	defer c.Unlock()
	l = len(c.list)
	if l == 0 || i < 0 || j < 0 || i >= l || j >= l {
		return
	}
	c.list[i], c.list[j] = c.list[j], c.list[i]
}
func (c *CachePool) Push(x interface{}) {
	c.Lock()
	defer c.Unlock()
	cell := x.(*CacheEntity)
	c.list = append(c.list, cell)
}
func (c *CachePool) Pop() interface{} {
	if len(c.list) == 0 {
		return nil
	}
	c.Lock()
	defer c.Unlock()
	if len(c.list) == 0 {
		return nil
	}
	l := len(c.list) - 1
	r := c.list[l]
	c.list = c.list[:l]
	return r
}

func (c *CachePool) Range(f func(data interface{}) (breaked bool)) interface{} {
	c.RLock()
	defer c.RUnlock()
	for _, v := range c.list {
		if f(v.data) {
			return v.data
		}
	}
	return nil
}

func (c *CachePool) Clear() {
	c.Lock()
	defer c.Unlock()
	c.list = c.list[:0]
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		timer:  time.NewTimer(10 * 60 * time.Second),
		cancel: make(chan bool, 1),
		pool:   NewCachePool(),
	}
}

type CacheManager struct {
	timer  *time.Timer
	cancel chan bool
	pool   *CachePool
}

func (c *CacheManager) Range(f func(data interface{}) (breaked bool)) interface{} {
	return c.pool.Range(f)
}

func (c *CacheManager) Push(data interface{}, ttl time.Duration) {
	Push(c.pool, &CacheEntity{
		data:    data,
		expires: time.Now().Add(ttl),
	})
	c.refresh()
}

func (c *CacheManager) Run() {
	go func() {
		for {
			select {
			case <-c.timer.C:
				Pop(c.pool)
				c.refresh()
			case <-c.cancel:
			}
		}
	}()
}

func (c *CacheManager) Stop() {
	c.cancel <- true
}

func (c *CacheManager) refresh() {
	var wait time.Duration = -1
	if c.pool.Head() != nil {
		entity := c.pool.Head()
		wait = entity.expires.Sub(time.Now())
		for wait < 0 {
			Pop(c.pool)
			if c.pool.Head() == nil {
				break
			}
			entity = c.pool.Head()
			wait = entity.expires.Sub(time.Now())
		}
	}
	if wait < -1 {
		wait = 10 * 60 * time.Second
	}
	c.timer.Reset(wait)
}
func (c *CacheManager) Clear() {
	c.pool.Clear()
}

var heapLock = &sync.Mutex{}
// Push pushes the element x onto the heap. The complexity is
// O(log(n)) where n = h.Len().
//
func Push(h heap.Interface, x interface{}) {
	heapLock.Lock()
	defer heapLock.Unlock()
	heap.Push(h, x)
}

// Pop removes the minimum element (according to Less) from the heap
// and returns it. The complexity is O(log(n)) where n = h.Len().
// It is equivalent to Remove(h, 0).
//
func Pop(h heap.Interface) interface{} {
	heapLock.Lock()
	defer heapLock.Unlock()
	return heap.Pop(h)
}
