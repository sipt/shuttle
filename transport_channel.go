package shuttle

import (
	"github.com/sipt/shuttle/pool"
	"io"
	"net/http"
	"time"
	"bytes"
	"bufio"
	"github.com/sipt/shuttle/util"
	"strings"
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
		if disposition := resp.Header.Get("Content-Disposition"); len(disposition) > 0 && strings.HasPrefix(disposition, "attachment") {
			//文件下载，不Dump
			err = resp.Write(to)
			if err != nil {
				if err != io.EOF {
					Logger.Errorf("ConnectID [%d] HttpChannel Transport [b]->c: %v", to.GetID(), err)
					return
				}
			}
			if h.allowDump {
				go func() {
					buffer := &bytes.Buffer{}
					err := resp.Write(buffer)
					if err == nil {
						dump.WriteResponse(h.id, buffer.Bytes())
						dump.Complete(h.id)
					}
				}()
			}
		} else {
			buffer := &bytes.Buffer{}
			err = resp.Write(buffer)
			if err != nil {
				if err != io.EOF {
					Logger.Errorf("ConnectID [%d] HttpChannel Transport [b]->c: %v", to.GetID(), err)
					return
				}
			}
			_, err = to.Write(buffer.Bytes())
			if err != nil {
				if err != io.EOF {
					Logger.Errorf("ConnectID [%d] HttpChannel Transport [b]->c: %v", to.GetID(), err)
					return
				}
			}
			if h.allowDump {
				go func() {
					dump.WriteResponse(h.id, buffer.Bytes())
					dump.Complete(h.id)
				}()
			}
		}
		boxChan <- &Box{h.id, RecordStatus, RecordStatusCompleted}
	}
}

func (h *HttpChannel) sendToServer(from, to IConn, first *http.Request) {
	var err error
	var b *bufio.Reader
	var respBuf []byte
	var buffer *bytes.Buffer
	for {
		respBuf = nil
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
			resp := RequestModify(h.req, h.isHttps)
			if resp != nil { // response mock ?
				buffer = &bytes.Buffer{}
				err = resp.Write(buffer)
				if err != nil {
					if err != io.EOF {
						Logger.Errorf("ConnectID [%d] HttpChannel Transport [req]->[b]: %v", to.GetID(), err)
						return
					}
				}
				respBuf = buffer.Bytes()
			}
		}
		if h.id == 0 {
			h.id = from.GetID()
		} else {
			h.id = util.NextID()
		}
		Logger.Debugf("[connID:%d] [reqID:%d] HttpChannel Transport c->[r]: %s", from.GetID(), h.id, h.req.URL.String())
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
		boxChan <- &Box{Op: RecordAppend, Value: &record}
		to.SetRecordID(record.ID)

		if content := h.req.Header.Get("Content-Type"); len(content) > 0 && strings.HasPrefix(content, "multipart/form-data") {
			//上传文件，不Dump
			if len(respBuf) == 0 {
				err = h.req.Write(to)
				if err != nil {
					if err != io.EOF {
						Logger.Errorf("ConnectID [%d] HttpChannel Transport [r]->s: %v", to.GetID(), err)
						return
					}
				}
			} else {
				_, err = from.Write(respBuf)
				if err != nil {
					if err != io.EOF {
						Logger.Errorf("ConnectID [%d] HttpChannel Transport [b]->c: %v", to.GetID(), err)
						return
					}
				}
			}
			buffer = &bytes.Buffer{}
			h.req.Write(buffer)
		} else {
			buffer = &bytes.Buffer{}
			err = h.req.Write(buffer)
			if err != nil {
				if err != io.EOF {
					Logger.Errorf("ConnectID [%d] HttpChannel Transport [req]->[b]: %v", to.GetID(), err)
					return
				}
			}
			if len(respBuf) == 0 {
				_, err = to.Write(buffer.Bytes())
				if err != nil {
					if err != io.EOF {
						Logger.Errorf("ConnectID [%d] HttpChannel Transport [b]->s: %v", to.GetID(), err)
						return
					}
				}
			} else {
				_, err = from.Write(respBuf)
				if err != nil {
					if err != io.EOF {
						Logger.Errorf("ConnectID [%d] HttpChannel Transport [b]->c: %v", to.GetID(), err)
						return
					}
				}
			}
		}
		if h.allowDump {
			go func(id int64, reqBuf, respBuf []byte) {
				dump.InitDump(id)
				writer := bytes.NewBuffer(pool.GetBuf()[:0])
				writer.Write(reqBuf)
				dump.WriteRequest(id, writer.Bytes())
				if len(respBuf) > 0 {
					dump.WriteResponse(id, respBuf)
					dump.Complete(id)
				}
			}(h.id, buffer.Bytes(), respBuf)
		}
	}
	return
}
