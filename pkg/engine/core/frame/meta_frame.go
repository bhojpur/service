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
	"strconv"
	"time"

	"github.com/bhojpur/service/pkg/engine/codec"
)

// MetaFrame is a Bhojpur Service encoded bytes, SeqID is a fixed value of TYPE_ID_TRANSACTION.
// used for describes metadata for a DataFrame.
type MetaFrame struct {
	tid string
}

// NewMetaFrame creates a new MetaFrame instance.
func NewMetaFrame() *MetaFrame {
	return &MetaFrame{
		tid: strconv.FormatInt(time.Now().Unix(), 10),
	}
}

// SetTransactinID set the transaction ID.
func (m *MetaFrame) SetTransactionID(transactionID string) {
	m.tid = transactionID
}

// TransactionID returns transactionID
func (m *MetaFrame) TransactionID() string {
	return m.tid
}

// Encode implements Frame.Encode method.
func (m *MetaFrame) Encode() []byte {
	meta := codec.NewNodePacketEncoder(int(byte(TagOfMetaFrame)))

	transactionID := codec.NewPrimitivePacketEncoder(int(byte(TagOfTransactionID)))
	transactionID.SetStringValue(m.tid)

	meta.AddPrimitivePacket(transactionID)
	return meta.Encode()
}

// DecodeToMetaFrame decode a MetaFrame instance from given buffer.
func DecodeToMetaFrame(buf []byte) (*MetaFrame, error) {
	//nodeBlock := codec.NodePacket{}
	nodeBlock, _, err := codec.DecodeNodePacket(buf)
	if err != nil {
		return nil, err
	}

	meta := &MetaFrame{}
	for _, v := range nodeBlock.PrimitivePackets {
		val, _ := v.ToUTF8String()
		meta.tid = val
		break
	}

	return meta, nil
}
