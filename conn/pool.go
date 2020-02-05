package conn

import (
	"container/list"
	"sync"
)

type Link struct {
	Input  ICtxConn
	Output ICtxConn
}

func (l *Link) include(conn ICtxConn) bool {
	return l.Input != nil && l.Input.GetConnID() == conn.GetConnID() || l.Output != nil && l.Output.GetConnID() == conn.GetConnID()
}

func (l *Link) remove(conn ICtxConn) bool {
	if l.Input != nil && l.Input.GetConnID() == conn.GetConnID() {
		l.Input = nil
		return true
	} else if l.Output != nil && l.Output.GetConnID() == conn.GetConnID() {
		l.Output = nil
		return true
	}
	return false
}

func (l *Link) isEmpty() bool {
	return l.Input == nil && l.Output == nil
}

var pool = list.New()
var lock = &sync.RWMutex{}

func PushInputConn(conn ICtxConn) {
	lock.Lock()
	defer lock.Unlock()
	pool.PushBack(&Link{Input: conn})
}
func PushOutputConn(conn ICtxConn) {
	lock.Lock()
	defer lock.Unlock()
	pool.PushBack(&Link{Output: conn})
}
func Remove(conn ICtxConn) {
	lock.Lock()
	defer lock.Unlock()
	for e := pool.Front(); e != nil; e = e.Next() {
		if l := e.Value.(*Link); l.remove(conn) {
			if l.isEmpty() {
				pool.Remove(e)
			}
			break
		}
	}
}
func PoolRange(f func(*Link) bool) {
	lock.RLock()
	defer lock.RUnlock()
	for e := pool.Front(); e != nil; e = e.Next() {
		if f(e.Value.(*Link)) {
			break
		}
	}
}
