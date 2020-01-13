package stream

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:   2048,
	WriteBufferSize:  2048,
	HandshakeTimeout: 5 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func InitAPI(e *gin.Engine) {
	e.GET("/ws/traffic", trafficHandleFunc)
}

func trafficHandleFunc(c *gin.Context) {
	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	go func() {
		defer conn.Close()
		for {
			typ, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if typ == websocket.CloseMessage {
				return
			}
		}
	}()
	for {
		err = conn.WriteJSON(gin.H{"up": up, "down": down})
		if err != nil {
			return
		}
		<-ticker.C
	}
}
