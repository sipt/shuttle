package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle/conf/logger"
	"github.com/sipt/shuttle/conf/model"

	capi "github.com/sipt/shuttle/cmd/api"
	tws "github.com/sipt/shuttle/conn/stream/traffic"
	dnsapi "github.com/sipt/shuttle/dns/api"
	gapi "github.com/sipt/shuttle/group/api"
	rapi "github.com/sipt/shuttle/rule/api"
	sapi "github.com/sipt/shuttle/server/api"
)

func init() {
	e.Use(gin.LoggerWithWriter(logger.Std))
	// api
	sapi.InitAPI(e)   // init server api
	gapi.InitAPI(e)   // init group api
	dnsapi.InitAPI(e) // init dns api
	rapi.InitAPI(e)   // init rule api
	capi.InitAPI(e)   // cmd api

	// ws
	tws.InitAPI(e)
}

func ApplyConfig(conf *model.Config) (closer func(), err error) {
	if len(conf.Controller.Addr) == 0 {
		return func() {}, nil
	}
	closer, err = InitEngine(conf.Controller.Addr, conf.Controller.Params)
	return
}
