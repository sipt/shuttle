package listener

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/conn"
	"github.com/sipt/shuttle/constant/typ"
	"github.com/sirupsen/logrus"
)

func init() {
	Register("udp", newUDPListener)
}

func newUDPListener(addr string) (func(context.Context, typ.HandleFunc), error) {
	l, err := net.Listen("udp", addr)
	logrus.WithField("addr", "udp://"+l.Addr().String()).Info("udp listen starting")
	if err != nil {
		return nil, errors.Errorf("listen udp://[%s] failed", addr)
	}
	return func(ctx context.Context, handle typ.HandleFunc) {
		go func() {
			for {
				c, err := l.Accept()
				if err != nil {
					// TODO call error center
					logrus.WithField("addr", addr).WithError(err).Errorf("[udp] listener accept failed")
					return
				}
				handle(conn.NewConn(c, ctx))
			}
		}()
		<-ctx.Done()
		_ = l.Close()
	}, nil
}
