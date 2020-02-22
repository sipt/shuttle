package dump

import (
	"context"
	"io"

	"github.com/sipt/shuttle/events/record"

	"github.com/sipt/shuttle/conn"
	"github.com/sipt/shuttle/conn/stream"
	"github.com/sipt/shuttle/constant"
	"github.com/sipt/shuttle/constant/typ"
	"github.com/sipt/shuttle/util"
	"github.com/sirupsen/logrus"
)

func init() {
	stream.RegisterStream("data-dump", newHTTPDump)
}

func checkAllowDump(c context.Context, protocol ...string) bool {
	if !allowDump {
		return false
	}
	var p string
	if len(protocol) > 0 && len(protocol[0]) > 0 {
		p = protocol[0]
	} else {
		var ok bool
		p, ok = c.Value(constant.KeyProtocol).(string)
		if !ok {
			return false
		}
	}
	if p == constant.ProtocolHTTP || p == constant.ProtocolSOCKS_HTTP {
		return true
	} else if p == constant.ProtocolHTTPS || p == constant.ProtocolSOCKS_HTTPS {
		if reqInfo, ok := c.Value(constant.KeyRequestInfo).(typ.RequestInfo); ok {
			return mitmIsEnabled(reqInfo.Domain())
		}
	}
	return false
}

func newHTTPDump(ctx context.Context, runtime typ.Runtime, _ map[string]string) (typ.DecorateFunc, error) {
	err := applyRuntime(ctx, runtime)
	if err != nil {
		return nil, err
	}
	go AutoSave(ctx)
	record.RegisterClearCallback(func() error {
		return ClearFiles()
	})
	return func(c conn.ICtxConn) conn.ICtxConn {
		p, ok := c.Value(constant.KeyProtocol).(string)
		if !ok {
			return c
		}
		if !checkAllowDump(c, p) {
			return c
		}
		if p == constant.ProtocolHTTPS || p == constant.ProtocolSOCKS_HTTPS {
			lc, err := Mitm(c)
			if err != nil {
				logrus.WithError(err).Error("call mitm failed")
				_ = c.Close()
			}
			c = lc
		}
		rc := &httpDumpConn{
			ICtxConn: c,
		}
		var wt io.WriterTo
		var rf io.ReaderFrom
		if _, ok := c.(io.WriterTo); ok {
			wt = &recordTrafficConnWithWriteTo{
				ICtxConn: rc,
				WriterTo: &writeTo{
					ICtxConn: c,
				},
			}
		}
		if _, ok := c.(io.ReaderFrom); ok {
			rf = &recordTrafficConnWithReadFrom{
				ICtxConn: rc,
				ReaderFrom: &readFrom{
					ICtxConn: c,
				},
			}
		}
		switch {
		case wt == nil && rf == nil:
			return rc
		case wt != nil && rf == nil:
			return wt.(conn.ICtxConn)
		case wt == nil && rf != nil:
			return rf.(conn.ICtxConn)
		default:
			return &recordTrafficConnWithWriteToAndReadFrom{
				ICtxConn:   rc,
				WriterTo:   wt,
				ReaderFrom: rf,
			}
		}
	}, nil
}

func applyRuntime(ctx context.Context, runtime typ.Runtime) error {
	allowDump, _ = runtime.Get("dump").(bool)
	dumpPath, _ := runtime.Get("dump_path").(string)
	err := InitDumpStorage(dumpPath)
	if err != nil {
		return err
	}
	mitmEnabled, _ := runtime.Get("mitm").(bool)
	domains, _ := runtime.Get("domains").([]interface{})
	keyEncode, _ := runtime.Get("key").(string)
	caEncode, _ := runtime.Get("ca").(string)
	err = InitMITM(keyEncode, caEncode, mitmEnabled, util.InterfaceSliceToStringSlice(domains))
	if err != nil {
		return err
	}
	if len(keyEncode) == 0 || len(caEncode) == 0 {
		keyEncode, caEncode, err = GenerateCA()
		if err != nil {
			logrus.WithError(err).Error("generate ca failed")
		} else {
			if err = runtime.Set("key", keyEncode); err != nil {
				logrus.WithError(err).Error("set generate key failed")
			}
			if err = runtime.Set("ca", caEncode); err != nil {
				logrus.WithError(err).Error("set generate ca failed")
			}
		}
	}
	return nil
}

type httpDumpConn struct {
	conn.ICtxConn
}

func (h *httpDumpConn) Read(b []byte) (n int, err error) {
	n, err = h.ICtxConn.Read(b)
	if n > 0 {
		if id, ok := recordID(h); ok {
			err := SaveRequest(id, b[:n])
			if err != nil {
				logrus.WithField("record_id", id).WithError(err).Error("[data_dump] save request failed")
			}
		}
	}
	return
}

func (h *httpDumpConn) Write(b []byte) (n int, err error) {
	n, err = h.ICtxConn.Write(b)
	if n > 0 {
		if id, ok := recordID(h); ok {
			err := SaveResponse(id, b[:n])
			if err != nil {
				logrus.WithField("record_id", id).WithError(err).Error("[data_dump] save response failed")
			}
		}
	}
	return
}

func (h *httpDumpConn) Close() (err error) {
	err = h.ICtxConn.Close()
	if id, ok := recordID(h); ok {
		err := CloseFiles(id)
		if err != nil {
			logrus.WithField("record_id", id).WithError(err).Error("[data_dump] close files failed")
		}
	}
	return err
}

type recordTrafficConnWithWriteTo struct {
	conn.ICtxConn
	io.WriterTo
}

type recordTrafficConnWithReadFrom struct {
	conn.ICtxConn
	io.ReaderFrom
}
type recordTrafficConnWithWriteToAndReadFrom struct {
	conn.ICtxConn
	io.WriterTo
	io.ReaderFrom
}

type writeTo struct {
	conn.ICtxConn
}

func (r *writeTo) WriteTo(w io.Writer) (n int64, err error) {
	wr := &writer{Writer: w, ICtxConn: r.ICtxConn}
	n, err = r.ICtxConn.(io.WriterTo).WriteTo(wr)
	return n, err
}

type readFrom struct {
	conn.ICtxConn
}

func (r *readFrom) ReadFrom(re io.Reader) (n int64, err error) {
	rr := &reader{Reader: re, ICtxConn: r.ICtxConn}
	n, err = r.ICtxConn.(io.ReaderFrom).ReadFrom(rr)
	return n, err
}

type writer struct {
	io.Writer
	conn.ICtxConn
}

func (w *writer) Write(b []byte) (n int, err error) {
	n, err = w.Writer.Write(b)
	if n > 0 {
		if id, ok := recordID(w); ok {
			err := SaveRequest(id, b[:n])
			if err != nil {
				logrus.WithField("record_id", id).WithError(err).Error("[data_dump] save request failed")
			}
		}
	}
	return n, err
}

type reader struct {
	io.Reader
	conn.ICtxConn
}

func (r *reader) Read(b []byte) (n int, err error) {
	n, err = r.Reader.Read(b)
	if n > 0 {
		if id, ok := recordID(r); ok {
			err := SaveResponse(id, b[:n])
			if err != nil {
				logrus.WithField("record_id", id).WithError(err).Error("[data_dump] save response failed")
			}
		}
	}
	return n, err
}

func recordID(ctx context.Context) (int64, bool) {
	req, ok := ctx.Value(constant.KeyRequestInfo).(typ.RequestInfo)
	if !ok {
		return 0, false
	}
	return req.ID(), true
}
