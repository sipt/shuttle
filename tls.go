package shuttle

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/sipt/shuttle/config"
	connect "github.com/sipt/shuttle/conn"
	"math/big"
	"reflect"
	"time"
	"unsafe"
)

var ca *x509.Certificate
var caBytes []byte
var key *rsa.PrivateKey

type IMITMConfig interface {
	GetMITM() *config.Mitm
	SetMITM(*config.Mitm)
}

func ApplyMITMConfig(config IMITMConfig) error {
	mitm := config.GetMITM()
	if mitm == nil {
		return nil
	}
	caBytes, err := base64.RawStdEncoding.DecodeString(mitm.CA)
	if err != nil {
		return err
	}
	keyBytes, err := base64.RawStdEncoding.DecodeString(mitm.Key)
	if err != nil {
		return err
	}
	ca, key, err = LoadCA(caBytes, keyBytes)

	MitMRules = mitm.Rules
	return err
}
func GetCACert() []byte {
	l := len(caBytes)
	if l == 0 {
		return nil
	}
	bak := make([]byte, l)
	copy(bak, caBytes)
	return bak
}

func GenerateCA() (mitm *config.Mitm, err error) {
	key, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}
	names := pkix.Name{
		Organization: []string{"Shuttle"},
		CommonName:   "Shuttle Generated CA",
		Names: []pkix.AttributeTypeAndValue{
			{
				Type:  []int{2, 5, 4, 10},
				Value: "Shuttle",
			},
			{
				Type:  []int{2, 5, 4, 3},
				Value: "Shuttle Generated CA",
			},
		},
	}
	template := &x509.Certificate{
		Version:               1,
		SerialNumber:          big.NewInt(1),
		Subject:               names,
		Issuer:                names,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(5, 0, 0),
		BasicConstraintsValid: true,
		IsCA:        true,
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}

	caBytes, err = x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return
	}
	ca, err = x509.ParseCertificate(caBytes)
	if err != nil {
		return
	}

	// Generate cert
	certBuffer := bytes.Buffer{}
	if err = pem.Encode(&certBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: caBytes}); err != nil {
		return
	}

	// Generate key
	keyBuffer := bytes.Buffer{}
	if err = pem.Encode(&keyBuffer, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}); err != nil {
		return
	}

	mitm = &config.Mitm{
		CA:    base64.RawStdEncoding.EncodeToString(certBuffer.Bytes()),
		Key:   base64.RawStdEncoding.EncodeToString(keyBuffer.Bytes()),
		Rules: MitMRules,
	}
	return
}

func LoadCA(caPem, keyPem []byte) (*x509.Certificate, *rsa.PrivateKey, error) {
	ca, _ := pem.Decode(caPem)
	if ca == nil {
		return nil, nil, errors.New("CA load failed")
	}
	caBytes = ca.Bytes
	caCert, err := x509.ParseCertificate(ca.Bytes)
	if err != nil {
		return nil, nil, err
	}
	key, _ := pem.Decode(keyPem)
	if key == nil {
		return nil, nil, errors.New("key load failed")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(key.Bytes)
	return caCert, privateKey, err
}

func makeCert(cert *x509.Certificate) ([]byte, error) {
	derBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}
	return derBytes, nil
}

func Mimt(lc, sc connect.IConn) (connect.IConn, connect.IConn, error) {
	if ca == nil {
		return nil, nil, errors.New("please first generate CA")
	}
	conf := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true,
	}
	lcID, scID := lc.GetID(), sc.GetID()
	scTls := tls.Client(sc, conf)
	err := scTls.Handshake()
	if err != nil {
		return nil, nil, fmt.Errorf("tls hand shake: %v", err)
	}
	//解析sc证书
	rt := reflect.TypeOf(scTls).Elem()
	filed, ok := rt.FieldByName("peerCertificates")
	var cert *x509.Certificate
	if ok {
		ptr := (uintptr)(unsafe.Pointer(scTls))
		cert = (*(*[]*x509.Certificate)(unsafe.Pointer(ptr + filed.Offset)))[0]
	}
	sc, err = connect.DefaultDecorateForTls(scTls, connect.TCP, scID)
	if err != nil {
		return nil, nil, err
	}
	derCert, err := makeCert(cert)
	if err != nil {
		return nil, nil, err
	}
	conf = &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true,
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{derCert},
				PrivateKey:  key,
			},
		},
	}
	lcTls := tls.Server(lc, conf)
	lc, err = connect.DefaultDecorateForTls(lcTls, connect.TCP, lcID)
	return lc, sc, err
}
