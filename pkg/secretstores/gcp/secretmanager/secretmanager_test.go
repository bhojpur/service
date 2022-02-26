package secretmanager

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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/secretstores"
	"github.com/bhojpur/service/pkg/utils/logger"
)

func TestInit(t *testing.T) {
	m := secretstores.Metadata{}
	sm := NewSecreteManager(logger.NewLogger("test"))
	t.Run("Init with Wrong metadata", func(t *testing.T) {
		m.Properties = map[string]string{
			"type":                        "service_account",
			"project_id":                  "a",
			"private_key_id":              "a",
			"private_key":                 "a",
			"client_email":                "a",
			"client_id":                   "a",
			"auth_uri":                    "a",
			"token_uri":                   "a",
			"auth_provider_x509_cert_url": "a",
			"client_x509_cert_url":        "a",
		}

		err := sm.Init(m)
		assert.NotNil(t, err)
		assert.Equal(t, err, fmt.Errorf("failed to setup secretmanager client: google: could not parse key: private key should be a PEM or plain PKCS1 or PKCS8; parse error: asn1: syntax error: truncated tag or length"))
	})

	t.Run("Init with missing `type` metadata", func(t *testing.T) {
		m.Properties = map[string]string{
			"dummy": "a",
		}
		err := sm.Init(m)
		assert.NotNil(t, err)
		assert.Equal(t, err, fmt.Errorf("missing property `type` in metadata"))
	})

	t.Run("Init with missing `project_id` metadata", func(t *testing.T) {
		m.Properties = map[string]string{
			"type": "service_account",
		}
		err := sm.Init(m)
		assert.NotNil(t, err)
		assert.Equal(t, err, fmt.Errorf("missing property `project_id` in metadata"))
	})
}

func TestGetSecret(t *testing.T) {
	sm := NewSecreteManager(logger.NewLogger("test"))

	t.Run("Get Secret - without Init", func(t *testing.T) {
		v, err := sm.GetSecret(secretstores.GetSecretRequest{Name: "test"})
		assert.NotNil(t, err)
		assert.Equal(t, err, fmt.Errorf("client is not initialized"))
		assert.Equal(t, secretstores.GetSecretResponse{Data: nil}, v)
	})

	t.Run("Get Secret - with wrong Init", func(t *testing.T) {
		m := secretstores.Metadata{
			Properties: map[string]string{
				"type":                        "service_account",
				"project_id":                  "a",
				"private_key_id":              "a",
				"private_key":                 "a",
				"client_email":                "a",
				"client_id":                   "a",
				"auth_uri":                    "a",
				"token_uri":                   "a",
				"auth_provider_x509_cert_url": "a",
				"client_x509_cert_url":        "a",
			},
		}
		sm.Init(m)
		v, err := sm.GetSecret(secretstores.GetSecretRequest{Name: "test"})
		assert.NotNil(t, err)
		assert.Equal(t, secretstores.GetSecretResponse{Data: nil}, v)
	})
}

func TestBulkGetSecret(t *testing.T) {
	sm := NewSecreteManager(logger.NewLogger("test"))

	t.Run("Bulk Get Secret - without Init", func(t *testing.T) {
		v, err := sm.BulkGetSecret(secretstores.BulkGetSecretRequest{})
		assert.NotNil(t, err)
		assert.Equal(t, err, fmt.Errorf("client is not initialized"))
		assert.Equal(t, secretstores.BulkGetSecretResponse{Data: nil}, v)
	})

	t.Run("Bulk Get Secret - with wrong Init", func(t *testing.T) {
		m := secretstores.Metadata{
			Properties: map[string]string{
				"type":                        "service_account",
				"project_id":                  "a",
				"private_key_id":              "a",
				"private_key":                 "a",
				"client_email":                "a",
				"client_id":                   "a",
				"auth_uri":                    "a",
				"token_uri":                   "a",
				"auth_provider_x509_cert_url": "a",
				"client_x509_cert_url":        "a",
			},
		}
		sm.Init(m)
		v, err := sm.BulkGetSecret(secretstores.BulkGetSecretRequest{})
		assert.NotNil(t, err)
		assert.Equal(t, secretstores.BulkGetSecretResponse{Data: nil}, v)
	})
}
