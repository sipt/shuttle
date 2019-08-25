package listener

import (
	"net"

	"github.com/sipt/shuttle/conn"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func init() {
	Register("tcp", newTCPListener)
}

func newTCPListener(addr string) (func(HandleFunc) error, error) {
	return func(handle HandleFunc) error {
		l, err := net.Listen("tcp", addr)
		logrus.WithField("addr", "tcp://"+addr).Info("tcp listen starting")
		if err != nil {
			return errors.Errorf("listen tcp://[%s] failed", addr)
		}
		for {
			c, err := l.Accept()
			if err != nil {
				// TODO call error center
				logrus.WithField("addr", addr).WithError(err).Errorf("[tcp] listener accept failed")
			}
			handle(conn.WrapConn(c))
		}
	}, nil

}
