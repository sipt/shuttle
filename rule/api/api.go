package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle/controller/model"
	"github.com/sipt/shuttle/global/namespace"
)

func InitAPI(e *gin.Engine) {
	e.GET("/api/rule/mode", func(c *gin.Context) {
		np := namespace.NamespaceWithContext(c)
		c.JSON(http.StatusOK, &model.Response{
			Code: 0,
			Data: np.Mode(),
		})
	})

	e.PUT("/api/rule/mode/:mode", func(c *gin.Context) {
		np := namespace.NamespaceWithContext(c)
		mode := c.Param("mode")
		if len(mode) == 0 {
			c.JSON(http.StatusBadRequest, &model.Response{
				Code:    1,
				Message: "mode is empty",
			})
			return
		}
		np.SetMode(mode)
		c.JSON(http.StatusOK, &model.Response{})
	})
}
