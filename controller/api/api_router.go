package api

import (
	"github.com/gin-gonic/gin"
	. "github.com/sipt/shuttle/constant"
)

func APIRoute(router *gin.RouterGroup, eventChan chan *EventObj) {
	//dns
	router.GET("/dns", DNSCacheList)
	router.DELETE("/dns", ClearDNSCache)

	//records
	router.GET("/records", GetRecords)
	router.DELETE("/records", ClearRecords)

	//dump
	dump := router.Group("/dump")
	{
		dump.POST("/allow", SetAllowDump)
		dump.GET("/allow", GetAllowDump)
		dump.GET("/data/:conn_id", DumpRequest)
		dump.GET("/large/:conn_id", DumpLarge)
	}

	//cert
	router.POST("/cert", GenerateCert)
	router.GET("/cert", DownloadCert)

	//server
	router.GET("/servers", ServerList)
	router.POST("/server/select", SelectServer)
	router.POST("/server/select/refresh", SelectRefresh)

	//general
	router.GET("/system/proxy/enable", EnableSystemProxy)
	router.GET("/system/proxy/disable", DisableSystemProxy)
	router.POST("/shutdown", NewShutdown(eventChan))
	router.POST("/reload", ReloadConfig(eventChan))
	router.GET("/mode", GetConnMode)
	router.POST("/mode/:mode", SetConnMode)
	router.GET("/upgrade/check", CheckUpdate)
	router.POST("/upgrade", NewUpgrade(eventChan))

	//ws
	router.GET("/ws/records", func(ctx *gin.Context) {
		WsHandler(ctx.Writer, ctx.Request)
	})
	router.GET("/ws/speed", func(ctx *gin.Context) {
		WsSpeedHandler(ctx.Writer, ctx.Request)
	}) // 时速
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
