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
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"golang.org/x/sync/errgroup"
)

func Test_Connectable_IterableChannel_Single(t *testing.T) {
	defer goleak.VerifyNone(t)
	ch := make(chan Item, 10)
	go func() {
		ch <- Of(1)
		ch <- Of(2)
		ch <- Of(3)
		close(ch)
	}()
	obs := &ObservableImpl{
		iterable: newChannelIterable(ch, WithPublishStrategy()),
	}
	testConnectableSingle(t, obs)
}

func Test_Connectable_IterableChannel_Composed(t *testing.T) {
	defer goleak.VerifyNone(t)
	ch := make(chan Item, 10)
	go func() {
		ch <- Of(1)
		ch <- Of(2)
		ch <- Of(3)
		close(ch)
	}()
	obs := &ObservableImpl{
		iterable: newChannelIterable(ch, WithPublishStrategy()),
	}
	testConnectableComposed(t, obs)
}

func Test_Connectable_IterableChannel_Disposed(t *testing.T) {
	defer goleak.VerifyNone(t)
	ch := make(chan Item, 10)
	go func() {
		ch <- Of(1)
		ch <- Of(2)
		ch <- Of(3)
		close(ch)
	}()
	obs := &ObservableImpl{
		iterable: newChannelIterable(ch, WithPublishStrategy()),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, disposable := obs.Connect(ctx)
	disposable()
	time.Sleep(50 * time.Millisecond)
	Assert(ctx, t, obs, IsEmpty())
}

func Test_Connectable_IterableChannel_WithoutConnect(t *testing.T) {
	defer goleak.VerifyNone(t)
	ch := make(chan Item, 10)
	go func() {
		ch <- Of(1)
		ch <- Of(2)
		ch <- Of(3)
		close(ch)
	}()
	obs := &ObservableImpl{
		iterable: newChannelIterable(ch, WithPublishStrategy()),
	}
	testConnectableWithoutConnect(t, obs)
}

func Test_Connectable_IterableCreate_Single(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := &ObservableImpl{
		iterable: newCreateIterable([]Producer{func(_ context.Context, ch chan<- Item) {
			ch <- Of(1)
			ch <- Of(2)
			ch <- Of(3)
			cancel()
		}}, WithPublishStrategy(), WithContext(ctx)),
	}
	testConnectableSingle(t, obs)
}

func Test_Connectable_IterableCreate_Composed(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := &ObservableImpl{
		iterable: newCreateIterable([]Producer{func(_ context.Context, ch chan<- Item) {
			ch <- Of(1)
			ch <- Of(2)
			ch <- Of(3)
			cancel()
		}}, WithPublishStrategy(), WithContext(ctx)),
	}
	testConnectableComposed(t, obs)
}

func Test_Connectable_IterableCreate_Disposed(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := &ObservableImpl{
		iterable: newCreateIterable([]Producer{func(_ context.Context, ch chan<- Item) {
			ch <- Of(1)
			ch <- Of(2)
			ch <- Of(3)
			cancel()
		}}, WithPublishStrategy(), WithContext(ctx)),
	}
	obs.Connect(ctx)
	_, cancel2 := context.WithTimeout(context.Background(), 550*time.Millisecond)
	defer cancel2()
	time.Sleep(50 * time.Millisecond)
	Assert(ctx, t, obs, IsEmpty())
}

func Test_Connectable_IterableCreate_WithoutConnect(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := &ObservableImpl{
		iterable: newCreateIterable([]Producer{func(_ context.Context, ch chan<- Item) {
			ch <- Of(1)
			ch <- Of(2)
			ch <- Of(3)
			cancel()
		}}, WithBufferedChannel(3), WithPublishStrategy(), WithContext(ctx)),
	}
	testConnectableWithoutConnect(t, obs)
}

func Test_Connectable_IterableDefer_Single(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := &ObservableImpl{
		iterable: newDeferIterable([]Producer{func(_ context.Context, ch chan<- Item) {
			ch <- Of(1)
			ch <- Of(2)
			ch <- Of(3)
			cancel()
		}}, WithBufferedChannel(3), WithPublishStrategy(), WithContext(ctx)),
	}
	testConnectableSingle(t, obs)
}

func Test_Connectable_IterableDefer_Composed(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := &ObservableImpl{
		iterable: newDeferIterable([]Producer{func(_ context.Context, ch chan<- Item) {
			ch <- Of(1)
			ch <- Of(2)
			ch <- Of(3)
			cancel()
		}}, WithBufferedChannel(3), WithPublishStrategy(), WithContext(ctx)),
	}
	testConnectableComposed(t, obs)
}

func Test_Connectable_IterableJust_Single(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := &ObservableImpl{
		iterable: newJustIterable(1, 2, 3)(WithPublishStrategy(), WithContext(ctx)),
	}
	testConnectableSingle(t, obs)
}

func Test_Connectable_IterableJust_Composed(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := &ObservableImpl{
		iterable: newJustIterable(1, 2, 3)(WithPublishStrategy(), WithContext(ctx)),
	}
	testConnectableComposed(t, obs)
}

func Test_Connectable_IterableRange_Single(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := &ObservableImpl{
		iterable: newRangeIterable(1, 3, WithPublishStrategy(), WithContext(ctx)),
	}
	testConnectableSingle(t, obs)
}

func Test_Connectable_IterableRange_Composed(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := &ObservableImpl{
		iterable: newRangeIterable(1, 3, WithPublishStrategy(), WithContext(ctx)),
	}
	testConnectableComposed(t, obs)
}

func Test_Connectable_IterableSlice_Single(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := &ObservableImpl{iterable: newSliceIterable([]Item{Of(1), Of(2), Of(3)},
		WithPublishStrategy(), WithContext(ctx))}
	testConnectableSingle(t, obs)
}

func Test_Connectable_IterableSlice_Composed(t *testing.T) {
	defer goleak.VerifyNone(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	obs := &ObservableImpl{iterable: newSliceIterable([]Item{Of(1), Of(2), Of(3)},
		WithPublishStrategy(), WithContext(ctx))}
	testConnectableComposed(t, obs)
}

func testConnectableSingle(t *testing.T, obs Observable) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	eg, _ := errgroup.WithContext(ctx)

	expected := []interface{}{1, 2, 3}

	nbConsumers := 3
	wg := sync.WaitGroup{}
	wg.Add(nbConsumers)
	// Before Connect() is called we create multiple observers
	// We check all observers receive the same items
	for i := 0; i < nbConsumers; i++ {
		eg.Go(func() error {
			observer := obs.Observe(WithContext(ctx))
			wg.Done()
			got, err := collect(ctx, observer)
			if err != nil {
				return err
			}
			if !reflect.DeepEqual(got, expected) {
				return fmt.Errorf("expected: %v, got: %v", expected, got)
			}
			return nil
		})
	}

	wg.Wait()
	obs.Connect(ctx)
	assert.NoError(t, eg.Wait())
}

func testConnectableComposed(t *testing.T, obs Observable) {
	obs = obs.Map(func(_ context.Context, i interface{}) (interface{}, error) {
		return i.(int) + 1, nil
	}, WithPublishStrategy())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	eg, _ := errgroup.WithContext(ctx)

	expected := []interface{}{2, 3, 4}

	nbConsumers := 3
	wg := sync.WaitGroup{}
	wg.Add(nbConsumers)
	// Before Connect() is called we create multiple observers
	// We check all observers receive the same items
	for i := 0; i < nbConsumers; i++ {
		eg.Go(func() error {
			observer := obs.Observe(WithContext(ctx))
			wg.Done()

			got, err := collect(ctx, observer)
			if err != nil {
				return err
			}
			if !reflect.DeepEqual(got, expected) {
				return fmt.Errorf("expected: %v, got: %v", expected, got)
			}
			return nil
		})
	}

	wg.Wait()
	obs.Connect(ctx)
	assert.NoError(t, eg.Wait())
}

func testConnectableWithoutConnect(t *testing.T, obs Observable) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	Assert(ctx, t, obs, IsEmpty())
}
