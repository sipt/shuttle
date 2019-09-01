package controller

import (
	"github.com/sipt/shuttle/conf/model"
	sapi "github.com/sipt/shuttle/server/api"
)

func ApplyConfig(conf *model.Config) (closer func(), err error) {
	sapi.InitAPI(e) // inti server api
	if len(conf.Controller.Addr) == 0 {
		return func() {}, nil
	}
	closer, err = InitEngine(conf.Controller.Addr, conf.Controller.Params)
	return
}
