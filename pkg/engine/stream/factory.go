package stream

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

	"github.com/bhojpur/service/pkg/reactive"
)

// Factory creates the rx.Stream from several sources.
type Factory interface {
	// FromChannel creates a new Stream from a channel.
	FromChannel(ctx context.Context, channel chan interface{}) Stream

	// FromItems creates a new Stream from items.
	FromItems(ctx context.Context, items []interface{}) Stream
}

type factoryImpl struct {
}

// NewFactory creates a new Rx factory.
func NewFactory() Factory {
	return &factoryImpl{}
}

// FromChannel creates a new Stream from a channel.
func (fac *factoryImpl) FromChannel(ctx context.Context, channel chan interface{}) Stream {
	f := func(ctx context.Context, next chan reactive.Item) {
		defer close(next)

		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-channel:
				if !ok {
					return
				}

				switch item := item.(type) {
				default:
					Of(item).SendContext(ctx, next)
				case error:
					reactive.Error(item).SendContext(ctx, next)
				}
			}
		}
	}
	return CreateObservable(ctx, f)
}

// FromItems creates a new Stream from items.
func (fac *factoryImpl) FromItems(ctx context.Context, items []interface{}) Stream {
	next := make(chan reactive.Item)
	go func() {
		for _, item := range items {
			next <- Of(item)
		}
	}()

	return ConvertObservable(ctx, reactive.FromChannel(next))
}
