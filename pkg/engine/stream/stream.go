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
	"time"

	"github.com/bhojpur/service/pkg/reactive"
	"github.com/cenkalti/backoff/v4"
)

// Stream is the interface for Reactive Stream.
type Stream interface {
	reactive.Iterable

	// PipeBackToProcessor write the DataFrame with a specified DataID.
	PipeBackToProcessor(dataID byte) Stream

	// RawBytes get the raw bytes in Stream which receives from Bhojpur Service-Processor.
	RawBytes() Stream

	// StdOut writes the value as standard output.
	StdOut(opts ...reactive.Option) Stream

	// AuditTime ignores values for duration milliseconds, then only emits the most recent value.
	AuditTime(milliseconds uint32, opts ...reactive.Option) Stream

	// DefaultIfEmptyWithTime emits a default value if didn't receive any values for duration milliseconds.
	DefaultIfEmptyWithTime(milliseconds uint32, defaultValue interface{}, opts ...reactive.Option) Stream

	// All determines whether all items emitted by an Observable meet some criteria
	All(predicate reactive.Predicate, opts ...reactive.Option) Stream

	// AverageFloat32 calculates the average of numbers emitted by an Observable and emits the average float32.
	AverageFloat32(opts ...reactive.Option) Stream

	// AverageFloat64 calculates the average of numbers emitted by an Observable and emits the average float64.
	AverageFloat64(opts ...reactive.Option) Stream

	// AverageInt calculates the average of numbers emitted by an Observable and emits the average int.
	AverageInt(opts ...reactive.Option) Stream

	// AverageInt8 calculates the average of numbers emitted by an Observable and emits the average int8.
	AverageInt8(opts ...reactive.Option) Stream

	// AverageInt16 calculates the average of numbers emitted by an Observable and emits the average int16.
	AverageInt16(opts ...reactive.Option) Stream

	// AverageInt32 calculates the average of numbers emitted by an Observable and emits the average int32.
	AverageInt32(opts ...reactive.Option) Stream

	// AverageInt64 calculates the average of numbers emitted by an Observable and emits the average int64.
	AverageInt64(opts ...reactive.Option) Stream

	// BackOffRetry implements a backoff retry if a source Observable sends an error, resubscribe to it in the hopes that it will complete without error.
	// Cannot be run in parallel.
	BackOffRetry(backOffCfg backoff.BackOff, opts ...reactive.Option) Stream

	// BufferWithCount returns an Observable that emits buffers of items it collects
	// from the source Observable.
	// The resulting Observable emits buffers every skip items, each containing a slice of count items.
	// When the source Observable completes or encounters an error,
	// the resulting Observable emits the current buffer and propagates
	// the notification from the source Observable.
	BufferWithCount(count int, opts ...reactive.Option) Stream

	// BufferWithTime returns an Observable that emits buffers of items it collects from the source
	// Observable. The resulting Observable starts a new buffer periodically, as determined by the
	// timeshift argument. It emits each buffer after a fixed timespan, specified by the timespan argument.
	// When the source Observable completes or encounters an error, the resulting Observable emits
	// the current buffer and propagates the notification from the source Observable.
	BufferWithTime(milliseconds uint32, opts ...reactive.Option) Stream

	// BufferWithTimeOrCount returns an Observable that emits buffers of items it collects from the source
	// Observable either from a given count or at a given time interval.
	BufferWithTimeOrCount(milliseconds uint32, count int, opts ...reactive.Option) Stream

	// Connect instructs a connectable Observable to begin emitting items to its subscribers.
	Connect(ctx context.Context) (context.Context, reactive.Disposable)

	// Contains determines whether an Observable emits a particular item or not.
	Contains(equal reactive.Predicate, opts ...reactive.Option) Stream

	// Count counts the number of items emitted by the source Observable and emit only this value.
	Count(opts ...reactive.Option) Stream

	// Debounce only emits an item from an Observable if a particular timespan has passed without it emitting another item.
	Debounce(milliseconds uint32, opts ...reactive.Option) Stream

	// DefaultIfEmpty returns an Observable that emits the items emitted by the source
	// Observable or a specified default item if the source Observable is empty.
	DefaultIfEmpty(defaultValue interface{}, opts ...reactive.Option) Stream

	// Distinct suppresses duplicate items in the original Observable and returns
	// a new Observable.
	Distinct(apply reactive.Func, opts ...reactive.Option) Stream

	// DistinctUntilChanged suppresses consecutive duplicate items in the original Observable.
	// Cannot be run in parallel.
	DistinctUntilChanged(apply reactive.Func, opts ...reactive.Option) Stream

	// DoOnCompleted registers a callback action that will be called once the Observable terminates.
	DoOnCompleted(completedFunc reactive.CompletedFunc, opts ...reactive.Option) reactive.Disposed

	// DoOnError registers a callback action that will be called if the Observable terminates abnormally.
	DoOnError(errFunc reactive.ErrFunc, opts ...reactive.Option) reactive.Disposed

	// DoOnNext registers a callback action that will be called on each item emitted by the Observable.
	DoOnNext(nextFunc reactive.NextFunc, opts ...reactive.Option) reactive.Disposed

	// ElementAt emits only item n emitted by an Observable.
	// Cannot be run in parallel.
	ElementAt(index uint, opts ...reactive.Option) Stream

	// Error returns the eventual Observable error.
	// This method is blocking.
	Error(opts ...reactive.Option) error

	// Errors returns an eventual list of Observable errors.
	// This method is blocking
	Errors(opts ...reactive.Option) []error

	// Filter emits only those items from an Observable that pass a predicate test.
	Filter(apply reactive.Predicate, opts ...reactive.Option) Stream

	// Find emits the first item passing a predicate then complete.
	Find(find reactive.Predicate, opts ...reactive.Option) Stream

	// First returns new Observable which emit only first item.
	// Cannot be run in parallel.
	First(opts ...reactive.Option) Stream

	// FirstOrDefault returns new Observable which emit only first item.
	// If the observable fails to emit any items, it emits a default value.
	// Cannot be run in parallel.
	FirstOrDefault(defaultValue interface{}, opts ...reactive.Option) Stream

	// FlatMap transforms the items emitted by an Observable into Observables, then flatten the emissions from those into a single Observable.
	FlatMap(apply reactive.ItemToObservable, opts ...reactive.Option) Stream

	// ForEach subscribes to the Observable and receives notifications for each element.
	ForEach(nextFunc reactive.NextFunc, errFunc reactive.ErrFunc, completedFunc reactive.CompletedFunc, opts ...reactive.Option) reactive.Disposed

	// GroupBy divides an Observable into a set of Observables that each emit a different group of items from the original Observable, organized by key.
	GroupBy(length int, distribution func(reactive.Item) int, opts ...reactive.Option) Stream

	// GroupByDynamic divides an Observable into a dynamic set of Observables that each emit GroupedObservable from the original Observable, organized by key.
	GroupByDynamic(distribution func(reactive.Item) string, opts ...reactive.Option) Stream

	// IgnoreElements ignores all items emitted by the source ObservableSource except for the errors.
	// Cannot be run in parallel.
	IgnoreElements(opts ...reactive.Option) Stream

	// Join combines items emitted by two Observables whenever an item from one Observable is emitted during
	// a time window defined according to an item emitted by the other Observable.
	// The time is extracted using a timeExtractor function.
	Join(joiner reactive.Func2, right reactive.Observable, timeExtractor func(interface{}) time.Time, windowInMS uint32, opts ...reactive.Option) Stream

	// Last returns a new Observable which emit only last item.
	// Cannot be run in parallel.
	Last(opts ...reactive.Option) Stream

	// LastOrDefault returns a new Observable which emit only last item.
	// If the observable fails to emit any items, it emits a default value.
	// Cannot be run in parallel.
	LastOrDefault(defaultValue interface{}, opts ...reactive.Option) Stream

	// Map transforms the items emitted by an Observable by applying a function to each item.
	Map(apply reactive.Func, opts ...reactive.Option) Stream

	// Marshal transforms the items emitted by an Observable by applying a marshalling to each item.
	Marshal(marshaller Marshaller, opts ...reactive.Option) Stream

	// Max determines and emits the maximum-valued item emitted by an Observable according to a comparator.
	Max(comparator reactive.Comparator, opts ...reactive.Option) Stream

	// Min determines and emits the minimum-valued item emitted by an Observable according to a comparator.
	Min(comparator reactive.Comparator, opts ...reactive.Option) Stream

	// OnErrorResumeNext instructs an Observable to pass control to another Observable rather than invoking
	// onError if it encounters an error.
	OnErrorResumeNext(resumeSequence reactive.ErrorToObservable, opts ...reactive.Option) Stream

	// OnErrorReturn instructs an Observable to emit an item (returned by a specified function)
	// rather than invoking onError if it encounters an error.
	OnErrorReturn(resumeFunc reactive.ErrorFunc, opts ...reactive.Option) Stream

	// OnErrorReturnItem instructs on Observable to emit an item if it encounters an error.
	OnErrorReturnItem(resume interface{}, opts ...reactive.Option) Stream

	// Reduce applies a function to each item emitted by an Observable, sequentially, and emit the final value.
	Reduce(apply reactive.Func2, opts ...reactive.Option) Stream

	// Repeat returns an Observable that repeats the sequence of items emitted by the source Observable
	// at most count times, at a particular frequency.
	// Cannot run in parallel.
	Repeat(count int64, milliseconds uint32, opts ...reactive.Option) Stream

	// Retry retries if a source Observable sends an error, resubscribe to it in the hopes that it will complete without error.
	// Cannot be run in parallel.
	Retry(count int, shouldRetry func(error) bool, opts ...reactive.Option) Stream

	// Run creates an Observer without consuming the emitted items.
	Run(opts ...reactive.Option) reactive.Disposed

	// Sample returns an Observable that emits the most recent items emitted by the source
	// Iterable whenever the input Iterable emits an item.
	Sample(iterable reactive.Iterable, opts ...reactive.Option) Stream

	// Scan apply a Func2 to each item emitted by an Observable, sequentially, and emit each successive value.
	// Cannot be run in parallel.
	Scan(apply reactive.Func2, opts ...reactive.Option) Stream

	// SequenceEqual emits true if an Observable and the input Observable emit the same items,
	// in the same order, with the same termination state. Otherwise, it emits false.
	SequenceEqual(iterable reactive.Iterable, opts ...reactive.Option) Stream

	// Send sends the items to a given channel.
	Send(output chan<- reactive.Item, opts ...reactive.Option)

	// Serialize forces an Observable to make serialized calls and to be well-behaved.
	Serialize(from int, identifier func(interface{}) int, opts ...reactive.Option) Stream

	// Skip suppresses the first n items in the original Observable and
	// returns a new Observable with the rest items.
	// Cannot be run in parallel.
	Skip(nth uint, opts ...reactive.Option) Stream

	// SkipLast suppresses the last n items in the original Observable and
	// returns a new Observable with the rest items.
	// Cannot be run in parallel.
	SkipLast(nth uint, opts ...reactive.Option) Stream

	// SkipWhile discard items emitted by an Observable until a specified condition becomes false.
	// Cannot be run in parallel.
	SkipWhile(apply reactive.Predicate, opts ...reactive.Option) Stream

	// StartWith emits a specified Iterable before beginning to emit the items from the source Observable.
	StartWith(iterable reactive.Iterable, opts ...reactive.Option) Stream

	// SumFloat32 calculates the average of float32 emitted by an Observable and emits a float32.
	SumFloat32(opts ...reactive.Option) Stream

	// SumFloat64 calculates the average of float64 emitted by an Observable and emits a float64.
	SumFloat64(opts ...reactive.Option) Stream

	// SumInt64 calculates the average of integers emitted by an Observable and emits an int64.
	SumInt64(opts ...reactive.Option) Stream

	// Take emits only the first n items emitted by an Observable.
	// Cannot be run in parallel.
	Take(nth uint, opts ...reactive.Option) Stream

	// TakeLast emits only the last n items emitted by an Observable.
	// Cannot be run in parallel.
	TakeLast(nth uint, opts ...reactive.Option) Stream

	// TakeUntil returns an Observable that emits items emitted by the source Observable,
	// checks the specified predicate for each item, and then completes when the condition is satisfied.
	// Cannot be run in parallel.
	TakeUntil(apply reactive.Predicate, opts ...reactive.Option) Stream

	// TakeWhile returns an Observable that emits items emitted by the source ObservableSource so long as each
	// item satisfied a specified condition, and then completes as soon as this condition is not satisfied.
	// Cannot be run in parallel.
	TakeWhile(apply reactive.Predicate, opts ...reactive.Option) Stream

	// TimeInterval converts an Observable that emits items into one that emits indications of the amount of time elapsed between those emissions.
	TimeInterval(opts ...reactive.Option) Stream

	// Timestamp attaches a timestamp to each item emitted by an Observable indicating when it was emitted.
	Timestamp(opts ...reactive.Option) Stream

	// ToMap convert the sequence of items emitted by an Observable
	// into a map keyed by a specified key function.
	// Cannot be run in parallel.
	ToMap(keySelector reactive.Func, opts ...reactive.Option) Stream

	// ToMapWithValueSelector convert the sequence of items emitted by an Observable
	// into a map keyed by a specified key function and valued by another
	// value function.
	// Cannot be run in parallel.
	ToMapWithValueSelector(keySelector, valueSelector reactive.Func, opts ...reactive.Option) Stream

	// ToSlice collects all items from an Observable and emit them in a slice and an optional error.
	// Cannot be run in parallel.
	ToSlice(initialCapacity int, opts ...reactive.Option) ([]interface{}, error)

	// Unmarshal transforms the items emitted by an Observable by applying an unmarshalling to each item.
	Unmarshal(unmarshaller Unmarshaller, factory func() interface{}, opts ...reactive.Option) Stream

	// WindowWithCount periodically subdivides items from an Observable into Observable windows of a given size and emit these windows
	// rather than emitting the items one at a time.
	WindowWithCount(count int, opts ...reactive.Option) Stream

	// WindowWithTime periodically subdivides items from an Observable into Observables based on timed windows
	// and emit them rather than emitting the items one at a time.
	WindowWithTime(milliseconds uint32, opts ...reactive.Option) Stream

	// WindowWithTimeOrCount periodically subdivides items from an Observable into Observables based on timed windows or a specific size
	// and emit them rather than emitting the items one at a time.
	WindowWithTimeOrCount(milliseconds uint32, count int, opts ...reactive.Option) Stream

	// ZipFromIterable merges the emissions of an Iterable via a specified function
	// and emit single items for each combination based on the results of this function.
	ZipFromIterable(iterable reactive.Iterable, processor reactive.Func2, opts ...reactive.Option) Stream

	// SlidingWindowWithCount buffers the data in the specified sliding window size, the buffered data can be processed in the handler func.
	// It returns the orginal data to Stream, not the buffered slice.
	SlidingWindowWithCount(windowSize int, slideSize int, handler Handler, opts ...reactive.Option) Stream

	// SlidingWindowWithTime buffers the data in the specified sliding window time in milliseconds, the buffered data can be processed in the handler func.
	// It returns the orginal data to Stream, not the buffered slice.
	SlidingWindowWithTime(windowTimeInMS uint32, slideTimeInMS uint32, handler Handler, opts ...reactive.Option) Stream

	// // ZipMultiObservers subscribes multi Bhojpur Service observers, zips the values into a slice and calls the processor callback when all keys are observed.
	// ZipMultiObservers(observers []KeyObserveFunc, processor func(items []interface{}) (interface{}, error)) Stream
}

// // KeyObserveFunc is a pair of subscribed key and onObserve callback.
// type KeyObserveFunc struct {
// 	Key       byte
// 	OnObserve decoder.OnObserveFunc
// }

type (
	// Marshaller defines a marshaller type (interface{} to []byte).
	Marshaller func(interface{}) ([]byte, error)
	// Unmarshaller defines an unmarshaller type ([]byte to interface).
	Unmarshaller func([]byte, interface{}) error
)
