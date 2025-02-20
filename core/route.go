package core

import (
	"github.com/xjasonlyu/tun2socks/v2/core/option"

	"gvisor.dev/gvisor/pkg/tcpip"
	"gvisor.dev/gvisor/pkg/tcpip/header"
	"gvisor.dev/gvisor/pkg/tcpip/stack"
)

func withRouteTable(enableIPv6 bool, nicID tcpip.NICID) option.Option {
	if enableIPv6 {
		return func(s *stack.Stack) error {
			s.SetRouteTable([]tcpip.Route{
				{
					Destination: header.IPv4EmptySubnet,
					NIC:         nicID,
				},
				{
					Destination: header.IPv6EmptySubnet,
					NIC:         nicID,
				},
			})
			return nil
		}
	} else {
		return func(s *stack.Stack) error {
			s.SetRouteTable([]tcpip.Route{
				{
					Destination: header.IPv4EmptySubnet,
					NIC:         nicID,
				},
			})
			return nil
		}
	}

}
