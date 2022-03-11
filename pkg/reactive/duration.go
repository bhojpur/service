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
	"time"

	"github.com/stretchr/testify/mock"
)

// Infinite represents an infinite wait time
var Infinite int64 = -1

// Duration represents a duration
type Duration interface {
	duration() time.Duration
}

type duration struct {
	d time.Duration
}

func (d *duration) duration() time.Duration {
	return d.d
}

// WithDuration is a duration option
func WithDuration(d time.Duration) Duration {
	return &duration{
		d: d,
	}
}

var tick = struct{}{}

type causalityDuration struct {
	fs []execution
}

type execution struct {
	f      func()
	isTick bool
}

func timeCausality(elems ...interface{}) (context.Context, Observable, Duration) {
	ch := make(chan Item, 1)
	fs := make([]execution, len(elems)+1)
	ctx, cancel := context.WithCancel(context.Background())
	for i, elem := range elems {
		i := i
		elem := elem
		if elem == tick {
			fs[i] = execution{
				f:      func() {},
				isTick: true,
			}
		} else {
			switch elem := elem.(type) {
			default:
				fs[i] = execution{
					f: func() {
						ch <- Of(elem)
					},
					isTick: false,
				}
			case error:
				fs[i] = execution{
					f: func() {
						ch <- Error(elem)
					},
					isTick: false,
				}
			}
		}
	}
	fs[len(elems)] = execution{
		f: func() {
			cancel()
		},
		isTick: false,
	}
	return ctx, FromChannel(ch), &causalityDuration{fs: fs}
}

func (d *causalityDuration) duration() time.Duration {
	pop := d.fs[0]
	pop.f()
	d.fs = d.fs[1:]
	if pop.isTick {
		return time.Nanosecond
	}
	return time.Minute
}

type mockDuration struct {
	mock.Mock
}

func (m *mockDuration) duration() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}
