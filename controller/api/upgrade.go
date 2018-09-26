package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle"
	"github.com/sipt/shuttle/extension/config"
	"github.com/sipt/shuttle/upgrade"
	"os"
	"path/filepath"
)

var latest string
var url string

func CheckUpdate(ctx *gin.Context) {
	var err error
	latest, url, err = upgrade.CheckUpgrade(shuttle.ShuttleVersion)
	if err != nil {
		ctx.JSON(500, Response{
			Code: 1, Message: err.Error(),
		})
		return
	}
	ctx.JSON(200, Response{
		Code: 0,
		Data: map[string]string{
			"Current": shuttle.ShuttleVersion,
			"Latest":  latest,
			"URL":     url,
		},
	})
}

func NewUpgrade(upgradeSignal chan string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		if len(latest) == 0 {
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
		upgradeSignal <- "shuttle.zip"
	}
}
