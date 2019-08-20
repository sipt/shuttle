package dns

import (
	"fmt"
	"net"
	"testing"
)

func TestResolveDomain(t *testing.T) {
	fmt.Println(net.LookupIP("www.baidu.com"))
}
