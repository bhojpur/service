package inmemory

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

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/pubsub"
	"github.com/bhojpur/service/pkg/utils/logger"
)

func TestNewInMemoryBus(t *testing.T) {
	bus := New(logger.NewLogger("test"))
	bus.Init(pubsub.Metadata{})

	ch := make(chan []byte)
	bus.Subscribe(pubsub.SubscribeRequest{Topic: "demo"}, func(ctx context.Context, msg *pubsub.NewMessage) error {
		return publish(ch, msg)
	})

	bus.Publish(&pubsub.PublishRequest{Data: []byte("ABCD"), Topic: "demo"})
	assert.Equal(t, "ABCD", string(<-ch))
}

func TestMultipleSubscribers(t *testing.T) {
	bus := New(logger.NewLogger("test"))
	bus.Init(pubsub.Metadata{})

	ch1 := make(chan []byte)
	ch2 := make(chan []byte)
	bus.Subscribe(pubsub.SubscribeRequest{Topic: "demo"}, func(ctx context.Context, msg *pubsub.NewMessage) error {
		return publish(ch1, msg)
	})

	bus.Subscribe(pubsub.SubscribeRequest{Topic: "demo"}, func(ctx context.Context, msg *pubsub.NewMessage) error {
		return publish(ch2, msg)
	})

	bus.Publish(&pubsub.PublishRequest{Data: []byte("ABCD"), Topic: "demo"})

	assert.Equal(t, "ABCD", string(<-ch1))
	assert.Equal(t, "ABCD", string(<-ch2))
}

func TestRetry(t *testing.T) {
	bus := New(logger.NewLogger("test"))
	bus.Init(pubsub.Metadata{})

	ch := make(chan []byte)
	i := -1

	bus.Subscribe(pubsub.SubscribeRequest{Topic: "demo"}, func(ctx context.Context, msg *pubsub.NewMessage) error {
		i++
		if i < 5 {
			return errors.New("if at first you don't succeed")
		}

		return publish(ch, msg)
	})

	bus.Publish(&pubsub.PublishRequest{Data: []byte("ABCD"), Topic: "demo"})
	assert.Equal(t, "ABCD", string(<-ch))
	assert.Equal(t, 5, i)
}

func publish(ch chan []byte, msg *pubsub.NewMessage) error {
	go func() { ch <- msg.Data }()

	return nil
}
