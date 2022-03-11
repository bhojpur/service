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
	"fmt"
	"sync"
	"time"

	"github.com/bhojpur/service/pkg/engine/core/frame"
	"github.com/bhojpur/service/pkg/engine/logger"
	"github.com/bhojpur/service/pkg/reactive"
	"github.com/cenkalti/backoff/v4"
)

// Of creates an item from a value.
func Of(i interface{}) reactive.Item {
	return reactive.Item{V: i}
}

// StreamImpl is the implementation of Stream interface.
type StreamImpl struct {
	ctx        context.Context
	observable reactive.Observable
}

// appendContinueOnError appends the "ContinueOnError" to options
func appendContinueOnError(ctx context.Context, opts ...reactive.Option) []reactive.Option {
	options := append(opts, reactive.WithErrorStrategy(reactive.ContinueOnError))
	return append(options, reactive.WithContext(ctx))
}

// All determines whether all items emitted by an Observable meet some criteria
func (s *StreamImpl) All(predicate reactive.Predicate, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.All(predicate, opts...).Observe(), opts...)}
}

// AverageFloat32 calculates the average of numbers emitted by an Observable and emits the average float32.
func (s *StreamImpl) AverageFloat32(opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.AverageFloat32(opts...).Observe(), opts...)}
}

// AverageFloat64 calculates the average of numbers emitted by an Observable and emits the average float64.
func (s *StreamImpl) AverageFloat64(opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.AverageFloat64(opts...).Observe(), opts...)}
}

// AverageInt calculates the average of numbers emitted by an Observable and emits the average int.
func (s *StreamImpl) AverageInt(opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.AverageInt(opts...).Observe(), opts...)}
}

// AverageInt8 calculates the average of numbers emitted by an Observable and emits the≤ average int8.
func (s *StreamImpl) AverageInt8(opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.AverageInt8(opts...).Observe(), opts...)}
}

// AverageInt16 calculates the average of numbers emitted by an Observable and emits the average int16.
func (s *StreamImpl) AverageInt16(opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.AverageInt16(opts...).Observe(), opts...)}
}

// AverageInt32 calculates the average of numbers emitted by an Observable and emits the average int32.
func (s *StreamImpl) AverageInt32(opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.AverageInt32(opts...).Observe(), opts...)}
}

// AverageInt64 calculates the average of numbers emitted by an Observable and emits the average int64.
func (s *StreamImpl) AverageInt64(opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.AverageInt64(opts...).Observe(), opts...)}
}

// BackOffRetry implements a backoff retry if a source Observable sends an error, resubscribe to it in the hopes that it will complete without error.
// Cannot be run in parallel.
func (s *StreamImpl) BackOffRetry(backOffCfg backoff.BackOff, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.BackOffRetry(backOffCfg, opts...).Observe(), opts...)}
}

// BufferWithCount returns an Observable that emits buffers of items it collects
// from the source Observable.
// The resulting Observable emits buffers every skip items, each containing a slice of count items.
// When the source Observable completes or encounters an error,
// the resulting Observable emits the current buffer and propagates
// the notification from the source Observable.
func (s *StreamImpl) BufferWithCount(count int, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.BufferWithCount(count, opts...).Observe(), opts...)}
}

func getRxDuration(milliseconds uint32) reactive.Duration {
	return reactive.WithDuration(time.Duration(milliseconds) * time.Millisecond)
}

// BufferWithTime returns an Observable that emits buffers of items it collects from the source
// Observable. The resulting Observable starts a new buffer periodically, as determined by the
// timeshift argument. It emits each buffer after a fixed timespan, specified by the timespan argument.
// When the source Observable completes or encounters an error, the resulting Observable emits
// the current buffer and propagates the notification from the source Observable.
func (s *StreamImpl) BufferWithTime(milliseconds uint32, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.BufferWithTime(getRxDuration(milliseconds), opts...).Observe(), opts...)}
}

// BufferWithTimeOrCount returns an Observable that emits buffers of items it collects from the source
// Observable either from a given count or at a given time interval.
func (s *StreamImpl) BufferWithTimeOrCount(milliseconds uint32, count int, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.BufferWithTimeOrCount(getRxDuration(milliseconds), count, opts...).Observe(), opts...)}
}

