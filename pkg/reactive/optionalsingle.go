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

// OptionalSingleEmpty is the constant returned when an OptionalSingle is empty.
var OptionalSingleEmpty = Item{}

// OptionalSingle is an optional single.
type OptionalSingle interface {
	Iterable
	Get(opts ...Option) (Item, error)
	Map(apply Func, opts ...Option) OptionalSingle
	Run(opts ...Option) Disposed
}

// OptionalSingleImpl implements OptionalSingle.
type OptionalSingleImpl struct {
	parent   context.Context
	iterable Iterable
}

// Get returns the item or reactive.OptionalEmpty. The error returned is if the context has been cancelled.
// This method is blocking.
func (o *OptionalSingleImpl) Get(opts ...Option) (Item, error) {
	option := parseOptions(opts...)
	ctx := option.buildContext(o.parent)

	observe := o.Observe(opts...)
	for {
		select {
		case <-ctx.Done():
			return Item{}, ctx.Err()
		case v, ok := <-observe:
			if !ok {
				return OptionalSingleEmpty, nil
			}
			return v, nil
		}
	}
}

// Map transforms the items emitted by an OptionalSingle by applying a function to each item.
func (o *OptionalSingleImpl) Map(apply Func, opts ...Option) OptionalSingle {
	return optionalSingle(o.parent, o, func() operator {
		return &mapOperatorOptionalSingle{apply: apply}
	}, false, true, opts...)
}

// Observe observes an OptionalSingle by returning its channel.
func (o *OptionalSingleImpl) Observe(opts ...Option) <-chan Item {
	return o.iterable.Observe(opts...)
}

type mapOperatorOptionalSingle struct {
	apply Func
}

func (op *mapOperatorOptionalSingle) next(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	res, err := op.apply(ctx, item.V)
	if err != nil {
		dst <- Error(err)
		operatorOptions.stop()
		return
	}
	dst <- Of(res)
}

func (op *mapOperatorOptionalSingle) err(ctx context.Context, item Item, dst chan<- Item, operatorOptions operatorOptions) {
	defaultErrorFuncOperator(ctx, item, dst, operatorOptions)
}

func (op *mapOperatorOptionalSingle) end(_ context.Context, _ chan<- Item) {
}

func (op *mapOperatorOptionalSingle) gatherNext(_ context.Context, item Item, dst chan<- Item, _ operatorOptions) {
	switch item.V.(type) {
	case *mapOperatorOptionalSingle:
		return
	}
	dst <- item
}

// Run creates an observer without consuming the emitted items.
func (o *OptionalSingleImpl) Run(opts ...Option) Disposed {
	dispose := make(chan struct{})
	option := parseOptions(opts...)
	ctx := option.buildContext(o.parent)

	go func() {
		defer close(dispose)
		observe := o.Observe(opts...)
		for {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-observe:
				if !ok {
					return
				}
			}
		}
	}()

	return dispose
}
