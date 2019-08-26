package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/sipt/shuttle/conf"
	"github.com/sipt/shuttle/constant"
	"github.com/sipt/shuttle/global"
	"github.com/sipt/shuttle/inbound"
	"github.com/sipt/shuttle/rule"
	"github.com/sipt/shuttle/server"
	"github.com/sirupsen/logrus"

	connpkg "github.com/sipt/shuttle/conn"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	params := map[string]string{"path": "config1.toml"}
	config, err := conf.LoadConfig(ctx, "file", "toml", params, func() {
		fmt.Println("config file change")
	})
	if err != nil {
		panic(err)
	}
	err = conf.ApplyConfig(ctx, config)
	if err != nil {
		panic(err)
	}
	err = inbound.ApplyConfig(config, handle)
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	cancel()
}

func handle(conn connpkg.ICtxConn) {
	requseInfo := conn.Value(constant.KeyRequestInfo).(rule.Info)
	// host, _, _ := net.SplitHostPort(connpkg.RemoteAddr().String())
	namespace := global.NamespaceWithName()
	profile := namespace.Profile()
	rule := profile.RuleHandle()(conn, requseInfo)
	if len(requseInfo.IP()) == 0 {
		answer := profile.DNSHandle()(conn, requseInfo.Domain())
		if answer != nil {
			requseInfo.SetIP(answer.CurrentIP)
		}
	}
	var s server.IServer
	g := profile.Group()[rule.Proxy]
	if g == nil {
		s = profile.Server()[rule.Proxy]
	} else {
		s = g.Server()
	}
	sc, err := s.Dial(conn, conn.RemoteAddr().Network(), requseInfo, connpkg.DefaultDial)
	if err != nil {
		logrus.WithField("proxy", rule.Proxy).WithError(err).Errorf("remote to server failed")
		return
	}
	transefer(conn, sc)
}

func transefer(from, to connpkg.ICtxConn) {
	go io.Copy(from, to)
	go io.Copy(to, from)
}
