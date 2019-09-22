package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/sipt/shuttle/conf"
	"github.com/sipt/shuttle/conf/logger"
	"github.com/sipt/shuttle/constant"
	"github.com/sipt/shuttle/constant/typ"
	"github.com/sipt/shuttle/controller"
	"github.com/sipt/shuttle/global"
	"github.com/sipt/shuttle/global/namespace"
	"github.com/sipt/shuttle/inbound"
	"github.com/sipt/shuttle/pkg/debug"
	"github.com/sipt/shuttle/server"
	"github.com/sirupsen/logrus"

	connpkg "github.com/sipt/shuttle/conn"
	closepkg "github.com/sipt/shuttle/pkg/close"
	rulepkg "github.com/sipt/shuttle/rule"
)

var Path = flag.String("c", os.Getenv("CONFIG_PATH"), "config file Path")
var Encoding = flag.String("e", os.Getenv("ENCODING"), "config file Encoding")
var LogPath = flag.String("logpath", os.Getenv("LOGGER_PATH"), "logger file")

func Start() error {
	logrus.SetLevel(logrus.DebugLevel)
	err := logger.ConfigOutput(*LogPath)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	params := map[string]string{"path": *Path}
	config, err := conf.LoadConfig(ctx, "file", "toml", params, func() {
		fmt.Println("config file change")
	})
	if err != nil {
		logrus.WithError(err).Error("load config failed")
		return err
	}
	l, err := logrus.ParseLevel(config.General.LoggerLevel)
	if err != nil {
		l = logrus.DebugLevel
	}
	logger.ConfigLogger(l)
	if err != nil {
		logrus.WithError(err).Error("load config failed")
		return err
	}
	err = conf.ApplyConfig(ctx, config)
	if err != nil {
		logrus.WithError(err).Error("apply config failed")
		return err
	}
	closer, err := controller.ApplyConfig(config)
	if err != nil {
		logrus.WithError(err).Error("start controller failed")
		return err
	}
	err = inbound.ApplyConfig(ctx, config, handle())
	if err != nil {
		logrus.WithError(err).Error("start inbound failed")
		return err
	}

	logrus.Info("server starting...")
	closepkg.AppendCloser(func() error {
		cancel()
		return nil
	})
	closepkg.AppendCloser(func() error {
		closer()
		return nil
	})
	return nil
}

func handle() typ.HandleFunc {
	handle := outboundHandle()
	handle = ruleHandle(handle)
	handle = namespaceHandle(handle)
	handle = recoverHandle(handle)
	return handle
}

func namespaceHandle(next typ.HandleFunc) typ.HandleFunc {
	return func(conn connpkg.ICtxConn) {
		if p, ok := conn.Value(constant.KeyProfile).(*global.Profile); !ok || p == nil {
			np := namespace.NamespaceWithContext(conn)
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
		if reqInfo.Network() == "udp" {
			rule = profile.UDPRuleHandle()(conn, reqInfo)
		}
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
			_ = lc.Close()
			return
		}
		lc = profile.BeforeStream()(lc)
		connpkg.Transfer(lc, sc)
	}
}

func recoverHandle(next typ.HandleFunc) typ.HandleFunc {
	return func(lc connpkg.ICtxConn) {
		defer func() {
			if e := recover(); e != nil {
				logrus.WithField("error", e).WithField("stack", string(debug.Stack(3))).
					Error("stacktrace from panic")
			}
		}()
		next(lc)
	}
}
