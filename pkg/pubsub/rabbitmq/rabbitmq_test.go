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
	"context"
	"errors"
	"testing"
	"time"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/pubsub"
	"github.com/bhojpur/service/pkg/utils/logger"
)

func newBroker() *rabbitMQInMemoryBroker {
	return &rabbitMQInMemoryBroker{
		buffer: make(chan amqp.Delivery, 2),
	}
}

func newRabbitMQTest(broker *rabbitMQInMemoryBroker) pubsub.PubSub {
	return &rabbitMQ{
		declaredExchanges: make(map[string]bool),
		stopped:           false,
		logger:            logger.NewLogger("test"),
		connectionDial: func(host string) (rabbitMQConnectionBroker, rabbitMQChannelBroker, error) {
			broker.connectCount++

			return broker, broker, nil
		},
	}
}

func TestNoHost(t *testing.T) {
	broker := newBroker()
	pubsubRabbitMQ := newRabbitMQTest(broker)
	err := pubsubRabbitMQ.Init(pubsub.Metadata{})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing RabbitMQ host")
}

func TestNoConsumer(t *testing.T) {
	broker := newBroker()
	pubsubRabbitMQ := newRabbitMQTest(broker)
	metadata := pubsub.Metadata{
		Properties: map[string]string{
			metadataHostKey: "anyhost",
		},
	}
	err := pubsubRabbitMQ.Init(metadata)
	assert.NoError(t, err)
	err = pubsubRabbitMQ.Subscribe(pubsub.SubscribeRequest{}, nil)
	assert.Contains(t, err.Error(), "consumerID is required for subscriptions")
}

func TestConcurrencyMode(t *testing.T) {
	t.Run("parallel", func(t *testing.T) {
		broker := newBroker()
		pubsubRabbitMQ := newRabbitMQTest(broker)
		metadata := pubsub.Metadata{
			Properties: map[string]string{
				metadataHostKey:       "anyhost",
				metadataConsumerIDKey: "consumer",
				pubsub.ConcurrencyKey: string(pubsub.Parallel),
			},
		}
		err := pubsubRabbitMQ.Init(metadata)
		assert.Nil(t, err)
		assert.Equal(t, pubsub.Parallel, pubsubRabbitMQ.(*rabbitMQ).metadata.concurrency)
	})

	t.Run("single", func(t *testing.T) {
		broker := newBroker()
		pubsubRabbitMQ := newRabbitMQTest(broker)
		metadata := pubsub.Metadata{
			Properties: map[string]string{
				metadataHostKey:       "anyhost",
				metadataConsumerIDKey: "consumer",
				pubsub.ConcurrencyKey: string(pubsub.Single),
			},
		}
		err := pubsubRabbitMQ.Init(metadata)
		assert.Nil(t, err)
		assert.Equal(t, pubsub.Single, pubsubRabbitMQ.(*rabbitMQ).metadata.concurrency)
	})

	t.Run("default", func(t *testing.T) {
		broker := newBroker()
		pubsubRabbitMQ := newRabbitMQTest(broker)
		metadata := pubsub.Metadata{
			Properties: map[string]string{
				metadataHostKey:       "anyhost",
				metadataConsumerIDKey: "consumer",
			},
		}
		err := pubsubRabbitMQ.Init(metadata)
		assert.Nil(t, err)
		assert.Equal(t, pubsub.Parallel, pubsubRabbitMQ.(*rabbitMQ).metadata.concurrency)
	})
}

