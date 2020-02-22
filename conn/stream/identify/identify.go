package identify

import (
	"context"
	"net/http"
	"strings"

	"github.com/sipt/shuttle/conn"
	"github.com/sipt/shuttle/conn/stream"
	"github.com/sipt/shuttle/conn/stream/dump"
	"github.com/sipt/shuttle/constant"
	"github.com/sipt/shuttle/constant/typ"
	"github.com/sipt/shuttle/pkg/pool"
	"github.com/sirupsen/logrus"
)

func init() {
	stream.RegisterStream("protocol-identify", newProtocolIdentify)
}

func newProtocolIdentify(ctx context.Context, _ typ.Runtime, _ map[string]string) (typ.DecorateFunc, error) {
	return func(c conn.ICtxConn) conn.ICtxConn {
		if c.Value(constant.KeyProtocol).(string) != constant.ProtocolSOCKS {
			return c
		}
		b := pool.GetBuf()
		n, err := c.Read(b)
		if err != nil {
			return &bufferConn{err: err}
		}
		reqInfo := c.Value(constant.KeyRequestInfo).(typ.RequestInfo)
		if isHTTP(b[:n]) > 0 {
			if isWebsocket(b[:n]) {
				c.WithValue(constant.KeyProtocol, constant.ProtocolSOCKS_Websocket)
			} else {
				c.WithValue(constant.KeyProtocol, constant.ProtocolSOCKS_HTTP)
			}
		} else if IsTls(b[:n]) {
			if reqInfo.Port() == 443 {
				c.WithValue(constant.KeyProtocol, constant.ProtocolSOCKS_HTTPS)
			} else {
				c.WithValue(constant.KeyProtocol, constant.ProtocolSOCKS_TLS)
			}
		}
		{
			err := dump.InitFiles(c, reqInfo.ID())
			if err != nil {
				logrus.WithField("record_id", reqInfo.ID()).WithError(err).Error("[data_dump] init files failed")
			}
		}
		return &bufferConn{
			buffer:   b[:n],
			ICtxConn: c,
		}
	}, nil
}

func isHTTP(b []byte) int {
	switch b[0] {
	case 'G':
		if string(b[:len(http.MethodGet)]) == http.MethodGet {
			return 2
		}
	case 'H':
		if string(b[:len(http.MethodHead)]) == http.MethodHead {
			return 1
		}
	case 'P':
		if string(b[:len(http.MethodPost)]) == http.MethodPost {
			return 1
		} else if string(b[:len(http.MethodPut)]) == http.MethodPut {
			return 1
		} else if string(b[:len(http.MethodPatch)]) == http.MethodPatch {
			return 1
		}
	case 'D':
		if string(b[:len(http.MethodDelete)]) == http.MethodDelete {
			return 1
		}
	case 'C':
		if string(b[:len(http.MethodConnect)]) == http.MethodConnect {
			return 1
		}
	case 'O':
		if string(b[:len(http.MethodOptions)]) == http.MethodOptions {
			return 1
		}
	case 'T':
		if string(b[:len(http.MethodTrace)]) == http.MethodTrace {
			return 1
		}
	default:
		return 0
	}
	return 0
}

func isWebsocket(b []byte) bool {
	return strings.Index(string(b), "Upgrade: websocket") > -1
}

type bufferConn struct {
	buffer []byte
	offset int
	conn.ICtxConn
	err error
}

func (i *bufferConn) Read(b []byte) (n int, err error) {
	if i.err != nil {
		return 0, i.err
	}
	if len(i.buffer) > i.offset {
		n = copy(b, i.buffer[i.offset:])
		i.offset += n
		return
	} else if i.buffer != nil {
		pool.PutBuf(i.buffer)
		i.buffer = nil
	}
	return i.ICtxConn.Read(b)
}