// Connect instructs a connectable Observable to begin emitting items to its subscribers.
func (s *StreamImpl) Connect(ctx context.Context) (context.Context, reactive.Disposable) {
	return s.observable.Connect(ctx)
}

// Contains determines whether an Observable emits a particular item or not.
func (s *StreamImpl) Contains(equal reactive.Predicate, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.Contains(equal, opts...).Observe(), opts...)}
}

// Count counts the number of items emitted by the source Observable and emit only this value.
func (s *StreamImpl) Count(opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.Count(opts...).Observe(), opts...)}
}

// Debounce only emits an item from an Observable if a particular timespan has passed without it emitting another item.
func (s *StreamImpl) Debounce(milliseconds uint32, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.Debounce(getRxDuration(milliseconds), opts...).Observe(), opts...)}
}

// DefaultIfEmpty returns an Observable that emits the items emitted by the source
// Observable or a specified default item if the source Observable is empty.
func (s *StreamImpl) DefaultIfEmpty(defaultValue interface{}, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.DefaultIfEmpty(defaultValue, opts...).Observe(), opts...)}
}

// Distinct suppresses duplicate items in the original Observable and returns
// a new Observable.
func (s *StreamImpl) Distinct(apply reactive.Func, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.Distinct(apply, opts...).Observe(), opts...)}
}

// DistinctUntilChanged suppresses consecutive duplicate items in the original Observable.
// Cannot be run in parallel.
func (s *StreamImpl) DistinctUntilChanged(apply reactive.Func, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.DistinctUntilChanged(apply, opts...).Observe(), opts...)}
}

// DoOnCompleted registers a callback action that will be called once the Observable terminates.
func (s *StreamImpl) DoOnCompleted(completedFunc reactive.CompletedFunc, opts ...reactive.Option) reactive.Disposed {
	opts = appendContinueOnError(s.ctx, opts...)
	return s.observable.DoOnCompleted(completedFunc, opts...)
}

// DoOnError registers a callback action that will be called if the Observable terminates abnormally.
func (s *StreamImpl) DoOnError(errFunc reactive.ErrFunc, opts ...reactive.Option) reactive.Disposed {
	opts = appendContinueOnError(s.ctx, opts...)
	return s.observable.DoOnError(errFunc, opts...)
}

// DoOnNext registers a callback action that will be called on each item emitted by the Observable.
func (s *StreamImpl) DoOnNext(nextFunc reactive.NextFunc, opts ...reactive.Option) reactive.Disposed {
	opts = appendContinueOnError(s.ctx, opts...)
	return s.observable.DoOnNext(nextFunc, opts...)
}

// Error returns the eventual Observable error.
// This method is blocking.
func (s *StreamImpl) Error(opts ...reactive.Option) error {
	opts = appendContinueOnError(s.ctx, opts...)
	return s.observable.Error(opts...)
}

// Errors returns an eventual list of Observable errors.
// This method is blocking
func (s *StreamImpl) Errors(opts ...reactive.Option) []error {
	opts = appendContinueOnError(s.ctx, opts...)
	return s.observable.Errors(opts...)
}

// Filter emits only those items from an Observable that pass a predicate test.
func (s *StreamImpl) Filter(apply reactive.Predicate, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.Filter(apply, opts...).Observe(), opts...)}
}

// ElementAt emits only item n emitted by an Observable.
// Cannot be run in parallel.
func (s *StreamImpl) ElementAt(index uint, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.ElementAt(index, opts...).Observe(), opts...)}
}

// Find emits the first item passing a predicate then complete.
func (s *StreamImpl) Find(find reactive.Predicate, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.Find(find, opts...).Observe(), opts...)}
}

// First returns new Observable which emit only first item.
// Cannot be run in parallel.
func (s *StreamImpl) First(opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.First(opts...).Observe(), opts...)}
}

// FirstOrDefault returns new Observable which emit only first item.
// If the observable fails to emit any items, it emits a default value.
// Cannot be run in parallel.
func (s *StreamImpl) FirstOrDefault(defaultValue interface{}, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.FirstOrDefault(defaultValue, opts...).Observe(), opts...)}
}

// FlatMap transforms the items emitted by an Observable into Observables, then flatten the emissions from those into a single Observable.
func (s *StreamImpl) FlatMap(apply reactive.ItemToObservable, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.FlatMap(apply, opts...).Observe(), opts...)}
}

