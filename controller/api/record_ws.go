package api

import (
	"github.com/gorilla/websocket"
	"time"
	"net/http"
	"sync"
	"github.com/sipt/shuttle/util"
	"github.com/sipt/shuttle"
)

func init() {
	shuttle.RegisterPusher(func(v interface{}) {
		wsCenter.RangeSend(v)
	})
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
			c.index[k] = c.index[len(c.index)]
			c.index = c.index[:len(c.index)-1]
			c.conns[k] = c.conns[len(c.conns)]
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
		shuttle.Logger.Errorf("Failed to set websocket upgrade: %v", err)
		return
	}
	index, _ := util.IW.NextId()
	wsCenter.Add(index, conn)
	for {
		_, _, err = conn.ReadMessage()
		if err != nil {
			wsCenter.Remove(index)
			break
		}
	}
}
