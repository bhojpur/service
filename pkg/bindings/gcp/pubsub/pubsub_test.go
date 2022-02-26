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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/utils/logger"
)

func TestInit(t *testing.T) {
	m := bindings.Metadata{}
	m.Properties = map[string]string{
		"auth_provider_x509_cert_url": "https://auth", "auth_uri": "https://auth", "client_x509_cert_url": "https://cert", "client_email": "test@test.com", "client_id": "id", "private_key": "****",
		"private_key_id": "key_id", "project_id": "project1", "token_uri": "https://token", "type": "serviceaccount", "topic": "t1", "subscription": "s1",
	}
	ps := GCPPubSub{logger: logger.NewLogger("test")}
	b, err := ps.parseMetadata(m)
	assert.Nil(t, err)

	var pubsubMeta pubSubMetadata
	err = json.Unmarshal(b, &pubsubMeta)
	assert.Nil(t, err)

	assert.Equal(t, "s1", pubsubMeta.Subscription)
	assert.Equal(t, "t1", pubsubMeta.Topic)
	assert.Equal(t, "https://auth", pubsubMeta.AuthProviderCertURL)
	assert.Equal(t, "https://auth", pubsubMeta.AuthURI)
	assert.Equal(t, "https://cert", pubsubMeta.ClientCertURL)
	assert.Equal(t, "test@test.com", pubsubMeta.ClientEmail)
	assert.Equal(t, "id", pubsubMeta.ClientID)
	assert.Equal(t, "****", pubsubMeta.PrivateKey)
	assert.Equal(t, "key_id", pubsubMeta.PrivateKeyID)
	assert.Equal(t, "project1", pubsubMeta.ProjectID)
	assert.Equal(t, "https://token", pubsubMeta.TokenURI)
	assert.Equal(t, "serviceaccount", pubsubMeta.Type)
}
