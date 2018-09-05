package ssstream

import (
	"crypto/cipher"
	"encoding/binary"
	"github.com/sipt/shuttle/pool"
	"golang.org/x/crypto/salsa20/salsa"
)

func init() {
	registerStreamCiphers("salsa20", &salsa20{32, 8})
}

type salsa20 struct {
	keyLen int
	ivLen  int
}

func (a *salsa20) KeyLen() int {
	return a.keyLen
}
func (a *salsa20) IVLen() int {
	return a.ivLen
}
func (a *salsa20) NewEncrypter(key, iv []byte) (cipher.Stream, error) {
	var c salsaStreamCipher
	copy(c.nonce[:], iv[:8])
	copy(c.key[:], key[:32])
	return &c, nil
}
func (a *salsa20) NewDecrypter(key, iv []byte) (cipher.Stream, error) {
	var c salsaStreamCipher
	copy(c.nonce[:], iv[:8])
	copy(c.key[:], key[:32])
	return &c, nil
}

type salsaStreamCipher struct {
	nonce   [8]byte
	key     [32]byte
	counter int
}

func (c *salsaStreamCipher) XORKeyStream(dst, src []byte) {
	var buf []byte
	padLen := c.counter % 64
	dataSize := len(src) + padLen
	if cap(dst) >= dataSize {
		buf = dst[:dataSize]
	} else if pool.BufferSize >= dataSize {
		buf = pool.GetBuf()
		defer pool.PutBuf(buf)
		buf = buf[:dataSize]
	} else {
		buf = make([]byte, dataSize)
	}

	var subNonce [16]byte
	copy(subNonce[:], c.nonce[:])
	binary.LittleEndian.PutUint64(subNonce[len(c.nonce):], uint64(c.counter/64))

	// It's difficult to avoid data copy here. src or dst maybe slice from
	// Conn.Read/Write, which can't have padding.
	copy(buf[padLen:], src[:])
	salsa.XORKeyStream(buf, buf, &subNonce, &c.key)
	copy(dst, buf[padLen:])

	c.counter += len(src)
}
