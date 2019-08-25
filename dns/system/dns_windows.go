//+build windows

package dns

import (
	"encoding/json"
	"fmt"
	"unsafe"

	"github.com/kbinani/win"
)

func GetNetworkParams() []string {
	fixedInfo := win.FIXED_INFO_W2KSP1{}
	bufLen := uint32(unsafe.Sizeof(fixedInfo))
	reply := win.GetNetworkParams(&fixedInfo, &bufLen)
	fmt.Println("call syscall3: ", reply)
	data, err := json.Marshal(fixedInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	return nil
}
