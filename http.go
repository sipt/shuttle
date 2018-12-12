package shuttle

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"context"
	"github.com/sipt/shuttle/config"
	connect "github.com/sipt/shuttle/conn"
	"github.com/sipt/shuttle/log"
	"github.com/sipt/shuttle/proxy"
	rule2 "github.com/sipt/shuttle/rule"
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

// init MitMRules
func SetMitMRules(rs []string) {
	MitMRules = rs
}
func GetMitMRules() []string { // For controller API
	return MitMRules
}
func AppendMitMRules(r string) { // For controller API
	MitMRules = append(MitMRules, r)
	conf := config.CurrentConfig()
	conf.Mitm.Rules = MitMRules
	config.SaveConfig(config.CurrentConfigFile(), conf)
}
func RemoveMitMRules(r string) { // For controller API
	for i, v := range MitMRules {
		if v == r {
			MitMRules[i] = MitMRules[len(MitMRules)-1]
			MitMRules = MitMRules[:len(MitMRules)-1]
			conf := config.CurrentConfig()
			conf.Mitm.Rules = MitMRules
			config.SaveConfig(config.CurrentConfigFile(), conf)
			return
		}
	}
}

func HandleHTTP(co net.Conn) {
	log.Logger.Debug("start conn.IConn wrap net.Con")
	conn, err := connect.NewDefaultConn(co, connect.TCP)
	if err != nil {
		log.Logger.Errorf("[HTTP] shuttle.IConn wrap net.Conn failed: %v", err)
		return
	}
	log.Logger.Debugf("[HTTP] [ID:%d] shuttle.IConn wrap net.Conn success", conn.GetID())
	log.Logger.Debugf("[HTTP] [ID:%d] start read http request", conn.GetID())
	//prepare request
	hreq, err := prepareRequest(conn)
	if err != nil {
		if err != io.EOF {
			log.Logger.Errorf("[HTTP] [ID:%d] prepareRequest failed: %s", conn.GetID(), err.Error())
		}
		return
	}

	//switch hreq.Proto {
	//case "HTTP/2":
	//	ProxyHTTP2()
	//case "HTTP/1.1":
	if hreq.URL.Scheme == HTTP { // HTTP
		ProxyHTTP(conn, hreq)
	} else { // HTTPS
		ProxyHTTPS(conn, hreq)
	}
	//}
}

func ProxyHTTP(lc connect.IConn, hreq *http.Request) {
	HttpTransport(lc, nil, allowDump, hreq)
}
func ProxyHTTPS(lc connect.IConn, hreq *http.Request) {
	// Handshake
	_, err := lc.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	if err != nil {
		log.Logger.Errorf("[HTTPS] [ID:%d] reply https-CONNECT failed: %s", lc.GetID(), err.Error())
		lc.Close()
		return
	}
	domain := hreq.URL.Hostname()
	rule, server, sc, err := ConnectFilter(hreq, lc.GetID())
	record := &Record{
		Protocol: HTTPS,
		Created:  time.Now(),
		Status:   RecordStatusActive,
		URL:      hreq.URL.String(),
		Proxy:    server,
		Rule:     rule,
	}
	if hreq.URL.Scheme == "" {
		record.URL = "https:" + record.URL
	}
	if err != nil {
		if err == ErrorReject {
			record.Status = RecordStatusReject
		} else {
			record.Status = RecordStatusFailed
			record.Rule = rule2.FailedRule
			record.Proxy = proxy.FailedServer
		}
		record.ID = util.NextID()
		boxChan <- &Box{Op: RecordAppend, Value: record}
		return
	}
	// MitM
	mitm := false
	if allowMitm {
		for _, v := range MitMRules {
			if v == "*" { // 通配
				log.Logger.Debugf("[HTTPS] [ID:%d] MitM RuleFilter [%s] use [%s]", lc.GetID(), domain, v)
				mitm = true
				break
			} else if v == domain { // 全区配
				log.Logger.Debugf("[HTTPS] [ID:%d] MitM RuleFilter [%s] use [%s]", lc.GetID(), domain, v)
				mitm = true
				break
			} else if v[0] == '*' && strings.HasSuffix(domain, v[1:]) { // 后缀匹配
				log.Logger.Debugf("[HTTPS] [ID:%d] MitM RuleFilter [%s] use [%s]", lc.GetID(), domain, v)
				mitm = true
				break
			}
		}
	}
	//MitM Decorate
	if mitm {
		log.Logger.Debugf("[HTTPS] [ID:%d] MitM Decorate", lc.GetID())
		lct, sct, err := Mimt(lc, sc)
		if err != nil {
			log.Logger.Error("[HTTPS] [ID:%d] MitM failed: %s", lc.GetID(), err.Error())
			record.Status = RecordStatusFailed
			boxChan <- &Box{Op: RecordAppend, Value: record}
			lc.Close()
			sc.Close()
			return
		}
		lc, sc = lct, sct
		ctx := sc.Context()
		ctx = context.WithValue(ctx, "rule", rule)
		ctx = context.WithValue(ctx, "server", server)
		sc.SetContext(ctx)
		HttpTransport(lc, sc, allowDump, nil)
		return
	}

	record.ID = util.NextID()
	boxChan <- &Box{Op: RecordAppend, Value: record}
	sc.SetRecordID(record.ID)
	direct := &DirectChannel{}
	direct.Transport(lc, sc)
	boxChan <- &Box{record.ID, RecordStatus, RecordStatusCompleted}
}

func ProxyHTTP2() {

}

func prepareRequest(conn connect.IConn) (*http.Request, error) {
	br := bufio.NewReader(conn)
	hreq, err := http.ReadRequest(br)
	if err != nil {
		return nil, err
	}
	log.Logger.Debugf("[ID:%d] [HTTP/HTTPS] %s:%s", conn.GetID(), hreq.URL.Hostname(), hreq.URL.Port())
	return hreq, nil
}

func StrToUint16(v string) (i uint16, err error) {
	r, err := strconv.ParseUint(v, 10, 2*8)
	if err == nil {
		i = uint16(r)
	}
	return
}

func IsPass(host, port, ip string) bool {
	if host == ControllerDomain {
		return true
	}
	if (host == "localhost" || host == "127.0.0.1" || ip == "127.0.0.1") && ControllerPort == port {
		return true
	}
	return false
}
