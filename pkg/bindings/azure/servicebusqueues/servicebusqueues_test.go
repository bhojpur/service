package servicebusqueues

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

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/metadata"
	"github.com/bhojpur/service/pkg/utils/logger"
)

func TestParseMetadata(t *testing.T) {
	var oneSecondDuration time.Duration = time.Second

	testCases := []struct {
		name                     string
		properties               map[string]string
		expectedConnectionString string
		expectedQueueName        string
		expectedTTL              time.Duration
	}{
		{
			name:                     "ConnectionString and queue name",
			properties:               map[string]string{"connectionString": "connString", "queueName": "queue1"},
			expectedConnectionString: "connString",
			expectedQueueName:        "queue1",
			expectedTTL:              AzureServiceBusDefaultMessageTimeToLive,
		},
		{
			name:                     "Empty TTL",
			properties:               map[string]string{"connectionString": "connString", "queueName": "queue1", metadata.TTLMetadataKey: ""},
			expectedConnectionString: "connString",
			expectedQueueName:        "queue1",
			expectedTTL:              AzureServiceBusDefaultMessageTimeToLive,
		},
		{
			name:                     "With TTL",
			properties:               map[string]string{"connectionString": "connString", "queueName": "queue1", metadata.TTLMetadataKey: "1"},
			expectedConnectionString: "connString",
			expectedQueueName:        "queue1",
			expectedTTL:              oneSecondDuration,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			m := bindings.Metadata{}
			m.Properties = tt.properties
			a := NewAzureServiceBusQueues(logger.NewLogger("test"))
			meta, err := a.parseMetadata(m)
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedConnectionString, meta.ConnectionString)
			assert.Equal(t, tt.expectedQueueName, meta.QueueName)
			assert.Equal(t, tt.expectedTTL, meta.ttl)
		})
	}
}

func TestParseMetadataWithInvalidTTL(t *testing.T) {
	testCases := []struct {
		name       string
		properties map[string]string
	}{
		{
			name:       "Whitespaces TTL",
			properties: map[string]string{"connectionString": "connString", "queueName": "queue1", metadata.TTLMetadataKey: "  "},
		},
		{
			name:       "Negative ttl",
			properties: map[string]string{"connectionString": "connString", "queueName": "queue1", metadata.TTLMetadataKey: "-1"},
		},
		{
			name:       "Non-numeric ttl",
			properties: map[string]string{"connectionString": "connString", "queueName": "queue1", metadata.TTLMetadataKey: "abc"},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			m := bindings.Metadata{}
			m.Properties = tt.properties

			a := NewAzureServiceBusQueues(logger.NewLogger("test"))
			_, err := a.parseMetadata(m)
			assert.NotNil(t, err)
		})
	}
}

func TestParseMetadataConnectionStringAndNamespaceNameExclusivity(t *testing.T) {
	testCases := []struct {
		name                     string
		properties               map[string]string
		expectedConnectionString string
		expectedNamespaceName    string
		expectedQueueName        string
		expectedErr              bool
	}{
		{
			name:                     "ConnectionString and queue name",
			properties:               map[string]string{"connectionString": "connString", "queueName": "queue1"},
			expectedConnectionString: "connString",
			expectedNamespaceName:    "",
			expectedQueueName:        "queue1",
			expectedErr:              false,
		},
		{
			name:                     "Empty TTL",
			properties:               map[string]string{"namespaceName": "testNamespace", "queueName": "queue1", metadata.TTLMetadataKey: ""},
			expectedConnectionString: "",
			expectedNamespaceName:    "testNamespace",
			expectedQueueName:        "queue1",
			expectedErr:              false,
		},
		{
			name:                     "With TTL",
			properties:               map[string]string{"connectionString": "connString", "namespaceName": "testNamespace", "queueName": "queue1", metadata.TTLMetadataKey: "1"},
			expectedConnectionString: "",
			expectedNamespaceName:    "",
			expectedQueueName:        "queue1",
			expectedErr:              true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			m := bindings.Metadata{}
			m.Properties = tt.properties
			a := NewAzureServiceBusQueues(logger.NewLogger("test"))
			meta, err := a.parseMetadata(m)
			if tt.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, tt.expectedConnectionString, meta.ConnectionString)
				assert.Equal(t, tt.expectedQueueName, meta.QueueName)
				assert.Equal(t, tt.expectedNamespaceName, meta.NamespaceName)
			}
		})
	}
}
