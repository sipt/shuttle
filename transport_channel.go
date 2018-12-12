package shuttle

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	connect "github.com/sipt/shuttle/conn"
	"github.com/sipt/shuttle/log"
	"github.com/sipt/shuttle/pool"
	"github.com/sipt/shuttle/proxy"
	rule2 "github.com/sipt/shuttle/rule"
	"github.com/sipt/shuttle/util"
)

type DirectChannel struct{}

func (d *DirectChannel) Transport(lc, sc connect.IConn) {
	errChan := make(chan error, 2)
	go func() {
		defer Recover(func() {
			select {
			case errChan <- nil:
			default:
			}
		})
		d.send(sc, lc, errChan)
		<-errChan
	}()
	go func() {
		defer Recover(func() {
			select {
			case errChan <- nil:
			default:
			}
		})
		d.send(lc, sc, errChan)

	}()
	<-errChan
	lc.Close()
	sc.Close()
}

func (d *DirectChannel) send(from, to connect.IConn, errChan chan error) {
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
				log.Logger.Errorf("[ID:%d] [DirectChannel] DirectChannel Transport: %v", from.GetID(), err)
			}
			errChan <- err
			return
		}
		n, err = to.Write(buf[:n])
		if err != nil {
			if err != io.EOF && !strings.Contains(err.Error(), "use of closed network connection") {
				log.Logger.Error("[ID:%d] [DirectChannel] DirectChannel Transport: %v", to.GetID(), err)
			}
			errChan <- err
			return
		}
	}
}

func HttpTransport(lc, sc connect.IConn, allowDump bool, first *http.Request) {
	h := &HttpChannel{
		allowDump: allowDump,
		isHttps:   first == nil,
	}
	h.Transport(lc, sc, first)
}

type HttpChannel struct {
	allowDump bool
	isHttps   bool
}

func (h *HttpChannel) Transport(lc, sc connect.IConn, first *http.Request) (err error) {
	var (
		oldHreq, hreq *http.Request
		lcBuf         = bufio.NewReader(lc)
		scBuf         *bufio.Reader
		resp          *http.Response
		rule          *rule2.Rule
		server        *proxy.Server
		passed        bool // inner request
		scid          int64
	)
	if sc != nil {
		scBuf = bufio.NewReader(sc)
		ctx := sc.Context()
		rule, _ = ctx.Value("rule").(*rule2.Rule)
		server, _ = ctx.Value("server").(*proxy.Server)
	}
	defer func() {
		lc.Close()
		if sc != nil {
			sc.Close()
		}
	}()
	for {
		// read from client
		if hreq == nil && first != nil {
			hreq = first
		} else {
			oldHreq = hreq
			hreq, err = http.ReadRequest(lcBuf)
			if err != nil {
				if err != io.EOF {
					log.Logger.Errorf("[ID:%d] [HttpChannel] HttpChannel Transport c->[hreq]: %v", sc.GetID(), err)
				}
				return
			}
		}

		//request update
		resp = RequestModify(hreq, h.isHttps)
		passed = IsPass(hreq.URL.Hostname(), hreq.URL.Hostname(), hreq.URL.Port())
		// Record
		record := &Record{
			ID:      util.NextID(),
			URL:     hreq.URL.String(),
			Status:  RecordStatusActive,
			Created: time.Now(),
			Dumped:  h.allowDump && !passed,
			Rule:    rule,
			Proxy:   server,
		}
		if h.isHttps {
			record.Protocol = HTTPS
		} else {
			record.Protocol = HTTP + "(" + hreq.Method + ")"
		}
		if hreq.URL.Host == "" {
			if h.isHttps {
				record.URL = "https://" + hreq.Host + record.URL
			} else {
				record.URL = "http://" + hreq.Host + record.URL
			}
		} else if hreq.URL.Scheme == "" {
			if h.isHttps {
				record.URL = "https:" + record.URL
			} else {
				record.URL = "http:" + record.URL
			}
		}
		log.Logger.Debugf("[ID:%d] [HttpChannel] [reqID:%d] HttpChannel Transport c->[hreq]: %s", lc.GetID(), record.ID, record.URL)

		// rule RuleFilter
		if resp == nil && (sc == nil || (oldHreq != nil && hreq.URL.Host != oldHreq.URL.Host)) {
			if sc != nil {
				sc.Close()
			}
			rule, server, sc, err = ConnectFilter(hreq, lc.GetID())
			record.Rule = rule
			record.Proxy = server
			if err != nil {
				if err == ErrorReject {
					record.Status = RecordStatusReject
				} else {
					record.Status = RecordStatusFailed
					record.Rule = rule2.FailedRule
					record.Proxy = proxy.FailedServer
				}
				record.Dumped = false
				if !passed {
					boxChan <- &Box{Op: RecordAppend, Value: record, ID: record.ID}
				}
				return
			}
			if scBuf == nil {
				scBuf = bufio.NewReader(sc)
			} else {
				scBuf.Reset(sc)
			}
			scid = sc.GetID()
		} else if resp != nil {
			record.Rule = rule2.MockRule
			record.Proxy = proxy.MockServer
		}
		if !passed {
			boxChan <- &Box{Op: RecordAppend, Value: record, ID: record.ID}
		}
		if sc != nil {
			log.Logger.Debugf("[ID:%d] [HttpChannel] [reqID:%d] HttpChannel Transport send record to boxChan", scid, record.ID)
			sc.SetRecordID(record.ID)
		}

		// dump
		var dumpWriter io.Writer
		if !passed && h.allowDump {
			dump.InitDump(record.ID)
			dumpWriter = ToWriter(func(b []byte) (int, error) {
				return dump.WriteRequest(record.ID, b)
			})
		}
		// 分流器
		var shunt *Shunt
		if resp != nil {
			// response mock, set nil to server conn
			shunt = NewShunt(nil, dumpWriter)
		} else {
			shunt = NewShunt(sc, dumpWriter)
		}

		err = hreq.Write(shunt)
		if err != nil {
			if err != io.EOF {
				log.Logger.Errorf("[ID:%d] [HttpChannel] HttpChannel Transport [hreq]->s: %v", scid, err)
				return
			}
		}
		//response mock ?
		if resp != nil {
			// write response to client
			err = h.writeResponse(resp, lc, record.ID, h.allowDump && !passed)
			if err != nil {
				return
			}
			continue
		}
		//==================
		//Read response
		//==================
		resp, err = http.ReadResponse(scBuf, nil)
		if err != nil {
			if err != io.EOF {
				log.Logger.Errorf("[ID:%d] [HttpChannel] HttpChannel Transport s->[b]: %v", scid, err)
			}
			return
		}
		log.Logger.Debugf("[ID:%d] [HttpChannel] HttpChannel Transport return s->[b]", scid)
		ResponseModify(hreq, resp, h.isHttps)
		err = h.writeResponse(resp, lc, record.ID, h.allowDump && !passed)
		if err != nil {
			return
		}
	}
	return
}