// ForEach subscribes to the Observable and receives notifications for each element.
func (s *StreamImpl) ForEach(nextFunc reactive.NextFunc, errFunc reactive.ErrFunc, completedFunc reactive.CompletedFunc, opts ...reactive.Option) reactive.Disposed {
	opts = appendContinueOnError(s.ctx, opts...)
	return s.observable.ForEach(nextFunc, errFunc, completedFunc, opts...)
}

// IgnoreElements ignores all items emitted by the source ObservableSource except for the errors.
// Cannot be run in parallel.
func (s *StreamImpl) IgnoreElements(opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.IgnoreElements(opts...).Observe(), opts...)}
}

// Join combines items emitted by two Observables whenever an item from one Observable is emitted during
// a time window defined according to an item emitted by the other Observable.
// The time is extracted using a timeExtractor function.
func (s *StreamImpl) Join(joiner reactive.Func2, right reactive.Observable, timeExtractor func(interface{}) time.Time, windowInMS uint32, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.Join(joiner, right, timeExtractor, getRxDuration(windowInMS), opts...).Observe(), opts...)}
}

// GroupBy divides an Observable into a set of Observables that each emit a different group of items from the original Observable, organized by key.
func (s *StreamImpl) GroupBy(length int, distribution func(reactive.Item) int, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.GroupBy(length, distribution, opts...).Observe(), opts...)}
}

// GroupByDynamic divides an Observable into a dynamic set of Observables that each emit GroupedObservable from the original Observable, organized by key.
func (s *StreamImpl) GroupByDynamic(distribution func(reactive.Item) string, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.GroupByDynamic(distribution, opts...).Observe(), opts...)}
}

// Last returns a new Observable which emit only last item.
// Cannot be run in parallel.
func (s *StreamImpl) Last(opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.Last(opts...).Observe(), opts...)}
}

// LastOrDefault returns a new Observable which emit only last item.
// If the observable fails to emit any items, it emits a default value.
// Cannot be run in parallel.
func (s *StreamImpl) LastOrDefault(defaultValue interface{}, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.LastOrDefault(defaultValue, opts...).Observe(), opts...)}
}

// Map transforms the items emitted by an Observable by applying a function to each item.
func (s *StreamImpl) Map(apply reactive.Func, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.Map(apply, opts...).Observe(), opts...)}
}

// Marshal transforms the items emitted by an Observable by applying a marshalling to each item.
func (s *StreamImpl) Marshal(marshaller Marshaller, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)

	return s.Map(func(_ context.Context, i interface{}) (interface{}, error) {
		return marshaller(i)
	}, opts...)
}

// Unmarshal transforms the items emitted by an Observable by applying an unmarshalling to each item.
func (s *StreamImpl) Unmarshal(unmarshaller Unmarshaller, factory func() interface{}, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)

	return s.Map(func(_ context.Context, i interface{}) (interface{}, error) {
		v := factory()
		err := unmarshaller(i.([]byte), v)
		if err != nil {
			return nil, err
		}
		return v, nil
	}, opts...)
}

// Max determines and emits the maximum-valued item emitted by an Observable according to a comparator.
func (s *StreamImpl) Max(comparator reactive.Comparator, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.Max(comparator, opts...).Observe(), opts...)}
}

// Min determines and emits the minimum-valued item emitted by an Observable according to a comparator.
func (s *StreamImpl) Min(comparator reactive.Comparator, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.Min(comparator, opts...).Observe(), opts...)}
}

// OnErrorResumeNext instructs an Observable to pass control to another Observable rather than invoking
// onError if it encounters an error.
func (s *StreamImpl) OnErrorResumeNext(resumeSequence reactive.ErrorToObservable, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.OnErrorResumeNext(resumeSequence, opts...).Observe(), opts...)}
}

// OnErrorReturn instructs an Observable to emit an item (returned by a specified function)
// rather than invoking onError if it encounters an error.
func (s *StreamImpl) OnErrorReturn(resumeFunc reactive.ErrorFunc, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.OnErrorReturn(resumeFunc, opts...).Observe(), opts...)}
}

