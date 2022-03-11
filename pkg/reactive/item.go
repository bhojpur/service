package reactive

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
	"reflect"
	"time"
)

type (
	// Item is a wrapper having either a value or an error.
	Item struct {
		V interface{}
		E error
	}

	// TimestampItem attach a timestamp to an item.
	TimestampItem struct {
		Timestamp time.Time
		V         interface{}
	}

	// CloseChannelStrategy indicates a strategy on whether to close a channel.
	CloseChannelStrategy uint32
)

const (
	// LeaveChannelOpen indicates to leave the channel open after completion.
	LeaveChannelOpen CloseChannelStrategy = iota
	// CloseChannel indicates to close the channel open after completion.
	CloseChannel
)

// Of creates an item from a value.
func Of(i interface{}) Item {
	return Item{V: i}
}

// Error creates an item from an error.
func Error(err error) Item {
	return Item{E: err}
}

// SendItems is an utility function that send a list of interface{} and indicate a strategy on whether to close
// the channel once the function completes.
func SendItems(ctx context.Context, ch chan<- Item, strategy CloseChannelStrategy, items ...interface{}) {
	if strategy == CloseChannel {
		defer close(ch)
	}
	send(ctx, ch, items...)
}

func send(ctx context.Context, ch chan<- Item, items ...interface{}) {
	for _, currentItem := range items {
		switch item := currentItem.(type) {
		default:
			rt := reflect.TypeOf(item)
			switch rt.Kind() {
			default:
				Of(item).SendContext(ctx, ch)
			case reflect.Chan:
				in := reflect.ValueOf(currentItem)
				for {
					v, ok := in.Recv()
					if !ok {
						return
					}
					currentItem := v.Interface()
					switch item := currentItem.(type) {
					default:
						Of(item).SendContext(ctx, ch)
					case error:
						Error(item).SendContext(ctx, ch)
					}
				}
			case reflect.Slice:
				s := reflect.ValueOf(currentItem)
				for i := 0; i < s.Len(); i++ {
					send(ctx, ch, s.Index(i).Interface())
				}
			}
		case error:
			Error(item).SendContext(ctx, ch)
		}
	}
}

// Error checks if an item is an error.
func (i Item) Error() bool {
	return i.E != nil
}

// SendBlocking sends an item and blocks until it is sent.
func (i Item) SendBlocking(ch chan<- Item) {
	ch <- i
}

// SendContext sends an item and blocks until it is sent or a context canceled.
// It returns a boolean to indicate whether the item was sent.
func (i Item) SendContext(ctx context.Context, ch chan<- Item) bool {
	select {
	case <-ctx.Done(): // Context's done channel has the highest priority
		return false
	default:
		select {
		case <-ctx.Done():
			return false
		case ch <- i:
			return true
		}
	}
}

// SendNonBlocking sends an item without blocking.
// It returns a boolean to indicate whether the item was sent.
func (i Item) SendNonBlocking(ch chan<- Item) bool {
	select {
	default:
		return false
	case ch <- i:
		return true
	}
}
