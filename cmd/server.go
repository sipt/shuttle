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
	params := map[string]string{"path": "/Users/sipt/workspace/go/shuttle/cmd/config1.toml"}
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
	logrus.Debug("start transefer")
	go io.Copy(from, to)
	io.Copy(to, from)
}

func Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return copyBuffer(dst, src, nil)
}

// copyBuffer is the actual implementation of Copy and CopyBuffer.
// if buf is nil, one is allocated.
func copyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	if buf == nil {
		size := 32 * 1024
		if l, ok := src.(*io.LimitedReader); ok && int64(size) > l.N {
			if l.N < 1 {
				size = 1
			} else {
				size = int(l.N)
			}
		}
		buf = make([]byte, size)
	}
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			fmt.Println(string(buf[:nr]))
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
