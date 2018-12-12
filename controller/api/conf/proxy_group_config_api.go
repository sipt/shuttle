package conf

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle/config"
	"github.com/sipt/shuttle/proxy"
	"github.com/sipt/shuttle/rule"
)

func GetProxyGroups(ctx *gin.Context) {
	ctx.JSON(200, &Response{
		Data: proxy.GetGroupExternals(),
	})
}

func SetProxyGroups(ctx *gin.Context) {
	conf := config.CurrentConfig()
	newConf := &config.Config{}
	*newConf = *conf
	data := make(map[string][]string)
	err := ctx.BindJSON(&data)
	if err != nil {
		ctx.JSON(500, &Response{Code: 1, Message: err.Error()})
		return
	}
	newConf.SetProxyGroup(data)
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
	conf.SetProxyGroup(data)
	err = config.SaveConfig(config.CurrentConfigFile(), conf)
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

func GetProxyGroup(ctx *gin.Context) {
	name := ctx.Query("name")
	if len(name) == 0 {
		ctx.JSON(200, &Response{
			Code:    1,
			Message: fmt.Sprintf("name is empty"),
		})
		return
	}
	groups := proxy.GetGroupExternals(name)
	if len(groups) > 0 {
		ctx.JSON(200, &Response{
			Data: groups[0],
		})
		return
	}
	ctx.JSON(500, &Response{
		Code:    1,
		Message: fmt.Sprintf("%s not found", name),
	})
}

func AddProxyGroup(ctx *gin.Context) {
	data := &ProxyRequest{}
	err := ctx.BindJSON(&data)
	if err != nil {
		ctx.JSON(500, &Response{Code: 1, Message: err.Error()})
		return
	}
	if len(data.Name) == 0 {
		ctx.JSON(500, &Response{Code: 1, Message: "ProxyGroup Name is empty"})
		return
	}
	err = proxy.AddGroup(data.Name, data.VS)
	if err != nil {
		ctx.JSON(500, &Response{Code: 1, Message: err.Error()})
		return
	}
	conf := config.CurrentConfig()
	conf.GetProxyGroup()[data.Name] = data.VS
	err = config.SaveConfig(config.CurrentConfigFile(), conf)
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

func EditProxyGroup(ctx *gin.Context) {
	data := &ProxyRequest{}
	err := ctx.BindJSON(&data)
	if err != nil {
		ctx.JSON(500, &Response{Code: 1, Message: err.Error()})
		return
	}
	err = proxy.EditGroup(data.Name, data.VS)
	if err != nil {
		ctx.JSON(500, &Response{Code: 1, Message: err.Error()})
		return
	}
	conf := config.CurrentConfig()
	conf.GetProxyGroup()[data.Name] = data.VS
	err = config.SaveConfig(config.CurrentConfigFile(), conf)
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

func RemoveProxyGroup(ctx *gin.Context) {
	name := ctx.Query("name")
	effects, deletes, err := proxy.RemoveGroup(name)
	if err != nil {
		ctx.JSON(500, &Response{
			Code:    1,
			Message: err.Error(),
		})
		return
	}
	conf := config.CurrentConfig()
	delete(conf.GetProxyGroup(), name)
	err = config.SaveConfig(config.CurrentConfigFile(), conf)
	if err != nil {
		ctx.JSON(500, &Response{
			Code:    1,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(200, &Response{
		Data: struct {
			Effects []string `json:"effects"`
			Deletes []string `json:"deletes"`
		}{effects, deletes},
	})
	return
}
