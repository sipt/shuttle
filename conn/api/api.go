package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle/conn"
)

func InitAPI(e *gin.Engine) {
	e.GET("/api/conns", connectionsHandleFunc)
}

type LinkResponse struct {
	Input  *ConnEntity `json:"input"`
	Output *ConnEntity `json:"output"`
}

type ConnEntity struct {
	ID         int64  `json:"id"`
	LocalAddr  string `json:"local_addr"`
	RemoteAddr string `json:"remote_addr"`
}

func connectionsHandleFunc(c *gin.Context) {
	list := make([]*LinkResponse, 0)
	conn.PoolRange(func(link *conn.Link) bool {
		l := &LinkResponse{}
		if link.Input != nil {
			l.Input = &ConnEntity{
				ID:         link.Input.GetConnID(),
				LocalAddr:  link.Input.LocalAddr().String(),
				RemoteAddr: link.Input.RemoteAddr().String(),
			}
		}
		if link.Output != nil {
			l.Output = &ConnEntity{
				ID:         link.Output.GetConnID(),
				LocalAddr:  link.Output.LocalAddr().String(),
				RemoteAddr: link.Output.RemoteAddr().String(),
			}
		}
		list = append(list, l)
		return false
	})
	c.JSON(200, list)
}
