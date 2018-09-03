package ssstream

import (
	"crypto/cipher"
	"github.com/Yawning/chacha20"
)

func init() {
	registerStreamCiphers("chacha20", &_chacha20{32, 8})
	registerStreamCiphers("chacha20-ietf", &_chacha20{32, 12})
}

type _chacha20 struct {
	keyLen int
	ivLen  int
}

func (a *_chacha20) KeyLen() int {
	return a.keyLen
}
func (a *_chacha20) IVLen() int {
	return a.ivLen
}
func (a *_chacha20) NewEncrypter(key, iv []byte) (cipher.Stream, error) {
	return chacha20.NewCipher(key, iv)
}
func (a *_chacha20) NewDecrypter(key, iv []byte) (cipher.Stream, error) {
	return chacha20.NewCipher(key, iv)
}
