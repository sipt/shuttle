package main

import (
	"net"
	"github.com/sipt/shuttle"
	_ "github.com/sipt/shuttle/ciphers"
	_ "github.com/sipt/shuttle/selector"
)

func main() {
	err := shuttle.InitConfig()
	if err != nil {
		panic(err)
	}
	//go HandleUDP()
	go HandleHTTP()
	HandleSocks5()
}

func HandleSocks5() {
	var port = "8080"
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	shuttle.Logger.Info("Listen to [SOCKS]: ", port)
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

func HandleHTTP() {
	var port = "8081"
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	shuttle.Logger.Info("Listen to [HTTP/HTTPS]: ", port)
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
