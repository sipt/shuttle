package stream

import (
	"context"
	"io"

	"github.com/sipt/shuttle/conn"
	"github.com/sipt/shuttle/conn/stream"
	"github.com/sipt/shuttle/constant"
	"github.com/sipt/shuttle/constant/typ"
	"github.com/sipt/shuttle/events"
	"github.com/sipt/shuttle/events/record"
)

func init() {
	stream.RegisterStream("record-traffic", newRecordTrafficMetrics)
}

func newRecordTrafficMetrics(ctx context.Context, _ map[string]string) (stream.DecorateFunc, error) {
	return func(c conn.ICtxConn) conn.ICtxConn {
		rc := &recordTrafficConn{
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

type recordTrafficConn struct {
	conn.ICtxConn
}

func (t *recordTrafficConn) Read(b []byte) (n int, err error) {
	n, err = t.ICtxConn.Read(b)
	if n > 0 {
		if reqInfo, ok := t.Value(constant.KeyRequestInfo).(typ.RequestInfo); ok {
			events.Bus <- &events.Event{
				Typ: record.UpdateRecordUpEvent,
				Value: &record.RecordEntity{
					ID: reqInfo.ID(),
					Up: int64(n),
				},
			}
		}
	}
	return
}

func (t *recordTrafficConn) Write(b []byte) (n int, err error) {
	n, err = t.ICtxConn.Write(b)
	if n > 0 {
		if reqInfo, ok := t.Value(constant.KeyRequestInfo).(typ.RequestInfo); ok {
			events.Bus <- &events.Event{
				Typ: record.UpdateRecordDownEvent,
				Value: &record.RecordEntity{
					ID:   reqInfo.ID(),
					Down: int64(n),
				},
			}
		}
	}
	return
}

func (t *recordTrafficConn) Close() (err error) {
	err = t.ICtxConn.Close()
	if reqInfo, ok := t.Value(constant.KeyRequestInfo).(typ.RequestInfo); ok {
		events.Bus <- &events.Event{
			Typ: record.UpdateRecordStatusEvent,
			Value: &record.RecordEntity{
				ID:     reqInfo.ID(),
				Status: record.CompletedStatus,
			},
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
	if reqInfo, ok := w.Value(constant.KeyRequestInfo).(typ.RequestInfo); ok {
		events.Bus <- &events.Event{
			Typ: record.UpdateRecordUpEvent,
			Value: &record.RecordEntity{
				ID: reqInfo.ID(),
				Up: int64(n),
			},
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
	if reqInfo, ok := r.Value(constant.KeyRequestInfo).(typ.RequestInfo); ok {
		events.Bus <- &events.Event{
			Typ: record.UpdateRecordDownEvent,
			Value: &record.RecordEntity{
				ID:   reqInfo.ID(),
				Down: int64(n),
			},
		}
	}
	return n, err
}
