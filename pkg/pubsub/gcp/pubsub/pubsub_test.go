package pubsub

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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/pubsub"
)

func TestInit(t *testing.T) {
	t.Run("metadata is correct with explicit creds", func(t *testing.T) {
		m := pubsub.Metadata{}
		m.Properties = map[string]string{
			"projectId":               "superproject",
			"authProviderX509CertUrl": "https://authcerturl",
			"authUri":                 "https://auth",
			"clientX509CertUrl":       "https://cert",
			"clientEmail":             "test@test.com",
			"clientId":                "id",
			"privateKey":              "****",
			"privateKeyId":            "key_id",
			"identityProjectId":       "project1",
			"tokenUri":                "https://token",
			"type":                    "serviceaccount",
			"enableMessageOrdering":   "true",
		}
		b, err := createMetadata(m)
		assert.Nil(t, err)

		assert.Equal(t, "https://authcerturl", b.AuthProviderCertURL)
		assert.Equal(t, "https://auth", b.AuthURI)
		assert.Equal(t, "https://cert", b.ClientCertURL)
		assert.Equal(t, "test@test.com", b.ClientEmail)
		assert.Equal(t, "id", b.ClientID)
		assert.Equal(t, "****", b.PrivateKey)
		assert.Equal(t, "key_id", b.PrivateKeyID)
		assert.Equal(t, "project1", b.IdentityProjectID)
		assert.Equal(t, "https://token", b.TokenURI)
		assert.Equal(t, "serviceaccount", b.Type)
		assert.Equal(t, true, b.EnableMessageOrdering)
	})

	t.Run("metadata is correct with implicit creds", func(t *testing.T) {
		m := pubsub.Metadata{}
		m.Properties = map[string]string{
			"projectId": "superproject",
		}

		b, err := createMetadata(m)
		assert.Nil(t, err)

		assert.Equal(t, "superproject", b.ProjectID)
		assert.Equal(t, "service_account", b.Type)
	})

	t.Run("missing project id", func(t *testing.T) {
		m := pubsub.Metadata{}
		m.Properties = map[string]string{}
		_, err := createMetadata(m)
		assert.Error(t, err)
	})
}
