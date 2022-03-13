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

	svcsvr "github.com/bhojpur/service/pkg/engine"
	"github.com/bhojpur/service/pkg/engine/core/frame"
	"github.com/bhojpur/service/pkg/engine/logger"
)

// Runtime is the Reactive Stream serverless runtime engine.
type Runtime struct {
	rawBytesChan chan interface{}
	sfn          svcsvr.StreamFunction
	stream       Stream
}

// NewRuntime creates a new Reactive Stream serverless runtime engine instance.
func NewRuntime(sfn svcsvr.StreamFunction) *Runtime {
	return &Runtime{
		rawBytesChan: make(chan interface{}),
		sfn:          sfn,
	}
}

// RawByteHandler is the Stream Handler for RawBytes.
func (r *Runtime) RawByteHandler(req []byte) (byte, []byte) {
	go func() {
		r.rawBytesChan <- req
	}()

	// observe the data from Reactive Stream.
	for item := range r.stream.Observe() {
		if item.Error() {
			logger.Errorf("[Reactive Handler] Handler got an error, err=%v", item.E)
			continue
		}

		if item.V == nil {
			logger.Warnf("[Reactive Handler] the returned data is nil.")
			continue
		}

		res, ok := (item.V).(frame.PayloadFrame)
		if !ok {
			logger.Warnf("[Reactive Handler] the data is not a frame.PayloadFrame, won't send it to Bhojpur Service-Processor.")
			continue
		}

		logger.Infof("[RawByteHandler] Send data with [tag=%#x] to Bhojpur Service-Processor.", res.Tag)
		return res.Tag, res.Carriage
	}

	// return empty data by default, the new data from Reactive Stream will be returned in `Pipe` function.
	return 0x0, nil
}

// PipeHandler processes data sequentially.
func (r *Runtime) PipeHandler(in <-chan []byte, out chan<- *frame.PayloadFrame) {
	for {
		select {
		case req := <-in:
			r.rawBytesChan <- req
		case item := <-r.stream.Observe():
			if item.Error() {
				logger.Errorf("[Reactive PipeHandler] Handler got an error, err=%v", item.E)
				continue
			}

			if item.V == nil {
				logger.Warnf("[Reactive PipeHandler] the returned data is nil.")
				continue
			}

			res, ok := (item.V).(frame.PayloadFrame)
			if !ok {
				logger.Warnf("[Reactive PipeHandler] the data is not a frame.PayloadFrame, won't send it to Bhojpur Service-Processor.")
				continue
			}

			logger.Infof("[Reactive PipeHandler] Send data with [tag=%#x] to Bhojpur Service-Processor.", res.Tag)
			out <- &res
		}
	}
}

// Pipe the Reactive Handler with Reactive Stream.
func (r *Runtime) Pipe(rxHandler func(rxstream Stream) Stream) {
	fac := NewFactory()
	// create a Reactive Stream from raw bytes channel.
	rxstream := fac.FromChannel(context.Background(), r.rawBytesChan)

	// run Reactive Handler and get a new Reactive Stream.
	r.stream = rxHandler(rxstream)
}
