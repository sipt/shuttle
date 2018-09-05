package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle"
	"fmt"
)

func NewShutdown(shutdownSignal chan bool) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.JSON(200, Response{})
		shutdownSignal <- true
	}
}

func ReloadConfig(reloadConfigSignal chan bool) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.JSON(200, Response{})
		reloadConfigSignal <- true
	}
}

func GetConnMode(ctx *gin.Context) {
	ctx.JSON(200, Response{
		Data: shuttle.GetConnMode(),
	})
}

func SetConnMode(ctx *gin.Context) {
	value := ctx.Param("mode")
	err := shuttle.SetConnMode(value)
	if err != nil {
		ctx.JSON(500, Response{
			Code:    1,
			Message: err.Error(),
		})
	}
	GetConnMode(ctx)
}

func Speed(ctx *gin.Context) {
	up, down := shuttle.CurrentSpeed()

	ctx.JSON(200, Response{
		Data: struct {
			UpSpeed   string `json:"up_speed"`
			DownSpeed string `json:"down_speed"`
		}{
			UpSpeed:   fmt.Sprintf("%s/s", capacityConversion(up)),
			DownSpeed: fmt.Sprintf("%s/s", capacityConversion(down)),
		},
	})
}
func capacityConversion(v int) string {
	unit := "B"
	t := v
	if n := t / 1024; n >= 1 {
		unit = "KB"
		t = n
		if n := t / 1024; n >= 1 {
			unit = "MB"
			t = n
			if n := t / 1024; n >= 1 {
				unit = "GB"
				t = n
			}
		}
	}
	return fmt.Sprintf("%d%s", t, unit)
}
