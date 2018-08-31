package shuttle

import (
	"net"
	"net/http"
	"strconv"
	"bufio"
	"strings"
	"time"
	"github.com/sipt/shuttle/util"
)

const (
	HTTP  = "http"
	HTTPS = "https"
)

var allowMitm = false
var allowDump = false
var MitMRules []string

func SetAllowMitm(b bool) {
	allowMitm = b
}
func SetAllowDump(b bool) {
	allowDump = b
}

func GetAllowMitm() bool {
	return allowMitm
}
func GetAllowDump() bool {
	return allowDump
}

func SetMitMRules(rs []string) {
	MitMRules = rs
}
func GetMitMRules() []string {
	return MitMRules
}

func HandleHTTP(co net.Conn) {
	Logger.Debug("start shuttle.IConn wrap net.Con")
	conn, err := NewDefaultConn(co, TCP)
	if err != nil {
		Logger.Errorf("shuttle.IConn wrap net.Conn failed: %v", err)
		return
	}
	Logger.Debugf("shuttle.IConn wrap net.Con success [ID:%d]", conn.GetID())
	Logger.Debugf("[ID:%d] start read http request", conn.GetID())
	//prepare request
	req, hreq, err := prepareRequest(conn)
	if err != nil {
		Logger.Error("prepareRequest failed: ", err)
		return
	}

	//request modify Or mock ?
	respBuf, err := RequestModifyOrMock(req, hreq, hreq.URL.Scheme == HTTP)
	if err != nil {
		Logger.Errorf("[ID:%d] request modify or mock failed: %v", conn.GetID(), err)
	}
	if len(respBuf) > 0 {
		conn.Write(respBuf)
		return
	}

	//inner controller domain
	if req.Addr == ControllerDomain {
		port, err := strconv.ParseUint(controllerPort, 10, 16)
		if err == nil {
			req.IP = []byte{127, 0, 0, 1}
			req.Port = uint16(port)
		}
	}

	//filter by Rules and DNS
	rule, s, err := FilterByReq(req)
	if err != nil {
		Logger.Error("ConnectToServer failed [", req.Host(), "] err: ", err)
	}

	//connect to server
	sc, err := s.Conn(req)
	if err != nil {
		if err == ErrorReject {
			Logger.Debugf("Reject [%s]", req.Target)
			recordChan <- &Record{
				ID:       util.NextID(),
				Protocol: req.Protocol,
				Created:  time.Now(),
				Proxy:    s,
				Status:   RecordStatusReject,
				URL:      req.Target,
				Rule:     rule,
			}
		} else {
			Logger.Error("ConnectToServer failed [", req.Host(), "] err: ", err)
		}
		return
	}
	Logger.Debugf("Bind [client-local](%d) [local-server](%d)", conn.GetID(), sc.GetID())
	if req.Protocol == ProtocolHttps {
		_, err = conn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		if err != nil {
			Logger.Error("reply https-CONNECT failed: ", err)
			conn.Close()
			sc.Close()
			return
		}
	}
	//todo 白名单判断
	if req.Addr == ControllerDomain {
		lc, err := TimerDecorate(conn, DefaultTimeOut, -1)
		if err != nil {
			Logger.Error("Timer Decorate net.Conn failed: ", err)
			lc = conn
		}
		hreq.Write(sc)
		direct := &DirectChannel{}
		direct.Transport(lc, sc)
		return
	}
	//MitM filter
	mitm := false
	if allowMitm && len(MitMRules) > 0 && req.Protocol == ProtocolHttps {
		for _, v := range MitMRules {
			if v == "*" { // 通配
				Logger.Debugf("[ID:%d] [HTTP/HTTPS] MitM filter [%s] use [%s]", conn.GetID(), req.Addr, v)
				mitm = true
				break
			} else if v == req.Addr { // 全区配
				Logger.Debugf("[ID:%d] [HTTP/HTTPS] MitM filter [%s] use [%s]", conn.GetID(), req.Addr, v)
				mitm = true
				break
			} else if v[0] == '*' && strings.HasSuffix(req.Addr, v[1:]) { // 后缀匹配
				Logger.Debugf("[ID:%d] [HTTP/HTTPS] MitM filter [%s] use [%s]", conn.GetID(), req.Addr, v)
				mitm = true
				break
			}
		}
		//MitM Decorate
		if mitm {
			Logger.Debugf("[ID:%d] [HTTP/HTTPS] MitM Decorate", conn.GetID())
			lct, sct, err := Mimt(conn, sc, req)
			if err != nil {
				Logger.Error("[HTTPS] MitM failed: ", err)
				conn.Close()
				sc.Close()
				return
			}
			conn, sc = lct, sct
		}
	}

	record := &Record{
		ID:       sc.GetID(),
		Protocol: req.Protocol,
		Created:  time.Now(),
		Proxy:    s,
		Status:   RecordStatusActive,
		URL:      req.Target,
		Rule:     rule,
	}
	lc, err := TimerDecorate(conn, DefaultTimeOut, -1)
	if err != nil {
		Logger.Error("Timer Decorate net.Conn failed: ", err)
		lc = conn
	}
	//Dump Decorate
	if mitm {
		HttpTransport(lc, sc, record, allowDump, nil)
		return
	} else if req.Protocol == ProtocolHttp {
		HttpTransport(lc, sc, record, allowDump, hreq)
		return
	} else {
		recordChan <- record
	}

	direct := &DirectChannel{}
	direct.Transport(lc, sc)
}

func prepareRequest(conn IConn) (*Request, *http.Request, error) {
	br := bufio.NewReader(conn)
	hreq, err := http.ReadRequest(br)
	if err != nil {
		return nil, nil, err
	}
	Logger.Debugf("[ID:%d] [HTTP/HTTPS] %s:%s", conn.GetID(), hreq.URL.Hostname(), hreq.URL.Port())
	req := &Request{
		Ver:    socksVer5,
		Cmd:    CmdTCP,
		Atyp:   AddrTypeDomain,
		ConnID: conn.GetID(),
	}
	if hreq.URL.Scheme == HTTP {
		req.Protocol = ProtocolHttp
		if req.Port == 0 {
			req.Port = 80
		}
	} else if hreq.Method == http.MethodConnect {
		req.Protocol = ProtocolHttps
		if req.Port == 0 {
			req.Port = 443
		}
	}
	return req, hreq, nil
}

func strToUint16(v string) (i uint16, err error) {
	r, err := strconv.ParseUint(v, 10, 2*8)
	if err == nil {
		i = uint16(r)
	}
	return
}