// write response in connection
func (h *HttpChannel) writeResponse(resp *http.Response, to connect.IConn, recordID int64, allowDump bool) (err error) {
	var dumpWriter io.Writer
	if allowDump {
		dumpWriter = ToWriter(func(b []byte) (int, error) {
			return dump.WriteResponse(recordID, b)
		})
	}
	// 分流器
	shunt := NewShunt(to, dumpWriter)
	err = resp.Write(shunt)
	if err != nil && err != io.EOF {
		log.Logger.Errorf("[ID:%d] [HttpChannel] HttpChannel Transport [b]->c: %v", to.GetID(), err)
	} else {
		log.Logger.Debugf("[ID:%d] [HttpChannel] HttpChannel Transport return [b]->c", to.GetID())
	}
	if allowDump {
		go func() {
			dump.Complete(recordID)
		}()
	}
	if err == nil || err == io.EOF {
		boxChan <- &Box{recordID, RecordStatus, RecordStatusCompleted}
	} else {
		boxChan <- &Box{recordID, RecordStatus, RecordStatusReject}
	}
	return
}

func HostName(req *http.Request) (host string) {
	if req.URL != nil {
		host = req.URL.Hostname()
	}
	if len(host) == 0 {
		host = req.Header.Get("Host")
	}
	return
}

func ConnectFilter(hreq *http.Request, connID int64) (rule *rule2.Rule, server *proxy.Server, conn connect.IConn, err error) {
	req := &HttpRequest{
		network:  connect.TCP,
		domain:   HostName(hreq),
		connID:   connID,
		port:     hreq.URL.Port(),
		protocol: hreq.URL.Scheme,
	}
	if len(req.protocol) == 0 {
		req.protocol = HTTPS
	}
	if len(net.ParseIP(req.domain)) > 0 {
		req.ip = req.domain
		req.domain = ""
	}
	rule, server, err = FilterByReq(req)
	if err != nil {
		log.Logger.Errorf("[HTTP] [ID:%d] ConnectToServer failed [%s] err: %s", connID, req.Host(), err)
		return
	}

	log.Logger.Debugf("[HTTP] [ID:%d] Start connect to Server [%s] [%s]", connID, req.Host(), server.Name)
	conn, err = server.Conn(req)
	if err != nil {
		if err == ErrorReject {
			log.Logger.Debugf("Reject [%s]", req.Host())
		} else {
			log.Logger.Errorf("[HTTP] [ID:%d] Connect to Server [%s] failed [%s] err: %s",
				connID, server.Name, req.Host(), err.Error())
			return
		}
	} else {
		log.Logger.Infof("[HTTP] [ClientConnID:%d] Bind to Server [ServerConnID:%d]", connID, conn.GetID())
	}
	return
}