func TestPublishAndSubscribe(t *testing.T) {
	broker := newBroker()
	pubsubRabbitMQ := newRabbitMQTest(broker)
	metadata := pubsub.Metadata{
		Properties: map[string]string{
			metadataHostKey:       "anyhost",
			metadataConsumerIDKey: "consumer",
		},
	}
	err := pubsubRabbitMQ.Init(metadata)
	assert.Nil(t, err)
	assert.Equal(t, 1, broker.connectCount)
	assert.Equal(t, 0, broker.closeCount)

	topic := "mytopic"

	messageCount := 0
	lastMessage := ""
	processed := make(chan bool)
	handler := func(ctx context.Context, msg *pubsub.NewMessage) error {
		messageCount++
		lastMessage = string(msg.Data)
		processed <- true

		return nil
	}

	err = pubsubRabbitMQ.Subscribe(pubsub.SubscribeRequest{Topic: topic}, handler)
	assert.Nil(t, err)

	err = pubsubRabbitMQ.Publish(&pubsub.PublishRequest{Topic: topic, Data: []byte("hello world")})
	assert.Nil(t, err)
	<-processed
	assert.Equal(t, 1, messageCount)
	assert.Equal(t, "hello world", lastMessage)

	err = pubsubRabbitMQ.Publish(&pubsub.PublishRequest{Topic: topic, Data: []byte("foo bar")})
	assert.Nil(t, err)
	<-processed
	assert.Equal(t, 2, messageCount)
	assert.Equal(t, "foo bar", lastMessage)
}

func TestPublishReconnect(t *testing.T) {
	broker := newBroker()
	pubsubRabbitMQ := newRabbitMQTest(broker)
	metadata := pubsub.Metadata{
		Properties: map[string]string{
			metadataHostKey:       "anyhost",
			metadataConsumerIDKey: "consumer",
		},
	}
	err := pubsubRabbitMQ.Init(metadata)
	assert.Nil(t, err)
	assert.Equal(t, 1, broker.connectCount)
	assert.Equal(t, 0, broker.closeCount)

	topic := "othertopic"

	messageCount := 0
	lastMessage := ""
	processed := make(chan bool)
	handler := func(ctx context.Context, msg *pubsub.NewMessage) error {
		messageCount++
		lastMessage = string(msg.Data)
		processed <- true

		return nil
	}

	err = pubsubRabbitMQ.Subscribe(pubsub.SubscribeRequest{Topic: topic}, handler)
	assert.Nil(t, err)

	err = pubsubRabbitMQ.Publish(&pubsub.PublishRequest{Topic: topic, Data: []byte("hello world")})
	assert.Nil(t, err)
	<-processed
	assert.Equal(t, 1, messageCount)
	assert.Equal(t, "hello world", lastMessage)

	err = pubsubRabbitMQ.Publish(&pubsub.PublishRequest{Topic: topic, Data: []byte(errorChannelConnection)})
	assert.NotNil(t, err)
	assert.Equal(t, 1, messageCount)
	assert.Equal(t, "hello world", lastMessage)
	// Check that reconnection happened
	assert.Equal(t, 3, broker.connectCount) // three counts - one initial connection plus 2 reconnect attempts
	assert.Equal(t, 4, broker.closeCount)   // four counts - one for connection, one for channel , times 2 reconnect attempts

	err = pubsubRabbitMQ.Publish(&pubsub.PublishRequest{Topic: topic, Data: []byte("foo bar")})
	assert.Nil(t, err)
	<-processed
	assert.Equal(t, 2, messageCount)
	assert.Equal(t, "foo bar", lastMessage)
}

func TestPublishReconnectAfterClose(t *testing.T) {
	broker := newBroker()
	pubsubRabbitMQ := newRabbitMQTest(broker)
	metadata := pubsub.Metadata{
		Properties: map[string]string{
			metadataHostKey:       "anyhost",
			metadataConsumerIDKey: "consumer",
		},
	}
	err := pubsubRabbitMQ.Init(metadata)
	assert.Nil(t, err)
	assert.Equal(t, 1, broker.connectCount)
	assert.Equal(t, 0, broker.closeCount)

	topic := "mytopic2"

	messageCount := 0
	lastMessage := ""
	processed := make(chan bool)
	handler := func(ctx context.Context, msg *pubsub.NewMessage) error {
		messageCount++
		lastMessage = string(msg.Data)
		processed <- true

		return nil
	}

	err = pubsubRabbitMQ.Subscribe(pubsub.SubscribeRequest{Topic: topic}, handler)
	assert.Nil(t, err)

	err = pubsubRabbitMQ.Publish(&pubsub.PublishRequest{Topic: topic, Data: []byte("hello world")})
	assert.Nil(t, err)
	<-processed
	assert.Equal(t, 1, messageCount)
	assert.Equal(t, "hello world", lastMessage)

	// Close PubSub
	err = pubsubRabbitMQ.Close()
	assert.Nil(t, err)
	assert.Equal(t, 2, broker.closeCount) // two counts - one for connection, one for channel

	err = pubsubRabbitMQ.Publish(&pubsub.PublishRequest{Topic: topic, Data: []byte(errorChannelConnection)})
	assert.NotNil(t, err)
	assert.Equal(t, 1, messageCount)
	assert.Equal(t, "hello world", lastMessage)
	// Check that reconnection did not happened
	assert.Equal(t, 1, broker.connectCount)
	assert.Equal(t, 2, broker.closeCount) // two counts - one for connection, one for channel
}

