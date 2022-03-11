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

import "context"

type (
	operatorOptions struct {
		stop          func()
		resetIterable func(Iterable)
	}

	// Comparator defines a func that returns an int:
	// - 0 if two elements are equals
	// - A negative value if the first argument is less than the second
	// - A positive value if the first argument is greater than the second
	Comparator func(interface{}, interface{}) int
	// ItemToObservable defines a function that computes an observable from an item.
	ItemToObservable func(Item) Observable
	// ErrorToObservable defines a function that transforms an observable from an error.
	ErrorToObservable func(error) Observable
	// Func defines a function that computes a value from an input value.
	Func func(context.Context, interface{}) (interface{}, error)
	// Func2 defines a function that computes a value from two input values.
	Func2 func(context.Context, interface{}, interface{}) (interface{}, error)
	// FuncN defines a function that computes a value from N input values.
	FuncN func(...interface{}) interface{}
	// ErrorFunc defines a function that computes a value from an error.
	ErrorFunc func(error) interface{}
	// Predicate defines a func that returns a bool from an input value.
	Predicate func(interface{}) bool
	// Marshaller defines a marshaller type (interface{} to []byte).
	Marshaller func(interface{}) ([]byte, error)
	// Unmarshaller defines an unmarshaller type ([]byte to interface).
	Unmarshaller func([]byte, interface{}) error
	// Producer defines a producer implementation.
	Producer func(ctx context.Context, next chan<- Item)
	// Supplier defines a function that supplies a result from nothing.
	Supplier func(ctx context.Context) Item
	// Disposed is a notification channel indicating when an Observable is closed.
	Disposed <-chan struct{}
	// Disposable is a function to be called in order to dispose a subscription.
	Disposable context.CancelFunc

	// NextFunc handles a next item in a stream.
	NextFunc func(interface{})
	// ErrFunc handles an error in a stream.
	ErrFunc func(error)
	// CompletedFunc handles the end of a stream.
	CompletedFunc func()
)

// BackpressureStrategy is the backpressure strategy type.
type BackpressureStrategy uint32

const (
	// Block blocks until the channel is available.
	Block BackpressureStrategy = iota
	// Drop drops the message.
	Drop
)

// OnErrorStrategy is the Observable error strategy.
type OnErrorStrategy uint32

const (
	// StopOnError is the default error strategy.
	// An operator will stop processing items on error.
	StopOnError OnErrorStrategy = iota
	// ContinueOnError means an operator will continue processing items after an error.
	ContinueOnError
)

// ObservationStrategy defines the strategy to consume from an Observable.
type ObservationStrategy uint32

const (
	// Lazy is the default observation strategy, when an Observer subscribes.
	Lazy ObservationStrategy = iota
	// Eager means consuming as soon as the Observable is created.
	Eager
)
