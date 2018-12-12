package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle/config"
	. "github.com/sipt/shuttle/constant"
	"github.com/sipt/shuttle/extension/network"
	"github.com/sipt/shuttle/rule"
	"strings"
)

func EnableSystemProxy(ctx *gin.Context) {
	g := config.CurrentConfig().General
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

func NewShutdown(eventChan chan *EventObj) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.JSON(200, Response{})
		eventChan <- EventShutdown
	}
}

func ReloadConfig(eventChan chan *EventObj) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.JSON(200, Response{})
		eventChan <- EventReloadConfig
	}
}

func GetConnMode(ctx *gin.Context) {
	ctx.JSON(200, Response{
		Data: rule.GetConnMode(),
	})
}

func SetConnMode(ctx *gin.Context) {
	value := ctx.Param("mode")
	value = strings.ToUpper(value)
	err := rule.SetConnMode(value)
	if err != nil {
		ctx.JSON(500, Response{
			Code:    1,
			Message: err.Error(),
		})
	}
	GetConnMode(ctx)
}
