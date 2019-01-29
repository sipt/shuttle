package api

import (
	"github.com/gin-gonic/gin"
	conf "github.com/sipt/shuttle/config"
	"github.com/sipt/shuttle/upgrade"
)

var latest string
var url string
var status string

func CheckUpdate(ctx *gin.Context) {
	var err error
	latest, url, status, err = upgrade.CheckUpgrade(conf.ShuttleVersion)
	if err != nil {
		ctx.JSON(500, Response{
			Code: 1, Message: err.Error(),
		})
		return
	}
	ctx.JSON(200, Response{
		Code: 0,
		Data: map[string]string{
			"Current": conf.ShuttleVersion,
			"Latest":  latest,
			"URL":     url,
			"Status":  status,
		},
	})
}
