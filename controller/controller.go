package controller

import (
	"github.com/sipt/shuttle/conf/model"

	dnsapi "github.com/sipt/shuttle/dns/api"
	gapi "github.com/sipt/shuttle/group/api"
	sapi "github.com/sipt/shuttle/server/api"
)

func ApplyConfig(conf *model.Config) (closer func(), err error) {
	sapi.InitAPI(e)   // init server api
	gapi.InitAPI(e)   // init group api
	dnsapi.InitAPI(e) // init dns api
	if len(conf.Controller.Addr) == 0 {
		return func() {}, nil
	}
	closer, err = InitEngine(conf.Controller.Addr, conf.Controller.Params)
	return
}
