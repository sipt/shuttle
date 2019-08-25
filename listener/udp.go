package listener

import (
	"net"

	"github.com/sipt/shuttle/conn"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func init() {
	Register("udp", newUDPListener)
}

func newUDPListener(addr string) (func(HandleFunc) error, error) {
	return func(handle HandleFunc) error {
		l, err := net.Listen("udp", addr)
		logrus.WithField("addr", "udp://"+addr).Info("udp listen starting")
		if err != nil {
			return errors.Errorf("listen udp://[%s] failed", addr)
		}
		for {
			c, err := l.Accept()
			if err != nil {
				// TODO call error center
				logrus.WithField("addr", addr).WithError(err).Errorf("[udp] listener accept failed")
			}
			handle(conn.WrapConn(c))
		}
	}, nil

}
