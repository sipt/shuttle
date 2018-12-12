package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle/dns"
)

func DNSCacheList(ctx *gin.Context) {
	ctx.JSON(200, &Response{
		Data: dns.DNSCacheList(),
	})
}
func ClearDNSCache(ctx *gin.Context) {
	dns.ClearDNSCache()
	ctx.JSON(200, &Response{})
}

