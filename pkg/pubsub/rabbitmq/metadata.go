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
	"strconv"
	"time"

	"github.com/streadway/amqp"

	"github.com/bhojpur/service/pkg/pubsub"
)

type metadata struct {
	consumerID       string
	host             string
	durable          bool
	enableDeadLetter bool
	deleteWhenUnused bool
	autoAck          bool
	requeueInFailure bool
	deliveryMode     uint8 // Transient (0 or 1) or Persistent (2)
	prefetchCount    uint8 // Prefetch deactivated if 0
	reconnectWait    time.Duration
	concurrency      pubsub.ConcurrencyMode
	maxLen           int64
	maxLenBytes      int64
	exchangeKind     string
}

// createMetadata creates a new instance from the pubsub metadata.
func createMetadata(pubSubMetadata pubsub.Metadata) (*metadata, error) {
	result := metadata{
		durable:          true,
		deleteWhenUnused: true,
		autoAck:          false,
		reconnectWait:    time.Duration(defaultReconnectWaitSeconds) * time.Second,
		exchangeKind:     fanoutExchangeKind,
	}

	if val, found := pubSubMetadata.Properties[metadataHostKey]; found && val != "" {
		result.host = val
	} else {
		return &result, fmt.Errorf("%s missing RabbitMQ host", errorMessagePrefix)
	}

	if val, found := pubSubMetadata.Properties[metadataConsumerIDKey]; found && val != "" {
		result.consumerID = val
	}

	if val, found := pubSubMetadata.Properties[metadataDeliveryModeKey]; found && val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			if intVal < 0 || intVal > 2 {
				return &result, fmt.Errorf("%s invalid RabbitMQ delivery mode, accepted values are between 0 and 2", errorMessagePrefix)
			}
			result.deliveryMode = uint8(intVal)
		}
	}

	if val, found := pubSubMetadata.Properties[metadataDurable]; found && val != "" {
		if boolVal, err := strconv.ParseBool(val); err == nil {
			result.durable = boolVal
		}
	}

	if val, found := pubSubMetadata.Properties[metadataEnableDeadLetter]; found && val != "" {
		if boolVal, err := strconv.ParseBool(val); err == nil {
			result.enableDeadLetter = boolVal
		}
	}

	if val, found := pubSubMetadata.Properties[metadataDeleteWhenUnusedKey]; found && val != "" {
		if boolVal, err := strconv.ParseBool(val); err == nil {
			result.deleteWhenUnused = boolVal
		}
	}

	if val, found := pubSubMetadata.Properties[metadataAutoAckKey]; found && val != "" {
		if boolVal, err := strconv.ParseBool(val); err == nil {
			result.autoAck = boolVal
		}
	}

	if val, found := pubSubMetadata.Properties[metadataRequeueInFailureKey]; found && val != "" {
		if boolVal, err := strconv.ParseBool(val); err == nil {
			result.requeueInFailure = boolVal
		}
	}

	if val, found := pubSubMetadata.Properties[metadataReconnectWaitSeconds]; found && val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			result.reconnectWait = time.Duration(intVal) * time.Second
		}
	}

	if val, found := pubSubMetadata.Properties[metadataPrefetchCount]; found && val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			result.prefetchCount = uint8(intVal)
		}
	}

	if val, found := pubSubMetadata.Properties[metadataMaxLen]; found && val != "" {
		if intVal, err := strconv.ParseInt(val, 10, 64); err == nil {
			result.maxLen = intVal
		}
	}

	if val, found := pubSubMetadata.Properties[metadataMaxLenBytes]; found && val != "" {
		if intVal, err := strconv.ParseInt(val, 10, 64); err == nil {
			result.maxLenBytes = intVal
		}
	}

	if val, found := pubSubMetadata.Properties[metadataExchangeKind]; found && val != "" {
		if exchangeKindValid(val) {
			result.exchangeKind = val
		} else {
			return &result, fmt.Errorf("%s invalid RabbitMQ exchange kind %s", errorMessagePrefix, val)
		}
	}

	c, err := pubsub.Concurrency(pubSubMetadata.Properties)
	if err != nil {
		return &result, err
	}
	result.concurrency = c

	return &result, nil
}

func (m *metadata) formatQueueDeclareArgs(origin amqp.Table) amqp.Table {
	if origin == nil {
		origin = amqp.Table{}
	}
	if m.maxLen > 0 {
		origin[argMaxLength] = m.maxLen
	}
	if m.maxLenBytes > 0 {
		origin[argMaxLengthBytes] = m.maxLenBytes
	}

	return origin
}

func exchangeKindValid(kind string) bool {
	return kind == amqp.ExchangeFanout || kind == amqp.ExchangeTopic || kind == amqp.ExchangeDirect || kind == amqp.ExchangeHeaders
}
