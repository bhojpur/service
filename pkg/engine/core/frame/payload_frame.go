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

// PayloadFrame is a Bhojpur Service encoded bytes, Tag is a fixed value TYPE_ID_PAYLOAD_FRAME
// the Len is the length of Val. Val is also a Bhojpur Service encoded PrimitivePacket, storing
// raw bytes as user's data
type PayloadFrame struct {
	Tag      byte
	Carriage []byte
}

// NewPayloadFrame creates a new PayloadFrame with a given TagID of user's data
func NewPayloadFrame(tag byte) *PayloadFrame {
	return &PayloadFrame{
		Tag: tag,
	}
}

// SetCarriage sets the user's raw data
func (m *PayloadFrame) SetCarriage(buf []byte) *PayloadFrame {
	m.Carriage = buf
	return m
}

// Encode to Bhojpur Service encoded bytes
func (m *PayloadFrame) Encode() []byte {
	carriage := codec.NewPrimitivePacketEncoder(int(m.Tag))
	carriage.SetBytesValue(m.Carriage)

	payload := codec.NewNodePacketEncoder(int(byte(TagOfPayloadFrame)))
	payload.AddPrimitivePacket(carriage)

	return payload.Encode()
}

// DecodeToPayloadFrame decodes Bhojpur Service encoded bytes to PayloadFrame
func DecodeToPayloadFrame(buf []byte) (*PayloadFrame, error) {
	//nodeBlock := codec.NodePacket{}
	nodeBlock, _, err := codec.DecodeNodePacket(buf)
	if err != nil {
		return nil, err
	}

	payload := &PayloadFrame{}
	for _, v := range nodeBlock.PrimitivePackets {
		payload.Tag = v.SeqID()
		payload.Carriage = v.ToBytes()
		break
	}

	return payload, nil
}
