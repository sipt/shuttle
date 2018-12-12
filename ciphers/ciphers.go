package ciphers

import (
	"fmt"
	"github.com/sipt/shuttle/ciphers/ssaead"
	"github.com/sipt/shuttle/ciphers/ssstream"
	connect "github.com/sipt/shuttle/conn"
)

type ConnDecorate func(password string, conn connect.IConn) (connect.IConn, error)

//加密装饰
func CipherDecorate(password, method string, conn connect.IConn) (connect.IConn, error) {
	d := ssstream.GetStreamCiphers(method)
	if d != nil {
		return d(password, conn)
	}
	d = ssaead.GetAEADCiphers(method)
	if d != nil {
		return d(password, conn)
	}
	return nil, fmt.Errorf("[SS Cipher] not support : %s", method)
}
