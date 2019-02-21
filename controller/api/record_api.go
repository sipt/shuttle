package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle"
	"github.com/sipt/shuttle/storage"
)

func GetClients(ctx *gin.Context) {
	ctx.JSON(200, &Response{
		Data: storage.Keys(),
	})
}

func GetRecords(ctx *gin.Context) {
	response := Response{}
	clientID := ctx.Query("client_id")
	if len(clientID) == 0 {
		response.Code = 1
		response.Message = "client_id is empty!"
		ctx.JSON(500, response)
		return
	}
	ctx.JSON(200, &Response{
		Data: storage.Get(clientID),
	})
}
func ClearRecords(ctx *gin.Context) {
	clientID := ctx.Query("client_id")
	if len(clientID) == 0 {
		storage.Clear()
	} else {
		storage.Clear(clientID)
	}
	dump := shuttle.GetDump()
	if dump != nil {
		dump.Clear()
	}
	ctx.JSON(200, &Response{})
}
