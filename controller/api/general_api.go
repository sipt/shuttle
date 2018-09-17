package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle"
	"strings"
)

func NewShutdown(shutdownSignal chan bool) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.JSON(200, Response{})
		shutdownSignal <- true
	}
}

func ReloadConfig(reloadConfigSignal chan bool) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.JSON(200, Response{})
		reloadConfigSignal <- true
	}
}

func GetConnMode(ctx *gin.Context) {
	ctx.JSON(200, Response{
		Data: shuttle.GetConnMode(),
	})
}

func SetConnMode(ctx *gin.Context) {
	value := ctx.Param("mode")
	value = strings.ToUpper(value)
	err := shuttle.SetConnMode(value)
	if err != nil {
		ctx.JSON(500, Response{
			Code:    1,
			Message: err.Error(),
		})
	}
	GetConnMode(ctx)
}
