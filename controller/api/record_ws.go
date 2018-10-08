package api

import (
	"github.com/gorilla/websocket"
	"time"
	"net/http"
	"sync"
	"github.com/sipt/shuttle/util"
	"github.com/sipt/shuttle/log"
	"github.com/sipt/shuttle"
)

func init() {
	pushTimeBuffer = &PushTimeBuffer{
		ticker: time.NewTicker(time.Second),
		buffer: make([]interface{}, 0, 8),
	}
	go pushTimeBuffer.Run()
	shuttle.RegisterPusher(func(v interface{}) {
		pushTimeBuffer.Push(v)
	})
}

var pushTimeBuffer *PushTimeBuffer

type PushTimeBuffer struct {
	ticker *time.Ticker
	buffer []interface{}
	sync.RWMutex
}

func (p *PushTimeBuffer) Push(v interface{}) {
	p.RLock()
	p.buffer = append(p.buffer, v)
	p.RUnlock()
}

func (p *PushTimeBuffer) Run() {
	for {
		<-p.ticker.C
		if len(p.buffer) > 0 {
			var buf []interface{}
			p.Lock()
			if len(p.buffer) > 0 {
				buf = p.buffer
				p.buffer = make([]interface{}, 0, 8)
			}
			p.Unlock()
			if len(buf) > 0 {
				wsCenter.RangeSend(buf)
			}
		}
	}
}

var wsCenter = &ConnCenter{
	conns: make([]*websocket.Conn, 0, 8),
	index: make([]int64, 0, 8),
}

type ConnCenter struct {
	conns []*websocket.Conn
	index []int64
	sync.RWMutex
}

func (c *ConnCenter) RangeSend(v interface{}) {
	c.RLock()
	for _, conn := range c.conns {
		conn.WriteJSON(v)
	}
	c.RUnlock()
}

func (c *ConnCenter) Add(i int64, conn *websocket.Conn) {
	c.Lock()
	c.index = append(c.index, i)
	c.conns = append(c.conns, conn)
	c.Unlock()
}

func (c *ConnCenter) Remove(i int64) {
	c.RLock()
	for k, v := range c.index {
		if v == i {
			if len(c.index) == 1 {
				c.index = c.index[:0]
				c.conns = c.conns[:0]
				break
			}
			c.index[k] = c.index[len(c.index)-1]
			c.index = c.index[:len(c.index)-1]
			c.conns[k] = c.conns[len(c.conns)-1]
			c.conns = c.conns[:len(c.conns)-1]
			break
		}
	}
	c.RUnlock()
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:   2048,
	WriteBufferSize:  2048,
	HandshakeTimeout: 5 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Logger.Errorf("[Shuttle-Controller] Failed to set websocket upgrade: %v", err)
		return
	}
	index, _ := util.IW.NextId()
	wsCenter.Add(index, conn)
	for {
		t, _, err := conn.ReadMessage()
		if t == -1 || err != nil {
			wsCenter.Remove(index)
			conn.Close()
			break
		}
	}
}
