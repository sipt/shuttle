package storage

import (
	"container/list"
	"sync"
)

func NewLRUList(count int) *LRUList {
	return &LRUList{
		count:   count,
		list:    &list.List{},
		RWMutex: &sync.RWMutex{},
	}
}

type LRUList struct {
	count int
	list  *list.List
	*sync.RWMutex
}

func (l *LRUList) PushBack(v interface{}) {
	l.Lock()
	defer l.Unlock()
	l.list.PushBack(v)
	if l.list.Len() > l.count {
		l.list.Remove(l.list.Front())
	}
}

func (l *LRUList) Clear() {
	l.Lock()
	defer l.Unlock()
	l.list.Init()
}

func (l *LRUList) Range(f func(v interface{}) bool) {
	l.RLock()
	defer l.RUnlock()
	for e := l.list.Front(); e != nil; e = e.Next() {
		if f(e.Value) {
			break
		}
	}
}
func (l *LRUList) RangeForUpdate(f func(v interface{}) bool) {
	l.Lock()
	defer l.Unlock()
	for e := l.list.Front(); e != nil; e = e.Next() {
		if f(e.Value) {
			break
		}
	}
}
