package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle"
)

func DNSCacheList(ctx *gin.Context) {
	ctx.JSON(200, &Response{
		Data: shuttle.DNSCacheList(),
	})
}
func ClearDNSCache(ctx *gin.Context) {
	shuttle.ClearDNSCache()
	ctx.JSON(200, &Response{})
}
