package inbound

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/sipt/shuttle/constant"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/listener"
	"github.com/sirupsen/logrus"

	connpkg "github.com/sipt/shuttle/conn"
)

const (
	ParamsKeyAuthType = "auth_type"

	AuthTypeBasic  = "basic"
	AuthTypeBearer = "bearer"

	ParamsKeyUser     = "user"
	ParamsKeyPassword = "password"
	ParamsKeyToken    = "token"
)

func init() {
	Register("http", newHTTPInbound)
}

func newHTTPInbound(addr string, params map[string]string) (listen func(context.Context, listener.HandleFunc), err error) {
	authType, ok := params[ParamsKeyAuthType]
	authFunc := func(r *http.Request) bool { return true }
	if ok {
		switch authType {
		case AuthTypeBasic:
			authFunc, err = newBasicAuth(params)
		case AuthTypeBearer:
			authFunc, err = newBearerAuth(params)
		default:
			err = errors.Errorf("[http.Inbound] is not support")
		}
		if err != nil {
			return
		}
	}
	dial, err := listener.Get("tcp", addr)
	if err != nil {
		return nil, err
	}
	logrus.WithField("addr", "http://"+addr).Info("http listen starting")
	return func(ctx context.Context, handle listener.HandleFunc) {
		dial(ctx, newHttpHandleFunc(authFunc, handle))
	}, nil
}

func newHttpHandleFunc(authFunc func(r *http.Request) bool, handle listener.HandleFunc) func(conn connpkg.ICtxConn) {
	if authFunc == nil {
		authFunc = func(r *http.Request) bool { return true }
	}
	return func(conn connpkg.ICtxConn) {
		defer conn.Close()
		logrus.Debug("start dial HTTP/HTTPS connection")
		for {
			req, err := http.ReadRequest(bufio.NewReader(conn))
			if err != nil {
				if err == io.EOF {
					return
				}
				logrus.WithError(err).Error("[http.Inbound] parse to http request failed")
				_ = conn.Close()
				return
			}
			if !authFunc(req) {
				resp := &http.Response{
					StatusCode: http.StatusProxyAuthRequired,
				}
				err = resp.Write(conn)
				if err != nil {
					logrus.WithError(err).Error("[http.Inbound] write to response failed")
					_ = conn.Close()
					return
				}
				return
			}
			if req.Method == http.MethodConnect {
				c, err := httpsHandshake(req, conn)
				if err != nil {
					logrus.WithError(err).Error("[http.Inbound] https handshake failed")
					_ = conn.Close()
					return
				}
				handle(c)
				return
			} else {
				c, err := httpHandshake(req, conn)
				if err != nil {
					logrus.WithError(err).Error("[http.Inbound] http handshake failed")
					_ = conn.Close()
					return
				}
				handle(c)
			}
		}
	}
}

func newBasicAuth(params map[string]string) (func(*http.Request) bool, error) {
	user := params[ParamsKeyUser]
	if len(user) == 0 {
		return nil, errors.New("[user] is empty")
	}
	password := params[ParamsKeyPassword]
	if len(password) == 0 {
		return nil, errors.New("[password] is empty")
	}
	authorization := "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", user, password)))
	return func(req *http.Request) bool {
		return req.Header.Get("Proxy-Authorization") == authorization
	}, nil
}

func newBearerAuth(params map[string]string) (func(*http.Request) bool, error) {
	token := params[ParamsKeyToken]
	if len(token) == 0 {
		return nil, errors.New("[token] is empty")
	}
	authorization := "Bearer " + token
	return func(req *http.Request) bool {
		return req.Header.Get("Proxy-Authorization") == authorization
	}, nil
}

func httpHandshake(req *http.Request, c connpkg.ICtxConn) (connpkg.ICtxConn, error) {
	ctxReq := &request{
		uri:    req.URL.String(),
		domain: req.URL.Hostname(),
	}
	var err error
	if port := req.URL.Port(); port != "" {
		ctxReq.port, err = strconv.Atoi(port)
		if err != nil {
			return nil, errors.Errorf("port [%s] is error: %s", port, err.Error())
		}
	} else {
		ctxReq.SetPort(80)
	}
	hc := &httpConn{
		header:   &bytes.Buffer{},
		body:     req.Body,
		ICtxConn: connpkg.NewConn(c, context.WithValue(c, constant.KeyRequestInfo, ctxReq)),
		Mutex:    &sync.Mutex{},
	}
	req.Body = nil
	err = req.Write(hc.header)
	if err != nil {
		return nil, errors.Wrapf(err, "read request failed")
	}
	return hc, nil
}

func httpsHandshake(req *http.Request, c connpkg.ICtxConn) (connpkg.ICtxConn, error) {
	_, err := c.Write([]byte(fmt.Sprintf("%s 200 Connection established\r\n\r\n", req.Proto)))
	if err != nil {
		return nil, errors.Wrapf(err, "https handshake failed")
	}
	ctxReq := &request{
		uri:    req.URL.String(),
		domain: req.URL.Hostname(),
	}
	if port := req.URL.Port(); port != "" {
		ctxReq.port, err = strconv.Atoi(port)
		if err != nil {
			return nil, errors.Errorf("port [%s] is error: %s", port, err.Error())
		}
	} else {
		ctxReq.SetPort(443)
	}
	type witchContext interface {
		WithContext(ctx context.Context)
	}
	c.(witchContext).WithContext(context.WithValue(c, constant.KeyRequestInfo, ctxReq))
	return c, nil
}

type httpConn struct {
	header *bytes.Buffer
	body   io.ReadCloser
	connpkg.ICtxConn
	*sync.Mutex
}

func (h *httpConn) Read(b []byte) (int, error) {
	if h.header.Len() > 0 {
		h.Lock()
		defer h.Unlock()
		return h.header.Read(b)
	}
	return h.body.Read(b)
}

func (h *httpConn) Close() error {
	return h.body.Close()
}

type request struct {
	network     string
	domain      string
	uri         string
	ip          net.IP
	port        int
	countryCode string
}

func (r *request) Network() string {
	return r.network
}
func (r *request) Domain() string {
	return r.domain
}
func (r *request) URI() string {
	return r.uri
}
func (r *request) IP() net.IP {
	return r.ip
}
func (r *request) CountryCode() string {
	return r.countryCode
}
func (r *request) Port() int {
	return r.port
}
func (r *request) SetIP(in net.IP) {
	r.ip = in
}
func (r *request) SetPort(in int) {
	r.port = in
}
func (r *request) SetCountryCode(in string) {
	r.countryCode = in
}
