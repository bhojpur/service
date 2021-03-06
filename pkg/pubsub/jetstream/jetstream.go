package jetstream

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
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"

	"github.com/bhojpur/service/pkg/pubsub"
	"github.com/bhojpur/service/pkg/utils/logger"
	"github.com/bhojpur/service/pkg/utils/retry"
)

type jetstreamPubSub struct {
	nc   *nats.Conn
	jsc  nats.JetStreamContext
	l    logger.Logger
	meta metadata

	ctx           context.Context
	ctxCancel     context.CancelFunc
	backOffConfig retry.Config
}

func NewJetStream(logger logger.Logger) pubsub.PubSub {
	return &jetstreamPubSub{l: logger}
}

func (js *jetstreamPubSub) Init(metadata pubsub.Metadata) error {
	var err error
	js.meta, err = parseMetadata(metadata)
	if err != nil {
		return err
	}

	var opts []nats.Option
	opts = append(opts, nats.Name(js.meta.name))

	// Set nats.UserJWT options when jwt and seed key is provided.
	if js.meta.jwt != "" && js.meta.seedKey != "" {
		opts = append(opts, nats.UserJWT(func() (string, error) {
			return js.meta.jwt, nil
		}, func(nonce []byte) ([]byte, error) {
			return sigHandler(js.meta.seedKey, nonce)
		}))
	}

	js.nc, err = nats.Connect(js.meta.natsURL, opts...)
	if err != nil {
		return err
	}
	js.l.Debugf("Connected to nats at %s", js.meta.natsURL)

	js.jsc, err = js.nc.JetStream()
	if err != nil {
		return err
	}

	js.ctx, js.ctxCancel = context.WithCancel(context.Background())

	// Default retry configuration is used if no backOff properties are set.
	if err := retry.DecodeConfigWithPrefix(
		&js.backOffConfig,
		metadata.Properties,
		"backOff"); err != nil {
		return err
	}

	js.l.Debug("JetStream initialization complete")

	return nil
}

func (js *jetstreamPubSub) Features() []pubsub.Feature {
	return nil
}

func (js *jetstreamPubSub) Publish(req *pubsub.PublishRequest) error {
	js.l.Debugf("Publishing topic %v with data: %v", req.Topic, req.Data)
	_, err := js.jsc.Publish(req.Topic, req.Data)

	return err
}

func (js *jetstreamPubSub) Subscribe(req pubsub.SubscribeRequest, handler pubsub.Handler) error {
	var opts []nats.SubOpt

	if v := js.meta.durableName; v != "" {
		opts = append(opts, nats.Durable(v))
	}

	if v := js.meta.startTime; !v.IsZero() {
		opts = append(opts, nats.StartTime(v))
	} else if v := js.meta.startSequence; v > 0 {
		opts = append(opts, nats.StartSequence(v))
	} else if js.meta.deliverAll {
		opts = append(opts, nats.DeliverAll())
	} else {
		opts = append(opts, nats.DeliverLast())
	}

	if js.meta.flowControl {
		opts = append(opts, nats.EnableFlowControl())
	}

	natsHandler := func(m *nats.Msg) {
		jsm, err := m.Metadata()
		if err != nil {
			// If we get an error, then we don't have a valid JetStream
			// message.
			js.l.Error(err)

			return
		}

		operation := func() error {
			js.l.Debugf("Processing JetStream message %s/%d", m.Subject,
				jsm.Sequence)
			opErr := handler(js.ctx, &pubsub.NewMessage{
				Topic: req.Topic,
				Data:  m.Data,
				Metadata: map[string]string{
					"Topic": m.Subject,
				},
			})
			if opErr != nil {
				return opErr
			}

			return m.Ack()
		}
		notify := func(nerr error, d time.Duration) {
			js.l.Errorf("Error processing JetStream message: %s/%d. Retrying...",
				m.Subject, jsm.Sequence)
		}
		recovered := func() {
			js.l.Infof("Successfully processed JetStream message after it previously failed: %s/%d",
				m.Subject, jsm.Sequence)
		}
		backOff := js.backOffConfig.NewBackOffWithContext(js.ctx)

		err = retry.NotifyRecover(operation, backOff, notify, recovered)
		if err != nil && !errors.Is(err, context.Canceled) {
			js.l.Errorf("Error processing message and retries are exhausted:  %s/%d.",
				m.Subject, jsm.Sequence)
		}
	}

	var err error
	if queue := js.meta.queueGroupName; queue != "" {
		js.l.Debugf("nats: subscribed to subject %s with queue group %s",
			req.Topic, js.meta.queueGroupName)
		_, err = js.jsc.QueueSubscribe(req.Topic, queue, natsHandler, opts...)
	} else {
		js.l.Debugf("nats: subscribed to subject %s", req.Topic)
		_, err = js.jsc.Subscribe(req.Topic, natsHandler, opts...)
	}

	return err
}

func (js *jetstreamPubSub) Close() error {
	js.ctxCancel()

	return js.nc.Drain()
}

// Handle nats signature request for challenge response authentication.
func sigHandler(seedKey string, nonce []byte) ([]byte, error) {
	kp, err := nkeys.FromSeed([]byte(seedKey))
	if err != nil {
		return nil, err
	}
	// Wipe our key on exit.
	defer kp.Wipe()

	sig, _ := kp.Sign(nonce)
	return sig, nil
}
