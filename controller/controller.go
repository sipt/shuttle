package controller

import "github.com/sipt/shuttle/conf/model"

func ApplyConfig(conf *model.Config) (closer func(), err error) {
	if len(conf.Controller.Addr) == 0 {
		return func() {}, nil
	}
	closer, err = InitEngine(conf.Controller.Addr, conf.Controller.Params)
	return
}
