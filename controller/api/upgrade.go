package api

import (
	"github.com/gin-gonic/gin"
	conf "github.com/sipt/shuttle/config"
	. "github.com/sipt/shuttle/constant"
	"github.com/sipt/shuttle/extension/config"
	"github.com/sipt/shuttle/upgrade"
	"os"
	"path/filepath"
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

func NewUpgrade(eventChan chan *EventObj) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		if status == upgrade.VersionEqual || status == upgrade.VersionGreater {
			ctx.JSON(500, Response{
				Code: 1, Message: "You're up-to-date!",
			})
			return
		}
		path := filepath.Join(config.HomeDir, "Downloads", "shuttle.zip")
		os.Remove(path)
		err := upgrade.DownloadFile(path, url)
		if err != nil {
			ctx.JSON(500, Response{
				Code: 1, Message: err.Error(),
			})
			return
		}
		ctx.JSON(200, Response{
			Code: 0, Message: "success",
		})
		eventChan <- EventUpgrade.SetData("shuttle.zip")
	}
}
