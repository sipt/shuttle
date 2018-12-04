package conf

import (
	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle/config"
	"github.com/sipt/shuttle/proxy"
	"github.com/sipt/shuttle/rule"
)

func GetProxy(ctx *gin.Context) {
	conf := config.CurrentConfig()
	ctx.JSON(200, &Response{
		Data: conf.GetProxy(),
	})
}

func SetProxy(ctx *gin.Context) {
	conf := config.CurrentConfig()
	newConf := &config.Config{}
	*newConf = *conf
	data := make(map[string][]string)
	err := ctx.BindJSON(&data)
	if err != nil {
		ctx.JSON(500, &Response{Code: 1, Message: err.Error()})
		return
	}
	newConf.SetProxy(data)
	err = proxy.ApplyConfig(newConf)
	if err != nil {
		ctx.JSON(500, &Response{Code: 1, Message: err.Error()})
		return
	}
	err = rule.ApplyConfig(conf)
	if err != nil {
		ctx.JSON(500, &Response{Code: 1, Message: err.Error()})
		return
	}
	conf.SetProxy(data)
	config.SaveConfig(config.CurrentConfigFile(), conf)
	if err != nil {
		ctx.JSON(500, &Response{
			Code:    1,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(200, &Response{})
	return
}
