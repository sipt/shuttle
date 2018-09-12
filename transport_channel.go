package shuttle

import (
	"github.com/sipt/shuttle/pool"
	"io"
	"net/http"
	"time"
	"bufio"
	"strings"
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
		// @fix 空数据返回引发断连
		//if n == 0 {
		//errChan <- nil
		//return
		//}
		if err != nil {
			if err != io.EOF && !strings.Contains(err.Error(), "use of closed network connection") {
				Logger.Errorf("ConnectID [%d] DirectChannel Transport: %v", from.GetID(), err)
			}
			errChan <- err
			return
		}
		n, err = to.Write(buf[:n])
		if err != nil {
			if err != io.EOF && !strings.Contains(err.Error(), "use of closed network connection") {
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
		isHttps:   first == nil,
	}
	h.Transport(lc, sc, first)
}

type HttpChannel struct {
	id        int64
	req       *http.Request
	urlStr    string
	allowDump bool
	template  *Record
	isHttps   bool
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
		boxChan <- &Box{h.id, RecordStatus, RecordStatusCompleted}
	}
}

func (h *HttpChannel) sendToClient(from, to IConn) {
	for {
		buf := bufio.NewReader(from)
		resp, err := http.ReadResponse(buf, h.req)
		if err != nil {
			if err != io.EOF && !strings.Contains(err.Error(), "use of closed network connection") {
				Logger.Errorf("ConnectID [%d] HttpChannel Transport s->[b]: %v", from.GetID(), err)
			}
			return
		}
		Logger.Debugf("ConnectID [%d] HttpChannel Transport return s->[b]", to.GetID())
		ResponseModify(h.req, resp, h.isHttps)
		err = h.writeResponse(resp, to)
		if err != nil {
			return
		}
	}
}

func (h *HttpChannel) sendToServer(from, to IConn, first *http.Request) {
	var err error
	var b *bufio.Reader
	var resp *http.Response
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
			//request update
			resp = RequestModify(h.req, h.isHttps)
		}
		h.id = util.NextID()
		Logger.Debugf("[connID:%d] [reqID:%d] HttpChannel Transport c->[req]: %s", from.GetID(), h.id, h.req.URL.String())
		record := *h.template
		record.ID = h.id
		record.URL = h.req.URL.String()
		if h.req.URL.Host == "" {
			if h.isHttps {
				record.URL = "https://" + h.req.Host + record.URL
			} else {
				record.URL = "http://" + h.req.Host + record.URL
			}
		}
		record.Status = RecordStatusActive
		record.Created = time.Now()
		record.Dumped = h.allowDump
		boxChan <- &Box{Op: RecordAppend, Value: &record, ID: record.ID}
		Logger.Debugf("[reqID:%d] HttpChannel Transport send record to boxChan", h.id)
		to.SetRecordID(record.ID)
		var dumpWriter io.Writer
		if h.allowDump {
			dump.InitDump(h.id)
			dumpWriter = ToWriter(func(b []byte) (int, error) {
				return dump.WriteRequest(h.id, b)
			})
		}
		// 分流器
		var shunt *Shunt
		if resp != nil {
			// response mock, set nil to server conn
			shunt = NewShunt(nil, dumpWriter)
		} else {
			shunt = NewShunt(to, dumpWriter)
		}

		err = h.req.Write(shunt)
		if err != nil {
			if err != io.EOF {
				Logger.Errorf("ConnectID [%d] HttpChannel Transport [req]->c: %v", to.GetID(), err)
				return
			}
		}
		//response mock ?
		if resp != nil {
			// write response to client
			err = h.writeResponse(resp, from)
			if err != nil {
				return
			}
		}
	}
	return
}

// write response in connection
func (h *HttpChannel) writeResponse(resp *http.Response, to IConn) (err error) {
	var dumpWriter io.Writer
	if h.allowDump {
		dumpWriter = ToWriter(func(b []byte) (int, error) {
			return dump.WriteResponse(h.id, b)
		})
	}
	// 分流器
	shunt := NewShunt(to, dumpWriter)
	err = resp.Write(shunt)
	if err != nil {
		if err != io.EOF {
			Logger.Errorf("ConnectID [%d] HttpChannel Transport [b]->c: %v", to.GetID(), err)
			return
		}
	}
	if h.allowDump {
		go func() {
			dump.Complete(h.id)
		}()
	}
	boxChan <- &Box{h.id, RecordStatus, RecordStatusCompleted}
	h.id = 0
	return
}
