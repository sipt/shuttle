package ssstream

import (
	"crypto/cipher"
	"crypto/aes"
)

func init() {
	registerStreamCiphers("aes-128-ctr", &aes_ctr{16, 16})
	registerStreamCiphers("aes-192-ctr", &aes_ctr{24, 16})
	registerStreamCiphers("aes-256-ctr", &aes_ctr{32, 16})
}

type aes_ctr struct {
	keyLen int
	ivLen  int
}

func (a *aes_ctr) KeyLen() int {
	return a.keyLen
}
func (a *aes_ctr) IVLen() int {
	return a.ivLen
}
func (a *aes_ctr) NewEncrypter(key, iv []byte) (cipher.Stream, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewCTR(block, iv), nil
}
func (a *aes_ctr) NewDecrypter(key, iv []byte) (cipher.Stream, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewCTR(block, iv), nil
}
