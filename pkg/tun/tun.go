package tun

import (
	"context"
	"fmt"
	"net"
	"net/netip"

	"github.com/sirupsen/logrus"
	"github.com/slackhq/nebula/config"
	"github.com/slackhq/nebula/overlay"
	"gvisor.dev/gvisor/pkg/buffer"
	"gvisor.dev/gvisor/pkg/tcpip/link/channel"
	"gvisor.dev/gvisor/pkg/tcpip/network/ipv4"
	"gvisor.dev/gvisor/pkg/tcpip/stack"
)

const NetCidr = "198.18.0.1/8"

type TunDevice struct {
	Device   overlay.Device
	Listener *Listener

	s      *stack.Stack
	linkEP *channel.Endpoint
}

type Listener struct {
	UdpListener UdpListener
	TcpListener TcpListener
}

func OpenTun(ctx context.Context) (*TunDevice, error) {
	// 创建 logger
	logger := logrus.WithField("method", "open-tun").Logger

	// 创建配置
	c := config.NewC(logger)

	// 设置 TUN 设备配置
	c.Settings = map[string]any{
		"tun": map[string]any{
			"mtu": 4000,
		},
	}

	tunCidr, err := netip.ParsePrefix(NetCidr)
	if err != nil {
		logger.WithError(err).Error("Failed to parse TUN CIDR")
		return nil, fmt.Errorf("Failed to parse TUN CIDR: %v", err)
	}

	// 创建 TUN 设备
	logger.Infof("Creating TUN device with CIDR: %s", tunCidr)
	tun, err := overlay.NewDeviceFromConfig(c, logger, []netip.Prefix{tunCidr}, 1)
	if err != nil {
		logger.WithError(err).Error("Failed to create TUN device")
		return nil, fmt.Errorf("Failed to create TUN device: %v", err)
	}
	logger.Info("Activating TUN device...")
	err = tun.Activate()
	if err != nil {
		logger.Errorf("Failed to activate TUN device: %v", err)
		return nil, fmt.Errorf("Failed to activate TUN device: %v", err)
	}

	linkEP := channel.New(1024, 4000, "")
	// 创建栈
	_, ipNet, err := net.ParseCIDR(NetCidr)
	if err != nil {
		logger.WithError(err).Error("Failed to parse TUN CIDR")
		return nil, fmt.Errorf("Failed to parse TUN CIDR: %v", err)
	}
	s, err := CreateStack(linkEP, []net.IPNet{*ipNet})
	if err != nil {
		logger.WithError(err).Error("Failed to create stack")
		return nil, fmt.Errorf("Failed to create stack: %v", err)
	}

	// 3. goroutine: TUN → netstack
	go func() {
		buf := make([]byte, 2000)
		for {
			n, err := tun.Read(buf)
			if err != nil {
				logger.Errorf("tun read err: %v", err)
				break
			}
			vv := stack.NewPacketBuffer(stack.PacketBufferOptions{
				Payload: buffer.MakeWithData(buf[:n]),
			})
			linkEP.InjectInbound(ipv4.ProtocolNumber, vv)
		}
	}()

	// // 4. goroutine: netstack → TUN
	go func() {
		for {
			pkt := linkEP.ReadContext(ctx)
			if pkt == nil {
				logger.Info("linkEP read nil, exit")
				return
			}
			vv := pkt.ToView()
			data := vv.AsSlice()
			_, err := tun.Write(data)
			if err != nil {
				logger.Errorf("tun write err: %v", err)
			}
			pkt.DecRef()
		}
	}()

	var listener = new(Listener)

	// 5. 启动 UDP 代理 - 拦截所有UDP连接
	udpListener, err := UdpForward(ctx, s, 1024)
	if err != nil {
		logger.Errorf("udp forwarder error: %v", err)
		return nil, fmt.Errorf("udp forwarder error: %v", err)
	}
	logger.Info("UDP proxy started - intercepting all UDP connections")
	listener.UdpListener = udpListener

	// 6. 启动 TCP 代理 - 拦截所有TCP连接
	tcpListener, err := TcpForward(ctx, s, 1024)
	if err != nil {
		logger.Errorf("tcp forwarder error: %v", err)
		return nil, fmt.Errorf("tcp forwarder error: %v", err)
	}
	logger.Info("TCP proxy started - intercepting all TCP connections")
	listener.TcpListener = tcpListener

	return &TunDevice{
		Device:   tun,
		Listener: listener,
		s:        s,
		linkEP:   linkEP,
	}, nil
}

func (t *TunDevice) Close() error {
	t.Device.Close()
	t.linkEP.Close()
	t.s.Close()
	return nil
}
