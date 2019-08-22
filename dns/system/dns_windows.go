//+build windows

package dns

import (
	"fmt"
	"github.com/kbinani/win"
	"unsafe"
)

func GetNetworkParams() []string {
	fixedInfo := win.FIXED_INFO_W2KSP1{}
	bufLen := uint32(unsafe.Sizeof(fixedInfo))
	reply := win.GetNetworkParams(&fixedInfo, &bufLen)
	fmt.Println("call syscall3: ", reply)
	server := &fixedInfo.DnsServerList
	for server != nil {
		fmt.Println(server.IpAddress.String, server.IpAddress.String)
		server = server.Next
	}
	return nil
}
