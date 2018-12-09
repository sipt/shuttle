package conf

import (
	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle/config"
	"github.com/sipt/shuttle/dns"
)

type DNSConfig struct {
	Servers  []string   `json:"servers"`
	LocalDNS [][]string `json:"local_dns"`
}

func GetDNSConfig(ctx *gin.Context) {
	conf := config.CurrentConfig()
	ctx.JSON(200, &Response{
		Data: &DNSConfig{
			Servers:  conf.GetDNSServers(),
			LocalDNS: conf.GetLocalDNS(),
		},
	})
}

func SetDNSConfig(ctx *gin.Context) {
	body := &DNSConfig{}
	err := ctx.BindJSON(body)
	if err != nil {
		ctx.JSON(500, &Response{
			Code:    1,
			Message: err.Error(),
		})
		return
	}

	if len(body.Servers) <= 0 {
		ctx.JSON(500, &Response{
			Code:    1,
			Message: "DNS Servers is empty",
		})
		return
	}
	var conf = &config.Config{}
	*conf = *(config.CurrentConfig())
	if len(body.LocalDNS) > 0 {
		conf.SetLocalDNS(body.LocalDNS)
	}
	conf.SetDNSServers(body.Servers)
	err = dns.ApplyConfig(conf)
	if err != nil {
		ctx.JSON(500, &Response{
			Code:    1,
			Message: err.Error(),
		})
		return
	}
	config.CurrentConfig().SetDNSServers(body.Servers)
	config.CurrentConfig().SetLocalDNS(body.LocalDNS)
	err = config.SaveConfig(config.CurrentConfigFile(), config.CurrentConfig())
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
