package rocketmq

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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/pubsub"
	"github.com/bhojpur/service/pkg/utils/logger"
)

func getTestMetadata() map[string]string {
	return map[string]string{
		"nameServer":         "127.0.0.1:9876",
		"consumerGroup":      "bhojpur.rocketmq.producer",
		"accessKey":          "RocketMQ",
		"secretKey":          "12345",
		"consumerBatchSize":  "1",
		"consumerThreadNums": "2",
		"retries":            "2",
	}
}

func TestParseRocketMQMetadata(t *testing.T) {
	t.Run("correct metadata", func(t *testing.T) {
		meta := getTestMetadata()
		_, err := parseRocketMQMetaData(pubsub.Metadata{Properties: meta})
		assert.Nil(t, err)
	})

	t.Run("correct init", func(t *testing.T) {
		meta := getTestMetadata()
		r := NewRocketMQ(logger.NewLogger("test"))
		err := r.Init(pubsub.Metadata{Properties: meta})
		assert.Nil(t, err)
	})

	t.Run("setup producer missing nameserver", func(t *testing.T) {
		meta := getTestMetadata()
		delete(meta, "nameServer")
		r := NewRocketMQ(logger.NewLogger("test"))
		err := r.Init(pubsub.Metadata{Properties: meta})
		assert.Nil(t, err)
		req := &pubsub.PublishRequest{
			Data:       []byte("hello"),
			PubsubName: "rocketmq",
			Topic:      "test",
			Metadata:   map[string]string{},
		}
		err = r.Publish(req)
		assert.NotNil(t, err)
	})

	t.Run("subscribe illegal type", func(t *testing.T) {
		meta := getTestMetadata()
		r := NewRocketMQ(logger.NewLogger("test"))
		err := r.Init(pubsub.Metadata{Properties: meta})
		assert.Nil(t, err)

		req := pubsub.SubscribeRequest{
			Topic: "test",
			Metadata: map[string]string{
				metadataRocketmqType: "incorrect type",
			},
		}
		handler := func(ctx context.Context, msg *pubsub.NewMessage) error {
			return nil
		}
		err = r.Subscribe(req, handler)
		assert.NotNil(t, err)
	})
}