// OnErrorReturnItem instructs on Observable to emit an item if it encounters an error.
func (s *StreamImpl) OnErrorReturnItem(resume interface{}, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.OnErrorReturnItem(resume, opts...).Observe(), opts...)}
}

// Reduce applies a function to each item emitted by an Observable, sequentially, and emit the final value.
func (s *StreamImpl) Reduce(apply reactive.Func2, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.Reduce(apply, opts...).Observe(), opts...)}
}

// Repeat returns an Observable that repeats the sequence of items emitted by the source Observable
// at most count times, at a particular frequency.
// Cannot run in parallel.
func (s *StreamImpl) Repeat(count int64, milliseconds uint32, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.Repeat(count, getRxDuration(milliseconds), opts...).Observe(), opts...)}
}

// Retry retries if a source Observable sends an error, resubscribe to it in the hopes that it will complete without error.
// Cannot be run in parallel.
func (s *StreamImpl) Retry(count int, shouldRetry func(error) bool, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.Retry(count, shouldRetry, opts...).Observe(), opts...)}
}

// Run creates an Observer without consuming the emitted items.
func (s *StreamImpl) Run(opts ...reactive.Option) reactive.Disposed {
	opts = appendContinueOnError(s.ctx, opts...)
	return s.observable.Run(opts...)
}

// Sample returns an Observable that emits the most recent items emitted by the source
// Iterable whenever the input Iterable emits an item.
func (s *StreamImpl) Sample(iterable reactive.Iterable, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.Sample(iterable, opts...).Observe(), opts...)}
}

// Scan apply a Func2 to each item emitted by an Observable, sequentially, and emit each successive value.
// Cannot be run in parallel.
func (s *StreamImpl) Scan(apply reactive.Func2, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.Scan(apply, opts...).Observe(), opts...)}
}

// Send sends the items to a given channel.
func (s *StreamImpl) Send(output chan<- reactive.Item, opts ...reactive.Option) {
	opts = appendContinueOnError(s.ctx, opts...)
	s.observable.Send(output, opts...)
}

// SequenceEqual emits true if an Observable and the input Observable emit the same items,
// in the same order, with the same termination state. Otherwise, it emits false.
func (s *StreamImpl) SequenceEqual(iterable reactive.Iterable, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.SequenceEqual(iterable, opts...).Observe(), opts...)}
}

// Serialize forces an Observable to make serialized calls and to be well-behaved.
func (s *StreamImpl) Serialize(from int, identifier func(interface{}) int, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.Serialize(from, identifier, opts...).Observe(), opts...)}
}

// Skip suppresses the first n items in the original Observable and
// returns a new Observable with the rest items.
// Cannot be run in parallel.
func (s *StreamImpl) Skip(nth uint, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.Skip(nth, opts...).Observe(), opts...)}
}

// SkipLast suppresses the last n items in the original Observable and
// returns a new Observable with the rest items.
// Cannot be run in parallel.
func (s *StreamImpl) SkipLast(nth uint, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.SkipLast(nth, opts...).Observe(), opts...)}
}

// SkipWhile discard items emitted by an Observable until a specified condition becomes false.
// Cannot be run in parallel.
func (s *StreamImpl) SkipWhile(apply reactive.Predicate, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.SkipWhile(apply, opts...).Observe(), opts...)}
}

// StartWith emits a specified Iterable before beginning to emit the items from the source Observable.
func (s *StreamImpl) StartWith(iterable reactive.Iterable, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.StartWith(iterable, opts...).Observe(), opts...)}
}

// SumFloat32 calculates the average of float32 emitted by an Observable and emits a float32.
func (s *StreamImpl) SumFloat32(opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.SumFloat32(opts...).Observe(), opts...)}
}

// SumFloat64 calculates the average of float64 emitted by an Observable and emits a float64.
func (s *StreamImpl) SumFloat64(opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.SumFloat64(opts...).Observe(), opts...)}
}

// SumInt64 calculates the average of integers emitted by an Observable and emits an int64.
func (s *StreamImpl) SumInt64(opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.SumInt64(opts...).Observe(), opts...)}
}

// Take emits only the first n items emitted by an Observable.
// Cannot be run in parallel.
func (s *StreamImpl) Take(nth uint, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.Take(nth, opts...).Observe(), opts...)}
}

