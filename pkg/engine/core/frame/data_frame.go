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

// DataFrame defines the data structure carried with user's data
// transferring within Bhojpur Service
type DataFrame struct {
	metaFrame    *MetaFrame
	payloadFrame *PayloadFrame
}

// NewDataFrame create `DataFrame` with a transactionID string,
// consider change transactionID to UUID type later
func NewDataFrame() *DataFrame {
	data := &DataFrame{
		metaFrame: NewMetaFrame(),
	}
	return data
}

// Type gets the type of Frame.
func (d *DataFrame) Type() Type {
	return TagOfDataFrame
}

// Tag return the tag of carriage data.
func (d *DataFrame) Tag() byte {
	return d.payloadFrame.Tag
}

// SetCarriage set user's raw data in `DataFrame`
func (d *DataFrame) SetCarriage(tag byte, carriage []byte) {
	d.payloadFrame = NewPayloadFrame(tag).SetCarriage(carriage)
}

// GetCarriage return user's raw data in `DataFrame`
func (d *DataFrame) GetCarriage() []byte {
	return d.payloadFrame.Carriage
}

// TransactionID return transactionID string
func (d *DataFrame) TransactionID() string {
	return d.metaFrame.TransactionID()
}

// SetTransactionID set transactionID string
func (d *DataFrame) SetTransactionID(transactionID string) {
	d.metaFrame.SetTransactionID(transactionID)
}

// GetMetaFrame return MetaFrame.
func (d *DataFrame) GetMetaFrame() *MetaFrame {
	return d.metaFrame
}

// GetDataTag return the Tag of user's data
func (d *DataFrame) GetDataTag() byte {
	return d.payloadFrame.Tag
}

// Encode return Bhojpur Service encoded bytes of `DataFrame`
func (d *DataFrame) Encode() []byte {
	data := codec.NewNodePacketEncoder(int(byte(d.Type())))
	// MetaFrame
	data.AddBytes(d.metaFrame.Encode())
	// PayloadFrame
	data.AddBytes(d.payloadFrame.Encode())

	return data.Encode()
}

// DecodeToDataFrame decode Bhojpur Service encoded bytes to `DataFrame`
func DecodeToDataFrame(buf []byte) (*DataFrame, error) {
	//packet := codec.NodePacket{}
	packet, _, err := codec.DecodeNodePacket(buf)
	if err != nil {
		return nil, err
	}

	data := &DataFrame{}

	metaBlock := packet.NodePackets[int(byte(TagOfMetaFrame))]
	meta, err := DecodeToMetaFrame(metaBlock.GetValBuf())
	if err != nil {
		return nil, err
	}
	data.metaFrame = meta

	payloadBlock := packet.NodePackets[int(byte(TagOfPayloadFrame))]
	payload, err := DecodeToPayloadFrame(payloadBlock.GetValBuf())
	if err != nil {
		return nil, err
	}
	data.payloadFrame = payload

	return data, nil
}
