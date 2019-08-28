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
	"github.com/sipt/shuttle/server"
	"github.com/sirupsen/logrus"

	connpkg "github.com/sipt/shuttle/conn"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
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

	logrus.Info("server starting...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	cancel()
}

func handle(conn connpkg.ICtxConn) {
	reqInfo := conn.Value(constant.KeyRequestInfo).(global.RequestInfo)
	// host, _, _ := net.SplitHostPort(connpkg.RemoteAddr().String())
	namespace := global.NamespaceWithName()
	profile := namespace.Profile()
	rule := profile.RuleHandle()(conn, reqInfo)
	if len(reqInfo.IP()) == 0 {
		answer := profile.DNSHandle()(conn, reqInfo.Domain())
		if answer != nil {
			reqInfo.SetIP(answer.CurrentIP)
		}
	}
	logrus.Infof("Match Rule [%s, %s, %s]", rule.Typ, rule.Value, rule.Proxy)
	var s server.IServer
	g := profile.Group()[rule.Proxy]
	if g == nil {
		s = profile.Server()[rule.Proxy]
	} else {
		s = g.Server()
	}
	sc, err := s.Dial(conn, reqInfo.Network(), reqInfo, connpkg.DefaultDial)
	if err != nil {
		logrus.WithField("proxy", rule.Proxy).WithError(err).Errorf("remote to server failed")
		return
	}
	logrus.Debug(reqInfo)
	transefer(conn, sc)
}

func transefer(from, to connpkg.ICtxConn) {
	go io.Copy(from, to)
	io.Copy(to, from)
}
