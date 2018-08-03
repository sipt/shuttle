package shuttle

import (
	"crypto/md5"
	"github.com/sipt/shuttle/pool"
	"crypto/cipher"
	"io"
	"crypto/rand"
	"fmt"
)

var ciphers = make(map[string]ICipher)

func RegisterCipher(method string, cipher ICipher) error {
	Logger.Debugf("[CONF] register cipher [%s]", method)
	ciphers[method] = cipher
	return nil
}

func CheckCipher(method string) bool {
	_, ok := ciphers[method]
	return ok
}

type ICipher interface {
	KeyLen() int
	IVLen() int
	NewEncrypter(key, iv []byte) (cipher.Stream, error)
	NewDecrypter(key, iv []byte) (cipher.Stream, error)
}

//加密装饰
func CipherDecorate(password, method string, conn IConn) (IConn, error) {
	cipher, ok := ciphers[method]
	if !ok || cipher == nil {
		return nil, fmt.Errorf("[Cipher] not support [%s]", method)
	}
	cipherConn := &cipherConn{
		IConn:  conn,
		key:    evpBytesToKey(password, cipher.KeyLen()),
		iv:     make([]byte, cipher.IVLen()),
		cipher: cipher,
	}
	if _, err := io.ReadFull(rand.Reader, cipherConn.iv); err != nil {
		return nil, err
	}
	//cipherConn.iv = []byte{203, 63, 174, 189, 49, 139, 9, 54, 159, 146, 88, 224, 30, 9, 12, 238}
	fmt.Println("iv -> ", cipherConn.iv)
	fmt.Println("key -> ", cipherConn.key)
	var err error
	cipherConn.Encrypter, err = cipher.NewEncrypter(cipherConn.key, cipherConn.iv)
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(cipherConn.iv)
	return cipherConn, err
}

func evpBytesToKey(password string, keyLen int) (key []byte) {
	const md5Len = 16

	cnt := (keyLen-1)/md5Len + 1
	m := make([]byte, cnt*md5Len)
	copy(m, MD5([]byte(password)))
	d := make([]byte, md5Len+len(password))
	start := 0
	for i := 1; i < cnt; i++ {
		start += md5Len
		copy(d, m[start-md5Len:start])
		copy(d[md5Len:], password)
		copy(m[start:], MD5(d))
	}
	return m[:keyLen]
}

func MD5(data []byte) []byte {
	hash := md5.New()
	hash.Write(data)
	return hash.Sum(nil)
}

type cipherConn struct {
	IConn
	key       []byte
	iv        []byte
	Encrypter cipher.Stream
	Decrypter cipher.Stream
	cipher    ICipher
}

func (c *cipherConn) Read(b []byte) (n int, err error) {
	switch c.GetNetwork() {
	case TCP:
		n, err = c.readTCP(b)
	case UDP:
		n, err = c.readUDP(b)
	}
	return
}

func (c *cipherConn) Write(b []byte) (n int, err error) {
	fmt.Println("berfore cipher: ", b)
	buf := pool.GetBuf()
	if len(buf) < len(b) {
		pool.PutBuf(buf)
		buf = make([]byte, len(b))
	} else {
		buf = buf[:len(b)]
		defer pool.PutBuf(buf)
	}
	c.Encrypter.XORKeyStream(buf, b)
	return c.IConn.Write(buf)
}

func (c *cipherConn) readTCP(b []byte) (n int, err error) {
	if c.Decrypter == nil {
		iv := make([]byte, c.cipher.IVLen())
		if _, err = c.IConn.Read(iv); err != nil {
			return
		}
		c.Decrypter, err = c.cipher.NewEncrypter(c.key, iv)
		if err != nil {
			return
		}
	}
	buf := pool.GetBuf()
	if len(buf) < len(b) {
		pool.PutBuf(buf)
		buf = make([]byte, len(b))
	}
	defer pool.PutBuf(buf)
	n, err = c.IConn.Read(buf)
	if err != nil {
		return
	}
	c.Decrypter.XORKeyStream(b[:n], buf[:n])
	return
}

func (c *cipherConn) readUDP(b []byte) (n int, err error) {
	buf := pool.GetBuf()
	if len(buf) < len(b)+c.cipher.IVLen() {
		pool.PutBuf(buf)
		buf = make([]byte, len(b)+c.cipher.IVLen())
	} else {
		buf = buf[:len(b)+c.cipher.IVLen()]
		defer pool.PutBuf(buf)
	}
	if n, err = c.IConn.Read(buf); err != nil {
		return
	}
	iv := buf[:c.cipher.IVLen()]
	buf = buf[c.cipher.IVLen():]
	n -= c.cipher.IVLen()

	c.Decrypter, err = c.cipher.NewDecrypter(c.key, iv)
	if err != nil {
		return
	}
	c.Decrypter.XORKeyStream(b[:n], buf[:n])
	fmt.Println("udp read: ", b[:n])
	return
}