// TakeLast emits only the last n items emitted by an Observable.
// Cannot be run in parallel.
func (s *StreamImpl) TakeLast(nth uint, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.TakeLast(nth, opts...).Observe(), opts...)}
}

// TakeUntil returns an Observable that emits items emitted by the source Observable,
// checks the specified predicate for each item, and then completes when the condition is satisfied.
// Cannot be run in parallel.
func (s *StreamImpl) TakeUntil(apply reactive.Predicate, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.TakeUntil(apply, opts...).Observe(), opts...)}
}

// TakeWhile returns an Observable that emits items emitted by the source ObservableSource so long as each
// item satisfied a specified condition, and then completes as soon as this condition is not satisfied.
// Cannot be run in parallel.
func (s *StreamImpl) TakeWhile(apply reactive.Predicate, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.TakeWhile(apply, opts...).Observe(), opts...)}
}

// TimeInterval converts an Observable that emits items into one that emits indications of the amount of time elapsed between those emissions.
func (s *StreamImpl) TimeInterval(opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.TimeInterval(opts...).Observe(), opts...)}
}

// Timestamp attaches a timestamp to each item emitted by an Observable indicating when it was emitted.
func (s *StreamImpl) Timestamp(opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.Timestamp(opts...).Observe(), opts...)}
}

// ToMap convert the sequence of items emitted by an Observable
// into a map keyed by a specified key function.
// Cannot be run in parallel.
func (s *StreamImpl) ToMap(keySelector reactive.Func, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.ToMap(keySelector, opts...).Observe(), opts...)}
}

// ToMapWithValueSelector convert the sequence of items emitted by an Observable
// into a map keyed by a specified key function and valued by another
// value function.
// Cannot be run in parallel.
func (s *StreamImpl) ToMapWithValueSelector(keySelector, valueSelector reactive.Func, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.ToMapWithValueSelector(keySelector, valueSelector, opts...).Observe(), opts...)}
}

// ToSlice collects all items from an Observable and emit them in a slice and an optional error.
// Cannot be run in parallel.
func (s *StreamImpl) ToSlice(initialCapacity int, opts ...reactive.Option) ([]interface{}, error) {
	opts = appendContinueOnError(s.ctx, opts...)
	return s.observable.ToSlice(initialCapacity, opts...)
}

// WindowWithCount periodically subdivides items from an Observable into Observable windows of a given size and emit these windows
// rather than emitting the items one at a time.
func (s *StreamImpl) WindowWithCount(count int, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.WindowWithCount(count, opts...).Observe(), opts...)}
}

// WindowWithTime periodically subdivides items from an Observable into Observables based on timed windows
// and emit them rather than emitting the items one at a time.
func (s *StreamImpl) WindowWithTime(milliseconds uint32, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.WindowWithTime(getRxDuration(milliseconds), opts...).Observe(), opts...)}
}

// WindowWithTimeOrCount periodically subdivides items from an Observable into Observables based on timed windows or a specific size
// and emit them rather than emitting the items one at a time.
func (s *StreamImpl) WindowWithTimeOrCount(milliseconds uint32, count int, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.WindowWithTimeOrCount(getRxDuration(milliseconds), count, opts...).Observe(), opts...)}
}

// ZipFromIterable merges the emissions of an Iterable via a specified function
// and emit single items for each combination based on the results of this function.
func (s *StreamImpl) ZipFromIterable(iterable reactive.Iterable, processor reactive.Func2, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	return &StreamImpl{ctx: s.ctx, observable: reactive.FromChannel(s.observable.ZipFromIterable(iterable, processor, opts...).Observe(), opts...)}
}

// Observe the items in Reactive Stream.
func (s *StreamImpl) Observe(opts ...reactive.Option) <-chan reactive.Item {
	opts = appendContinueOnError(s.ctx, opts...)
	return s.observable.Observe(opts...)
}

// DefaultIfEmptyWithTime emits a default value if didn't receive any values for duration milliseconds.
func (s *StreamImpl) DefaultIfEmptyWithTime(milliseconds uint32, defaultValue interface{}, opts ...reactive.Option) Stream {
	f := func(ctx context.Context, next chan reactive.Item) {
		defer close(next)
		observe := s.Observe(opts...)

		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-observe:
				if !ok {
					return
				}

				if !item.SendContext(ctx, next) {
					return
				}
			case <-time.After(time.Duration(milliseconds) * time.Millisecond):
				if !reactive.Of(defaultValue).SendContext(ctx, next) {
					return
				}
			}
		}
	}
	return CreateObservable(s.ctx, f, opts...)
}

