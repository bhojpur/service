package core

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
	"errors"
	"io"

	"github.com/bhojpur/service/pkg/engine/core/frame"
)

// FrameStream is the QUIC Stream with the minimum unit Frame.
type FrameStream struct {
	// Stream is a QUIC stream.
	stream io.ReadWriter
}

// NewFrameStream creates a new FrameStream.
func NewFrameStream(s io.ReadWriter) *FrameStream {
	return &FrameStream{
		stream: s,
	}
}

// ReadFrame reads next frame from QUIC stream.
func (fs *FrameStream) ReadFrame() (frame.Frame, error) {
	if fs.stream == nil {
		return nil, errors.New("core.ReadStream: stream can not be nil")
	}
	return ParseFrame(fs.stream)
}

// WriteFrame writes a frame into QUIC stream.
func (fs *FrameStream) WriteFrame(f frame.Frame) (int, error) {
	if fs.stream == nil {
		return 0, errors.New("core.WriteFrame: stream can not be nil")
	}
	return fs.stream.Write(f.Encode())
}
