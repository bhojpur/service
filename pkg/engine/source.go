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

	"github.com/bhojpur/service/pkg/engine/core"
	"github.com/bhojpur/service/pkg/engine/core/frame"
)

const (
	sourceLogPrefix = "\033[32m[bhojpur:source]\033[0m "
)

// Source is responsible for sending data to Bhojpur Service.
type Source interface {
	// Close will close the connection to Bhojpur Service-Processor.
	Close() error
	// Connect to Bhojpur Service-Processor.
	Connect() error
	// SetDataTag will set the tag of data when invoking Write().
	SetDataTag(tag uint8)
	// Write the data to downstream.
	Write(p []byte) (n int, err error)
	// WriteWithTag will write data with specified tag, default transactionID is epoch time.
	WriteWithTag(tag uint8, data []byte) error
}

// Bhojpur Service Data-Source
type dataSource struct {
	name              string
	processorEndpoint string
	client            *core.Client
	tag               uint8
}

var _ Source = &dataSource{}

// NewSource create a Bhojpur Service Data-Source
func NewSource(name string, opts ...Option) Source {
	options := NewOptions(opts...)
	client := core.NewClient(name, core.ClientTypeSource, options.ClientOptions...)

	return &dataSource{
		name:              name,
		processorEndpoint: options.ProcessorAddr,
		client:            client,
	}
}

// Write the data to downstream.
func (s *dataSource) Write(data []byte) (int, error) {
	return len(data), s.WriteWithTag(s.tag, data)
}

// SetDataTag will set the tag of data when invoking Write().
func (s *dataSource) SetDataTag(tag uint8) {
	s.tag = tag
}

// Close will close the connection to Bhojpur Service-Processor.
func (s *dataSource) Close() error {
	if err := s.client.Close(); err != nil {
		s.client.Logger().Errorf("%sClose(): %v", sourceLogPrefix, err)
		return err
	}
	s.client.Logger().Debugf("%s is closed", sourceLogPrefix)
	return nil
}

// Connect to Bhojpur Service-Processor.
func (s *dataSource) Connect() error {
	err := s.client.Connect(context.Background(), s.processorEndpoint)
	if err != nil {
		s.client.Logger().Errorf("%sConnect() error: %s", sourceLogPrefix, err)
	}
	return err
}

// WriteWithTag will write data with specified tag, default transactionID is epoch time.
func (s *dataSource) WriteWithTag(tag uint8, data []byte) error {
	s.client.Logger().Debugf("%sWriteWithTag: len(data)=%d, data=%# x", sourceLogPrefix, len(data), frame.Shortly(data))
	frame := frame.NewDataFrame()
	frame.SetCarriage(byte(tag), data)
	return s.client.WriteFrame(frame)
}
