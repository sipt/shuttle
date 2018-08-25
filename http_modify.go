package shuttle

import (
	"net/http"
	"regexp"
	"net/url"
	"os"
	"strconv"
	"net"
	"strings"
	"bytes"
	"time"
	"github.com/sipt/shuttle/util"
	"github.com/sipt/shuttle/pool"
)

const (
	ModifyTypeURL    = "URL"
	ModifyTypeHeader = "HEADER"
	ModifyTypeStatus = "STATUS"
	ModifyTypeBody   = "BODY"

	ModifyMock   = "MOCK"
	ModifyUpdate = "UPDATE"

	BodyFileDir = "./RespFiles/"
)

var reqPolicies []*ModifyPolicy
var respPolicies []*ModifyPolicy

func InitHttpModify(req []*ModifyPolicy, resp []*ModifyPolicy) {
	reqPolicies = req
	respPolicies = resp
}

func ClearHttpModify() {
	reqPolicies = nil
	respPolicies = nil
}

func RequestModifyOrMock(req *Request, hreq *http.Request, isHttps bool) (respBuf []byte, err error) {
	//request update
	resp := RequestModify(hreq, isHttps)
	req.Addr = hreq.URL.Hostname()
	req.IP = net.ParseIP(req.Addr)
	if port := hreq.URL.Port(); len(port) > 0 {
		req.Port, err = strToUint16(port)
		if err != nil {
			Logger.Error("http port error:" + port)
			return
		}
	}
	req.Target = hreq.URL.String()
	if strings.HasPrefix(req.Target, "//") {
		req.Target = req.Target[2:]
	}
	if resp != nil { // response mock ?
		buffer := &bytes.Buffer{}
		err = resp.Write(buffer)
		if err != nil {
			return
		}
		respBuf = buffer.Bytes()
		//mock record to storage
		id := util.NextID()
		recordChan <- &Record{
			ID:       id,
			Protocol: req.Protocol,
			Created:  time.Now(),
			Proxy:    &Server{Name: "MOCK"},
			Status:   RecordStatusCompleted,
			Dumped:   allowDump,
			URL:      req.Target,
			Rule:     &Rule{},
		}
		if allowDump {
			go func(id int64, respBuf []byte) {
				dump.InitDump(id)
				writer := bytes.NewBuffer(pool.GetBuf()[:0])
				hreq.Write(writer)
				dump.WriteRequest(id, writer.Bytes())
				dump.WriteResponse(id, respBuf)
				dump.Complete(id)
			}(id, respBuf)
		}
	}
	return
}

func RequestModify(req *http.Request, isHttps bool) *http.Response {
	if len(reqPolicies) == 0 {
		return nil
	}
	l := req.URL.String()
	if req.URL.Host == "" {
		if isHttps {
			l = "https://" + req.Host + l
		} else {
			l = "http://" + req.Host + l
		}
	}
	for _, v := range reqPolicies {
		if v.rex.MatchString(l) {
			switch v.Type {
			case ModifyMock:
				return modifyMock(v, req, isHttps)
			case ModifyUpdate:
				modifyUpdate(v, req, isHttps)
				return nil
			}
		}
	}
	return nil
}

func modifyMock(v *ModifyPolicy, req *http.Request, _ bool) *http.Response {
	resp := &http.Response{
		StatusCode:    200,
		Proto:         req.Proto,
		ContentLength: 0,
		Header:        make(http.Header),
	}
	for _, e := range v.MVs {
		var err error
		switch e.Type {
		case ModifyTypeHeader:
			Logger.Debugf("[Http Modify] [Mock] response set Header [%s:%s]", e.Key, e.Value)
			resp.Header.Set(e.Key, e.Value)
		case ModifyTypeBody:
			file, err := os.Open(BodyFileDir + e.Value)
			if err != nil {
				Logger.Errorf("[HTTP MODIFY] open mock file failed: %v", err)
				return nil
			}
			status, err := file.Stat()
			if err != nil {
				Logger.Errorf("[HTTP MODIFY] read mock file [FileInfo] failed: %v", err)
				return nil
			}
			Logger.Debugf("[Http Modify] [Mock] response set body [ContentLength:%d]", status.Size())
			resp.ContentLength = status.Size()
			resp.Body = file
		case ModifyTypeStatus:
			Logger.Debugf("[Http Modify] [Mock] response set body [Status:%s]", e.Value)
			resp.StatusCode, err = strconv.Atoi(e.Value)
			if err != nil {
				resp.StatusCode = 200
			}
		}
	}
	return resp
}

func modifyUpdate(v *ModifyPolicy, req *http.Request, isHttps bool) {
	for _, e := range v.MVs {
		switch e.Type {
		case ModifyTypeURL:
			l := req.URL.String()
			if req.URL.Host == "" {
				if isHttps {
					l = "https://" + req.Host + l
				} else {
					l = "http://" + req.Host + l
				}
			}
			l = v.rex.ReplaceAllString(l, e.Value)
			u, err := url.Parse(l)
			if err != nil {
				Logger.Errorf("[HTTP MODIFY] parse [%s] to url failed: %v", e.Value, err)
				return
			}
			req.Host = u.Host
			if req.URL.Scheme == "" {
				u.Scheme = ""
			}
			if req.Host == "" {
				u.Host = ""
			}
			if req.URL.Scheme != u.Scheme {
				Logger.Errorf("[HTTP MODIFY] not support [%s] to [%s]", req.URL.Scheme, u.Scheme)
				return
			}
			Logger.Debugf("[Http Modify] [Update] response set URL [%s]", e.Value)
			req.URL = u
		case ModifyTypeHeader:
			Logger.Debugf("[Http Modify] [Update] response set Header [%s:%s]", e.Key, e.Value)
			req.Header.Set(e.Key, e.Value)
		}
	}
}

func ResponseModify(req *http.Request, resp *http.Response, isHttps bool) {
	if len(respPolicies) == 0 {
		return
	}
	l := req.URL.String()
	if req.URL.Host == "" {
		if isHttps {
			l = "https://" + req.Host + l
		} else {
			l = "http://" + req.Host + l
		}
	}
	for _, v := range respPolicies {
		if v.rex.MatchString(l) {
			for _, e := range v.MVs {
				switch e.Type {
				case ModifyTypeHeader:
					Logger.Debugf("[Http Modify] [Update] response set Header [%s, %s]", e.Key, e.Value)
					resp.Header.Set(e.Key, e.Value)
				case ModifyTypeStatus:
					Logger.Debugf("[Http Modify] [Update] response  [Status:%s]", e.Value)
					code, err := strconv.Atoi(e.Value)
					if err == nil {
						resp.StatusCode = code
						resp.Status = ""
					}
				}
			}
		}
	}
}

type ModifyPolicy struct {
	Type   string
	UrlRex string
	rex    *regexp.Regexp
	MVs    []*ModifyValue
}

type ModifyValue struct {
	Type  string
	Key   string
	Value string
}
