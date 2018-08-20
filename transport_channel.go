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
	errChan := make(chan error, 2)
	go d.send(sc, lc, errChan)
	go d.send(lc, sc, errChan)
	<-errChan
	lc.Close()
	sc.Close()
}

func (d *DirectChannel) send(from, to IConn, errChan chan error) {
	var (
		buf []byte
		n   int
		err error
	)
	for {
		buf = pool.GetBuf()
		n, err = from.Read(buf)
		if n == 0 {
			errChan <- nil
			return
		}
		if err != nil {
			if err != io.EOF {
				Logger.Errorf("ConnectID [%d] DirectChannel Transport: %v", from.GetID(), err)
			}
			errChan <- err
			return
		}
		n, err = to.Write(buf[:n])
		if err != nil {
			if err != io.EOF {
				Logger.Error("ConnectID [%d] DirectChannel Transport: %v", to.GetID(), err)
			}
			errChan <- err
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
	if h.id != 0 {
		go storage.Put(h.id, RecordStatus, RecordStatusCompleted)
	}
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
	var isHttps = first == nil
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
		if h.req.URL.Host == "" {
			if isHttps {
				record.URL = "https://" + h.req.Host + record.URL
			} else {
				record.URL = "http://" + h.req.Host + record.URL
			}
		}
		record.Status = RecordStatusActive
		record.Created = time.Now()
		record.Dumped = h.allowDump
		recordChan <- &record
		err = h.req.Write(to)
		if h.allowDump {
			go func(id, oldID int64, req *http.Request) {
				if oldID != 0 && oldID != id {
					dump.Complete(oldID)
				}
				dump.InitDump(id)
				writer := bytes.NewBuffer(pool.GetBuf()[:0])
				req.Write(writer)
				dump.WriteRequest(id, writer.Bytes())
			}(h.id, h.oldID, h.req)
		}
		if h.oldID != 0 && h.oldID != h.id {
			go storage.Put(h.oldID, RecordStatus, RecordStatusCompleted)
			h.oldID = 0
		}
	}
	return
}
