package ssaead

import (
	"crypto/cipher"
	"crypto/sha1"
	connect "github.com/sipt/shuttle/conn"
	"io"

	"bytes"
	"crypto/md5"
	"crypto/rand"
	"github.com/sipt/shuttle/log"
	"github.com/sipt/shuttle/pool"
	"golang.org/x/crypto/hkdf"
)

var aeadCiphers = make(map[string]IAEADCipher)

func registerAEADCiphers(method string, c IAEADCipher) {
	aeadCiphers[method] = c
}

func GetAEADCiphers(method string) func(string, connect.IConn) (connect.IConn, error) {
	c, ok := aeadCiphers[method]
	if !ok {
		return nil
	}
	return func(password string, conn connect.IConn) (connect.IConn, error) {
		salt := make([]byte, c.SaltSize())
		if _, err := io.ReadFull(rand.Reader, salt); err != nil {
			return nil, err
		}
		sc := &aeadConn{
			IConn:       conn,
			IAEADCipher: c,
			key:         evpBytesToKey(password, c.KeySize()),
			wNonce:      make([]byte, c.NonceSize()),
			rNonce:      make([]byte, c.NonceSize()),
			readBuffer:  bytes.NewBuffer(pool.GetBuf()[:0]),
		}
		var err error
		sc.Encrypter, err = sc.NewEncrypter(sc.key, salt)
		_, err = conn.Write(salt)
		return sc, err
	}
}

const DataMaxSize = 0x3FFF

type IAEADCipher interface {
	KeySize() int
	SaltSize() int
	NonceSize() int
	NewEncrypter(key []byte, salt []byte) (cipher.AEAD, error)
	NewDecrypter(key []byte, salt []byte) (cipher.AEAD, error)
}

type aeadConn struct {
	connect.IConn
	IAEADCipher
	key        []byte
	rNonce     []byte
	wNonce     []byte
	readBuffer *bytes.Buffer
	Encrypter  cipher.AEAD
	Decrypter  cipher.AEAD
}

func (a *aeadConn) Read(b []byte) (n int, err error) {
	if a.readBuffer.Len() > 0 {
		n, err = a.readBuffer.Read(b)
		return
	}
	if a.Decrypter == nil {
		salt := make([]byte, a.SaltSize())
		if _, err = io.ReadFull(a.IConn, salt); err != nil {
			return
		}
		a.Decrypter, err = a.NewDecrypter(a.key, salt)
		if err != nil {
			log.Logger.Errorf("[AEAD Conn] init decrypter failed: %v", err)
			return 0, err
		}
	}
	var overHead = a.Decrypter.Overhead()
	buf := make([]byte, 2+overHead+DataMaxSize+overHead)
	dataBuf := buf[:2+a.Decrypter.Overhead()]
	_, err = io.ReadFull(a.IConn, dataBuf)
	if err != nil {
		return
	}

	_, err = a.Decrypter.Open(dataBuf[:0], a.rNonce, dataBuf, nil)
	increment(a.rNonce)
	if err != nil {
		return 0, err
	}

	size := (int(dataBuf[0])<<8 + int(dataBuf[1])) & DataMaxSize

	dataBuf = buf[:size+a.Decrypter.Overhead()]
	_, err = io.ReadFull(a.IConn, dataBuf)
	if err != nil {
		return 0, err
	}
	if len(b) >= size {
		n = size
		_, err = a.Decrypter.Open(b[:0], a.rNonce, dataBuf, nil)
	} else {
		_, err = a.Decrypter.Open(dataBuf[:0], a.rNonce, dataBuf, nil)
		if err == nil {
			n = copy(b, dataBuf[:len(b)])
			a.readBuffer.Write(dataBuf[n:size])
		}
	}
	increment(a.rNonce)
	return
}

func (a *aeadConn) Write(b []byte) (n int, err error) {
	r := bytes.NewBuffer(b)
	var rn int
	var overHead = a.Encrypter.Overhead()
	for {
		buf := make([]byte, 2+overHead+DataMaxSize+overHead)
		dataBuf := buf[2+overHead : 2+overHead+DataMaxSize]
		rn, err = r.Read(dataBuf)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		if rn > 0 {
			n += rn
			buf = buf[:2+overHead+rn+overHead]
			dataBuf = dataBuf[:rn]
			buf[0], buf[1] = byte(rn>>8), byte(rn&0xffff)
			a.Encrypter.Seal(buf[:0], a.wNonce, buf[:2], nil)
			increment(a.wNonce)

			a.Encrypter.Seal(dataBuf[:0], a.wNonce, dataBuf, nil)
			increment(a.wNonce)

			_, ew := a.IConn.Write(buf)
			if ew != nil {
				err = ew
				break
			}
		} else {
			break
		}
	}
	return n, err
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

func HKDF_SHA1(secret, salt, info, key []byte) error {
	_, err := io.ReadFull(hkdf.New(sha1.New, secret, salt, info), key)
	return err
}

func increment(b []byte) {
	for i := range b {
		b[i]++
		if b[i] != 0 {
			return
		}
	}
}
