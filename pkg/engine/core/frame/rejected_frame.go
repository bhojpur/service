package frame

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
	"github.com/bhojpur/service/pkg/engine/codec"
)

// RejectedFrame is a Bhojpur Service encoded bytes, Tag is a fixed value TYPE_ID_REJECTED_FRAME
type RejectedFrame struct{}

// NewRejectedFrame creates a new RejectedFrame with a given TagID of user's data
func NewRejectedFrame() *RejectedFrame {
	return &RejectedFrame{}
}

// Type gets the type of Frame.
func (m *RejectedFrame) Type() Type {
	return TagOfRejectedFrame
}

// Encode to Bhojpur Service encoded bytes
func (m *RejectedFrame) Encode() []byte {
	rejected := codec.NewNodePacketEncoder(int(byte(m.Type())))
	rejected.AddBytes(nil)

	return rejected.Encode()
}

// DecodeToRejectedFrame decodes Bhojpur Service encoded bytes to RejectedFrame
func DecodeToRejectedFrame(buf []byte) (*RejectedFrame, error) {
	//nodeBlock := codec.NodePacket{}
	_, _, err := codec.DecodeNodePacket(buf)
	if err != nil {
		return nil, err
	}
	return &RejectedFrame{}, nil
}
