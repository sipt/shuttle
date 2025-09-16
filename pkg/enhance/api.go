package enhance

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle/controller/model"
)

func InitAPI(e *gin.Engine) {
	e.GET("/api/enhance/status", getStatus)
	e.POST("/api/enhance/start", start)
	e.POST("/api/enhance/stop", stop)
}

func getStatus(c *gin.Context) {
	state := "unknown"
	switch EnhanceModeInstance.GetState() {
	case EnhanceModeStateRunning:
		state = "running"
	case EnhanceModeStateStopped:
		state = "stopped"
	}

	c.JSON(http.StatusOK, &model.Response[string]{
		Code:    0,
		Message: "success",
		Data:    state,
	})
}

func start(c *gin.Context) {
	err := EnhanceModeInstance.Start()
	if err != nil {
		c.JSON(http.StatusOK, &model.Response[string]{
			Code:    1,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &model.Response[string]{
		Code:    0,
		Message: "success",
	})
}

func stop(c *gin.Context) {
	err := EnhanceModeInstance.Stop()
	if err != nil {
		c.JSON(http.StatusOK, &model.Response[string]{
			Code:    1,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &model.Response[string]{
		Code:    0,
		Message: "success",
	})
}
