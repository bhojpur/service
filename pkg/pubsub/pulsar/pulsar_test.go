package pulsar

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
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/pubsub"
)

func TestParsePulsarMetadata(t *testing.T) {
	m := pubsub.Metadata{}
	m.Properties = map[string]string{
		"host":            "a",
		"enableTLS":       "false",
		"disableBatching": "true",
	}
	meta, err := parsePulsarMetadata(m)

	assert.Nil(t, err)
	assert.Equal(t, "a", meta.Host)
	assert.Equal(t, false, meta.EnableTLS)
	assert.Equal(t, true, meta.DisableBatching)
	assert.Equal(t, defaultTenant, meta.Tenant)
	assert.Equal(t, defaultNamespace, meta.Namespace)
}

func TestParsePublishMetadata(t *testing.T) {
	m := &pubsub.PublishRequest{}
	m.Metadata = map[string]string{
		"deliverAt":    "2021-08-31T11:45:02Z",
		"deliverAfter": "60s",
	}
	msg, err := parsePublishMetadata(m)
	assert.Nil(t, err)

	val, _ := time.ParseDuration("60s")
	assert.Equal(t, val, msg.DeliverAfter)
	assert.Equal(t, "2021-08-31T11:45:02Z",
		msg.DeliverAt.Format(time.RFC3339))
}

func TestMissingHost(t *testing.T) {
	m := pubsub.Metadata{}
	m.Properties = map[string]string{"host": ""}
	meta, err := parsePulsarMetadata(m)

	assert.Error(t, err)
	assert.Nil(t, meta)
	assert.Equal(t, "pulsar error: missing pulsar host", err.Error())
}

func TestInvalidTLSInput(t *testing.T) {
	m := pubsub.Metadata{}
	m.Properties = map[string]string{"host": "a", "enableTLS": "honk"}
	meta, err := parsePulsarMetadata(m)

	assert.Error(t, err)
	assert.Nil(t, meta)
	assert.Equal(t, "pulsar error: invalid value for enableTLS", err.Error())
}

func TestValidTenantAndNS(t *testing.T) {
	var (
		testTenant                = "testTenant"
		testNamespace             = "testNamespace"
		testTopic                 = "testTopic"
		expectPersistentResult    = "persistent://testTenant/testNamespace/testTopic"
		expectNonPersistentResult = "non-persistent://testTenant/testNamespace/testTopic"
	)
	m := pubsub.Metadata{}
	m.Properties = map[string]string{"host": "a", tenant: testTenant, namespace: testNamespace}

	t.Run("test vaild tenant and namespace", func(t *testing.T) {
		meta, err := parsePulsarMetadata(m)

		assert.Nil(t, err)
		assert.Equal(t, testTenant, meta.Tenant)
		assert.Equal(t, testNamespace, meta.Namespace)
	})

	t.Run("test persistent format topic", func(t *testing.T) {
		meta, err := parsePulsarMetadata(m)
		p := Pulsar{metadata: *meta}
		res := p.formatTopic(testTopic)

		assert.Nil(t, err)
		assert.Equal(t, true, meta.Persistent)
		assert.Equal(t, expectPersistentResult, res)
	})

	t.Run("test non-persistent format topic", func(t *testing.T) {
		m.Properties[persistent] = "false"
		meta, err := parsePulsarMetadata(m)
		p := Pulsar{metadata: *meta}
		res := p.formatTopic(testTopic)

		assert.Nil(t, err)
		assert.Equal(t, false, meta.Persistent)
		assert.Equal(t, expectNonPersistentResult, res)
	})
}
