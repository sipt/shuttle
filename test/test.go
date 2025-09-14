package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/sipt/shuttle/pkg/dns"
	"github.com/sipt/shuttle/pkg/tun"
)

func main() {
	// 1. 打开TUN接口
	listener, err := tun.OpenTun()
	if err != nil {
		log.Fatalf("Failed to open TUN: %v", err)
	}
	fmt.Println("TUN opened successfully")

	// 2. 创建固定IP的DNS处理器，所有DNS查询都返回 10.0.0.12
	dnsHandler := dns.NewFixedIPDNSHandler("10.0.0.12")
	dnsServer := dns.NewDNSServer(dnsHandler)

	fmt.Println("DNS server created with fixed IP: 10.0.0.12")

	// 3. 启动DNS服务器处理循环
	go func() {
		for {
			// 接受UDP连接
			conn, err := listener.UdpListener.Accept()
			if err != nil {
				log.Printf("Failed to accept UDP connection: %v", err)
				continue
			}

			fmt.Printf("New UDP connection: %s -> %s\n",
				conn.RemoteAddr(), conn.LocalAddr())

			// 检查是否是DNS查询（端口53）
			isDNS := false
			if localAddr := conn.LocalAddr(); localAddr != nil {
				if udpAddr, ok := localAddr.(*net.UDPAddr); ok {
					isDNS = udpAddr.Port == 53
				}
			}

			if isDNS {
				fmt.Printf("DNS query detected from %s\n", conn.RemoteAddr())

				// 处理DNS连接
				go handleDNSConnection(dnsServer, conn)
			} else {
				// 非DNS连接，简单回复并关闭
				fmt.Printf("Non-DNS UDP connection from %s\n", conn.RemoteAddr())
				conn.Write([]byte("Hello from TUN interface!"))
				conn.Close()
			}
		}
	}()

	fmt.Println("DNS server is running...")
	fmt.Println("All DNS queries will be resolved to 10.0.0.12")
	fmt.Println("Press Ctrl+C to stop")

	// 保持程序运行
	select {}
}

// handleDNSConnection 处理DNS连接
func handleDNSConnection(dnsServer *dns.DNSServer, conn tun.UdpConn) {
	defer conn.Close()

	// 创建包装器以适配net.PacketConn接口
	packetConn := &dnsPacketConn{
		conn:       conn,
		localAddr:  conn.LocalAddr(),
		remoteAddr: conn.RemoteAddr(),
	}

	// 设置超时
	timeout := time.NewTimer(30 * time.Second)
	defer timeout.Stop()

	// 在单独的goroutine中处理DNS请求
	done := make(chan error, 1)
	go func() {
		done <- dnsServer.ServeDNSOverUDP(packetConn)
	}()

	// 等待完成或超时
	select {
	case err := <-done:
		if err != nil {
			log.Printf("DNS server error: %v", err)
		}
	case <-timeout.C:
		log.Printf("DNS connection timeout for %s", conn.RemoteAddr())
	}
}

// dnsPacketConn 实现net.PacketConn接口，用于适配DNS服务器
type dnsPacketConn struct {
	conn       tun.UdpConn
	localAddr  net.Addr
	remoteAddr net.Addr
}

func (c *dnsPacketConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	n, err = c.conn.Read(p)
	if err != nil {
		return 0, nil, err
	}
	return n, c.remoteAddr, nil
}

func (c *dnsPacketConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	return c.conn.Write(p)
}

func (c *dnsPacketConn) Close() error {
	return c.conn.Close()
}

func (c *dnsPacketConn) LocalAddr() net.Addr {
	return c.localAddr
}

func (c *dnsPacketConn) SetDeadline(t time.Time) error {
	// TUN连接不支持deadline，返回nil
	return nil
}

func (c *dnsPacketConn) SetReadDeadline(t time.Time) error {
	// TUN连接不支持deadline，返回nil
	return nil
}

func (c *dnsPacketConn) SetWriteDeadline(t time.Time) error {
	// TUN连接不支持deadline，返回nil
	return nil
}
