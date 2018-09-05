package ssstream

import (
	"crypto/cipher"
	"golang.org/x/crypto/cast5"
)

func init() {
	registerStreamCiphers("cast5-cfb", &cast5_cfb{16, 8})
}

type cast5_cfb struct {
	keyLen int
	ivLen  int
}

func (a *cast5_cfb) KeyLen() int {
	return a.keyLen
}
func (a *cast5_cfb) IVLen() int {
	return a.ivLen
}
func (a *cast5_cfb) NewEncrypter(key, iv []byte) (cipher.Stream, error) {
	block, err := cast5.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewCFBEncrypter(block, iv), nil
}
func (a *cast5_cfb) NewDecrypter(key, iv []byte) (cipher.Stream, error) {
	block, err := cast5.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewCFBDecrypter(block, iv), nil
}