func TestSubscribeReconnect(t *testing.T) {
	broker := newBroker()
	pubsubRabbitMQ := newRabbitMQTest(broker)
	metadata := pubsub.Metadata{
		Properties: map[string]string{
			metadataHostKey:              "anyhost",
			metadataConsumerIDKey:        "consumer",
			metadataAutoAckKey:           "true",
			metadataReconnectWaitSeconds: "0",
			pubsub.ConcurrencyKey:        string(pubsub.Single),
		},
	}
	err := pubsubRabbitMQ.Init(metadata)
	assert.Nil(t, err)
	assert.Equal(t, 1, broker.connectCount)
	assert.Equal(t, 0, broker.closeCount)

	topic := "thetopic"

	messageCount := 0
	lastMessage := ""
	processed := make(chan bool)
	handler := func(ctx context.Context, msg *pubsub.NewMessage) error {
		messageCount++
		lastMessage = string(msg.Data)
		processed <- true

		return errors.New(errorChannelConnection)
	}

	err = pubsubRabbitMQ.Subscribe(pubsub.SubscribeRequest{Topic: topic}, handler)
	assert.Nil(t, err)

	err = pubsubRabbitMQ.Publish(&pubsub.PublishRequest{Topic: topic, Data: []byte("hello world")})
	assert.Nil(t, err)
	<-processed
	assert.Equal(t, 1, messageCount)
	assert.Equal(t, "hello world", lastMessage)

	err = pubsubRabbitMQ.Publish(&pubsub.PublishRequest{Topic: topic, Data: []byte("foo bar")})
	assert.Nil(t, err)
	<-processed
	assert.Equal(t, 2, messageCount)
	assert.Equal(t, "foo bar", lastMessage)

	// allow last reconnect completion
	time.Sleep(time.Second)

	// Check that reconnection happened
	assert.Equal(t, 3, broker.connectCount) // initial connect + 2 reconnects
	assert.Equal(t, 4, broker.closeCount)   // two counts for each connection closure - one for connection, one for channel
}

func createAMQPMessage(body []byte) amqp.Delivery {
	return amqp.Delivery{Body: body}
}

type rabbitMQInMemoryBroker struct {
	buffer chan amqp.Delivery

	connectCount int
	closeCount   int
}

func (r *rabbitMQInMemoryBroker) Qos(prefetchCount, prefetchSize int, global bool) error {
	return nil
}

func (r *rabbitMQInMemoryBroker) Publish(exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error {
	if string(msg.Body) == errorChannelConnection {
		return errors.New(errorChannelConnection)
	}

	r.buffer <- createAMQPMessage(msg.Body)

	return nil
}

func (r *rabbitMQInMemoryBroker) QueueDeclare(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return amqp.Queue{Name: name}, nil
}

func (r *rabbitMQInMemoryBroker) QueueBind(name string, key string, exchange string, noWait bool, args amqp.Table) error {
	return nil
}

func (r *rabbitMQInMemoryBroker) Consume(queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return r.buffer, nil
}

func (r *rabbitMQInMemoryBroker) Nack(tag uint64, multiple bool, requeue bool) error {
	return nil
}

func (r *rabbitMQInMemoryBroker) Ack(tag uint64, multiple bool) error {
	return nil
}

func (r *rabbitMQInMemoryBroker) ExchangeDeclare(name string, kind string, durable bool, autoDelete bool, internal bool, noWait bool, args amqp.Table) error {
	return nil
}

func (r *rabbitMQInMemoryBroker) Close() error {
	r.closeCount++

	return nil
}
