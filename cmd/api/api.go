package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle/controller/model"
	"github.com/sipt/shuttle/global/namespace"
)

const (
	StatusStarting = "starting"
	StatusRunning  = "running"
	StatusStopped  = "stopped"
)

var Status = StatusStopped
var StartFunc func() error
var CloseFunc func() error
var CheckConfig func() error

func InitAPI(e *gin.Engine) {
	e.GET("/api/status", getStatus)
	e.GET("/api/inbounds", inbound)
	e.PUT("/api/reload", reload)
}

func getStatus(c *gin.Context) {
	c.JSON(http.StatusOK, &model.Response{
		Code: 0,
		Data: Status,
	})
}

func inbound(c *gin.Context) {
	// TODO: support namespace
	np := namespace.NamespaceWithName()
	conf := np.Profile().Config()
	type listener struct {
		Name string `json:"name"`
		Typ  string `json:"typ"`
		Addr string `json:"addr"`
	}
	inbounds := make([]*listener, 0, len(conf.Listener))
	for _, v := range conf.Listener {
		inbounds = append(inbounds, &listener{
			Name: v.Name,
			Typ:  v.Typ,
			Addr: v.Addr,
		})
	}
	c.JSON(http.StatusOK, &model.Response{
		Code: 0,
		Data: inbounds,
	})
}

func reload(c *gin.Context) {
	err := CheckConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, &model.Response{
			Code:    1,
			Message: fmt.Sprintf("load config file failed: %s", err.Error()),
		})
	}
	c.JSON(http.StatusOK, &model.Response{
		Code: 0,
		Data: "success",
	})
	go func() {
		_ = CloseFunc()
		_ = StartFunc()
	}()
}
