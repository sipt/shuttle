package conn

import (
	"container/list"
	"net"
	"strconv"
	"sync"
	"time"
)

func Init() {
	connMap = &ConnMap{
		m: make(map[string]IConnPool),
	}
}

func GetPool(key string) IConnPool {
	return connMap.Get(key)
}

func RemovePool(key string) {
	connMap.Remove(key)
}

var connMap *ConnMap

type ConnMap struct {
	m map[string]IConnPool
	sync.RWMutex
}

func (c *ConnMap) Get(key string) IConnPool {
	c.Lock()
	defer c.Unlock()
	if p, ok := c.m[key]; ok {
		return p
	} else {
		p = NewPool()
		c.m[key] = p
		return p
	}
}
func (c *ConnMap) Remove(key string) {
	c.Lock()
	defer c.Unlock()
	delete(c.m, key)
}

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
	Put(lc, sc IConn, keys ...string)
	Close(clientConnID, serverConnID int64, keys ...string)
	List(keys ...string) []Connection
	Clear(keys ...string)
	Replace(clientConnID int64, sc IConn, keys ...string)
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

func (p *Pool) Put(lc, sc IConn, keys ...string) {
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
	l.PushBack(&Connection{
		ClientConn: lc,
		ClientLine: ParseConnLine(lc),
		ServerConn: sc,
		ServerLine: ParseConnLine(sc),
	})
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

func (p *Pool) Replace(clientConnID int64, sc IConn, keys ...string) {
	p.RLock()
	key := defaultPoolKey
	if len(keys) > 0 && len(keys[0]) > 0 {
		key = keys[0]
	}
	conns, ok := p.connMap[key]
	if !ok {
		return
	}
	conns.Lock()
	defer conns.Unlock()
	p.RUnlock()
	for n := conns.Front(); n != nil; n = n.Next() {
		c := n.Value.(*Connection)
		if c.ClientConn.GetID() == clientConnID {
			c.ServerConn = sc
			c.ServerLine = ParseConnLine(sc)
		}
	}
}

func ParseConnLine(c IConn) *ConnLine {
	cl := &ConnLine{
		Local: &Address{},
	}
	host, port, err := net.SplitHostPort(c.LocalAddr().String())
	if err != nil {
		cl.Local.IPAddr = net.ParseIP(host)
		cl.Local.Port, _ = strconv.Atoi(port)
	}
	host, port, err = net.SplitHostPort(c.RemoteAddr().String())
	if err != nil {
		cl.Local.IPAddr = net.ParseIP(host)
		cl.Local.Port, _ = strconv.Atoi(port)
	}
	return cl
}
