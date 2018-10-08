package main

import (
	"github.com/sipt/shuttle"
	_ "github.com/sipt/shuttle/ciphers"
	"github.com/sipt/shuttle/controller"
	"github.com/sipt/shuttle/extension/config"
	"github.com/sipt/shuttle/extension/network"
	"github.com/sipt/shuttle/log"
	_ "github.com/sipt/shuttle/protocol"
	_ "github.com/sipt/shuttle/selector"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"syscall"
	"time"
)

var (
	ShutdownSignal     = make(chan bool, 1)
	UpgradeSignal      = make(chan string, 1)
	StopSocksSignal    = make(chan bool, 1)
	StopHTTPSignal     = make(chan bool, 1)
	ReloadConfigSignal = make(chan bool, 1)
)

func configPath() (fullPath string, err error) {
	var configFile = "shuttle.yaml"
	dir := config.ShuttleHomeDir
	fullPath = dir + string(os.PathSeparator) + configFile
	rc, err := os.Open(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		} else {
			return
		}
	} else {
		rc.Close()
		return
	}
	cc, err := os.Open("shuttle.yaml")
	if err != nil {
		return
	}
	defer cc.Close()
	// not exist
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return
	}
	//
	dc, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer dc.Close()
	_, err = io.Copy(dc, cc)
	return
}

func main() {
	configPath, err := configPath()
	if err != nil {
		log.Logger.Errorf("[PANIC] [ConfigPath] %s", err.Error())
		return
	}
	general, err := shuttle.InitConfig(configPath)
	if err != nil {
		log.Logger.Errorf("[PANIC] [InitConfig] %s", err.Error())
		return
	}
	var geoIPDB = "GeoLite2-Country.mmdb"
	err = shuttle.InitGeoIP(geoIPDB)
	if err != nil {
		log.Logger.Errorf("[PANIC] [InitGeoIP] %s", err.Error())
		return
	}
	// 启动api控制
	go controller.StartController(general.ControllerInterface, general.ControllerPort,
		ShutdownSignal,     // shutdown program
		ReloadConfigSignal, // reload config
		UpgradeSignal,      // upgrade
		general.LogLevel,
	)
	//go HandleUDP()
	go HandleHTTP(general.HttpPort, general.HttpInterface, StopSocksSignal)
	go HandleSocks5(general.SocksPort, general.SocksInterface, StopHTTPSignal)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	if general.SetAsSystemProxy == "" || general.SetAsSystemProxy == shuttle.SetAsSystemProxyAuto {
		//enable system proxy
		EnableSystemProxy(general)
	}
	for {
		select {
		case fileName := <-UpgradeSignal:
			shutdown()
			log.Logger.Info("[Shuttle] is shutdown, for upgrade!")
			var name string
			if runtime.GOOS == "windows" {
				name = "upgrade"
			} else {
				name = "./upgrade"
			}
			cmd := exec.Command(name, "-f="+fileName)
			err = cmd.Start()
			if err != nil {
				ioutil.WriteFile(filepath.Join(config.ShuttleHomeDir, "logs", "error.log"), []byte(err.Error()), 0664)
			}
			ioutil.WriteFile(filepath.Join(config.ShuttleHomeDir, "logs", "end.log"), []byte("ending"), 0664)
			os.Exit(0)
		case <-ShutdownSignal:
			log.Logger.Info("[Shuttle] is shutdown, see you later!")
			shutdown()
			os.Exit(0)
			return
		case <-signalChan:
			log.Logger.Info("[Shuttle] is shutdown, see you later!")
			shutdown()
			os.Exit(0)
			return
		case <-ReloadConfigSignal:
			StopSocksSignal <- true
			StopHTTPSignal <- true
			general, err := shuttle.ReloadConfig()
			if err != nil {
				log.Logger.Error("Reload Config failed: ", err)
			}
			if general.SetAsSystemProxy == "" || general.SetAsSystemProxy == shuttle.SetAsSystemProxyAuto {
				//enable system proxy
				EnableSystemProxy(general)
			}
			go HandleHTTP(general.HttpPort, general.HttpInterface, StopSocksSignal)
			go HandleSocks5(general.SocksPort, general.SocksInterface, StopHTTPSignal)
		}
	}
}

func shutdown() {
	StopSocksSignal <- true
	StopHTTPSignal <- true
	//disable system proxy
	DisableSystemProxy()
	log.Logger.Close()
	shuttle.CloseGeoDB()
	time.Sleep(time.Second)
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
	log.Logger.Info("Listen to [SOCKS]: ", net.JoinHostPort(socksInterface, socksPort))
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
					log.Logger.Error("[HTTP/HTTPS]panic :", err)
					log.Logger.Error("[HTTP/HTTPS]stack :", debug.Stack())
					conn.Close()
				}
			}()
			log.Logger.Debug("[SOCKS]Accept tcp connection")
			shuttle.SocksHandle(conn)
		}()
	}
}
func HandleHTTP(httpPort, httpInterface string, stopHandle chan bool) {
	listener, err := net.Listen("tcp", net.JoinHostPort(httpInterface, httpPort))
	if err != nil {
		panic(err)
	}
	log.Logger.Info("Listen to [HTTP/HTTPS]: ", net.JoinHostPort(httpInterface, httpPort))

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
func HandleUDP() {
	var port = "8080"
	listener, err := net.Listen("udp", ":"+port)
	if err != nil {
		panic(err)
	}
	log.Logger.Info("Listen to [udp]: ", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Logger.Error(err)
			continue
		}
		go func() {
			log.Logger.Info("Accept tcp connection")
			shuttle.SocksHandle(conn)
		}()
	}
}
