package main

import (
	"net"
	"github.com/sipt/shuttle"
	_ "github.com/sipt/shuttle/ciphers"
	_ "github.com/sipt/shuttle/selector"
	_ "github.com/sipt/shuttle/protocol"
	"github.com/sipt/shuttle/controller"
	"time"
	"strings"
	"runtime/debug"
	"github.com/sipt/shuttle/extension/network"
	"os"
	"os/signal"
	"syscall"
)

var (
	ShutdownSignal     = make(chan bool, 1)
	StopSocksSignal    = make(chan bool, 1)
	StopHTTPSignal     = make(chan bool, 1)
	ReloadConfigSignal = make(chan bool, 1)
)

func main() {
	var configFile = "shuttle.yaml"
	var geoIPDB = "GeoLite2-Country.mmdb"
	general, err := shuttle.InitConfig(configFile)
	if err != nil {
		panic(err)
	}
	shuttle.InitGeoIP(geoIPDB)
	// 启动api控制
	go controller.StartController(general.ControllerInterface, general.ControllerPort,
		ShutdownSignal,     // shutdown program
		ReloadConfigSignal, // reload config
		general.LogLevel,
	)
	//go HandleUDP()
	go HandleHTTP(general.HttpPort, general.HttpInterface, StopSocksSignal)
	go HandleSocks5(general.SocksPort, general.SocksInterface, StopHTTPSignal)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	//enable system proxy
	EnableSystemProxy(general)
	for {
		select {
		case <-ShutdownSignal:
			StopSocksSignal <- true
			StopHTTPSignal <- true
			//disable system proxy
			DisableSystemProxy()
			time.Sleep(time.Second)
			shuttle.Logger.Info("[Shuttle] is shutdown, see you later!")
			return
		case <-signalChan:
			StopSocksSignal <- true
			StopHTTPSignal <- true
			//disable system proxy
			DisableSystemProxy()
			time.Sleep(time.Second)
			shuttle.Logger.Info("[Shuttle] is shutdown, see you later!")
			return
		case <-ReloadConfigSignal:
			StopSocksSignal <- true
			StopHTTPSignal <- true
			general, err := shuttle.ReloadConfig()
			if err != nil {
				shuttle.Logger.Error("Reload Config failed: ", err)
			}
			//enable system proxy
			EnableSystemProxy(general)
			go HandleHTTP(general.HttpPort, general.HttpInterface, StopSocksSignal)
			go HandleSocks5(general.SocksPort, general.SocksInterface, StopHTTPSignal)
		}
	}
}

func EnableSystemProxy(g *shuttle.General) {
	network.WebProxySwitch(true, "127.0.0.1", g.HttpPort)
	network.SecureWebProxySwitch(true, "127.0.0.1", g.HttpPort)
	network.SocksProxySwitch(true, "127.0.0.1", g.SocksPort)
}

func DisableSystemProxy() {
	network.WebProxySwitch(false)
	network.SecureWebProxySwitch(false)
	network.SocksProxySwitch(false)
}

func HandleSocks5(socksPort, socksInterface string, stopHandle chan bool) {
	listener, err := net.Listen("tcp", net.JoinHostPort(socksInterface, socksPort))
	if err != nil {
		panic(err)
	}
	shuttle.Logger.Info("Listen to [SOCKS]: ", net.JoinHostPort(socksInterface, socksPort))
	var shutdown = false
	go func() {
		if shutdown = <-stopHandle; shutdown {
			listener.Close()
			shuttle.Logger.Infof("close socks listener!")
		}
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			if shutdown && strings.Contains(err.Error(), "use of closed network connection") {
				shuttle.Logger.Info("Stopped HTTP/HTTPS Proxy goroutine...")
				return
			} else {
				shuttle.Logger.Error(err)
			}
			continue
		}
		go func() {
			defer func() {
				if err := recover(); err != nil {
					shuttle.Logger.Error("[HTTP/HTTPS]panic :", err)
					shuttle.Logger.Error("[HTTP/HTTPS]stack :", debug.Stack())
					conn.Close()
				}
			}()
			shuttle.Logger.Debug("[SOCKS]Accept tcp connection")
			shuttle.SocksHandle(conn)
		}()
	}
}
func HandleHTTP(httpPort, httpInterface string, stopHandle chan bool) {
	listener, err := net.Listen("tcp", net.JoinHostPort(httpInterface, httpPort))
	if err != nil {
		panic(err)
	}
	shuttle.Logger.Info("Listen to [HTTP/HTTPS]: ", net.JoinHostPort(httpInterface, httpPort))

	var shutdown = false
	go func() {
		if shutdown = <-stopHandle; shutdown {
			listener.Close()
			shuttle.Logger.Infof("close HTTP/HTTPS listener!")
		}
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			if shutdown && strings.Contains(err.Error(), "use of closed network connection") {
				shuttle.Logger.Info("Stopped HTTP/HTTPS Proxy goroutine...")
				return
			} else {
				shuttle.Logger.Error(err)
			}
			continue
		}
		go func() {
			defer func() {
				conn.Close()
				if err := recover(); err != nil {
					shuttle.Logger.Errorf("[HTTP/HTTPS]panic :%v", err)
					shuttle.Logger.Errorf("[HTTP/HTTPS]stack :%s", debug.Stack())
				}
			}()
			shuttle.Logger.Debug("[HTTP/HTTPS]Accept tcp connection")
			shuttle.HandleHTTP(conn)
		}()
	}
}
func HandleUDP() {
	var port = "8080"
	listener, err := net.Listen("udp", ":"+port)
	if err != nil {
		panic(err)
	}
	shuttle.Logger.Info("Listen to [udp]: ", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			shuttle.Logger.Error(err)
			continue
		}
		go func() {
			shuttle.Logger.Info("Accept tcp connection")
			shuttle.SocksHandle(conn)
		}()
	}
}