// StdOut writes the item as standard output.
func (s *StreamImpl) StdOut(opts ...reactive.Option) Stream {
	f := func(ctx context.Context, next chan reactive.Item) {
		defer close(next)
		observe := s.Observe(opts...)

		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-observe:
				if !ok {
					return
				}

				if !item.Error() {
					fmt.Println("[StdOut]: ", item.V)
				}

				if !item.SendContext(ctx, next) {
					return
				}
			}
		}
	}
	return CreateObservable(s.ctx, f, opts...)
}

// AuditTime ignores values for duration milliseconds, then only emits the most recent value.
func (s *StreamImpl) AuditTime(milliseconds uint32, opts ...reactive.Option) Stream {
	opts = appendContinueOnError(s.ctx, opts...)
	o := s.observable.BufferWithTime(getRxDuration(milliseconds), opts...).Map(func(_ context.Context, i interface{}) (interface{}, error) {
		return i.([]interface{})[len(i.([]interface{}))-1], nil
	}, opts...)
	return ConvertObservable(s.ctx, o)
}

// RawBytes get the raw bytes in Stream which receives from Bhojpur Service-Processor.
func (s *StreamImpl) RawBytes() Stream {
	panic("RawBytes()")
	// f := func(ctx context.Context, next chan reactive.Item) {
	// 	defer close(next)
	// 	observe := s.Observe()
	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			return
	// 		case item, ok := <-observe:
	// 			if !ok {
	// 				return
	// 			}
	// 			if item.Error() {
	// 				continue
	// 			}
	// 			bhojpurstream, ok := (item.V).(decoder.Observable)
	// 			if !ok {
	// 				logger.Error("[RawBytes] the type of item.V is not `decoder.Observable`")
	// 				return
	// 			}

	// 			bufCh := bhojpurstream.RawBytes()
	// 			go func() {
	// 				for buf := range bufCh {
	// 					logger.Debug("[RawBytes] get the raw bytes from Bhojpur Service-Processor.", "buf", buf)
	// 					Of(buf).SendContext(ctx, next)
	// 				}
	// 			}()
	// 		}
	// 	}
	// }
	// return CreateObservable(s.ctx, f)
}

// // ZipMultiObservers subscribes multi Bhojpur Service observers, zips the values into a slice and calls the processor callback when all keys are observed.
// func (s *StreamImpl) ZipMultiObservers(observers []KeyObserveFunc, processor func(items []interface{}) (interface{}, error)) Stream {
// 	count := len(observers)
// 	if count < 2 {
// 		return s.thrown(errors.New("[ZipMultiObservers] the number of observers must be >= 2"))
// 	}

// 	// the function to zip the values into a slice
// 	var zipObserveFunc = func(_ context.Context, a interface{}, b interface{}) (interface{}, error) {
// 		items, ok := a.([]interface{})
// 		if !ok {
// 			return []interface{}{a, b}, nil
// 		}

// 		items = append(items, b)
// 		return items, nil
// 	}

// 	// the function of the `ZipMultiObservers` operator
// 	f := func(ctx context.Context, next chan reactive.Item) {
// 		defer close(next)

// 		// prepare slices and maps
// 		keys := make([]byte, count)
// 		keyObserveMap := make(map[byte]decoder.OnObserveFunc, count)
// 		keyChans := make(map[byte]chan reactive.Item, count)
// 		keyObservables := make([]reactive.Observable, count)
// 		for i, item := range observers {
// 			keys[i] = item.Key
// 			keyObserveMap[item.Key] = item.OnObserve
// 			ch := make(chan reactive.Item)
// 			keyChans[item.Key] = ch
// 			keyObservables[i] = reactive.FromChannel(ch)
// 		}

// 		// zip all observables
// 		zipObservable := keyObservables[0]
// 		for i := 1; i < count; i++ {
// 			zipObservable = zipObservable.ZipFromIterable(keyObservables[i], zipObserveFunc)
// 		}

