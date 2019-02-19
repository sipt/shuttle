package conn

import (
	"container/list"
	"net"
	"sync"
	"time"
)

type Address struct {
	IPAddr net.IP
	Port   int
}

type ConnLine struct {
	Local  *Address
	Remote *Address
}

type Connection struct {
	ClientConn IConn
	ClientLine *ConnLine
	ServerConn IConn
	ServerLine *ConnLine
	CreateAt   time.Time
}

type IConnPool interface {
	Put(connection Connection, keys ...string)
	Close(clientConnID, serverConnID int64, keys ...string)
	List(keys ...string) []Connection
	Clear(keys ...string)
}

const defaultPoolKey = "_GLOBAL"

func NewPool() IConnPool {
	return &Pool{
		connMap: make(map[string]*List),
	}
}

type List struct {
	*list.List
	sync.RWMutex
}
type Pool struct {
	connMap map[string]*List
	sync.RWMutex
}

func (p *Pool) Put(connection Connection, keys ...string) {
	p.Lock()
	var key = defaultPoolKey
	if len(keys) > 0 {
		key = keys[0]
	}
	l, ok := p.connMap[key]
	if !ok {
		l = &List{List: list.New()}
		p.connMap[key] = l
	}
	l.Lock()
	p.Unlock()
	l.PushBack(&connection)
	l.Unlock()
}

func (p *Pool) Close(clientConnID, serverConnID int64, keys ...string) {
	p.RLock()
	var key = defaultPoolKey
	if len(keys) > 0 {
		key = keys[0]
	}
	l, ok := p.connMap[key]
	if ok {
		l.Lock()
		p.RUnlock()
		node := l.Front()
		for ; node != nil; node = node.Next() {
			c, ok := node.Value.(*Connection)
			if ok && c.ClientConn.GetID() == clientConnID &&
				c.ServerConn.GetID() == serverConnID {
				l.Remove(node)
				//close conn
				c, ok := node.Value.(*Connection)
				if ok {
					_ = c.ClientConn.Close()
					_ = c.ServerConn.Close()
				}
				break
			}
		}
		l.Unlock()
	} else {
		p.RUnlock()
	}
}

func (p *Pool) List(keys ...string) []Connection {
	p.RLock()
	var key = defaultPoolKey
	if len(keys) > 0 {
		key = keys[0]
	}
	l, ok := p.connMap[key]
	if !ok {
		p.RUnlock()
		return nil
	}
	l.RLock()
	p.RUnlock()
	conns := make([]Connection, l.Len())
	n := l.Front()
	for i := 0; n != nil; i++ {
		conns[i] = *(n.Value.(*Connection))
		n = n.Next()
	}
	l.RUnlock()
	return conns
}

func (p *Pool) Clear(keys ...string) {
	if len(keys) == 0 {
		//close all
		p.Lock()
		m := p.connMap
		p.connMap = make(map[string]*List)
		p.Unlock()
		for _, v := range m {
			n := v.Front()
			for ; n != nil; n = n.Next() {
				if c, ok := n.Value.(*Connection); ok {
					_ = c.ClientConn.Close()
					_ = c.ServerConn.Close()
				}
			}
		}
	} else {
		p.Lock()
		for _, v := range keys {
			l := p.connMap[v]
			l.RLock()
			n := l.Front()
			for ; n != nil; n = n.Next() {
				if c, ok := n.Value.(*Connection); ok {
					_ = c.ClientConn.Close()
					_ = c.ServerConn.Close()
				}
			}
			l.RUnlock()
			delete(p.connMap, v)
		}
		p.Unlock()
	}
}
