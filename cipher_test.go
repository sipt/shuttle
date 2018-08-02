package shuttle

import (
	"testing"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rc4"
	"io"
	"crypto/rand"
	"fmt"
)

func TestCheckCipher(t *testing.T) {
	cipher := ciphers["rc4-md5"]
	key := evpBytesToKey("07071818w", cipher.KeyLen())
	iv := make([]byte, cipher.IVLen())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	//cipherConn.iv = []byte{203, 63, 174, 189, 49, 139, 9, 54, 159, 146, 88, 224, 30, 9, 12, 238}
	fmt.Println("iv -> ", iv)
	fmt.Println("key -> ", key)
	var err error
	enc, err := cipher.NewEncrypter(key, iv)
	if err != nil {
		panic(err)
	}

	dec, _ := cipher.NewDecrypter(key, iv)
	src := []byte{3, 14, 119, 119, 119, 46, 103, 111, 111, 103, 108, 101, 46, 99, 111, 109, 80, 0}
	dst := make([]byte, len(src))
	enc.XORKeyStream(dst, src)
	fmt.Println(dst)

	dst2 := make([]byte, len(src))
	dec.XORKeyStream(dst2, dst)
	fmt.Println(dst2)

}

func init() {
	RegisterCipher("rc4-md5", &rc4_md5{16, 16})
}

type rc4_md5 struct {
	keyLen int
	ivLen  int
}

func (a *rc4_md5) KeyLen() int {
	return a.keyLen
}
func (a *rc4_md5) IVLen() int {
	return a.ivLen
}
func (a *rc4_md5) NewEncrypter(key, iv []byte) (cipher.Stream, error) {
	h := md5.New()
	h.Write(key)
	h.Write(iv)
	rc4key := h.Sum(nil)

	return rc4.NewCipher(rc4key)
}
func (a *rc4_md5) NewDecrypter(key, iv []byte) (cipher.Stream, error) {
	h := md5.New()
	h.Write(key)
	h.Write(iv)
	rc4key := h.Sum(nil)

	return rc4.NewCipher(rc4key)
}
