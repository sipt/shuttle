package shuttle

import (
	"encoding/binary"
	"errors"
	"github.com/sipt/shuttle/log"
	"github.com/sipt/shuttle/pool"
	"github.com/sipt/shuttle/util"
	"net"
	"strconv"
	"time"
)

const (
	socksVer5 = 0x05

	verIndex     = 0
	nMethodIndex = 1
	methodIndex  = 2

	cmdIndex       = 1
	rsvIndex       = cmdIndex + 1
	atypIndex      = rsvIndex + 1
	addrIndex      = atypIndex + 1
	domianLenIndex = atypIndex + 1
	ipv4PortIndex  = addrIndex + 4
	ipv6PortIndex  = addrIndex + 16
)

func SocksHandle(co net.Conn) {
	log.Logger.Debug("[SOCKS] start shuttle.IConn wrap net.Conn")
	conn, err := NewDefaultConn(co, TCP)
	if err != nil {
		log.Logger.Errorf("shuttle.IConn wrap net.Conn failed: %v", err)
		return
	}
	log.Logger.Debugf("[SOCKS] [ID:%d] shuttle.IConn wrap net.Conn success ", conn.GetID())
	log.Logger.Debugf("[SOCKS] [ID:%d] start handShake", conn.GetID())
	err = handShake(conn)
	if err != nil {
		log.Logger.Errorf("[SOCKS] [ID:%d] handShake failed: %s", conn.GetID(), err.Error())
		return
	}
	req, err := parseRequest(conn)
	if err != nil {
		log.Logger.Errorf("[SOCKS] [ID:%d] parseRequest failed: %s", conn.GetID(), err.Error())
		return
	}
	req.Protocol = ProtocolSocks
	req.Target = req.Host()
	_, err = conn.Write([]byte{socksVer5, 0x00, 0x00, AddrTypeIPv4, 0x00, 0x00, 0x00, 0x00, 0x08, 0x43})
	if err != nil {
		log.Logger.Errorf("[SOCKS] [ID:%d] send connection confirmation: %s", conn.GetID(), err.Error())
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

	//RuleFilter by Rules and DNS
	rule, s, err := FilterByReq(req)
	if err != nil {
		log.Logger.Errorf("[SOCKS] [ID:%d] ConnectToServer failed [%s] err: %s", conn.GetID(), req.Host(), err)
	}

	//connnet to server
	log.Logger.Debugf("[SOCKS] [ID:%d] Start connect to Server [%s]", conn.GetID(), s.Name)
	sc, err := s.Conn(req)
	if err != nil {
		if err == ErrorReject {
			log.Logger.Debugf("[SOCKS] [ID:%d] Reject [%s]", conn.GetID(), req.Target)
		} else {
			log.Logger.Error("[SOCKS] [ID:%d] ConnectToServer failed [%s] err: %s", conn.GetID(), req.Host(), err.Error())
		}
		return
	}
	log.Logger.Debugf("[SOCKS] [ID:%d] Server [%s] Connected success", conn.GetID(), s.Name)
	log.Logger.Debugf("[HTTP] [ClientConnID:%d] Bind to [ServerConnID:%d]", conn.GetID(), sc.GetID())
	// 白名单判断
	if IsPass(req.Addr, req.IP.String(), strconv.Itoa(int(req.Port))) {
		direct := &DirectChannel{}
		direct.Transport(conn, sc)
		return
	}

	id := util.NextID()
	sc.SetRecordID(id)
	boxChan <- &Box{Op: RecordAppend, Value: &Record{
		ID:       id,
		Protocol: req.Protocol,
		Created:  time.Now(),
		Proxy:    s,
		Status:   RecordStatusActive,
		URL:      req.Target,
		Rule:     rule,
	}, ID: id}
	direct := &DirectChannel{}
	direct.Transport(conn, sc)
	boxChan <- &Box{id, RecordStatus, RecordStatusCompleted}
}

//socks 握手
func handShake(conn net.Conn) error {
	buf := pool.GetBuf()
	_, err := conn.Read(buf)
	if err != nil {
		return err
	}
	if buf[verIndex] != socksVer5 {
		return errors.New("socks version not supported")
	}
	//todo get methods

	//return supported methods
	conn.Write([]byte{0x05, 0x00})
	pool.PutBuf(buf)
	return nil
}

//获取协议
func parseRequest(conn IConn) (*Request, error) {
	//+----+-----+-------+------+----------+----------+
	//|VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	//+----+-----+-------+------+----------+----------+
	//| 1  |  1  |   1   |  1   | Variable |    2     |
	//+----+-----+-------+------+----------+----------+
	//CMD 字段：command 的缩写，shadowsocks 只用到了：
	//    0x01：建立 TCP 连接
	//    0x03：关联 UDP 请求
	//RSV 字段：保留字段，值为 0x00；
	//ATYP 字段：address type 的缩写，取值为：
	//    0x01：IPv4
	//    0x03：域名
	//    0x04：IPv6
	//DST.ADDR 字段：destination address 的缩写，取值随 ATYP 变化：
	//    ATYP == 0x01：4 个字节的 IPv4 地址
	//    ATYP == 0x03：1 个字节表示域名长度，紧随其后的是对应的域名
	//    ATYP == 0x04：16 个字节的 IPv6 地址
	//DST.PORT 字段：目的服务器的端口。
	buf := pool.GetBuf()
	_, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	request := &Request{
		Ver:    uint8(buf[verIndex]),
		Cmd:    uint8(buf[cmdIndex]),
		Rsv:    uint8(buf[rsvIndex]),
		Atyp:   uint8(buf[atypIndex]),
		ConnID: conn.GetID(),
	}
	switch request.Atyp {
	case AddrTypeIPv4:
		request.IP = buf[atypIndex+1 : ipv4PortIndex]
		request.Port = binary.BigEndian.Uint16(buf[ipv4PortIndex : ipv4PortIndex+2])
		if request.Cmd == CmdUDP {
			request.Data = buf[ipv4PortIndex+2:]
		}
	case AddrTypeDomain:
		end := buf[domianLenIndex] + 1 + domianLenIndex
		request.Addr = string(buf[domianLenIndex+1 : end])
		request.Port = binary.BigEndian.Uint16(buf[end : end+2])
		if request.Cmd == CmdUDP {
			request.Data = buf[end+2:]
		}
	case AddrTypeIPv6:
		request.IP = buf[atypIndex+1 : ipv6PortIndex]
		request.Port = binary.BigEndian.Uint16(buf[ipv6PortIndex : ipv6PortIndex+2])
		if request.Cmd == CmdUDP {
			request.Data = buf[ipv6PortIndex+2:]
		}
	}
	if request.Cmd != CmdUDP {
		pool.PutBuf(buf) // 回收
	}
	return request, nil
}
