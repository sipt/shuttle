package shuttle

import (
	"net"
	"errors"
	"github.com/sipt/shuttle/pool"
	"encoding/binary"
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
	domianLenIndex = rsvIndex + 1
	ipv4PortIndex  = addrIndex + 4
	ipv6PortIndex  = addrIndex + 16
	addrTypeIPv4   = 0x01 //    0x01：IPv4
	addrTypeDomain = 0x03 //    0x03：域名
	addrTypeIPv6   = 0x04 //    0x04：IPv6
)

func SocksHandle(co net.Conn) {
	Logger.Debug("start shuttle.IConn wrap net.Con")
	conn, err := NewDefaultConn(co, TCP)
	if err != nil {
		Logger.Errorf("shuttle.IConn wrap net.Conn failed: %v", err)
		return
	}
	Logger.Debugf("shuttle.IConn wrap net.Con success [ID:%d]", conn.GetID())
	Logger.Debugf("[ID:%d] start handShake", conn.GetID())
	err = handShake(conn)
	if err != nil {
		Logger.Errorf("[%d] handShake failed: %v", conn.GetID(), err)
		return
	}
	req, err := parseRequest(conn)
	if err != nil {
		Logger.Error("parseRequest failed: ", err)
		return
	}
	req.Protocol = ProtocolSocks
	req.Target = req.Host()
	_, err = conn.Write([]byte{socksVer5, 0x00, 0x00, req.Atyp, 0x00, 0x00, 0x00, 0x00, 0x08, 0x43})
	if err != nil {
		Logger.Error("send connection confirmation:", err)
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
	direct := &DirectChannel{}
	direct.Transport(conn, sc)
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
	case addrTypeIPv4:
		request.IP = buf[atypIndex+1 : ipv4PortIndex]
		request.Port = binary.BigEndian.Uint16(buf[ipv4PortIndex : ipv4PortIndex+2])
		if request.Cmd == cmdUDP {
			request.Data = buf[ipv4PortIndex+2:]
		}
	case addrTypeDomain:
		end := buf[domianLenIndex] + 1 + domianLenIndex
		request.IP = buf[domianLenIndex+1 : end]
		request.Port = binary.BigEndian.Uint16(buf[end : end+2])
		if request.Cmd == cmdUDP {
			request.Data = buf[end+2:]
		}
	case addrTypeIPv6:
		request.IP = buf[atypIndex+1 : ipv6PortIndex]
		request.Port = binary.BigEndian.Uint16(buf[ipv6PortIndex : ipv6PortIndex+2])
		if request.Cmd == cmdUDP {
			request.Data = buf[ipv6PortIndex+2:]
		}
	}
	if request.Cmd != cmdUDP {
		pool.PutBuf(buf) // 回收
	}
	return request, nil
}