// 		observe := s.Observe()
// 		go func() {
// 			defer func() {
// 				for _, ch := range keyChans {
// 					close(ch)
// 				}
// 			}()

// 			for {
// 				select {
// 				case <-ctx.Done():
// 					return
// 				case item, ok := <-observe:
// 					if !ok {
// 						return
// 					}
// 					if item.Error() {
// 						continue
// 					}

// 					bhojpurstream := (item.V).(decoder.Observable)
// 					// subscribe multi keys
// 					bhojpurObservable := bhojpurstream.MultiSubscribe(keys...)
// 					go func() {
// 						// get the value when the key is observed
// 						kvCh := bhojpurObservable.OnMultiObserve(keyObserveMap)
// 						for {
// 							select {
// 							case <-ctx.Done():
// 								return
// 							case kv, ok := <-kvCh:
// 								if !ok {
// 									return
// 								}

// 								ch := keyChans[kv.Key]
// 								if ch != nil {
// 									ch <- reactive.Item{V: kv.Value}
// 								} else {
// 									ch <- reactive.Item{E: fmt.Errorf("[ZipMultiObservers] ch is not found for key %v", kv.Key)}
// 								}
// 							}
// 						}
// 					}()
// 				}
// 			}
// 		}()

// 		for {
// 			// observe the value from zipObservable
// 			for item := range zipObservable.Observe(reactive.WithErrorStrategy(reactive.ContinueOnError)) {
// 				if item.Error() {
// 					logger.Error("[ZipMultiObservers] observe the value failed.", "err", item.E)
// 					continue
// 				}

// 				items, ok := item.V.([]interface{})
// 				if !ok {
// 					logger.Error("[ZipMultiObservers] item.V is not a slice")
// 					continue
// 				}

// 				// call the processor callback
// 				v, err := processor(items)
// 				if err != nil {
// 					logger.Error("[ZipMultiObservers] processor func returns an err.", "err", err)
// 					continue
// 				}

// 				if !Of(v).SendContext(ctx, next) {
// 					return
// 				}
// 			}
// 		}
// 	}
// 	return CreateObservable(s.ctx, f)
// }

// PipeBackToProcessor sets a specified DataTag to bytes and will pipe it back to processor.
func (s *StreamImpl) PipeBackToProcessor(dataTag byte) Stream {
	f := func(ctx context.Context, next chan reactive.Item) {
		defer close(next)
		observe := s.Observe()

		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-observe:
				if !ok {
					return
				}

				if item.Error() {
					continue
				}

				buf, ok := (item.V).([]byte)
				if !ok {
					logger.Errorf("[PipeBackToProcessor] the data is not a []byte, won't send pass it to next.")
					continue
				}

				data := frame.PayloadFrame{
					Tag:      dataTag,
					Carriage: buf,
				}

				if !Of(data).SendContext(ctx, next) {
					return
				}
			}
		}
	}
	return CreateObservable(s.ctx, f)
}

// SlidingWindowWithCount buffers the data in the specified sliding window size, the buffered data can be processed in the handler func.
// It returns the orginal data to Stream, not the buffered slice.
func (s *StreamImpl) SlidingWindowWithCount(windowSize int, slideSize int, handler Handler, opts ...reactive.Option) Stream {
	if windowSize <= 0 {
		return s.thrown(errors.New("windowSize must be positive"))
	}
	if slideSize <= 0 {
		return s.thrown(errors.New("slideSize must be positive"))
	}

	f := func(ctx context.Context, next chan reactive.Item) {
		defer close(next)
		observe := s.Observe()

		windowCount := 0
		currentSlideCount := 0
		buf := make([]interface{}, windowSize)
		firstTimeSend := true
		mutex := sync.Mutex{}

		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-observe:
				if !ok {
					return
				}
				if item.Error() {
					continue
				}

				mutex.Lock()
				if windowCount < windowSize {
					buf[windowCount] = item.V
					windowCount++
				}

				if windowCount == windowSize {
					// start sliding
					currentSlideCount++

					// append slide item to buffer
					if !firstTimeSend {
						buf = append(buf[1:windowSize], item.V)
					}

					// reach slide size
					if currentSlideCount%slideSize == 0 {
						err := handler(buf)
						firstTimeSend = false
						if err != nil {
							reactive.Error(err).SendContext(ctx, next)
							return
						}
					}
				}
				mutex.Unlock()
				// immediately send the original item to downstream
				Of(item.V).SendContext(ctx, next)
			}
		}
	}
	return CreateObservable(s.ctx, f, opts...)
}

