package codec

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
	"fmt"
	"time"
)

var (
	// enabledTestPrintf set whether to print debug information in the test
	enabledTestPrintf = false
)

// observableTester use Observable to listen the node
type observableTester struct {
	observe       byte
	sourceChannel chan interface{}
	source        Observable
}

// newObservableTester creat a observableTester
func newObservableTester(observe byte) *observableTester {
	return &observableTester{observe: observe}
}

// testDecoder is a shortcut to perform decoding tests
func testDecoder(observe byte, buf []byte, callback func(v []byte) (interface{}, error)) {
	newObservableTester(observe).
		Init(callback).
		Write(buf).
		CloseWith(150)
}

// Init create a channel for testing
func (t *observableTester) Init(callback func(v []byte) (interface{}, error)) *observableTester {
	t.sourceChannel = make(chan interface{})
	subscribers := make([]chan interface{}, 0)

	t.source = &observableImpl{iterable: &iterableImpl{next: t.sourceChannel, subscribers: subscribers}}

	consumer := t.source.Subscribe(t.observe).OnObserve(callback)

	go func() {
		for c := range consumer {
			if c != 0 {
				//TODO: Why empty branch?
				testPrintf("TODO: Empty branch reached\n")
			}
		}
	}()

	return t
}

// Write is used to write data to the Channel
func (t *observableTester) Write(buf []byte) *observableTester {
	t.sourceChannel <- buf
	return t
}

// Close is used to close the Channel
func (t *observableTester) Close() {
	close(t.sourceChannel)
}

// Close is used to close the Channel with waiting time
func (t *observableTester) CloseWith(millisecond int64) {
	time.Sleep(time.Duration(millisecond) * time.Millisecond)
	close(t.sourceChannel)
}

// testPrintf print debug output in test cases
func testPrintf(format string, a ...interface{}) {
	if enabledTestPrintf {
		fmt.Printf(format, a...)
	}
}
