package shuttle

import (
	"github.com/sipt/shuttle/pool"
	"io"
	"net/http"
	"time"
	"bytes"
	"bufio"
	"github.com/sipt/shuttle/util"
)

type DirectChannel struct{}

func (d *DirectChannel) Transport(lc, sc IConn) {
	go d.send(sc, lc)
	d.send(lc, sc)
	lc.Close()
	sc.Close()
}

func (d *DirectChannel) send(from, to IConn) {
	var (
		buf []byte
		n   int
		err error
	)
	for {
		buf = pool.GetBuf()
		n, err = from.Read(buf)
		if n == 0 {
			return
		}
		if err != nil {
			if err != io.EOF {
				Logger.Errorf("ConnectID [%d] DirectChannel Transport: %v", from.GetID(), err)
			}
			return
		}
		n, err = to.Write(buf[:n])
		if err != nil {
			if err != io.EOF {
				Logger.Error("ConnectID [%d] DirectChannel Transport: %v", to.GetID(), err)
			}
			return
		}
	}
}

func HttpTransport(lc, sc IConn, template *Record, allowDump bool, first *http.Request) {
	h := &HttpChannel{
		template:  template,
		allowDump: allowDump,
	}
	h.Transport(lc, sc, first)
}

type HttpChannel struct {
	id, oldID int64
	req       *http.Request
	allowDump bool
	template  *Record
}

func (h *HttpChannel) Transport(lc, sc IConn, first *http.Request) {
	go h.sendToClient(sc, lc)
	h.sendToServer(lc, sc, first)

	if h.allowDump {
		go dump.Complete(h.id)
	}
	lc.Close()
	sc.Close()
}

func (h *HttpChannel) sendToClient(from, to IConn) {
	var (
		buf []byte
		n   int
		err error
	)
	for {
		buf = pool.GetBuf()
		n, err = from.Read(buf)
		if n == 0 {
			return
		}
		if err != nil {
			if err != io.EOF {
				Logger.Errorf("ConnectID [%d] HttpChannel Transport s->[b]: %v", from.GetID(), err)
			}
			return
		}
		n, err = to.Write(buf[:n])
		if err != nil {
			if err != io.EOF {
				Logger.Error("ConnectID [%d] HttpChannel Transport [b]->c: %v", to.GetID(), err)
			}
			return
		}
		if h.allowDump {
			go dump.WriteResponse(h.id, buf[:n])
		}
	}
}

func (h *HttpChannel) sendToServer(from, to IConn, first *http.Request) {
	var err error
	var b *bufio.Reader
	for {
		if first != nil {
			h.req = first
			first = nil
		} else {
			if b == nil {
				b = bufio.NewReader(from)
			}
			h.req, err = http.ReadRequest(b)
			if err != nil {
				if err != io.EOF {
					Logger.Errorf("ConnectID [%d] HttpChannel Transport c->[r]: %v", from.GetID(), err)
				}
				return
			}
		}
		if h.id == 0 {
			h.id = from.GetID()
		} else {
			h.oldID, h.id = h.id, util.NextID()
		}
		Logger.Debugf("[connID:%d] [reqID:%d] HttpChannel Transport c->[r]: %s", from.GetID(), h.id, h.req.URL.String())
		record := *h.template
		record.ID = h.id
		record.URL = h.req.URL.String()
		record.Status = RecordStatusActive
		record.Created = time.Now()
		recordChan <- &record
		err = h.req.Write(to)
		if h.allowDump {
			go func(id int64, req *http.Request) {
				if h.oldID != 0 && h.oldID != h.id {
					dump.Complete(h.oldID)
					h.oldID = 0
				}
				dump.InitDump(h.id)
				writer := bytes.NewBuffer(pool.GetBuf()[:0])
				req.Write(writer)
				dump.WriteRequest(h.id, writer.Bytes())
			}(h.id, h.req)
		}
	}
	return
}
