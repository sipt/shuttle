package dump

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
	"math/big"
	"strings"
	"time"

	"github.com/sipt/shuttle/constant"
	"github.com/sipt/shuttle/constant/typ"

	"github.com/sipt/shuttle/conn"
)

var (
	mitmEnabled = false
	domainRules []string
)
var ca *x509.Certificate
var caBytes []byte
var key *rsa.PrivateKey

func InitMITM(keyEncode, caEncode string, enable bool, domains []string) error {
	mitmEnabled = enable
	domainRules = domains
	if len(keyEncode) == 0 || len(caEncode) == 0 {
		return nil
	}
	caBytes, err := base64.RawStdEncoding.DecodeString(caEncode)
	if err != nil {
		return err
	}
	keyBytes, err := base64.RawStdEncoding.DecodeString(keyEncode)
	if err != nil {
		return err
	}
	ca, key, err = LoadCA(caBytes, keyBytes)
	return err
}

func mitmIsEnabled(domain string) bool {
	if !mitmEnabled {
		return false
	}
	for _, v := range domainRules {
		if v[0] == '*' {
			return strings.HasSuffix(domain, v[1:])
		} else {
			return v == domain
		}
	}
	return false
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

func GenerateCA() (keyEncode, caEncode string, err error) {
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
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
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

	keyEncode = base64.RawStdEncoding.EncodeToString(keyBuffer.Bytes())
	caEncode = base64.RawStdEncoding.EncodeToString(certBuffer.Bytes())
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

func Mitm(lc conn.ICtxConn) (conn.ICtxConn, error) {
	return mitm(lc)
}

func mitm(lc conn.ICtxConn) (conn.ICtxConn, error) {
	if ca == nil {
		return nil, errors.New("please first generate CA")
	}
	cert := &x509.Certificate{
		SignatureAlgorithm: 4,
		Version:            3,
		Subject: pkix.Name{
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
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 366),
	}
	cert.SerialNumber, _ = rand.Prime(rand.Reader, 128)

	// request info
	req := lc.Value(constant.KeyRequestInfo).(typ.RequestInfo)
	items := strings.Split(req.Domain(), ".")
	base := fmt.Sprintf("%s.%s", items[len(items)-2], items[len(items)-1])
	cert.DNSNames = []string{base, fmt.Sprintf("*.%s", base)}

	derCert, err := makeCert(cert)
	if err != nil {
		return nil, err
	}
	conf := &tls.Config{
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
	lc = conn.NewConn(lcTls, lc.GetContext())
	lc.WithValue(constant.KeyUseTLS, true)
	return lc, err
}
