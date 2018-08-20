package api

import (
	"github.com/gin-gonic/gin"
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
