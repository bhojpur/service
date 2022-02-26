package redis

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
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/pubsub"
	"github.com/bhojpur/service/pkg/utils/logger"
)

func getFakeProperties() map[string]string {
	return map[string]string{
		consumerID:   "fakeConsumer",
		enableTLS:    "true",
		maxLenApprox: "1000",
	}
}

func TestParseRedisMetadata(t *testing.T) {
	t.Run("metadata is correct", func(t *testing.T) {
		fakeProperties := getFakeProperties()

		fakeMetaData := pubsub.Metadata{
			Properties: fakeProperties,
		}

		// act
		m, err := parseRedisMetadata(fakeMetaData)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, fakeProperties[consumerID], m.consumerID)
		assert.Equal(t, int64(1000), m.maxLenApprox)
	})

	t.Run("consumerID is not given", func(t *testing.T) {
		fakeProperties := getFakeProperties()

		fakeMetaData := pubsub.Metadata{
			Properties: fakeProperties,
		}
		fakeMetaData.Properties[consumerID] = ""

		// act
		m, err := parseRedisMetadata(fakeMetaData)
		// assert
		assert.Error(t, errors.New("redis streams error: missing consumerID"), err)
		assert.Empty(t, m.consumerID)
	})
}

func TestProcessStreams(t *testing.T) {
	fakeConsumerID := "fakeConsumer"
	topicCount := 0
	messageCount := 0
	expectedData := "testData"

	var wg sync.WaitGroup
	wg.Add(3)

	fakeHandler := func(ctx context.Context, msg *pubsub.NewMessage) error {
		defer wg.Done()

		messageCount++
		if topicCount == 0 {
			topicCount = 1
		}

		// assert
		assert.Equal(t, expectedData, string(msg.Data))

		// return fake error to skip executing redis client command
		return errors.New("fake error")
	}

	// act
	testRedisStream := &redisStreams{logger: logger.NewLogger("test")}
	testRedisStream.ctx, testRedisStream.cancel = context.WithCancel(context.Background())
	testRedisStream.queue = make(chan redisMessageWrapper, 10)
	go testRedisStream.worker()
	testRedisStream.enqueueMessages(fakeConsumerID, fakeHandler, generateRedisStreamTestData(2, 3, expectedData))

	// Wait for the handler to finish processing
	wg.Wait()

	// assert
	assert.Equal(t, 1, topicCount)
	assert.Equal(t, 3, messageCount)
}

func generateRedisStreamTestData(topicCount, messageCount int, data string) []redis.XMessage {
	generateXMessage := func(id int) redis.XMessage {
		return redis.XMessage{
			ID: fmt.Sprintf("%d", id),
			Values: map[string]interface{}{
				"data": data,
			},
		}
	}

	xmessageArray := make([]redis.XMessage, messageCount)
	for i := range xmessageArray {
		xmessageArray[i] = generateXMessage(i)
	}

	return xmessageArray
}
