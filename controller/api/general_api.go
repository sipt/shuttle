package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle"
	"github.com/sipt/shuttle/extension/network"
)

func EnableSystemProxy(ctx *gin.Context) {
	g := shuttle.GetGeneralConfig()
	network.WebProxySwitch(true, "127.0.0.1", g.HttpPort)
	network.SecureWebProxySwitch(true, "127.0.0.1", g.HttpPort)
	network.SocksProxySwitch(true, "127.0.0.1", g.SocksPort)
	ctx.JSON(200, Response{})
}

func DisableSystemProxy(ctx *gin.Context) {
	network.WebProxySwitch(false)
	network.SecureWebProxySwitch(false)
	network.SocksProxySwitch(false)
	ctx.JSON(200, Response{})
}

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
	err := shuttle.SetConnMode(value)
	if err != nil {
		ctx.JSON(500, Response{
			Code:    1,
			Message: err.Error(),
		})
	}
	GetConnMode(ctx)
}
