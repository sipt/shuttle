package main

import (
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	"github.com/pkg/taptun"
)

func main() {
	//c, err := net.Dial("ip:udp", "8.8.8.8")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(c == nil)
	tun()
}

func UnPacketUDP(data []byte) (b []byte, srcIP, dstIP net.IP, srcPort, dstPort layers.UDPPort) {
	var ip4 layers.IPv4
	var udp layers.UDP
	var dns layers.DNS
	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeIPv4, &ip4, &udp, &dns)
	decoded := []gopacket.LayerType{}
	err := parser.DecodeLayers(data, &decoded)
	if err != nil {
		fmt.Println("DecodingLayerParser error:", err)
	}
	fmt.Printf("src[%s:%d] -> dst[%s:%d]\n", ip4.SrcIP, udp.SrcPort, ip4.DstIP, udp.DstPort)
	srcIP = net.ParseIP(ip4.SrcIP.String())
	dstIP = net.ParseIP(ip4.DstIP.String())
	b = udp.Payload
	srcPort, dstPort = udp.SrcPort, udp.DstPort
	return
}

func PacketUDP(b []byte, srcIP, dstIP net.IP, srcPort, dstPort layers.UDPPort) []byte {
	sbuf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
	ip := &layers.IPv4{
		Version:  4,
		DstIP:    srcIP,
		SrcIP:    dstIP,
		TTL:      64,
		Protocol: layers.IPProtocolUDP,
	}
	udp := &layers.UDP{
		SrcPort: dstPort,
		DstPort: srcPort,
	}
	err := udp.SetNetworkLayerForChecksum(ip)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	err = gopacket.SerializeLayers(sbuf, opts, ip, udp, gopacket.Payload(b))
	if err != nil {
		fmt.Println("SerializeLayers error:", err)
	}
	writeBuf := make([]byte, 2048)
	writeBuf[3] = 2
	n := copy(writeBuf[4:], sbuf.Bytes())
	return writeBuf[:n+4]
}

func tun() {
	tun, err := taptun.OpenTun()
	if err != nil {
		panic(err)
	}
	fmt.Println("create tun:", tun.String())
	c, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		panic(err)
	}
	fmt.Println(c == nil)
	for {
		b := make([]byte, 1500)
		n, err := tun.Read(b)
		if err != nil {
			panic(err)
		}
		data, srcIP, dstIP, srcPort, dstPort := UnPacketUDP(b[4:n])
		_, err = c.Write(data)
		if err != nil {
			panic(err)
		}
		fmt.Println(srcIP.String(), dstIP.String())
		n, err = c.Read(b)
		if err != nil {
			panic(err)
		}
		data = PacketUDP(b[:n], srcIP, dstIP, srcPort, dstPort)
		_, err = tun.Write(data)
		if err != nil {
			panic(err)
		}
	}
}

//
//type ITest interface {
//	Print()
//	CallPrint()
//}
//
//type Test1 struct {
//	*Test
//}
//
//func (t *Test1) Print() {
//	fmt.Println("test1")
//}
//
//type Test struct {
//	ITest
//}
//
//func (t *Test) CallPrint() {
//	t.Print()
//}
//func f(err error) {
//	if err != nil {
//		panic(err)
//	}
//}
//
//func sendSocks5() {
//	c, err := net.Dial("tcp", "127.0.0.1:9000")
//	f(err)
//	_, err = c.Write([]byte{5, 1, 0})
//	f(err)
//	b := make([]byte, 1500)
//	n, err := c.Read(b)
//	fmt.Println(b[:n])
//	_, err = c.Write([]byte{5, 3, 0, 1, 8, 8, 8, 8, 0, 53})
//	f(err)
//	n, err = c.Read(b)
//	f(err)
//	fmt.Println(b[:n])
//	addr := test(b[:n])
//	c, err = net.Dial("udp", addr)
//	f(err)
//	fmt.Println(c.LocalAddr().String())
//	_, err = c.Write([]byte{0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 212, 155, 3, 147, 1, 32, 0, 1, 0, 0, 0, 0, 0, 1, 3, 119, 119, 119, 5, 98, 97, 105, 100, 117, 3, 99, 111, 109, 0, 0, 1, 0, 1, 0, 0, 41, 16, 0, 0, 0, 0, 0, 0, 0})
//	f(err)
//	n, err = c.Read(b)
//	f(err)
//	fmt.Println(b[:n])
//}
//
//func tcp() {
//	l, err := net.Listen("tcp", ":53")
//	if err != nil {
//		panic(err)
//	}
//	for {
//		conn, _ := l.Accept()
//		go func() {
//			b := make([]byte, 2048)
//			for {
//				n, err := conn.Read(b)
//				if err != nil {
//					fmt.Println(err)
//					return
//				}
//				fmt.Println(b[:n])
//			}
//		}()
//	}
//}
//
//func udp() {
//	pc, err := net.ListenPacket("udp", ":53")
//	if err != nil {
//		panic(err)
//	}
//	for {
//		b := make([]byte, 1500)
//		n, addr, err := pc.ReadFrom(b)
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//		fmt.Println(addr.String(), b[:n])
//
//	}
//}
//
//func test(b []byte) string {
//	if b[0] != 5 {
//		panic(b[0])
//	}
//	if b[1]|b[2] != 0 {
//		panic(b[1])
//	}
//	l := 2
//	off := 4
//	dstAddr := &Addr{}
//	switch b[3] {
//	case AddrTypeIPv4:
//		l += net.IPv4len
//		dstAddr.IP = make(net.IP, net.IPv4len)
//	case AddrTypeIPv6:
//		l += net.IPv6len
//		dstAddr.IP = make(net.IP, net.IPv6len)
//	case AddrTypeFQDN:
//		l += int(b[4])
//		off = 5
//	default:
//		logrus.WithField("server", "udp").Debugf("ATYP [%x] unknown address type", b[3])
//		return ""
//	}
//	if len(b[off:]) < l {
//		logrus.WithField("server", "udp").Debugf("short cmd request")
//		return ""
//	}
//	if dstAddr.IP != nil {
//		copy(dstAddr.IP, b[off:])
//	} else {
//		dstAddr.Name = string(b[off : off+l-2])
//	}
//	dstAddr.Port = int(b[off+l-2])<<8 | int(b[off+l-1])
//	return dstAddr.String()
//}
