package web

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/static"
)

const indexHtml = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Shuttle</title>
  <base href="/">

  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="icon" type="image/x-icon" href="/assets/favicon.ico">
</head>
<body>
  <app-root></app-root>
<script type="text/javascript" src="runtime.js"></script><script type="text/javascript" src="polyfills.js"></script><script type="text/javascript" src="styles.js"></script><script type="text/javascript" src="vendor.js"></script><script type="text/javascript" src="main.js"></script></body>
</html>`

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
