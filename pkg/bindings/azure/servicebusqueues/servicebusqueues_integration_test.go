//go:build integration_test
// +build integration_test

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
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/metadata"
	"github.com/bhojpur/service/pkg/utils/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	// Environment variable containing the connection string to Azure Service Bus
	testServiceBusEnvKey = "APP_TEST_AZURE_SERVICEBUS"
	ttlInSeconds         = 5
)

func getTestServiceBusConnectionString() string {
	return os.Getenv(testServiceBusEnvKey)
}

type testQueueHandler struct {
	callback func(*servicebus.Message)
}

func (h testQueueHandler) Handle(ctx context.Context, message *servicebus.Message) error {
	h.callback(message)
	return message.Complete(ctx)
}

func getMessageWithRetries(queue *servicebus.Queue, maxDuration time.Duration) (*servicebus.Message, bool, error) {
	var receivedMessage *servicebus.Message

	queueHandler := testQueueHandler{
		callback: func(msg *servicebus.Message) {
			receivedMessage = msg
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), maxDuration)
	defer cancel()
	err := queue.ReceiveOne(ctx, queueHandler)
	if err != nil && err != context.DeadlineExceeded {
		return nil, false, err
	}

	return receivedMessage, receivedMessage != nil, nil
}

func TestQueueWithTTL(t *testing.T) {
	serviceBusConnectionString := getTestServiceBusConnectionString()
	assert.NotEmpty(t, serviceBusConnectionString, fmt.Sprintf("Azure ServiceBus connection string must set in environment variable '%s'", testServiceBusEnvKey))

	queueName := uuid.New().String()
	a := NewAzureServiceBusQueues(logger.NewLogger("test"))
	m := bindings.Metadata{}
	m.Properties = map[string]string{"connectionString": serviceBusConnectionString, "queueName": queueName, metadata.TTLMetadataKey: fmt.Sprintf("%d", ttlInSeconds)}
	err := a.Init(m)
	assert.Nil(t, err)

	// Assert thet queue was created with an time to live value
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(serviceBusConnectionString))
	assert.Nil(t, err)
	queue, err := ns.NewQueue(queueName)
	assert.Nil(t, err)

	qmr := ns.NewQueueManager()
	defer qmr.Delete(context.Background(), queueName)

	queueEntity, err := qmr.Get(context.Background(), queueName)
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("PT%dS", ttlInSeconds), *queueEntity.DefaultMessageTimeToLive)

	// Assert that if waited too long, we won't see any message
	const tooLateMsgContent = "too_late_msg"
	_, err = a.Invoke(&bindings.InvokeRequest{Data: []byte(tooLateMsgContent)})
	assert.Nil(t, err)

	time.Sleep(time.Second * (ttlInSeconds + 2))

	const maxGetDuration = ttlInSeconds * time.Second

	_, ok, err := getMessageWithRetries(queue, maxGetDuration)
	assert.Nil(t, err)
	assert.False(t, ok)

	// Getting before it is expired, should return it
	const testMsgContent = "test_msg"
	_, err = a.Invoke(&bindings.InvokeRequest{Data: []byte(testMsgContent)})
	assert.Nil(t, err)

	msg, ok, err := getMessageWithRetries(queue, maxGetDuration)
	assert.Nil(t, err)
	assert.True(t, ok)
	msgBody := string(msg.Data)
	assert.Equal(t, testMsgContent, msgBody)
	assert.NotNil(t, msg.TTL)
	assert.Equal(t, ttlInSeconds*time.Second, *msg.TTL)
}

func TestPublishingWithTTL(t *testing.T) {
	serviceBusConnectionString := getTestServiceBusConnectionString()
	assert.NotEmpty(t, serviceBusConnectionString, fmt.Sprintf("Azure ServiceBus connection string must set in environment variable '%s'", testServiceBusEnvKey))

	queueName := uuid.New().String()
	queueBinding1 := NewAzureServiceBusQueues(logger.NewLogger("test"))
	bindingMetadata := bindings.Metadata{}
	bindingMetadata.Properties = map[string]string{"connectionString": serviceBusConnectionString, "queueName": queueName}
	err := queueBinding1.Init(bindingMetadata)
	assert.Nil(t, err)

	// Assert thet queue was created with Azure default time to live value
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(serviceBusConnectionString))
	assert.Nil(t, err)

	queue, err := ns.NewQueue(queueName)
	assert.Nil(t, err)

	qmr := ns.NewQueueManager()
	defer qmr.Delete(context.Background(), queueName)

	queueEntity, err := qmr.Get(context.Background(), queueName)
	assert.Nil(t, err)
	const defaultAzureServiceBusMessageTimeToLive = "P14D"
	assert.Equal(t, defaultAzureServiceBusMessageTimeToLive, *queueEntity.DefaultMessageTimeToLive)

	const tooLateMsgContent = "too_late_msg"
	writeRequest := bindings.InvokeRequest{
		Data: []byte(tooLateMsgContent),
		Metadata: map[string]string{
			metadata.TTLMetadataKey: fmt.Sprintf("%d", ttlInSeconds),
		},
	}
	_, err = queueBinding1.Invoke(&writeRequest)
	assert.Nil(t, err)

	time.Sleep(time.Second * (ttlInSeconds + 2))

	const maxGetDuration = ttlInSeconds * time.Second

	_, ok, err := getMessageWithRetries(queue, maxGetDuration)
	assert.Nil(t, err)
	assert.False(t, ok)

	// Getting before it is expired, should return it
	queueBinding2 := NewAzureServiceBusQueues(logger.NewLogger("test"))
	err = queueBinding2.Init(bindingMetadata)
	assert.Nil(t, err)

	const testMsgContent = "test_msg"
	writeRequest = bindings.InvokeRequest{
		Data: []byte(testMsgContent),
		Metadata: map[string]string{
			metadata.TTLMetadataKey: fmt.Sprintf("%d", ttlInSeconds),
		},
	}
	_, err = queueBinding2.Invoke(&writeRequest)
	assert.Nil(t, err)

	msg, ok, err := getMessageWithRetries(queue, maxGetDuration)
	assert.Nil(t, err)
	assert.True(t, ok)
	msgBody := string(msg.Data)
	assert.Equal(t, testMsgContent, msgBody)
	assert.NotNil(t, msg.TTL)

	assert.Equal(t, ttlInSeconds*time.Second, *msg.TTL)
}