// SlidingWindowWithTime buffers the data in the specified sliding window time, the buffered data can be processed in the handler func.
// It returns the orginal data to Stream, not the buffered slice.
func (s *StreamImpl) SlidingWindowWithTime(windowTimeInMS uint32, slideTimeInMS uint32, handler Handler, opts ...reactive.Option) Stream {
	f := func(ctx context.Context, next chan reactive.Item) {
		observe := s.Observe()
		buf := make([]slidingWithTimeItem, 0)
		stop := make(chan struct{})
		firstTimeSend := true
		mutex := sync.Mutex{}

		checkBuffer := func() {
			mutex.Lock()
			// filter items by time
			updatedBuf := make([]slidingWithTimeItem, 0)
			availableItems := make([]interface{}, 0)
			t := time.Now().Add(-time.Duration(windowTimeInMS) * time.Millisecond)
			for _, item := range buf {
				if item.timestamp.After(t) || item.timestamp.Equal(t) {
					updatedBuf = append(updatedBuf, item)
					availableItems = append(availableItems, item.data)
				}
			}
			buf = updatedBuf

			// apply and send items
			if len(availableItems) != 0 {
				err := handler(availableItems)
				if err != nil {
					reactive.Error(err).SendContext(ctx, next)
					return
				}
			}
			firstTimeSend = false
			mutex.Unlock()
		}

		go func() {
			defer close(next)
			for {
				select {
				case <-stop:
					checkBuffer()
					return
				case <-ctx.Done():
					return
				case <-time.After(time.Duration(windowTimeInMS) * time.Millisecond):
					if firstTimeSend {
						checkBuffer()
					}
				case <-time.After(time.Duration(slideTimeInMS) * time.Millisecond):
					checkBuffer()
				}
			}
		}()

		for {
			select {
			case <-ctx.Done():
				close(stop)
				return
			case item, ok := <-observe:
				if !ok {
					close(stop)
					return
				}
				if item.Error() {
					continue
				} else {
					mutex.Lock()
					// buffer data
					buf = append(buf, slidingWithTimeItem{
						timestamp: time.Now(),
						data:      item.V,
					})
					mutex.Unlock()
				}
				// immediately send the original item to downstream
				Of(item.V).SendContext(ctx, next)
			}
		}
	}
	return CreateObservable(s.ctx, f, opts...)
}

type slidingWithTimeItem struct {
	timestamp time.Time
	data      interface{}
}

// Handler defines a function that handle the input value.
type Handler func(interface{}) error

func (s *StreamImpl) thrown(err error) Stream {
	next := make(chan reactive.Item, 1)
	next <- reactive.Error(err)
	defer close(next)
	return &StreamImpl{observable: reactive.FromChannel(next)}
}

// CreateObservable creates a new observable.
func CreateObservable(ctx context.Context, f func(ctx context.Context, next chan reactive.Item), opts ...reactive.Option) Stream {
	next := make(chan reactive.Item)
	if ctx == nil {
		ctx = context.Background()
	}
	go f(ctx, next)
	opts = appendContinueOnError(ctx, opts...)
	return &StreamImpl{ctx: ctx, observable: reactive.FromChannel(next, opts...)}
}

// CreateProcessorObservable creates a new observable with the capacity 100 for Processor.
func CreateProcessorObservable(ctx context.Context, f func(ctx context.Context, next chan reactive.Item), opts ...reactive.Option) Stream {
	next := make(chan reactive.Item, 100)
	if ctx == nil {
		ctx = context.Background()
	}
	go f(ctx, next)
	opts = appendContinueOnError(ctx, opts...)
	return &StreamImpl{ctx: ctx, observable: reactive.FromChannel(next, opts...)}
}

// ConvertObservable converts the observable to Reactive Stream.
func ConvertObservable(ctx context.Context, observable reactive.Observable) Stream {
	if ctx == nil {
		ctx = context.Background()
	}
	return &StreamImpl{ctx: ctx, observable: observable}
}
