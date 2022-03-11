package tls

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"time"
)

var isDev bool

// CreateServerTLSConfig creates server tls config.
func CreateServerTLSConfig(host string) (*tls.Config, error) {
	// development mode
	if isDev {
		tc, err := developmentTLSConfig(host)
		if err != nil {
			return nil, err
		}
		return tc, nil
	}
	// production mode
	// CA pool
	pool, err := getCACertPool()
	if err != nil {
		return nil, err
	}
	// server certificate
	tlsCert, err := getCertAndKey()
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{*tlsCert},
		ClientCAs:    pool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		NextProtos:   []string{"bhojpur"},
	}, nil
}

// CreateClientTLSConfig creates client TLS config.
func CreateClientTLSConfig() (*tls.Config, error) {
	// development mode
	if isDev {
		return &tls.Config{
			InsecureSkipVerify: true,
			NextProtos:         []string{"bhojpur"},
			ClientSessionCache: tls.NewLRUClientSessionCache(64),
		}, nil
	}
	// production mode
	pool, err := getCACertPool()
	if err != nil {
		return nil, err
	}

	tlsCert, err := getCertAndKey()
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{*tlsCert},
		RootCAs:            pool,
		NextProtos:         []string{"bhojpur"},
		ClientSessionCache: tls.NewLRUClientSessionCache(0),
	}, nil
}

func getCACertPool() (*x509.CertPool, error) {
	var err error
	var caCert []byte

	caCertPath := os.Getenv("BHOJPUR_TLS_CACERT_FILE")
	if len(caCertPath) == 0 {
		return nil, errors.New("tls: must provide CA certificate on production mode, you can configure this via environment variables: `BHOJPUR_TLS_CACERT_FILE`")
	}

	caCert, err = ioutil.ReadFile(caCertPath)
	if err != nil {
		return nil, err
	}

	if len(caCert) == 0 {
		return nil, errors.New("tls: cannot load CA cert")
	}

	pool := x509.NewCertPool()
	if ok := pool.AppendCertsFromPEM(caCert); !ok {
		return nil, errors.New("tls: cannot append CA cert to pool")
	}

	return pool, nil
}

func getCertAndKey() (*tls.Certificate, error) {
	var err error
	var cert, key []byte

	certPath := os.Getenv("BHOJPUR_TLS_CERT_FILE")
	keyPath := os.Getenv("BHOJPUR_TLS_KEY_FILE")
	if len(certPath) == 0 || len(keyPath) == 0 {
		return nil, errors.New("tls: must provide certificate on production mode, you can configure this via environment variables: `BHOJPUR_TLS_CERT_FILE` and `BHOJPUR_TLS_KEY_FILE`")
	}

	// certificate
	cert, err = ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	// private key
	key, err = ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	if len(cert) == 0 || len(key) == 0 {
		return nil, errors.New("tls: cannot load tls cert/key")
	}

	tlsCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	return &tlsCert, nil
}

// IsDev development mode
func IsDev() bool {
	return isDev
}

// developmentTLSConfig Setup a bare-bones TLS config for the server
func developmentTLSConfig(host ...string) (*tls.Config, error) {
	tlsCert, err := generateCertificate(host...)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates:       []tls.Certificate{tlsCert},
		ClientSessionCache: tls.NewLRUClientSessionCache(1),
		NextProtos:         []string{"bhojpur"},
	}, nil
}

func generateCertificate(host ...string) (tls.Certificate, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(time.Hour * 24 * 365)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return tls.Certificate{}, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Bhojpur Consulting"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
	}

	for _, h := range host {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	// create public key
	certOut := bytes.NewBuffer(nil)
	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return tls.Certificate{}, err
	}

	// create private key
	keyOut := bytes.NewBuffer(nil)
	b, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return tls.Certificate{}, err
	}
	err = pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: b})
	if err != nil {
		return tls.Certificate{}, err
	}

	return tls.X509KeyPair(certOut.Bytes(), keyOut.Bytes())
}

func init() {
	env := os.Getenv("BHOJPUR_SERVICE_ENV")
	isDev = len(env) == 0 || env != "production"
}
