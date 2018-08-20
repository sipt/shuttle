package selector

import (
	"testing"
	"github.com/sipt/shuttle"
	_ "github.com/sipt/shuttle/ciphers"
	"fmt"
	"net"
)

func TestRttSelector_Current(t *testing.T) {
	shuttle.InitDNS(nil, nil)
	shuttle.SetLeve("trace")
	var s = &shuttle.Server{
		Name:     "hk.b.cloudss.win",
		Host:     "42.200.227.97",
		Port:     "13819",
		Method:   "rc4-md5",
		Password: "07071818w",
	}

	conn, err := net.Dial("tcp", net.JoinHostPort(s.Host, s.Port))
	if err != nil {
		panic(err)
	}
	_, err = conn.Write([]byte{220, 149, 11, 49, 85, 116, 217, 251, 100, 252, 3, 210, 18, 155, 63, 126, 88, 104, 247, 122, 46, 102, 128, 168, 241, 184, 251, 28, 41, 165, 182, 25, 175, 148, 74, 168, 249, 38, 217, 131, 166, 2, 159, 7, 91, 134, 104, 210, 255, 98, 82, 197, 214, 217, 249, 248, 190, 187, 220, 66, 197, 169, 221, 207, 200, 222, 162, 108, 5, 58, 84, 249, 151, 175, 245, 236, 137, 115, 133, 114, 254, 127, 127, 55, 119, 232, 9, 29, 121, 28, 200, 215, 195, 16, 111, 145, 15, 59, 138, 65, 173, 17, 173, 218, 35, 159, 228, 21, 231, 5, 42, 97, 167, 175, 12, 56, 124, 180, 151, 178, 72, 206, 195, 218, 224, 84, 91, 62, 120, 154, 22, 213, 185, 105, 2, 206, 113, 128, 146, 59, 172, 176, 176, 30, 187, 10, 3, 202, 117, 103, 60, 50, 45, 112, 205, 78, 39, 2, 225})
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf[:n])
}
