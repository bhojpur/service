package rabbitmq

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

	"github.com/streadway/amqp"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/pubsub"
)

func getFakeProperties() map[string]string {
	props := map[string]string{}
	props[metadataHostKey] = "fakehost"
	props[metadataConsumerIDKey] = "fakeConsumerID"

	return props
}

func TestCreateMetadata(t *testing.T) {
	booleanFlagTests := []struct {
		in       string
		expected bool
	}{
		{"true", true},
		{"TRUE", true},
		{"false", false},
		{"FALSE", false},
	}

	t.Run("metadata is correct", func(t *testing.T) {
		fakeProperties := getFakeProperties()

		fakeMetaData := pubsub.Metadata{
			Properties: fakeProperties,
		}

		// act
		m, err := createMetadata(fakeMetaData)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, fakeProperties[metadataHostKey], m.host)
		assert.Equal(t, fakeProperties[metadataConsumerIDKey], m.consumerID)
		assert.Equal(t, false, m.autoAck)
		assert.Equal(t, false, m.requeueInFailure)
		assert.Equal(t, true, m.deleteWhenUnused)
		assert.Equal(t, false, m.enableDeadLetter)
		assert.Equal(t, uint8(0), m.deliveryMode)
		assert.Equal(t, uint8(0), m.prefetchCount)
		assert.Equal(t, int64(0), m.maxLen)
		assert.Equal(t, int64(0), m.maxLenBytes)
		assert.Equal(t, fanoutExchangeKind, m.exchangeKind)
	})

	t.Run("host is not given", func(t *testing.T) {
		fakeProperties := getFakeProperties()

		fakeMetaData := pubsub.Metadata{
			Properties: fakeProperties,
		}
		fakeMetaData.Properties[metadataHostKey] = ""

		// act
		m, err := createMetadata(fakeMetaData)

		// assert
		assert.EqualError(t, err, "rabbitmq pub/sub error: missing RabbitMQ host")
		assert.Empty(t, m.host)
		assert.Empty(t, m.consumerID)
	})

	invalidDeliveryModes := []string{"3", "10", "-1"}

	for _, deliveryMode := range invalidDeliveryModes {
		t.Run(fmt.Sprintf("deliveryMode value=%s", deliveryMode), func(t *testing.T) {
			fakeProperties := getFakeProperties()

			fakeMetaData := pubsub.Metadata{
				Properties: fakeProperties,
			}
			fakeMetaData.Properties[metadataDeliveryModeKey] = deliveryMode

			// act
			m, err := createMetadata(fakeMetaData)

			// assert
			assert.EqualError(t, err, "rabbitmq pub/sub error: invalid RabbitMQ delivery mode, accepted values are between 0 and 2")
			assert.Equal(t, fakeProperties[metadataHostKey], m.host)
			assert.Equal(t, fakeProperties[metadataConsumerIDKey], m.consumerID)
			assert.Equal(t, uint8(0), m.deliveryMode)
		})
	}

	t.Run("deliveryMode is set", func(t *testing.T) {
		fakeProperties := getFakeProperties()

		fakeMetaData := pubsub.Metadata{
			Properties: fakeProperties,
		}
		fakeMetaData.Properties[metadataDeliveryModeKey] = "2"

		// act
		m, err := createMetadata(fakeMetaData)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, fakeProperties[metadataHostKey], m.host)
		assert.Equal(t, fakeProperties[metadataConsumerIDKey], m.consumerID)
		assert.Equal(t, uint8(2), m.deliveryMode)
	})

	t.Run("invalid concurrency", func(t *testing.T) {
		fakeProperties := getFakeProperties()

		fakeMetaData := pubsub.Metadata{
			Properties: fakeProperties,
		}
		fakeMetaData.Properties[pubsub.ConcurrencyKey] = "a"

		// act
		_, err := createMetadata(fakeMetaData)

		// assert
		assert.Error(t, err)
	})

	t.Run("prefetchCount is set", func(t *testing.T) {
		fakeProperties := getFakeProperties()

		fakeMetaData := pubsub.Metadata{
			Properties: fakeProperties,
		}
		fakeMetaData.Properties[metadataPrefetchCount] = "1"

		// act
		m, err := createMetadata(fakeMetaData)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, fakeProperties[metadataHostKey], m.host)
		assert.Equal(t, fakeProperties[metadataConsumerIDKey], m.consumerID)
		assert.Equal(t, uint8(1), m.prefetchCount)
	})

	t.Run("maxLen and maxLenBytes is set", func(t *testing.T) {
		fakeProperties := getFakeProperties()

		fakeMetaData := pubsub.Metadata{
			Properties: fakeProperties,
		}
		fakeMetaData.Properties[metadataMaxLen] = "1"
		fakeMetaData.Properties[metadataMaxLenBytes] = "2000000"

		// act
		m, err := createMetadata(fakeMetaData)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, fakeProperties[metadataHostKey], m.host)
		assert.Equal(t, fakeProperties[metadataConsumerIDKey], m.consumerID)
		assert.Equal(t, int64(1), m.maxLen)
		assert.Equal(t, int64(2000000), m.maxLenBytes)
	})

	for _, tt := range booleanFlagTests {
		t.Run(fmt.Sprintf("autoAck value=%s", tt.in), func(t *testing.T) {
			fakeProperties := getFakeProperties()

			fakeMetaData := pubsub.Metadata{
				Properties: fakeProperties,
			}
			fakeMetaData.Properties[metadataAutoAckKey] = tt.in

			// act
			m, err := createMetadata(fakeMetaData)

			// assert
			assert.NoError(t, err)
			assert.Equal(t, fakeProperties[metadataHostKey], m.host)
			assert.Equal(t, fakeProperties[metadataConsumerIDKey], m.consumerID)
			assert.Equal(t, tt.expected, m.autoAck)
		})
	}

	for _, tt := range booleanFlagTests {
		t.Run(fmt.Sprintf("requeueInFailure value=%s", tt.in), func(t *testing.T) {
			fakeProperties := getFakeProperties()

			fakeMetaData := pubsub.Metadata{
				Properties: fakeProperties,
			}
			fakeMetaData.Properties[metadataRequeueInFailureKey] = tt.in

			// act
			m, err := createMetadata(fakeMetaData)

			// assert
			assert.NoError(t, err)
			assert.Equal(t, fakeProperties[metadataHostKey], m.host)
			assert.Equal(t, fakeProperties[metadataConsumerIDKey], m.consumerID)
			assert.Equal(t, tt.expected, m.requeueInFailure)
		})
	}

	for _, tt := range booleanFlagTests {
		t.Run(fmt.Sprintf("deleteWhenUnused value=%s", tt.in), func(t *testing.T) {
			fakeProperties := getFakeProperties()

			fakeMetaData := pubsub.Metadata{
				Properties: fakeProperties,
			}
			fakeMetaData.Properties[metadataDeleteWhenUnusedKey] = tt.in

			// act
			m, err := createMetadata(fakeMetaData)

			// assert
			assert.NoError(t, err)
			assert.Equal(t, fakeProperties[metadataHostKey], m.host)
			assert.Equal(t, fakeProperties[metadataConsumerIDKey], m.consumerID)
			assert.Equal(t, tt.expected, m.deleteWhenUnused)
		})
	}

	for _, tt := range booleanFlagTests {
		t.Run(fmt.Sprintf("durable value=%s", tt.in), func(t *testing.T) {
			fakeProperties := getFakeProperties()

			fakeMetaData := pubsub.Metadata{
				Properties: fakeProperties,
			}
			fakeMetaData.Properties[metadataDurable] = tt.in

			// act
			m, err := createMetadata(fakeMetaData)

			// assert
			assert.NoError(t, err)
			assert.Equal(t, fakeProperties[metadataHostKey], m.host)
			assert.Equal(t, fakeProperties[metadataConsumerIDKey], m.consumerID)
			assert.Equal(t, tt.expected, m.durable)
		})
	}

	for _, tt := range booleanFlagTests {
		t.Run(fmt.Sprintf("enableDeadLetter value=%s", tt.in), func(t *testing.T) {
			fakeProperties := getFakeProperties()

			fakeMetaData := pubsub.Metadata{
				Properties: fakeProperties,
			}
			fakeMetaData.Properties[metadataEnableDeadLetter] = tt.in

			// act
			m, err := createMetadata(fakeMetaData)

			// assert
			assert.NoError(t, err)
			assert.Equal(t, fakeProperties[metadataHostKey], m.host)
			assert.Equal(t, fakeProperties[metadataConsumerIDKey], m.consumerID)
			assert.Equal(t, tt.expected, m.enableDeadLetter)
		})
	}
	validExchangeKind := []string{amqp.ExchangeDirect, amqp.ExchangeTopic, amqp.ExchangeFanout, amqp.ExchangeHeaders}

	for _, exchangeKind := range validExchangeKind {
		t.Run(fmt.Sprintf("exchangeKind value=%s", exchangeKind), func(t *testing.T) {
			fakeProperties := getFakeProperties()

			fakeMetaData := pubsub.Metadata{
				Properties: fakeProperties,
			}
			fakeMetaData.Properties[metadataExchangeKind] = exchangeKind

			// act
			m, err := createMetadata(fakeMetaData)

			// assert
			assert.NoError(t, err)
			assert.Equal(t, fakeProperties[metadataHostKey], m.host)
			assert.Equal(t, fakeProperties[metadataConsumerIDKey], m.consumerID)
			assert.Equal(t, exchangeKind, m.exchangeKind)
		})
	}

	t.Run("exchangeKind is invalid", func(t *testing.T) {
		fakeProperties := getFakeProperties()

		fakeMetaData := pubsub.Metadata{
			Properties: fakeProperties,
		}
		fakeMetaData.Properties[metadataExchangeKind] = "invalid"

		// act
		_, err := createMetadata(fakeMetaData)

		// assert
		assert.Error(t, err)
	})
}
