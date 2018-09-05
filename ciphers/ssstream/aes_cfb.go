package ssstream

import (
	"crypto/cipher"
	"crypto/aes"
)

func init() {
	registerStreamCiphers("aes-128-cfb", &aes_cfb{16, 16})
	registerStreamCiphers("aes-192-cfb", &aes_cfb{24, 16})
	registerStreamCiphers("aes-256-cfb", &aes_cfb{32, 16})
}

type aes_cfb struct {
	keyLen int
	ivLen  int
}

func (a *aes_cfb) KeyLen() int {
	return a.keyLen
}
func (a *aes_cfb) IVLen() int {
	return a.ivLen
}
func (a *aes_cfb) NewEncrypter(key, iv []byte) (cipher.Stream, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewCFBEncrypter(block, iv), nil
}
func (a *aes_cfb) NewDecrypter(key, iv []byte) (cipher.Stream, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewCFBDecrypter(block, iv), nil
}
