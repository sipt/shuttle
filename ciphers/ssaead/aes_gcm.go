package ssaead

import (
	"crypto/cipher"
	"crypto/aes"
)

func init() {
	registerAEADCiphers("aes-256-gcm", &aesGcm{32, 32, 12, 16})
	registerAEADCiphers("aes-192-gcm", &aesGcm{24, 24, 12, 16})
	registerAEADCiphers("aes-128-gcm", &aesGcm{16, 16, 12, 16})
}

type aesGcm struct {
	keySize   int
	saltSize  int
	nonceSize int
	tagSize   int
}

func (a *aesGcm) KeySize() int {
	return a.keySize
}

func (a *aesGcm) SaltSize() int {
	return a.saltSize
}

func (a *aesGcm) NonceSize() int {
	return a.nonceSize
}

func (a *aesGcm) NewEncrypter(key []byte, salt []byte) (cipher.AEAD, error) {
	subkey := make([]byte, a.KeySize())
	HKDF_SHA1(key, salt, []byte("ss-subkey"), subkey)
	blk, err := aes.NewCipher(subkey)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(blk)
}

func (a *aesGcm) NewDecrypter(key []byte, salt []byte) (cipher.AEAD, error) {
	subkey := make([]byte, a.KeySize())
	HKDF_SHA1(key, salt, []byte("ss-subkey"), subkey)
	blk, err := aes.NewCipher(subkey)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(blk)
}
