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
	"errors"
	"testing"
	"time"

	"github.com/bhojpur/service/pkg/reactive"
	"github.com/stretchr/testify/assert"
)

// HELPER FUNCTIONS

// // Reference:
// func channelValue(ctx context.Context, items ...interface{}) chan reactive.Item {
// 	next := make(chan reactive.Item)
// 	go func() {
// 		for _, item := range items {
// 			switch item := item.(type) {
// 			default:
// 				reactive.Of(item).SendContext(ctx, next)
// 			case error:
// 				reactive.Error(item).SendContext(ctx, next)
// 			}
// 		}
// 		close(next)
// 	}()
// 	return next
// }

// func newStream(ctx context.Context, items ...interface{}) Stream {
// 	return &StreamImpl{
// 		observable: reactive.FromChannel(channelValue(ctx, items...)),
// 	}
// }

func toStream(obs reactive.Observable) Stream {
	return &StreamImpl{observable: obs}
}

// TESTS

var testStream = toStream(reactive.Defer([]reactive.Producer{func(_ context.Context, ch chan<- reactive.Item) {
	for i := 1; i <= 3; i++ {
		ch <- reactive.Of(i)
		time.Sleep(100 * time.Millisecond)
	}
}}))

func Test_DefaultIfEmptyWithTime_Empty(t *testing.T) {
	t.Run("0 milliseconds", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		st := toStream(reactive.Empty()).DefaultIfEmptyWithTime(0, 3)
		reactive.Assert(ctx, t, st, reactive.IsEmpty())
	})

	t.Run("100 milliseconds", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		obs := reactive.Timer(reactive.WithDuration(120 * time.Millisecond))
		st := toStream(obs).DefaultIfEmptyWithTime(100, 3)
		reactive.Assert(ctx, t, st, reactive.HasItem(3))
	})
}

func Test_DefaultIfEmptyWithTime_NotEmpty(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	st := testStream.DefaultIfEmptyWithTime(100, 3)
	reactive.Assert(ctx, t, st, reactive.HasItemsNoOrder(1, 3, 2, 3, 3, 3))
}

func Test_StdOut_Empty(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	st := toStream(reactive.Empty()).StdOut()
	reactive.Assert(ctx, t, st, reactive.IsEmpty())
}

func Test_StdOut_NotEmpty(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	st := testStream.StdOut()
	reactive.Assert(ctx, t, st, reactive.HasItems(1, 2, 3))
}

func Test_AuditTime(t *testing.T) {
	t.Run("0 milliseconds", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		st := testStream.AuditTime(0)
		reactive.Assert(ctx, t, st, reactive.HasItems(1, 2, 3))
	})

	t.Run("100 milliseconds", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		st := testStream.AuditTime(120)
		reactive.Assert(ctx, t, st, reactive.HasItems(2, 3))
	})

	t.Run("keep last", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		st := testStream.AuditTime(500)
		reactive.Assert(ctx, t, st, reactive.HasItem(3))
	})
}

type testStruct struct {
	ID   uint32 `bhojpur:"0x11"`
	Name string `bhojpur:"0x12"`
}

func Test_SlidingWindowWithCount(t *testing.T) {
	t.Run("window size = 1, slide size = 1, handler does nothing", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		st := testStream.SlidingWindowWithCount(1, 1, func(buf interface{}) error {
			return nil
		})
		reactive.Assert(ctx, t, st, reactive.HasItems(1, 2, 3))
	})

	t.Run("window size = 3, slide size = 3, handler sums elements in buf", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		st := testStream.SlidingWindowWithCount(3, 3, func(buf interface{}) error {
			slice, ok := buf.([]interface{})
			assert.Equal(t, true, ok)
			sum := 0
			for _, v := range slice {
				sum += v.(int)
			}
			assert.Equal(t, 6, sum)
			return nil
		})
		reactive.Assert(ctx, t, st, reactive.HasItems(1, 2, 3))
	})
}

func Test_SlidingWindowWithTime(t *testing.T) {
	t.Run("window size = 120ms, slide size = 120ms, handler does nothing", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		st := testStream.SlidingWindowWithTime(120, 120, func(buf interface{}) error {
			return nil
		})
		reactive.Assert(ctx, t, st, reactive.HasItems(1, 2, 3))
	})

	t.Run("window size = 360ms, slide size = 360ms, handler sums elements in buf", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		st := testStream.SlidingWindowWithTime(360, 360, func(buf interface{}) error {
			slice, ok := buf.([]interface{})
			assert.Equal(t, true, ok)
			sum := 0
			for _, v := range slice {
				sum += v.(int)
			}
			assert.Equal(t, 6, sum)
			return nil
		})
		reactive.Assert(ctx, t, st, reactive.HasItems(1, 2, 3))
	})
}

func Test_ContinueOnError(t *testing.T) {
	t.Run("ContinueOnError on a single operator by default", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		errFoo := errors.New("foo")
		defer cancel()
		obs := testStream.
			Map(func(_ context.Context, i interface{}) (interface{}, error) {
				if i == 2 {
					return nil, errFoo
				}
				return i, nil
			})
		reactive.Assert(ctx, t, obs, reactive.HasItems(1, 3), reactive.HasError(errFoo))
	})

	t.Run("ContinueOnError on Handler by default", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		errFoo := errors.New("foo")
		defer cancel()

		handler := func(stream Stream) Stream {
			stream = stream.
				Map(func(_ context.Context, i interface{}) (interface{}, error) {
					if i == 2 {
						return nil, errFoo
					}
					return i, nil
				})
			return stream
		}

		stream := handler(testStream)
		reactive.Assert(ctx, t, stream, reactive.HasItems(1, 3), reactive.HasError(errFoo))
	})
}
