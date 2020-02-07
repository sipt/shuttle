package dump

import (
	"context"
	"io"

	"github.com/sipt/shuttle/conn"
	"github.com/sipt/shuttle/conn/stream"
	"github.com/sipt/shuttle/constant"
	"github.com/sipt/shuttle/constant/typ"
	"github.com/sirupsen/logrus"
)

func init() {
	stream.RegisterStream("data-dump", newHTTPDump)
}

func checkAllowDump(c context.Context, protocol ...string) bool {
	if len(protocol) > 0 && len(protocol[0]) > 0 {
		p := protocol[0]
		return allowDump && (p == constant.ProtocolHTTP || p == constant.ProtocolHTTPS && mitm)
	} else {
		p, ok := c.Value(constant.KeyProtocol).(string)
		return ok && allowDump && (p == constant.ProtocolHTTP || p == constant.ProtocolHTTPS && mitm)
	}
}

func newHTTPDump(ctx context.Context, params map[string]string) (stream.DecorateFunc, error) {
	allowDump = params["enabled"] == "true"
	err := InitDumpStorage(params["dump_path"])
	if err != nil {
		return nil, err
	}
	go AutoSave(ctx)
	return func(c conn.ICtxConn) conn.ICtxConn {
		if !checkAllowDump(c) {
			return c
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
