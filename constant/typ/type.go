package typ

import (
	"net"
	"sync"

	"github.com/sipt/shuttle/conn"
)

type HandleFunc func(conn.ICtxConn)

type K struct {
	V string
}

type RequestInfo interface {
	ID() int64
	Network() string
	Domain() string
	URI() string
	IP() net.IP
	CountryCode() string
	Port() int
	SetIP(net.IP)
	SetPort(int)
	SetCountryCode(string)
}

type Runtime interface {
	sync.Locker
	Get(string) interface{}
	Set(string, interface{}) error
}
