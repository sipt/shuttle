package shuttle

import (
	"net"
	"net/http"
	"strconv"
	"errors"
	"bufio"
	"bytes"
	"github.com/sipt/shuttle/pool"
	"strings"
)

const (
	HTTP  = "http"
	HTTPS = "https"
)

var allowMitm = false
var allowDump = false

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

func HandleHTTP(co net.Conn) {
	Logger.Debug("start shuttle.IConn wrap net.Con")
	conn, err := NewDefaultConn(co, TCP)
	if err != nil {
		Logger.Errorf("shuttle.IConn wrap net.Conn failed: %v", err)
		return
	}
	Logger.Debugf("shuttle.IConn wrap net.Con success [ID:%d]", conn.GetID())
	Logger.Debugf("[ID:%d] start read http request", conn.GetID())
	req, err := prepareRequest(conn)
	if err != nil {
		Logger.Error("prepareRequest failed: ", err)
		return
	}
	sc, err := ConnectToServer(req)
	if err != nil {
		if err == ErrorReject {
			Logger.Debugf("Reject [%s]", req.Target)
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
	//Dump Decorate
	if allowDump {
		if req.Protocol == ProtocolHttps && allowMitm {
			//MITM
			lct, sct, err := Mimt(conn, sc, req)
			if err != nil {
				Logger.Error("[HTTPS] Mitm failed: ", err)
				conn.Close()
				sc.Close()
				return
			}
			conn, sc = lct, sct
		}
		if req.Protocol == ProtocolHttp || (req.Protocol == ProtocolHttps && allowMitm) {
			sc, err = DumperDecorate(sc)
			if err != nil {
				Logger.Error("DumperDecorate failed: ", err)
				conn.Close()
				sc.Close()
				return
			}
			dump.InitDump(sc.GetID())
		}
	}
	//http
	if req.Protocol == ProtocolHttp {
		_, err = sc.Write(req.Data)
		if err != nil {
			Logger.Error("send http request to ss-server failed: ", err)
			conn.Close()
			sc.Close()
			return
		}
	}

	direct := &DirectChannel{}
	direct.Transport(conn, sc)
}

func prepareRequest(conn IConn) (*Request, error) {
	br := bufio.NewReader(conn)
	hreq, err := http.ReadRequest(br)
	if err != nil {
		return nil, err
	}
	Logger.Debugf("[ID:%d] [HTTP/HTTPS] %s:%s", conn.GetID(), hreq.URL.Hostname(), hreq.URL.Port())
	req := &Request{
		Ver:    socksVer5,
		Cmd:    cmdTCP,
		Addr:   hreq.URL.Hostname(),
		Atyp:   addrTypeDomain,
		ConnID: conn.GetID(),
	}
	req.IP = net.ParseIP(req.Addr)
	if port := hreq.URL.Port(); len(port) > 0 {
		req.Port, err = strToUint16(port)
		if err != nil {
			return nil, errors.New("http port error:" + port)
		}
	}
	req.Target = hreq.URL.String()
	if strings.HasPrefix(req.Target, "//") {
		req.Target = req.Target[2:]
	}
	if hreq.URL.Scheme == HTTP {
		req.Protocol = ProtocolHttp
		if req.Port == 0 {
			req.Port = 80
		}
		buffer := bytes.NewBuffer(pool.GetBuf()[:0])
		hreq.Write(buffer)
		req.Data = buffer.Bytes()
	} else if hreq.Method == http.MethodConnect {
		req.Protocol = ProtocolHttps
		if req.Port == 0 {
			req.Port = 443
		}
	}

	return req, nil
}

func strToUint16(v string) (i uint16, err error) {
	r, err := strconv.ParseUint(v, 10, 2*8)
	if err == nil {
		i = uint16(r)
	}
	return
}
