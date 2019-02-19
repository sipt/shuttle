package storage

import (
	"container/list"
	"sync"
)

const EngineMemory = "memory"

func init() {
	Register(EngineMemory, func() IStorage {
		return NewMemoryStorage(DefaultMaxLength)
	})
}

func NewMemoryStorage(maxLength int) IStorage {
	return &memoryStorage{
		maxLength: maxLength,
		pool:      make(map[string]*List),
	}
}

type List struct {
	*list.List
	sync.RWMutex
}

type memoryStorage struct {
	sync.RWMutex
	maxLength int
	pool      map[string]*List
}

//添加记录到指定的key下
func (m *memoryStorage) Put(key string, value Record) {
	l, ok := m.pool[key]
	if !ok {
		m.Lock()
		if l, ok = m.pool[key]; !ok {
			l = &List{List: list.New()}
			m.pool[key] = l
		}
		m.Unlock()
	}
	l.Lock()
	l.PushBack(&value)
	for m.maxLength < l.Len() {
		l.Remove(l.Front())
	}
	l.Unlock()
}

//当前所有的Key数量
func (m *memoryStorage) Keys() []string {
	m.RLock()
	keys := make([]string, 0, len(m.pool))
	for k := range m.pool {
		keys = append(keys, k)
	}
	m.RUnlock()
	return keys
}

//key下的记录数量
func (m *memoryStorage) Count(key string) int {
	m.RLock()
	l, ok := m.pool[key]
	m.RUnlock()
	if ok {
		l.RLock()
		defer l.RUnlock()
		return l.Len()
	}
	return 0
}

//获取key的所有记录
func (m *memoryStorage) Get(key string) []Record {
	m.RLock()
	l, ok := m.pool[key]
	m.RUnlock()
	if ok {
		l.RLock()
		defer l.RUnlock()
		reply := make([]Record, 0, l.Len())
		node := l.Front()
		for node != nil {
			reply = append(reply, *(node.Value.(*Record)))
			node = node.Next()
		}
		return reply
	}
	return nil
}

// 设置每个Key的上限记录数
func (m *memoryStorage) SetLength(l int) {
	m.maxLength = l
}

// 清空
func (m *memoryStorage) Clear(keys ...string) {
	m.Lock()
	defer m.Unlock()
	if len(keys) == 0 {
		m.pool = make(map[string]*List)
	} else {
		for _, k := range keys {
			delete(m.pool, k)
		}
	}
}
