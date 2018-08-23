package shuttle

import (
	"net/http"
	"regexp"
	"net/url"
	"os"
	"strconv"
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

func RequestModify(req *http.Request) *http.Response {
	if len(reqPolicies) == 0 {
		return nil
	}
	for _, v := range reqPolicies {
		if v.rex.MatchString(req.URL.String()) {
			switch v.Type {
			case ModifyMock:
				return modifyMock(v, req)
			case ModifyUpdate:
				modifyUpdate(v, req)
				return nil
			}
		}
	}
	return nil
}

func modifyMock(v *ModifyPolicy, req *http.Request) *http.Response {
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
			Logger.Debugf("[Http Modify] [Mock] response set body [Status:%d]", e.Value)
			resp.StatusCode, err = strconv.Atoi(e.Value)
			if err != nil {
				resp.StatusCode = 200
			}
		}
	}
	return resp
}

func modifyUpdate(v *ModifyPolicy, req *http.Request) {
	for _, e := range v.MVs {
		switch e.Type {
		case ModifyTypeURL:
			u, err := url.Parse(e.Value)
			if err != nil {
				Logger.Errorf("[HTTP MODIFY] parse [%s] to url failed: %v", e.Value, err)
				return
			}
			if req.URL.Scheme != u.Scheme {
				Logger.Errorf("[HTTP MODIFY] not support [%s] to [%s]", req.URL.Scheme, u.Scheme)
				return
			}
			Logger.Debugf("[Http Modify] [Update] response set URL [%s]", e.Value)
			req.URL = u
			req.Host = u.Host
		case ModifyTypeHeader:
			Logger.Debugf("[Http Modify] [Update] response set Header [%s:%s]", e.Key, e.Value)
			req.Header.Set(e.Key, e.Value)
		}
	}
}

func ResponseModify(req *http.Request, resp *http.Response) {
	if len(respPolicies) == 0 {
		return
	}
	for _, v := range respPolicies {
		if v.rex.MatchString(req.URL.String()) {
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
