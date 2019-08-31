package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/sipt/shuttle/conf"
	"github.com/sipt/shuttle/constant"
	"github.com/sipt/shuttle/constant/typ"
	"github.com/sipt/shuttle/global"
	"github.com/sipt/shuttle/inbound"
	"github.com/sipt/shuttle/server"
	"github.com/sirupsen/logrus"

	connpkg "github.com/sipt/shuttle/conn"
	rulepkg "github.com/sipt/shuttle/rule"
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
	err = inbound.ApplyConfig(config, handle())
	if err != nil {
		panic(err)
	}

	logrus.Info("server starting...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	cancel()
}

func handle() typ.HandleFunc {
	handle := outboundHandle()
	handle = ruleHandle(handle)
	handle = namespaceHandle(handle)
	return handle
}

func namespaceHandle(next typ.HandleFunc) typ.HandleFunc {
	return func(conn connpkg.ICtxConn) {
		if p, ok := conn.Value(constant.KeyProfile).(*global.Profile); !ok || p == nil {
			np := global.NamespaceWithContext(conn)
			conn.WithValue(constant.KeyProfile, np.Profile())
		}
		next(conn)
	}
}

func ruleHandle(next typ.HandleFunc) typ.HandleFunc {
	return func(conn connpkg.ICtxConn) {
		reqInfo := conn.Value(constant.KeyRequestInfo).(global.RequestInfo)
		profile := conn.Value(constant.KeyProfile).(*global.Profile)

		rule := profile.RuleHandle()(conn, reqInfo)
		if len(reqInfo.IP()) == 0 {
			answer := profile.DNSHandle()(conn, reqInfo.Domain())
			if answer != nil {
				reqInfo.SetIP(answer.CurrentIP)
				reqInfo.SetCountryCode(answer.CurrentCountry)
			}
		}
		logrus.Infof("Match Rule [%s, %s, %s]", rule.Typ, rule.Value, rule.Proxy)
		conn.WithValue(constant.KeyRule, rule)
		next = profile.Filter()(next)
		next(conn)
	}
}

func outboundHandle() typ.HandleFunc {
	return func(lc connpkg.ICtxConn) {
		reqInfo := lc.Value(constant.KeyRequestInfo).(global.RequestInfo)
		rule := lc.Value(constant.KeyRule).(*rulepkg.Rule)
		profile := lc.Value(constant.KeyProfile).(*global.Profile)

		var s server.IServer
		g := profile.Group()[rule.Proxy]
		if g == nil {
			s = profile.Server()[rule.Proxy]
		} else {
			s = g.Server()
		}
		logrus.WithField("network", reqInfo.Network()).
			WithField("domain", reqInfo.Domain()).
			WithField("addr", fmt.Sprintf("%s:%d", reqInfo.IP(), reqInfo.Port())).
			WithField("country-code", reqInfo.CountryCode()).
			WithField("rule", rule.String()).
			Infof("URI: %s", reqInfo.URI())
		sc, err := s.Dial(lc, reqInfo.Network(), reqInfo, func(ctx context.Context, network string, addr, port string) (conn connpkg.ICtxConn, e error) {
			conn, err := connpkg.DefaultDial(ctx, network, addr, port)
			if err != nil {
				return nil, err
			}
			conn = profile.AfterStream()(conn)
			return conn, nil
		})
		if err != nil {
			logrus.WithField("proxy", rule.Proxy).WithError(err).Errorf("remote to server failed")
			return
		}
		lc = profile.BeforeStream()(lc)
		transefer(lc, sc)
	}
}

func transefer(from, to connpkg.ICtxConn) {
	go io.Copy(from, to)
	io.Copy(to, from)
}
