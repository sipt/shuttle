package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle/conn"
)

func ConnClients(ctx *gin.Context) {
	ctx.JSON(200, &Response{
		Data: conn.GetClients(),
	})
}

func ConnList(ctx *gin.Context) {
	clientID := ctx.Query("client_id")
	if len(clientID) == 0 {
		ctx.JSON(500, &Response{
			Code: 1, Message: "client_id is empty",
		})
		return
	}
	l := conn.GetPool(clientID)
	if l != nil {
		cs := l.List()
		var connList = make([]*struct {
			Left  *conn.ConnLine
			Right *conn.ConnLine
		}, len(cs))
		for k, v := range l.List() {
			connList[k] = &struct {
				Left  *conn.ConnLine
				Right *conn.ConnLine
			}{Left: v.ClientLine, Right: v.ServerLine}
		}
		ctx.JSON(200, &Response{
			Data: connList,
		})
	}
	ctx.JSON(200, &Response{Data: []struct{}{}})

}
