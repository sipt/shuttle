package conf

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle/config"
	"github.com/sipt/shuttle/proxy"
	"github.com/sipt/shuttle/rule"
)

func GetProxies(ctx *gin.Context) {
	ctx.JSON(200, &Response{
		Data: proxy.GetServerExternals(),
	})
}

func SetProxies(ctx *gin.Context) {
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

type ProxyRequest struct {
	Name string   `json:"name"`
	VS   []string `json:"vs"`
}

func GetProxy(ctx *gin.Context) {
	name := ctx.Query("name")
	if len(name) == 0 {
		ctx.JSON(200, &Response{
			Code:    1,
			Message: fmt.Sprintf("name is empty"),
		})
		return
	}
	proxyConf := config.CurrentConfig().GetProxy()[name]
	if len(proxyConf) > 0 {
		ctx.JSON(200, &Response{
			Data: &ProxyRequest{
				Name: name,
				VS:   proxyConf,
			},
		})
		return
	}
	ctx.JSON(500, &Response{
		Code:    1,
		Message: fmt.Sprintf("%s not found", name),
	})
}

func AddProxy(ctx *gin.Context) {
	data := &ProxyRequest{}
	err := ctx.BindJSON(&data)
	if err != nil {
		ctx.JSON(500, &Response{Code: 1, Message: err.Error()})
		return
	}
	if len(data.Name) == 0 {
		ctx.JSON(500, &Response{Code: 1, Message: "Proxy Name is empty"})
		return
	}
	err = proxy.AddProxy(data.Name, data.VS)
	if err != nil {
		ctx.JSON(500, &Response{Code: 1, Message: err.Error()})
		return
	}
	conf := config.CurrentConfig()
	conf.GetProxy()[data.Name] = data.VS
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

func EditProxy(ctx *gin.Context) {
	data := &ProxyRequest{}
	err := ctx.BindJSON(&data)
	if err != nil {
		ctx.JSON(500, &Response{Code: 1, Message: err.Error()})
		return
	}
	err = proxy.EditProxy(data.Name, data.VS)
	if err != nil {
		ctx.JSON(500, &Response{Code: 1, Message: err.Error()})
		return
	}
	conf := config.CurrentConfig()
	conf.GetProxy()[data.Name] = data.VS
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

func RemoveProxy(ctx *gin.Context) {
	name := ctx.Query("name")
	effects, deletes, err := proxy.RemoveProxy(name)
	if err != nil {
		ctx.JSON(500, &Response{
			Code:    1,
			Message: err.Error(),
		})
		return
	}
	conf := config.CurrentConfig()
	delete(conf.GetProxy(), name)
	groups := conf.GetProxyGroup()
	for _, v := range effects {
		for i, s := range groups[v] {
			if s == name {
				groups[v] = append(groups[v][:i], groups[v][i+1:]...)
				break
			}
		}
	}
	for _, v := range deletes {
		delete(groups, v)
	}
	conf.SetProxyGroup(groups)
	conf.GetProxy()
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

func Policies(ctx *gin.Context) {
	conf := config.CurrentConfig()
	proxy := conf.GetProxy()
	group := conf.GetProxyGroup()
	policies := make([]string, 2, 2+len(proxy)+len(group))
	policies[0] = "DIRECT"
	policies[1] = "REJECT"
	for k := range proxy {
		policies = append(policies, k)
	}
	for k := range group {
		policies = append(policies, k)
	}
	ctx.JSON(200, &Response{
		Data: policies,
	})
	return
}
