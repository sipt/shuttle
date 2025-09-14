package tun

import (
	"fmt"
	"net"

	"gvisor.dev/gvisor/pkg/tcpip"
	"gvisor.dev/gvisor/pkg/tcpip/header"
	"gvisor.dev/gvisor/pkg/tcpip/network/ipv4"
	"gvisor.dev/gvisor/pkg/tcpip/network/ipv6"
	"gvisor.dev/gvisor/pkg/tcpip/stack"
	"gvisor.dev/gvisor/pkg/tcpip/transport/icmp"
	"gvisor.dev/gvisor/pkg/tcpip/transport/tcp"
	"gvisor.dev/gvisor/pkg/tcpip/transport/udp"
)

func CreateStack(endpoint stack.LinkEndpoint, ipCidrs []net.IPNet) (*stack.Stack, error) {
	s := stack.New(stack.Options{
		NetworkProtocols: []stack.NetworkProtocolFactory{
			ipv4.NewProtocol,
			ipv6.NewProtocol,
		},
		TransportProtocols: []stack.TransportProtocolFactory{
			tcp.NewProtocol,
			udp.NewProtocol,
			icmp.NewProtocol4,
			icmp.NewProtocol6,
		},
	})

	if err := s.SetForwardingDefaultAndAllNICs(ipv4.ProtocolNumber, true); err != nil {
		return nil, fmt.Errorf("set ipv4 forwarding: %s", err)
	}
	if err := s.SetForwardingDefaultAndAllNICs(ipv6.ProtocolNumber, true); err != nil {
		return nil, fmt.Errorf("set ipv6 forwarding: %s", err)
	}

	nicID := s.NextNICID()
	e := s.CreateNIC(nicID, endpoint)
	if e != nil {
		return nil, fmt.Errorf("create nic: %s", e)
	}

	var routeTable []tcpip.Route
	for _, cidr := range ipCidrs {
		var addr tcpip.ProtocolAddress
		switch len(cidr.IP) {
		case net.IPv4len:
			addr.Protocol = ipv4.ProtocolNumber
			prefixLen, err := maskToPrefixLen(cidr.Mask)
			if err != nil {
				return nil, err
			}
			var ip4 [net.IPv4len]byte
			copy(ip4[:], cidr.IP)
			addr.AddressWithPrefix = tcpip.AddressWithPrefix{
				Address:   tcpip.AddrFrom4(ip4),
				PrefixLen: prefixLen,
			}
			routeTable = append(routeTable, tcpip.Route{
				Destination: header.IPv4EmptySubnet,
				NIC:         nicID,
			})
		case net.IPv6len:
			addr.Protocol = ipv6.ProtocolNumber
			prefixLen, err := maskToPrefixLen(cidr.Mask)
			var ip6 [net.IPv6len]byte
			copy(ip6[:], cidr.IP)
			if err != nil {
				return nil, err
			}
			addr.AddressWithPrefix = tcpip.AddressWithPrefix{
				Address:   tcpip.AddrFrom16(ip6),
				PrefixLen: prefixLen,
			}
			routeTable = append(routeTable, tcpip.Route{
				Destination: header.IPv6EmptySubnet,
				NIC:         nicID,
			})
		default:
			return nil, fmt.Errorf("invalid ip: %s", cidr.IP)
		}

		s.AddProtocolAddress(nicID, addr, stack.AddressProperties{PEB: stack.CanBePrimaryEndpoint})
	}
	s.SetRouteTable(routeTable)
	e = s.SetSpoofing(nicID, true)
	if e != nil {
		return nil, fmt.Errorf("set spoofing: %s", e)
	}
	e = s.SetPromiscuousMode(nicID, true)
	if e != nil {
		return nil, fmt.Errorf("set promiscuous mode: %s", e)
	}

	return s, nil
}

func maskToPrefixLen(mask net.IPMask) (int, error) {
	var prefixLen int
	for _, b := range mask {
		for b > 0 {
			if b&1 == 1 {
				prefixLen++
			}
			b >>= 1
		}
	}
	if prefixLen == 0 {
		return 0, fmt.Errorf("invalid mask: %s", mask)
	}
	return prefixLen, nil
}
