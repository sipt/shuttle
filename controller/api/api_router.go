package api

import "github.com/gin-gonic/gin"

func APIRoute(router *gin.RouterGroup) {
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
	}

	//cert
	router.POST("/cert", GenerateCert)
	router.GET("/cert", DownloadCert)

	//server
	router.GET("/servers", ServerList)
	router.POST("/server/select", SelectServer)
}

type Response struct {
	Code    int
	Message string
	Data    interface{} `json:"omitempty"`
}
