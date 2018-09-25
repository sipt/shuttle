package shuttle

import (
	"net"
	"net/http"
	"strconv"
	"bufio"
	"strings"
	"time"
	"github.com/sipt/shuttle/util"
	"github.com/sipt/shuttle/log"
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

// init MitMRules
func SetMitMRules(rs []string) {
	MitMRules = rs
}
func GetMitMRules() []string { // For controller API
	return MitMRules
}
func AppendMitMRules(r string) { // For controller API
	MitMRules = append(MitMRules, r)
	SetMimt(&Mitm{Rules: MitMRules})
}
func RemoveMitMRules(r string) { // For controller API
	for i, v := range MitMRules {
		if v == r {
			MitMRules[i] = MitMRules[len(MitMRules)-1]
			MitMRules = MitMRules[:len(MitMRules)-1]
			SetMimt(&Mitm{Rules: MitMRules})
			return
		}
	}
}

func HandleHTTP(co net.Conn) {
	log.Logger.Debug("start shuttle.IConn wrap net.Con")
	conn, err := NewDefaultConn(co, TCP)
	if err != nil {
		log.Logger.Errorf("[HTTP] shuttle.IConn wrap net.Conn failed: %v", err)
		return
	}
	log.Logger.Debugf("[HTTP] [ID:%d] shuttle.IConn wrap net.Conn success", conn.GetID())
	log.Logger.Debugf("[HTTP] [ID:%d] start read http request", conn.GetID())
	//prepare request
	req, hreq, err := prepareRequest(conn)
	if err != nil {
		log.Logger.Errorf("[HTTP] [ID:%d] prepareRequest failed: %s", conn.GetID(), err.Error())
		return
	}

	//request modify Or mock ?
	respBuf, err := RequestModifyOrMock(req, hreq, hreq.URL.Scheme == HTTP)
	if err != nil {
		log.Logger.Errorf("[HTTP] [ID:%d] request modify or mock failed: %s", conn.GetID(), err.Error())
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
		log.Logger.Errorf("[HTTP] [ID:%d] ConnectToServer failed [%s] err: %s", conn.GetID(), req.Host(), err)
	}

	//connect to server
	log.Logger.Infof("[HTTP] [ID:%d] Start connect to Server [%s]", conn.GetID(), s.Name)
	sc, err := s.Conn(req)
	if err != nil {
		if err == ErrorReject {
			log.Logger.Debugf("Reject [%s]", req.Target)
			boxChan <- &Box{
				Op: RecordAppend,
				Value: &Record{
					ID:       util.NextID(),
					Protocol: req.Protocol,
					Created:  time.Now(),
					Proxy:    s,
					Status:   RecordStatusReject,
					URL:      req.Target,
					Rule:     rule,
				},
			}
		} else {
			log.Logger.Error("[HTTP] [ID:%d] Connect to Server [%s] failed [%s] err: %s",
				conn.GetID(), s.Name, req.Host(), err.Error())
		}
		return
	}
	log.Logger.Infof("[HTTP] [ID:%d] Server [%s] Connected success", conn.GetID(), s.Name)
	log.Logger.Debugf("[HTTP] [ClientConnID:%d] Bind to [ServerConnID:%d]", conn.GetID(), sc.GetID())
	if req.Protocol == ProtocolHttps {
		_, err = conn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		if err != nil {
			log.Logger.Errorf("[HTTP] [ID:%d] reply https-CONNECT failed: %s", conn.GetID(), err.Error())
			conn.Close()
			sc.Close()
			return
		}
	}
	//todo 白名单判断
	if IsPass(req) {
		// @fix time out bug
		//lc, err := TimerDecorate(conn, DefaultTimeOut, -1)
		//if err != nil {
		//	log.Logger.Error("Timer Decorate net.Conn failed: ", err)
		//	lc = conn
		//}
		hreq.Write(sc)
		direct := &DirectChannel{}
		direct.Transport(conn, sc)
		return
	}
	//MitM filter
	mitm := false
	if allowMitm && len(MitMRules) > 0 && req.Protocol == ProtocolHttps {
		for _, v := range MitMRules {
			if v == "*" { // 通配
				log.Logger.Debugf("[HTTP] [ID:%d] MitM filter [%s] use [%s]", conn.GetID(), req.Addr, v)
				mitm = true
				break
			} else if v == req.Addr { // 全区配
				log.Logger.Debugf("[HTTP] [ID:%d] MitM filter [%s] use [%s]", conn.GetID(), req.Addr, v)
				mitm = true
				break
			} else if v[0] == '*' && strings.HasSuffix(req.Addr, v[1:]) { // 后缀匹配
				log.Logger.Debugf("[HTTP] [ID:%d] MitM filter [%s] use [%s]", conn.GetID(), req.Addr, v)
				mitm = true
				break
			}
		}
		//MitM Decorate
		if mitm {
			log.Logger.Debugf("[HTTP] [ID:%d] MitM Decorate", conn.GetID())
			lct, sct, err := Mimt(conn, sc, req)
			if err != nil {
				log.Logger.Error("[HTTP] [ID:%d] MitM failed: %s", conn.GetID(), err.Error())
				conn.Close()
				sc.Close()
				return
			}
			conn, sc = lct, sct
		}
	}

	record := &Record{
		Protocol: req.Protocol,
		Created:  time.Now(),
		Proxy:    s,
		Status:   RecordStatusActive,
		URL:      req.Target,
		Rule:     rule,
	}

	// @fix time out bug
	//lc, err := TimerDecorate(conn, -1, -1)
	//if err != nil {
	//	log.Logger.Error("Timer Decorate net.Conn failed: ", err)
	//	lc = conn
	//}
	lc := conn
	//Dump Decorate
	if mitm {
		HttpTransport(lc, sc, record, allowDump, nil)
		return
	} else if req.Protocol == ProtocolHttp {
		HttpTransport(lc, sc, record, allowDump, hreq)
		return
	}
	record.ID = util.NextID()
	boxChan <- &Box{Op: RecordAppend, Value: record}
	sc.SetRecordID(record.ID)
	direct := &DirectChannel{}
	direct.Transport(lc, sc)
	boxChan <- &Box{record.ID, RecordStatus, RecordStatusCompleted}
}

func prepareRequest(conn IConn) (*Request, *http.Request, error) {
	br := bufio.NewReader(conn)
	hreq, err := http.ReadRequest(br)
	if err != nil {
		return nil, nil, err
	}
	log.Logger.Debugf("[ID:%d] [HTTP/HTTPS] %s:%s", conn.GetID(), hreq.URL.Hostname(), hreq.URL.Port())
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

func IsPass(req *Request) bool {
	if req.Addr == ControllerDomain {
		return true
	}
	port, _ := strToUint16(controllerPort)
	if (req.Addr == "localhost" || req.Addr == "127.0.0.1" || req.IP.String() == "127.0.0.1") && req.Port == port {
		return true
	}
	return false
}
