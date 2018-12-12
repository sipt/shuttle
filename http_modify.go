package shuttle

import (
	"bytes"
	"fmt"
	"github.com/sipt/shuttle/config"
	"github.com/sipt/shuttle/log"
	"github.com/sipt/shuttle/pool"
	"github.com/sipt/shuttle/proxy"
	"github.com/sipt/shuttle/rule"
	"github.com/sipt/shuttle/util"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	ModifyTypeURL    = "URL"
	ModifyTypeHeader = "HEADER"
	ModifyTypeStatus = "STATUS"
	ModifyTypeBody   = "BODY"

	ModifyMock   = "MOCK"
	ModifyUpdate = "UPDATE"

	BodyFileDir = "RespFiles"
)

var reqPolicies []*ModifyPolicy
var respPolicies []*ModifyPolicy

type IHttpModifyConfig interface {
	GetHTTPMap() *config.HttpMap
}

func ApplyHTTPModifyConfig(config IHttpModifyConfig) (err error) {
	httpMap := config.GetHTTPMap()
	if httpMap != nil {
		if len(httpMap.ReqMap) > 0 {
			reqps := make([]*ModifyPolicy, len(httpMap.ReqMap))
			for i, v := range httpMap.ReqMap {
				switch v.Type {
				case ModifyMock, ModifyUpdate:
				default:
					return fmt.Errorf("resolve config file [Http-Map] [Req-Map] not support [%s]", v.Type)
				}
				reqps[i] = &ModifyPolicy{
					Type:   v.Type,
					UrlRex: v.UrlRex,
				}
				reqps[i].rex, err = regexp.Compile(v.UrlRex)
				if err != nil {
					return fmt.Errorf("resolve config file [Http-Map] [%s] failed: %v", v.UrlRex, err)
				}
				if len(v.Items) > 0 {
					reqps[i].MVs = make([]*ModifyValue, len(v.Items))
					for j, e := range v.Items {
						if len(e) != 3 {
							return fmt.Errorf("resolve config file [Http-Map] failed: %v, item's count must be 3", e)
						}
						switch e[0] {
						case ModifyTypeURL, ModifyTypeHeader, ModifyTypeStatus, ModifyTypeBody:
						default:
							return fmt.Errorf("resolve config file [Http-Map] [Req-Map] not support [%s]", v.Type)
						}
						reqps[i].MVs[j] = &ModifyValue{
							Type:  e[0],
							Key:   e[1],
							Value: e[2],
						}
					}
				}
			}
			reqPolicies = reqps
		}

		if len(httpMap.ReqMap) > 0 {
			respps := make([]*ModifyPolicy, len(httpMap.RespMap))
			for i, v := range httpMap.RespMap {
				switch v.Type {
				case ModifyMock, ModifyUpdate:
				default:
					return fmt.Errorf("resolve config file [Http-Map] [Resp-Map] not support [%s]", v.Type)
				}
				respps[i] = &ModifyPolicy{
					Type:   v.Type,
					UrlRex: v.UrlRex,
				}
				respps[i].rex, err = regexp.Compile(v.UrlRex)
				if err != nil {
					return fmt.Errorf("resolve config file [Http-Map] [%s] failed: %v", err)
				}
				if len(v.Items) > 0 {
					respps[i].MVs = make([]*ModifyValue, len(v.Items))
					for j, e := range v.Items {
						if len(e) != 3 {
							return fmt.Errorf("resolve config file [Http-Map] failed: %v, item's count must be 3", e)
						}
						switch e[0] {
						case ModifyTypeHeader, ModifyTypeStatus, ModifyTypeBody:
						default:
							return fmt.Errorf("resolve config file [Http-Map] [Req-Map] not support [%s]", v.Type)
						}
						respps[i].MVs[j] = &ModifyValue{
							Type:  e[0],
							Key:   e[1],
							Value: e[2],
						}
					}
				}
			}
			respPolicies = respps
		}
	}

	return
}

func RequestModifyOrMock(req *HttpRequest, hreq *http.Request, isHttps bool) (respBuf []byte, err error) {
	//request update
	resp := RequestModify(hreq, isHttps)
	req.domain = hreq.URL.Hostname()
	if net.ParseIP(req.domain) != nil {
		req.ip = req.domain
		req.domain = ""
	}
	req.port = hreq.URL.Port()
	req.target = hreq.URL.String()
	if strings.HasPrefix(req.target, "//") {
		req.target = req.target[2:]
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
		boxChan <- &Box{
			Op: RecordAppend,
			Value: &Record{
				ID:       id,
				Protocol: req.protocol,
				Created:  time.Now(),
				Proxy:    proxy.MockServer,
				Status:   RecordStatusCompleted,
				Dumped:   allowDump,
				URL:      req.target,
				Rule:     &rule.Rule{},
			},
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
			log.Logger.Debugf("[Http Modify] [Mock] response set Header [%s:%s]", e.Key, e.Value)
			resp.Header.Set(e.Key, e.Value)
		case ModifyTypeBody:
			file, err := os.Open(e.Value)
			if err != nil {
				log.Logger.Errorf("[HTTP MODIFY] open mock file failed: %v", err)
				return nil
			}
			status, err := file.Stat()
			if err != nil {
				log.Logger.Errorf("[HTTP MODIFY] read mock file [FileInfo] failed: %v", err)
				return nil
			}
			log.Logger.Debugf("[Http Modify] [Mock] response set body [ContentLength:%d]", status.Size())
			resp.ContentLength = status.Size()
			resp.Body = file
		case ModifyTypeStatus:
			log.Logger.Debugf("[Http Modify] [Mock] response set body [Status:%s]", e.Value)
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
				log.Logger.Errorf("[HTTP MODIFY] parse [%s] to url failed: %v", e.Value, err)
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
				log.Logger.Errorf("[HTTP MODIFY] not support [%s] to [%s]", req.URL.Scheme, u.Scheme)
				return
			}
			log.Logger.Debugf("[Http Modify] [Update] response set URL [%s]", e.Value)
			req.URL = u
		case ModifyTypeHeader:
			log.Logger.Debugf("[Http Modify] [Update] response set Header [%s:%s]", e.Key, e.Value)
			req.Header.Set(e.Key, e.Value)
		}
	}
}

func ResponseModify(req *http.Request, resp *http.Response, isHttps bool) {
	if len(respPolicies) == 0 {
		return
	}
	if req.URL == nil {
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
					log.Logger.Debugf("[Http Modify] [Update] response set Header [%s, %s]", e.Key, e.Value)
					resp.Header.Set(e.Key, e.Value)
				case ModifyTypeStatus:
					log.Logger.Debugf("[Http Modify] [Update] response  [Status:%s]", e.Value)
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
