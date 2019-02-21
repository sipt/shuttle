package cmd

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"github.com/sipt/shuttle"

	"github.com/sipt/shuttle/config"
	"github.com/sipt/shuttle/controller"
	"github.com/sipt/shuttle/dns"
	"github.com/sipt/shuttle/extension/network"
	"github.com/sipt/shuttle/log"

	_ "github.com/sipt/shuttle/ciphers"
	_ "github.com/sipt/shuttle/proxy/protocol"
	_ "github.com/sipt/shuttle/proxy/selector"
)

var (
	StopSocksSignal = make(chan bool, 1)
	StopHTTPSignal  = make(chan bool, 1)
)

func Run(logMode, logPath, configPath string) {
	var (
		conf *config.Config
		err  error
	)
	//init Logger
	if err = log.InitLogger(logMode, logPath); err != nil {
		fmt.Println(err.Error())
		return
	}
	if conf, err = loadConfig(configPath); err != nil {
		fmt.Println(err.Error())
		return
	}

	//event listen
	ListenEvent()

	// 启动api控制
	go controller.StartController(conf, eventChan)
	//go HandleUDP()
	go HandleHTTP(conf, StopHTTPSignal)
	go HandleSocks5(conf, StopSocksSignal)

	// Catch "Ctrl + C"
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	if conf.General.SetAsSystemProxy == "" || conf.General.SetAsSystemProxy == config.SetAsSystemProxyAuto {
		//enable system proxy
		EnableSystemProxy(conf)
	}
	fmt.Println("success")

	<-signalChan
	log.Logger.Info("[Shuttle] is shutdown, see you later!")
	shutdown(conf.General.SetAsSystemProxy)
	os.Exit(0)
	return
}

func shutdown(setAsSystemProxy string) {
	controller.ShutdownController()
	StopSocksSignal <- true
	StopHTTPSignal <- true
	if setAsSystemProxy == "" || setAsSystemProxy == config.SetAsSystemProxyAuto {
		//disable system proxy
		DisableSystemProxy()
	}
	log.Logger.Close()
	dns.CloseGeoDB()
	time.Sleep(time.Second)
}

func EnableSystemProxy(config IProxyConfig) {
	network.WebProxySwitch(true, "127.0.0.1", config.GetHTTPPort())
	network.SecureWebProxySwitch(true, "127.0.0.1", config.GetHTTPPort())
	network.SocksProxySwitch(true, "127.0.0.1", config.GetSOCKSPort())
}

func DisableSystemProxy() {
	network.WebProxySwitch(false)
	network.SecureWebProxySwitch(false)
	network.SocksProxySwitch(false)
}

type IProxyConfig interface {
	ISOCKSProxyConfig
	IHTTPProxyConfig
}

//SOCKS5 Proxy
type ISOCKSProxyConfig interface {
	GetSOCKSInterface() string
	SetSOCKSInterface(string)
	GetSOCKSPort() string
	SetSOCKSPort(string)
}

func HandleSocks5(config ISOCKSProxyConfig, stopHandle chan bool) {
	addr := net.JoinHostPort(config.GetSOCKSInterface(), config.GetSOCKSPort())
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	log.Logger.Info("Listen to [SOCKS]: ", addr)
	var shutdown = false
	go func() {
		if shutdown = <-stopHandle; shutdown {
			listener.Close()
			log.Logger.Infof("close socks listener!")
		}
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			if shutdown && strings.Contains(err.Error(), "use of closed network connection") {
				log.Logger.Info("Stopped HTTP/HTTPS Proxy goroutine...")
				return
			} else {
				log.Logger.Error(err)
			}
			continue
		}
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Logger.Errorf("[HTTP/HTTPS]panic :%v", err)
					log.Logger.Errorf("[HTTP/HTTPS]stack :%s", debug.Stack())
					conn.Close()
				}
			}()
			log.Logger.Debug("[SOCKS]Accept tcp connection")
			shuttle.SocksHandle(conn)
		}()
	}
}

//HTTP Proxy
type IHTTPProxyConfig interface {
	GetHTTPInterface() string
	SetHTTPInterface(string)
	GetHTTPPort() string
	SetHTTPPort(string)
}

func HandleHTTP(config IHTTPProxyConfig, stopHandle chan bool) {
	addr := net.JoinHostPort(config.GetHTTPInterface(), config.GetHTTPPort())
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	log.Logger.Info("Listen to [HTTP/HTTPS]: ", addr)

	var shutdown = false
	go func() {
		if shutdown = <-stopHandle; shutdown {
			listener.Close()
			log.Logger.Infof("close HTTP/HTTPS listener!")
		}
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			if shutdown && strings.Contains(err.Error(), "use of closed network connection") {
				log.Logger.Info("Stopped HTTP/HTTPS Proxy goroutine...")
				return
			} else {
				log.Logger.Error(err)
			}
			continue
		}
		go func() {
			defer func() {
				conn.Close()
				if err := recover(); err != nil {
					log.Logger.Errorf("[HTTP/HTTPS]panic :%v", err)
					log.Logger.Errorf("[HTTP/HTTPS]stack :%s", debug.Stack())
				}
			}()
			log.Logger.Debug("[HTTP/HTTPS]Accept tcp connection")
			shuttle.HandleHTTP(conn)
		}()
	}
}
