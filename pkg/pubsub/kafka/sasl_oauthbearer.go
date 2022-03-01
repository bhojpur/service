package kafka

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
	ctx "context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	ccred "golang.org/x/oauth2/clientcredentials"
)

type OAuthTokenSource struct {
	CachedToken   oauth2.Token
	Extensions    map[string]string
	TokenEndpoint oauth2.Endpoint
	ClientID      string
	ClientSecret  string
	Scopes        []string
	httpClient    *http.Client
	trustedCas    []*x509.Certificate
	skipCaVerify  bool
}

type AccessToken struct {
	// Token is the access token payload.
	Token string
	// Extensions is a optional map of arbitrary key-value pairs that can be
	// sent with the SASL/OAUTHBEARER initial client response. These values are
	// ignored by the SASL server if they are unexpected. This feature is only
	// supported by Kafka >= 2.1.0.
	Extensions map[string]string
}

func newOAuthTokenSource(oidcTokenEndpoint, oidcClientID, oidcClientSecret string, oidcScopes []string) OAuthTokenSource {
	return OAuthTokenSource{TokenEndpoint: oauth2.Endpoint{TokenURL: oidcTokenEndpoint}, ClientID: oidcClientID, ClientSecret: oidcClientSecret, Scopes: oidcScopes}
}

var tokenRequestTimeout, _ = time.ParseDuration("30s")

func (ts *OAuthTokenSource) addCa(caPem string) error {
	pemBytes := []byte(caPem)

	block, _ := pem.Decode(pemBytes)

	if block == nil || block.Type != "CERTIFICATE" {
		return fmt.Errorf("PEM data not valid or not of a valid type (CERTIFICATE)")
	}

	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("error parsing PEM certificate: %w", err)
	}

	if ts.trustedCas == nil {
		ts.trustedCas = make([]*x509.Certificate, 0)
	}
	ts.trustedCas = append(ts.trustedCas, caCert)

	return nil
}

func (ts *OAuthTokenSource) configureClient() {
	if ts.httpClient != nil {
		return
	}

	tlsConfig := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: ts.skipCaVerify, //nolint:gosec
	}

	if ts.trustedCas != nil {
		caPool, err := x509.SystemCertPool()
		if err != nil {
			caPool = x509.NewCertPool()
		}

		for _, c := range ts.trustedCas {
			caPool.AddCert(c)
		}
		tlsConfig.RootCAs = caPool
	}

	ts.httpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
}

func (ts *OAuthTokenSource) Token() (*AccessToken, error) {
	if ts.CachedToken.Valid() {
		return ts.asSaramaToken(), nil
	}

	if ts.TokenEndpoint.TokenURL == "" || ts.ClientID == "" || ts.ClientSecret == "" {
		return nil, fmt.Errorf("cannot generate token, OAuthTokenSource not fully configured")
	}

	oidcCfg := ccred.Config{ClientID: ts.ClientID, ClientSecret: ts.ClientSecret, Scopes: ts.Scopes, TokenURL: ts.TokenEndpoint.TokenURL, AuthStyle: ts.TokenEndpoint.AuthStyle}

	timeoutCtx, _ := ctx.WithTimeout(ctx.TODO(), tokenRequestTimeout) //nolint:govet

	ts.configureClient()

	timeoutCtx = ctx.WithValue(timeoutCtx, oauth2.HTTPClient, ts.httpClient)

	token, err := oidcCfg.Token(timeoutCtx)
	if err != nil {
		return nil, fmt.Errorf("error generating oauth2 token: %w", err)
	}

	ts.CachedToken = *token
	return ts.asSaramaToken(), nil
}

func (ts *OAuthTokenSource) asSaramaToken() *AccessToken {
	return &(AccessToken{Token: ts.CachedToken.AccessToken, Extensions: ts.Extensions})
}
