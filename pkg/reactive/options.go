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
	"runtime"

	"github.com/teivah/onecontext"
)

var emptyContext context.Context

// Option handles configurable options.
type Option interface {
	apply(*funcOption)
	toPropagate() bool
	isEagerObservation() bool
	getPool() (bool, int)
	buildChannel() chan Item
	buildContext(parent context.Context) context.Context
	getBackPressureStrategy() BackpressureStrategy
	getErrorStrategy() OnErrorStrategy
	isConnectable() bool
	isConnectOperation() bool
	isSerialized() (bool, func(interface{}) int)
}

type funcOption struct {
	f                    func(*funcOption)
	isBuffer             bool
	buffer               int
	ctx                  context.Context
	observation          ObservationStrategy
	pool                 int
	backPressureStrategy BackpressureStrategy
	onErrorStrategy      OnErrorStrategy
	propagate            bool
	connectable          bool
	connectOperation     bool
	serialized           func(interface{}) int
}

func (fdo *funcOption) toPropagate() bool {
	return fdo.propagate
}

func (fdo *funcOption) isEagerObservation() bool {
	return fdo.observation == Eager
}

func (fdo *funcOption) getPool() (bool, int) {
	return fdo.pool > 0, fdo.pool
}

func (fdo *funcOption) buildChannel() chan Item {
	if fdo.isBuffer {
		return make(chan Item, fdo.buffer)
	}
	return make(chan Item)
}

func (fdo *funcOption) buildContext(parent context.Context) context.Context {
	if fdo.ctx != nil && parent != nil {
		ctx, _ := onecontext.Merge(fdo.ctx, parent)
		return ctx
	}

	if fdo.ctx != nil {
		return fdo.ctx
	}
	if parent != nil {
		return parent
	}
	return context.Background()
}

func (fdo *funcOption) getBackPressureStrategy() BackpressureStrategy {
	return fdo.backPressureStrategy
}

func (fdo *funcOption) getErrorStrategy() OnErrorStrategy {
	return fdo.onErrorStrategy
}

func (fdo *funcOption) isConnectable() bool {
	return fdo.connectable
}

func (fdo *funcOption) isConnectOperation() bool {
	return fdo.connectOperation
}

func (fdo *funcOption) apply(do *funcOption) {
	fdo.f(do)
}

func (fdo *funcOption) isSerialized() (bool, func(interface{}) int) {
	if fdo.serialized == nil {
		return false, nil
	}
	return true, fdo.serialized
}

func newFuncOption(f func(*funcOption)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func parseOptions(opts ...Option) Option {
	o := new(funcOption)
	for _, opt := range opts {
		opt.apply(o)
	}
	return o
}

// WithBufferedChannel allows to configure the capacity of a buffered channel.
func WithBufferedChannel(capacity int) Option {
	return newFuncOption(func(options *funcOption) {
		options.isBuffer = true
		options.buffer = capacity
	})
}

// WithContext allows to pass a context.
func WithContext(ctx context.Context) Option {
	return newFuncOption(func(options *funcOption) {
		options.ctx = ctx
	})
}

// WithObservationStrategy uses the eager observation mode meaning consuming the items even without subscription.
func WithObservationStrategy(strategy ObservationStrategy) Option {
	return newFuncOption(func(options *funcOption) {
		options.observation = strategy
	})
}

// WithPool allows to specify an execution pool.
func WithPool(pool int) Option {
	return newFuncOption(func(options *funcOption) {
		options.pool = pool
	})
}

// WithCPUPool allows to specify an execution pool based on the number of logical CPUs.
func WithCPUPool() Option {
	return newFuncOption(func(options *funcOption) {
		options.pool = runtime.NumCPU()
	})
}

// WithBackPressureStrategy sets the back pressure strategy: drop or block.
func WithBackPressureStrategy(strategy BackpressureStrategy) Option {
	return newFuncOption(func(options *funcOption) {
		options.backPressureStrategy = strategy
	})
}

// WithErrorStrategy defines how an observable should deal with error.
// This strategy is propagated to the parent observable.
func WithErrorStrategy(strategy OnErrorStrategy) Option {
	return newFuncOption(func(options *funcOption) {
		options.onErrorStrategy = strategy
	})
}

// WithPublishStrategy converts an ordinary Observable into a connectable Observable.
func WithPublishStrategy() Option {
	return newFuncOption(func(options *funcOption) {
		options.connectable = true
	})
}

// Serialize forces an Observable to make serialized calls and to be well-behaved.
func Serialize(identifier func(interface{}) int) Option {
	return newFuncOption(func(options *funcOption) {
		options.serialized = identifier
	})
}

func connect() Option {
	return newFuncOption(func(options *funcOption) {
		options.connectOperation = true
	})
}
