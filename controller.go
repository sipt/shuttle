package shuttle

import (
	"net/http"
	"github.com/gin-gonic/gin/json"
	"fmt"
	"errors"
	"text/template"
	"time"
	"github.com/apaxa-go/helper/strconvh"
	"encoding/base64"
	"io"
	"bytes"
)

var errNotSupportMethod = errors.New("not support method")
var recordsTemplate *template.Template

func StartController(port string) error {
	var err error
	recordsTemplate, err = template.New("records.html").Funcs(template.FuncMap{
		"Format": func(now time.Time) interface{} {
			return now.Format("01-02 15:04:05")
		},
		"Policy": func(rule *Rule, server *Server) interface{} {
			return fmt.Sprintf("%s(%s %s)", server.Name, rule.Type, rule.Value)
		},
		"Add": func(left int, right int) int {
			return left + right
		},
	}).ParseFiles("template/records.html") // Parse template file.
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/api/dns/cache", dnsCacheHandler)
	http.HandleFunc("/api/records", apiRecordHandler)
	http.HandleFunc("/api/dump", apiDumpHandler)

	http.HandleFunc("/records", recordHandler)
	http.HandleFunc("/cert", certHandler)
	Logger.Info("[Controller] listen to :", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		Logger.Errorf("[Controller] controller server start failed: %v", err)
	}
	return err
}

// path :/api/dns/cache
func dnsCacheHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		bytes, err := json.Marshal(_CacheDNS.List())
		toClient(w, bytes, err)
	case http.MethodDelete:
		_CacheDNS.Clear()
		toClient(w, nil, nil)
	default:
		toClient(w, nil, errNotSupportMethod)
	}
}

// path: /api/records
func apiRecordHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		bytes, err := json.Marshal(storage.List())
		toClient(w, bytes, err)
	case http.MethodDelete:
		storage.Clear()
		toClient(w, nil, nil)
	default:
		toClient(w, nil, errNotSupportMethod)
	}
}

// path: /api/dump
func apiDumpHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		allow_dump := req.FormValue("allow_dump")
		switch allow_dump {
		case "true":
			allowDump = true
		case "false":
			allowDump = false
		default:
			toClient(w, nil, fmt.Errorf("allow_dump value must be [true/false]: %v", allow_dump))
		}
		allow_mitm := req.FormValue("allow_mitm")
		switch allow_mitm {
		case "true":
			mitm = true
		case "false":
			mitm = false
		default:
			toClient(w, nil, fmt.Errorf("allow_mitm value error: %v", allow_mitm))
		}
		toClient(w, nil, nil)
	case http.MethodGet:
		idStr := req.FormValue("id")
		id, err := strconvh.ParseInt64(idStr)
		if err != nil {
			toClient(w, nil, err)
		}
		r := storage.Get(id)
		if r == nil {
			toClient(w, nil, errors.New(idStr+" not exist"))
			return
		}
		if r.Status != RecordStatusCompleted {
			toClient(w, nil, errors.New(idStr+" not Completed"))
			return
		}
		reply := &struct {
			RequestData  string
			ResponseData string
		}{}
		data, err := dump.ReadRequest(id)
		if err != nil {
			toClient(w, nil, err)
			return
		}
		reply.RequestData = base64.RawStdEncoding.EncodeToString(data)
		data, err = dump.ReadResponse(id)
		if err != nil {
			toClient(w, nil, err)
			return
		}
		reply.ResponseData = base64.RawStdEncoding.EncodeToString(data)

		toClient(w, []byte(fmt.Sprintf(`{"Request":"%s", "Response":"%s"}`, reply.RequestData, reply.ResponseData)), err)
	default:
		toClient(w, nil, errNotSupportMethod)
	}
}

// path: /records
func recordHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		err := recordsTemplate.Execute(w, storage.List()) // merge.
		if err != nil {
			toClient(w, nil, err)
			Logger.Errorf("records template error:", err)
		}
	default:
		toClient(w, nil, errNotSupportMethod)
	}
}

// path: /cert
func certHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		if len(caBytes) == 0 {
			toClient(w, nil, errors.New("please generate CA"))
			return
		}
		w.Header().Add("Content-Type", "application/octet-stream")
		w.Header().Add("content-disposition", "attachment; filename=\"Shuttle.cer\"")
		_, err := io.Copy(w, bytes.NewBuffer(caBytes))
		if err != nil {
			toClient(w, nil, err)
			return
		}
		toClient(w, nil, nil)
	case http.MethodPost:
		err := generateCA()
		if err != nil {
			toClient(w, nil, err)
			return
		}
		toClient(w, nil, nil)
	default:
		toClient(w, nil, errNotSupportMethod)
	}
}

func toClient(w http.ResponseWriter, data []byte, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		if err == errNotSupportMethod {
			w.WriteHeader(405)
			return
		}
		w.WriteHeader(500)
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(fmt.Sprintf(`{"code":1, "message":"%s"}`, err)))
	} else if len(data) > 0 {
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf(`{"code":0, "data":%s}`, data)))
	} else {
		w.WriteHeader(200)
		w.Write([]byte(`{"code":0}`))
	}
}
