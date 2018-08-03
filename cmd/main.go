package main

import (
	"net"
	"github.com/sipt/shuttle"
	_ "github.com/sipt/shuttle/ciphers"
	_ "github.com/sipt/shuttle/selector"
)

func main() {
	general, err := shuttle.InitConfig("/Users/sipt/Documents/GOPATH/src/github.com/sipt/shuttle/.conf/sipt.yaml")
	if err != nil {
		panic(err)
	}
	//go HandleUDP()
	go HandleHTTP(general.HttpPort, general.HttpInterface)
	HandleSocks5(general.SocksPort, general.SocksInterface)
}

func HandleSocks5(socksPort, socksInterface string) {
	listener, err := net.Listen("tcp", net.JoinHostPort(socksInterface, socksPort))
	if err != nil {
		panic(err)
	}
	shuttle.Logger.Info("Listen to [SOCKS]: ", net.JoinHostPort(socksInterface, socksPort))
	for {
		conn, err := listener.Accept()
		if err != nil {
			shuttle.Logger.Error(err)
			continue
		}
		go func() {
			defer conn.Close()
			shuttle.Logger.Debug("[SOCKS]Accept tcp connection")
			shuttle.SocksHandle(conn)
		}()
	}
}

func HandleHTTP(httpPort, httpInterface string) {
	listener, err := net.Listen("tcp", net.JoinHostPort(httpInterface, httpPort))
	if err != nil {
		panic(err)
	}
	shuttle.Logger.Info("Listen to [HTTP/HTTPS]: ", net.JoinHostPort(httpInterface, httpPort))
	for {
		conn, err := listener.Accept()
		if err != nil {
			shuttle.Logger.Error(err)
			continue
		}
		go func() {
			defer conn.Close()
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
