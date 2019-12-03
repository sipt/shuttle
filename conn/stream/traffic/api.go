package stream

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

func InitAPI(e *gin.Engine) {
	e.GET("/ws/traffic", trafficHandleFunc)
}

func trafficHandleFunc(c *gin.Context) {
	handler := websocket.Handler(func(conn *websocket.Conn) {
		defer conn.Close()
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			_, err := conn.Write([]byte(fmt.Sprintf(`{"up": %d, "down":%d}`, up, down)))
			if err != nil {
				return
			}
			<-ticker.C
		}
	})
	handler.ServeHTTP(c.Writer, c.Request)
}
