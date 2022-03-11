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

import "github.com/bhojpur/service/pkg/engine/codec"

// AcceptedFrame is a Bhojpur Service encoded bytes, Tag is a fixed value TYPE_ID_ACCEPTED_FRAME
type AcceptedFrame struct{}

// NewAcceptedFrame creates a new AcceptedFrame with a given TagID of user's data
func NewAcceptedFrame() *AcceptedFrame {
	return &AcceptedFrame{}
}

// Type gets the type of Frame.
func (m *AcceptedFrame) Type() Type {
	return TagOfAcceptedFrame
}

// Encode to Bhojpur Service encoded bytes.
func (m *AcceptedFrame) Encode() []byte {
	accepted := codec.NewNodePacketEncoder(int(byte(m.Type())))
	accepted.AddBytes(nil)

	return accepted.Encode()
}

// DecodeToAcceptedFrame decodes Bhojpur Service encoded bytes to AcceptedFrame.
func DecodeToAcceptedFrame(buf []byte) (*AcceptedFrame, error) {
	//nodeBlock := codec.NodePacket{}
	_, _, err := codec.DecodeNodePacket(buf)
	if err != nil {
		return nil, err
	}
	return &AcceptedFrame{}, nil
}
