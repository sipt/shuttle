package web

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/static"
)

func WebRoute(e *gin.Engine) {
	e.Use(static.Serve("/", static.LocalFile("view/", false)))
	e.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(301, "/records")
	})
	e.GET("/records", index)
	e.GET("/dns-cache", index)
	e.GET("/servers", index)
	e.GET("/mitm", index)
}

func index(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.Status(200)
	ctx.File("./view/index.html")
}
