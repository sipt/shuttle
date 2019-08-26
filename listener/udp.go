package listener

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/conn"
	"github.com/sirupsen/logrus"
)

func init() {
	Register("udp", newUDPListener)
}

func newUDPListener(addr string) (func(context.Context, HandleFunc), error) {
	l, err := net.Listen("udp", addr)
	logrus.WithField("addr", "udp://"+addr).Info("udp listen starting")
	if err != nil {
		return nil, errors.Errorf("listen udp://[%s] failed", addr)
	}
	return func(ctx context.Context, handle HandleFunc) {
		go func() {
			for {
				c, err := l.Accept()
				if err != nil {
					// TODO call error center
					logrus.WithField("addr", addr).WithError(err).Errorf("[udp] listener accept failed")
					return
				}
				handle(conn.WrapConn(c))
			}
		}()
		<-ctx.Done()
		_ = l.Close()
	}, nil
}
