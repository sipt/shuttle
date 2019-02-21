package cmd

import (
	"github.com/sipt/shuttle"
	"github.com/sipt/shuttle/config"
	"github.com/sipt/shuttle/conn"
	"github.com/sipt/shuttle/controller"
	"github.com/sipt/shuttle/dns"
	"github.com/sipt/shuttle/log"
	"github.com/sipt/shuttle/proxy"
	"github.com/sipt/shuttle/rule"
	"github.com/sipt/shuttle/storage"
)

//load config
func loadConfig(configPath string) (conf *config.Config, err error) {
	//init Config
	conf, err = config.LoadConfig(configPath)
	if err != nil {
		return
	}
	//init Config Value
	shuttle.InitConfigValue(conf)
	//init DNS & GeoIP
	if err = dns.ApplyConfig(conf); err != nil {
		return
	}
	//init Logger
	if err = log.ApplyConfig(conf); err != nil {
		return
	}
	//init Proxy & ProxyGroup
	if err = proxy.ApplyConfig(conf); err != nil {
		return
	}
	//init Rule
	if err = rule.ApplyConfig(conf); err != nil {
		return
	}
	//init Record-Storage
	if err = storage.ApplyConfig(conf); err != nil {
		return
	}
	storage.Prepare()
	//init HttpMap
	if err = shuttle.ApplyHTTPModifyConfig(conf); err != nil {
		return
	}
	//init MITM
	if err = shuttle.ApplyMITMConfig(conf); err != nil {
		return
	}
	//init conn
	conn.Init()
	return
}

//reload config
func reloadConfig(configPath string, StopSocksSignal, StopHTTPSignal chan bool) (conf *config.Config, err error) {
	oldConf := config.CurrentConfig()
	conf, err = loadConfig(configPath)
	if err != nil {
		return
	}
	// controller
	if oldConf.GetControllerInterface() != conf.GetControllerInterface() ||
		oldConf.GetControllerPort() != conf.GetControllerPort() {
		//restart controller
		controller.ShutdownController()
		// 启动api控制
		go controller.StartController(conf, eventChan)
	}

	// http proxy
	if oldConf.GetHTTPInterface() != conf.GetHTTPInterface() ||
		oldConf.GetHTTPPort() != conf.GetHTTPPort() {
		//restart http proxy
		StopHTTPSignal <- true
		go HandleHTTP(conf, StopHTTPSignal)
	}

	// socks5 proxy
	if oldConf.GetSOCKSInterface() != conf.GetSOCKSInterface() ||
		oldConf.GetSOCKSPort() != conf.GetSOCKSPort() {
		//restart http proxy
		StopSocksSignal <- true
		go HandleSocks5(conf, StopSocksSignal)
	}
	return
}
