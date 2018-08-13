package api

import (
	"github.com/sipt/shuttle"
	"github.com/gin-gonic/gin"
)

func GetRecords(ctx *gin.Context) {
	ctx.JSON(200, &Response{
		Data: shuttle.GetRecords(),
	})
}
func ClearRecords(ctx *gin.Context) {
	shuttle.ClearRecords()
	dump := shuttle.GetDump()
	if dump != nil {
		dump.Clear()
	}
	ctx.JSON(200, &Response{})
}
