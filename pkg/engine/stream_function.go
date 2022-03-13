package engine

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

	engine "github.com/bhojpur/service/pkg/engine/core"
	"github.com/bhojpur/service/pkg/engine/core/frame"
)

const (
	streamFunctionLogPrefix = "\033[31m[bhojpur:streamfunc]\033[0m "
)

// StreamFunction defines serverless streaming functions.
type StreamFunction interface {
	// SetObserveDataTags set the data tag list that will be observed
	// Deprecated: use bhojpur.WithObserveDataTags instead
	SetObserveDataTags(tag ...byte)
	// SetHandler set the handler function, which accept the raw bytes data and return the tag & response
	SetHandler(fn engine.AsyncHandler) error
	// SetPipeHandler set the pipe handler function
	SetPipeHandler(fn engine.PipeHandler) error
	// Connect create a connection to the Processor
	Connect() error
	// Close will close the connection
	Close() error
	// Send a data to Processor.
	Write(tag byte, carriage []byte) error
}

// NewStreamFunction create a stream function.
func NewStreamFunction(name string, opts ...Option) StreamFunction {
	options := NewOptions(opts...)
	client := engine.NewClient(name, engine.ClientTypeStreamFunction, options.ClientOptions...)
	sfn := &streamFunction{
		name:              name,
		processorEndpoint: options.ProcessorAddr,
		client:            client,
		observeDataTags:   make([]byte, 0),
	}

	return sfn
}

var _ StreamFunction = &streamFunction{}

// streamFunction implements StreamFunction interface.
type streamFunction struct {
	name              string
	processorEndpoint string
	client            *engine.Client
	observeDataTags   []byte              // tag list that will be observed
	fn                engine.AsyncHandler // user's function which will be invoked when data arrived
	pfn               engine.PipeHandler
	pIn               chan []byte
	pOut              chan *frame.PayloadFrame
}

// SetObserveDataTags set the data tag list that will be observed.
// Deprecated: use bhojpur.WithObserveDataTags instead
func (s *streamFunction) SetObserveDataTags(tag ...byte) {
	s.client.SetObserveDataTags(tag...)
	s.client.Logger().Debugf("%sSetObserveDataTag(%v)", streamFunctionLogPrefix, s.observeDataTags)
}

// SetHandler set the handler function, which accept the raw bytes data and return the tag & response.
func (s *streamFunction) SetHandler(fn engine.AsyncHandler) error {
	s.fn = fn
	s.client.Logger().Debugf("%sSetHandler(%v)", streamFunctionLogPrefix, s.fn)
	return nil
}

func (s *streamFunction) SetPipeHandler(fn engine.PipeHandler) error {
	s.pfn = fn
	s.client.Logger().Debugf("%sSetHandler(%v)", streamFunctionLogPrefix, s.fn)
	return nil
}

// Connect create a connection to the Processor, when data arrvied, the data will be passed to the
// handler which setted by SetHandler method.
func (s *streamFunction) Connect() error {
	s.client.Logger().Debugf("%s Connect()", streamFunctionLogPrefix)
	// notify underlying network operations, when data with tag we observed arrived, invoke the func
	s.client.SetDataFrameObserver(func(data *frame.DataFrame) {
		s.client.Logger().Debugf("%sreceive DataFrame, tag=%# x, carraige=%# x", streamFunctionLogPrefix, data.Tag(), data.GetCarriage())
		s.onDataFrame(data.GetCarriage(), data.GetMetaFrame())
	})

	if s.pfn != nil {
		s.pIn = make(chan []byte)
		s.pOut = make(chan *frame.PayloadFrame)

		// handle user's pipe function
		go func() {
			s.pfn(s.pIn, s.pOut)
		}()

		// send user's pipe function outputs to Processor
		go func() {
			for {
				data := <-s.pOut
				if data != nil {
					s.client.Logger().Debugf("%spipe function send: tag=%#x, data=%# x", streamFunctionLogPrefix, data.Tag, data.Carriage)
					frame := frame.NewDataFrame()
					// todo: frame.SetTransactionID
					frame.SetCarriage(data.Tag, data.Carriage)
					s.client.WriteFrame(frame)
				}
			}
		}()
	}

	err := s.client.Connect(context.Background(), s.processorEndpoint)
	if err != nil {
		s.client.Logger().Errorf("%sConnect() error: %s", streamFunctionLogPrefix, err)
	}
	return err
}

// Close will close the connection.
func (s *streamFunction) Close() error {
	if s.pIn != nil {
		close(s.pIn)
	}

	if s.pOut != nil {
		close(s.pOut)
	}

	if s.client != nil {
		if err := s.client.Close(); err != nil {
			s.client.Logger().Errorf("%sClose(): %v", err)
			return err
		}
	}

	return nil
}

// when DataFrame we observed arrived, invoke the user's function
func (s *streamFunction) onDataFrame(data []byte, metaFrame *frame.MetaFrame) {
	s.client.Logger().Infof("%sonDataFrame ->[%s]", streamFunctionLogPrefix, s.name)

	if s.fn != nil {
		go func() {
			s.client.Logger().Debugf("%sexecute-start function: data[%d]=%# x", streamFunctionLogPrefix, len(data), frame.Shortly(data))
			// invoke serverless
			tag, resp := s.fn(data)
			s.client.Logger().Debugf("%sexecute-done function: tag=%#x, resp[%d]=%# x", streamFunctionLogPrefix, tag, len(resp), frame.Shortly(resp))
			// if resp is not nil, means the user's function has returned something, we should send it to the Processor
			if len(resp) != 0 {
				s.client.Logger().Debugf("%sstart WriteFrame(): tag=%#x, data[%d]=%# x", streamFunctionLogPrefix, tag, len(resp), frame.Shortly(resp))
				// build a DataFrame
				// TODO: seems we should implement a DeepCopy() of MetaFrame in the future
				frame := frame.NewDataFrame()
				// reuse transactionID
				frame.SetTransactionID(metaFrame.TransactionID())
				frame.SetCarriage(tag, resp)
				s.client.WriteFrame(frame)
			}
		}()
	} else if s.pfn != nil {
		s.client.Logger().Debugf("%spipe function receive: data[%d]=%# x", streamFunctionLogPrefix, len(data), data)
		s.pIn <- data
	} else {
		s.client.Logger().Warnf("%sStreamFunction is nil", streamFunctionLogPrefix)
	}
}

// Send a DataFrame to the Processor.
func (s *streamFunction) Write(tag byte, carriage []byte) error {
	frame := frame.NewDataFrame()
	frame.SetCarriage(tag, carriage)
	return s.client.WriteFrame(frame)
}
