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

	"github.com/asaskevich/EventBus"

	"github.com/bhojpur/service/pkg/pubsub"
	"github.com/bhojpur/service/pkg/utils/logger"
)

type bus struct {
	bus EventBus.Bus
	ctx context.Context
	log logger.Logger
}

func New(logger logger.Logger) pubsub.PubSub {
	return &bus{
		log: logger,
	}
}

func (a *bus) Close() error {
	return nil
}

func (a *bus) Features() []pubsub.Feature {
	return nil
}

func (a *bus) Init(metadata pubsub.Metadata) error {
	a.bus = EventBus.New()
	a.ctx = context.Background()

	return nil
}

func (a *bus) Publish(req *pubsub.PublishRequest) error {
	a.bus.Publish(req.Topic, a.ctx, req.Data)

	return nil
}

func (a *bus) Subscribe(req pubsub.SubscribeRequest, handler pubsub.Handler) error {
	return a.bus.Subscribe(req.Topic, func(ctx context.Context, data []byte) {
		for i := 0; i < 10; i++ {
			if err := handler(ctx, &pubsub.NewMessage{Data: data, Topic: req.Topic, Metadata: req.Metadata}); err != nil {
				a.log.Error(err)

				continue
			}

			return
		}
	})
}
