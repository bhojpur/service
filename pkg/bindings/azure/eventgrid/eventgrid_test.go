package eventgrid

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

	"github.com/bhojpur/service/pkg/bindings"
)

func TestParseMetadata(t *testing.T) {
	m := bindings.Metadata{}
	m.Properties = map[string]string{
		"tenantId":              "a",
		"subscriptionId":        "a",
		"clientId":              "a",
		"clientSecret":          "a",
		"subscriberEndpoint":    "a",
		"handshakePort":         "a",
		"scope":                 "a",
		"eventSubscriptionName": "a",
		"accessKey":             "a",
		"topicEndpoint":         "a",
	}

	eh := AzureEventGrid{}
	meta, err := eh.parseMetadata(m)

	assert.Nil(t, err)
	assert.Equal(t, "a", meta.TenantID)
	assert.Equal(t, "a", meta.SubscriptionID)
	assert.Equal(t, "a", meta.ClientID)
	assert.Equal(t, "a", meta.ClientSecret)
	assert.Equal(t, "a", meta.SubscriberEndpoint)
	assert.Equal(t, "a", meta.HandshakePort)
	assert.Equal(t, "a", meta.Scope)
	assert.Equal(t, "a", meta.EventSubscriptionName)
	assert.Equal(t, "a", meta.AccessKey)
	assert.Equal(t, "a", meta.TopicEndpoint)
}
