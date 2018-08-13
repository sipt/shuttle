package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle/controller/api"
)

func StartController(port string) {
	e := gin.Default()
	api.APIRoute(e.Group("/api"))
	e.Run(":" + port)
}
