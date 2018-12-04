package conf

import (
	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle/config"
	. "github.com/sipt/shuttle/constant"
	"github.com/sipt/shuttle/log"
	"net"
	"strconv"
	"sync"
)

type GeneralConfig struct {
	LogLevel            string `json:"log_level"`
	HttpPort            string `json:"http_port"`
	HttpInterface       string `json:"http_interface"`
	SocksPort           string `json:"socks_port"`
	SocksInterface      string `json:"socks_interface"`
	ControllerPort      string `json:"controller_port"`
	ControllerInterface string `json:"controller_interface"`
}

var general sync.RWMutex

func GetGeneralConfig(ctx *gin.Context) {
	general.RLock()
	defer general.RUnlock()
	conf := config.CurrentConfig()
	ctx.JSON(200, &Response{
		Data: &GeneralConfig{
			LogLevel:            conf.GetLogLevel(),
			HttpPort:            conf.GetHTTPPort(),
			HttpInterface:       conf.GetHTTPInterface(),
			SocksPort:           conf.GetSOCKSPort(),
			SocksInterface:      conf.GetSOCKSInterface(),
			ControllerPort:      conf.GetControllerPort(),
			ControllerInterface: conf.GetControllerInterface(),
		},
	})
}

func SetGeneralConfig(eventChan chan *EventObj) func(*gin.Context) {
	return func(ctx *gin.Context) {
		general.Lock()
		defer general.Unlock()
		body := &GeneralConfig{}
		err := ctx.BindJSON(body)
		if err != nil {
			ctx.JSON(500, &Response{
				Code:    1,
				Message: err.Error(),
			})
			return
		}
		oldConf := config.CurrentConfig()
		newConf := &config.Config{}
		*newConf = *oldConf
		changed := false
		if body.LogLevel != oldConf.GetLogLevel() {
			newConf.SetLogLevel(body.LogLevel)
			err = log.ApplyConfig(newConf)
			if err != nil {
				log.ApplyConfig(oldConf)
				ctx.JSON(500, &Response{Code: 1, Message: err.Error()})
				return
			}
			oldConf.SetLogLevel(body.LogLevel)
			changed = true
		}
		//http
		if body.HttpInterface != oldConf.GetHTTPInterface() || body.HttpPort != oldConf.GetHTTPPort() {
			if ip := net.ParseIP(body.HttpInterface); ip == nil {
				ctx.JSON(500, &Response{Code: 1, Message: "http_interface incorrect"})
				return
			}
			if _, err := strconv.Atoi(body.HttpPort); err != nil {
				ctx.JSON(500, &Response{Code: 1, Message: "http_port incorrect"})
				return
			}
			oldConf.SetHTTPInterface(body.HttpInterface)
			oldConf.SetHTTPPort(body.HttpPort)
			eventChan <- EventRestartHttpProxy
			changed = true
		}
		//socks
		if body.SocksInterface != oldConf.GetSOCKSInterface() || body.SocksPort != oldConf.GetSOCKSPort() {
			if ip := net.ParseIP(body.SocksInterface); ip == nil {
				ctx.JSON(500, &Response{Code: 1, Message: "socks_interface incorrect"})
				return
			}
			if _, err := strconv.Atoi(body.SocksPort); err != nil {
				ctx.JSON(500, &Response{Code: 1, Message: "socks_port incorrect"})
				return
			}
			oldConf.SetSOCKSInterface(body.SocksInterface)
			oldConf.SetSOCKSPort(body.SocksPort)
			eventChan <- EventRestartSocksProxy
			changed = true
		}
		//controller
		if body.ControllerInterface != oldConf.GetControllerInterface() ||
			body.ControllerPort != oldConf.GetControllerPort() {
			if ip := net.ParseIP(body.ControllerInterface); ip == nil {
				ctx.JSON(500, &Response{Code: 1, Message: "controller_interface incorrect"})
				return
			}
			if _, err := strconv.Atoi(body.ControllerPort); err != nil {
				ctx.JSON(500, &Response{Code: 1, Message: "controller_port incorrect"})
				return
			}
			oldConf.SetControllerInterface(body.ControllerInterface)
			oldConf.SetControllerPort(body.ControllerPort)
			eventChan <- EventRestartController
			changed = true
		}
		if changed {
			config.SaveConfig(config.CurrentConfigFile(), oldConf)
			if err != nil {
				ctx.JSON(500, &Response{
					Code:    1,
					Message: err.Error(),
				})
				return
			}
		}
		ctx.JSON(200, &Response{})
		return
	}
}
