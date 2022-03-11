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

	"github.com/bhojpur/service/pkg/utils/encoding"
)

// the minimal length of a packet is 2 bytes
const primitivePacketBufferMinimalLength = 2

// PrimitivePacket，Codec，TLV
type PrimitivePacket basePacket

// SeqID returns the key of primitive packet
func (p *PrimitivePacket) SeqID() byte {
	return p.tag.SeqID()
}

// String prints debug info
func (p *PrimitivePacket) String() string {
	return fmt.Sprintf("Tag=%#x, Length=%v, RawDataLength=%v, Raw=[%#x]", p.tag, p.length, len(p.valBuf), p.valBuf)
}

// ToInt32 parse raw as int32 value
func (p *PrimitivePacket) ToInt32() (int32, error) {
	var val int32
	codec := encoding.VarCodec{}
	err := codec.DecodePVarInt32(p.valBuf, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// ToUInt32 parse raw as int32 value
func (p *PrimitivePacket) ToUInt32() (uint32, error) {
	var val uint32
	codec := encoding.VarCodec{}
	err := codec.DecodePVarUInt32(p.valBuf, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// ToInt64 parse raw as int32 value
func (p *PrimitivePacket) ToInt64() (int64, error) {
	var val int64
	codec := encoding.VarCodec{}
	err := codec.DecodePVarInt64(p.valBuf, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// ToUInt64 parse raw as uint64 value
func (p *PrimitivePacket) ToUInt64() (uint64, error) {
	var val uint64
	codec := encoding.VarCodec{}
	err := codec.DecodePVarUInt64(p.valBuf, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// ToFloat32 parse raw as float32 value
func (p *PrimitivePacket) ToFloat32() (float32, error) {
	var val float32
	codec := encoding.VarCodec{Size: len(p.valBuf)}
	err := codec.DecodeVarFloat32(p.valBuf, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// ToFloat64 parse raw as float64 value
func (p *PrimitivePacket) ToFloat64() (float64, error) {
	var val float64
	codec := encoding.VarCodec{Size: len(p.valBuf)}
	err := codec.DecodeVarFloat64(p.valBuf, &val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// ToBool parse raw as bool value
func (p *PrimitivePacket) ToBool() (bool, error) {
	var val bool
	codec := encoding.VarCodec{Size: len(p.valBuf)}
	err := codec.DecodePVarBool(p.valBuf, &val)
	if err != nil {
		return false, err
	}
	return val, nil
}

// ToUTF8String parse raw data as string value
func (p *PrimitivePacket) ToUTF8String() (string, error) {
	return string(p.valBuf), nil
}

// ToBytes returns raw buffer data
func (p *PrimitivePacket) ToBytes() []byte {
	return p.valBuf
}
