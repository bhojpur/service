package hazelcast

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
	"strconv"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hazelcast/hazelcast-go-client"
	hazelcastCore "github.com/hazelcast/hazelcast-go-client"

	"github.com/bhojpur/service/pkg/pubsub"
	"github.com/bhojpur/service/pkg/utils/logger"
	"github.com/bhojpur/service/pkg/utils/retry"
)

const (
	hazelcastServers           = "hazelcastServers"
	hazelcastBackOffMaxRetries = "backOffMaxRetries"
)

type Hazelcast struct {
	client   hazelcast.Client
	logger   logger.Logger
	metadata metadata

	ctx     context.Context
	cancel  context.CancelFunc
	backOff backoff.BackOff
}

// NewHazelcastPubSub returns a new hazelcast pub-sub implementation.
func NewHazelcastPubSub(logger logger.Logger) pubsub.PubSub {
	return &Hazelcast{logger: logger}
}

func parseHazelcastMetadata(meta pubsub.Metadata) (metadata, error) {
	m := metadata{}
	if val, ok := meta.Properties[hazelcastServers]; ok && val != "" {
		m.hazelcastServers = val
	} else {
		return m, errors.New("hazelcast error: missing hazelcast servers")
	}

	if val, ok := meta.Properties[hazelcastBackOffMaxRetries]; ok && val != "" {
		backOffMaxRetriesInt, err := strconv.Atoi(val)
		if err != nil {
			return m, fmt.Errorf("hazelcast error: invalid backOffMaxRetries %s, %v", val, err)
		}
		m.backOffMaxRetries = backOffMaxRetriesInt
	}

	return m, nil
}

func (p *Hazelcast) Init(metadata pubsub.Metadata) error {
	m, err := parseHazelcastMetadata(metadata)
	if err != nil {
		return err
	}

	p.metadata = m
	hzConfig := hazelcast.NewConfig()

	servers := m.hazelcastServers
	hzConfig.Cluster.Network.SetAddresses(strings.Split(servers, ",")...)

	client, err := hazelcast.StartNewClientWithConfig(p.ctx, hzConfig)
	if err != nil {
		return fmt.Errorf("hazelcast error: failed to create new client, %v", err)
	} else {
		p.client = *client
	}

	p.ctx, p.cancel = context.WithCancel(context.Background())

	// TODO: Make the backoff configurable for constant or exponential
	b := backoff.NewConstantBackOff(5 * time.Second)
	p.backOff = backoff.WithContext(b, p.ctx)

	return nil
}

func (p *Hazelcast) Publish(req *pubsub.PublishRequest) error {
	topic, err := p.client.GetTopic(p.ctx, req.Topic)
	if err != nil {
		return fmt.Errorf("hazelcast error: failed to get topic for %s", req.Topic)
	}

	if err = topic.Publish(p.ctx, req.Data); err != nil {
		return fmt.Errorf("hazelcast error: failed to publish data, %v", err)
	}

	return nil
}

func (p *Hazelcast) Subscribe(req pubsub.SubscribeRequest, handler pubsub.Handler) error {
	topic, err := p.client.GetTopic(p.ctx, req.Topic)
	if err != nil {
		return fmt.Errorf("hazelcast error: failed to get topic for %s", req.Topic)
	}

	_, err = topic.AddMessageListener(p.ctx, &hazelcastMessageListener{p, topic.Name(), handler})
	if err != nil {
		return fmt.Errorf("hazelcast error: failed to add new listener, %v", err)
	}

	return nil
}

func (p *Hazelcast) Close() error {
	p.cancel()
	p.client.Shutdown(p.ctx)

	return nil
}

func (p *Hazelcast) Features() []pubsub.Feature {
	return nil
}

type hazelcastMessageListener struct {
	p             *Hazelcast
	topicName     string
	pubsubHandler pubsub.Handler
}

func (l *hazelcastMessageListener) OnMessage(event *hazelcast.MessagePublished) error {
	msg, ok := event.Value([]byte)
	if !ok {
		return errors.New("hazelcast error: cannot cast message to byte array")
	}

	if err := l.handleMessageObject(msg); err != nil {
		l.p.logger.Error("Failure processing Hazelcast message")

		return err
	}

	return nil
}

func (l *hazelcastMessageListener) handleMessageObject(message []byte) error {
	pubsubMsg := pubsub.NewMessage{
		Data:  message,
		Topic: l.topicName,
	}

	b := l.p.backOff
	if l.p.metadata.backOffMaxRetries >= 0 {
		b = backoff.WithMaxRetries(b, uint64(l.p.metadata.backOffMaxRetries))
	}

	return retry.NotifyRecover(func() error {
		l.p.logger.Debug("Processing Hazelcast message")

		return l.pubsubHandler(l.p.ctx, &pubsubMsg)
	}, b, func(err error, d time.Duration) {
		l.p.logger.Error("Error processing Hazelcast message. Retrying...")
	}, func() {
		l.p.logger.Info("Successfully processed Hazelcast message after it previously failed")
	})
}
