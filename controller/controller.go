package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle/conf/logger"
	"github.com/sipt/shuttle/conf/model"

	dnsapi "github.com/sipt/shuttle/dns/api"
	gapi "github.com/sipt/shuttle/group/api"
	rapi "github.com/sipt/shuttle/rule/api"
	sapi "github.com/sipt/shuttle/server/api"
)

func ApplyConfig(conf *model.Config) (closer func(), err error) {
	e.Use(gin.LoggerWithWriter(logger.Std))
	sapi.InitAPI(e)   // init server api
	gapi.InitAPI(e)   // init group api
	dnsapi.InitAPI(e) // init dns api
	rapi.InitAPI(e)   // init rule api
	if len(conf.Controller.Addr) == 0 {
		return func() {}, nil
	}
	closer, err = InitEngine(conf.Controller.Addr, conf.Controller.Params)
	return
}
