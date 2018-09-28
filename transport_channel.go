package shuttle

import (
	"bufio"
	"github.com/sipt/shuttle/log"
	"github.com/sipt/shuttle/pool"
	"github.com/sipt/shuttle/util"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type DirectChannel struct{}

func (d *DirectChannel) Transport(lc, sc IConn) {
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

func HttpTransport(lc, sc IConn, allowDump bool, first *http.Request) {
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

func (h *HttpChannel) Transport(lc, sc IConn, first *http.Request) (err error) {
	var (
		oldHreq, hreq *http.Request
		lcBuf         = bufio.NewReader(lc)
		scBuf         *bufio.Reader
		resp          *http.Response
		rule          *Rule
		server        *Server
		passed        bool // inner request
	)
	if sc != nil {
		scBuf = bufio.NewReader(sc)
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

			//request update
			resp = RequestModify(hreq, h.isHttps)
		}
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

		// rule filter
		if resp == nil && (sc == nil || hreq.URL.Host != oldHreq.URL.Host) {
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
		}
		if !passed {
			boxChan <- &Box{Op: RecordAppend, Value: record, ID: record.ID}
		}
		if sc != nil {
			log.Logger.Debugf("[ID:%d] [HttpChannel] [reqID:%d] HttpChannel Transport send record to boxChan", sc.GetID(), record.ID)
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
				log.Logger.Errorf("[ID:%d] [HttpChannel] HttpChannel Transport [hreq]->s: %v", sc.GetID(), err)
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
				log.Logger.Errorf("[ID:%d] [HttpChannel] HttpChannel Transport s->[b]: %v", sc.GetID(), err)
			}
			return
		}
		log.Logger.Debugf("[ID:%d] [HttpChannel] HttpChannel Transport return s->[b]", sc.GetID())
		ResponseModify(hreq, resp, h.isHttps)
		err = h.writeResponse(resp, lc, record.ID, h.allowDump && !passed)
		if err != nil {
			return
		}
	}
	return
}

// write response in connection
func (h *HttpChannel) writeResponse(resp *http.Response, to IConn, recordID int64, allowDump bool) (err error) {
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

func ConnectFilter(hreq *http.Request, connID int64) (rule *Rule, server *Server, conn IConn, err error) {
	req := &Request{}
	req.Addr = HostName(hreq)
	req.ConnID = connID
	ip := net.ParseIP(req.Addr)
	if ip == nil {
		req.Atyp = AddrTypeDomain
	} else if len(ip) == net.IPv4len {
		req.Atyp = AddrTypeIPv4
	} else {
		req.Atyp = AddrTypeIPv6
	}
	req.Cmd = CmdTCP
	port := hreq.URL.Port()
	if len(port) > 0 {
		req.Port, err = strToUint16(port)
		if err != nil {
			log.Logger.Errorf("[HTTP] [ID:%d] Port to int16 failed [%d] err: %s", connID, req.Port, err)
			return
		}
	}

	rule, server, err = FilterByReq(req)
	if err != nil {
		log.Logger.Errorf("[HTTP] [ID:%d] ConnectToServer failed [%s] err: %s", connID, req.Host(), err)
		return
	}

	if req.Port == 0 {
		if hreq.URL.Scheme == HTTP {
			req.Port = 80
		} else if hreq.URL.Scheme == HTTPS {
			req.Port = 443
		}
	}

	log.Logger.Infof("[HTTP] [ID:%d] Start connect to Server [%s]", connID, server.Name)
	conn, err = server.Conn(req)
	if err != nil {
		if err == ErrorReject {
			log.Logger.Debugf("Reject [%s]", req.Target)
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
