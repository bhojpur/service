package encoding

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
	"math"
	mbit "math/bits"
)

// SizeOfVarFloat32 return the buffer size after encoding value as VarFloat32
func SizeOfVarFloat32(value float32) int {
	return sizeOfVarFloat(uint64(math.Float32bits(value)), 4)
}

// EncodeVarFloat32 encode value as VarFloat32 to buffer
func (codec *VarCodec) EncodeVarFloat32(buffer []byte, value float32) error {
	return codec.encodeVarFloat(buffer, uint64(math.Float32bits(value)), 4)
}

// DecodeVarFloat32 decode to value as VarFloat32 from buffer
func (codec *VarCodec) DecodeVarFloat32(buffer []byte, value *float32) error {
	var bits = uint64(math.Float32bits(*value))
	var err = codec.decodeVarFloat(buffer, &bits, 4)
	*value = math.Float32frombits(uint32(bits))
	return err
}

// SizeOfVarFloat64 return the buffer size after encoding value as VarFloat32
func SizeOfVarFloat64(value float64) int {
	return sizeOfVarFloat(math.Float64bits(value), 8)
}

// EncodeVarFloat64 encode value as VarFloat64 to buffer
func (codec *VarCodec) EncodeVarFloat64(buffer []byte, value float64) error {
	return codec.encodeVarFloat(buffer, math.Float64bits(value), 8)
}

// DecodeVarFloat64 decode to value as VarFloat64 from buffer
func (codec *VarCodec) DecodeVarFloat64(buffer []byte, value *float64) error {
	var bits = math.Float64bits(*value)
	var err = codec.decodeVarFloat(buffer, &bits, 8)
	*value = math.Float64frombits(bits)
	return err
}

func sizeOfVarFloat(bits uint64, width int) int {
	const unit = 8            // bit width of encoding unit
	const mask = uint64(0xFF) // mask of encoding unit

	for s := 0; width > 1; s += unit {
		if bits&(mask<<s) != 0 {
			return width
		}
		width--
	}
	return 1
}

func (codec *VarCodec) encodeVarFloat(buffer []byte, bits uint64, width int) error {
	if codec == nil || codec.Size == 0 {
		return errors.New("nothing to encode")
	}

	const unit = 8 // bit width of encoding unit
	var gap, mask = codec.sizeOfGap(width)

	for (codec.Size & mask) > 0 {
		if codec.Ptr >= len(buffer) {
			return ErrBufferInsufficient
		}
		codec.Size--
		buffer[codec.Ptr] = byte(bits >> ((codec.Size&mask + gap) * unit))
		codec.Ptr++
	}

	codec.Size = 0
	return nil
}

func (codec *VarCodec) decodeVarFloat(buffer []byte, bits *uint64, width int) error {
	if codec == nil || codec.Size == 0 {
		return errors.New("nothing to decode")
	}

	const unit = 8 // bit width of encoding unit
	var gap, mask = codec.sizeOfGap(width)

	for (codec.Size & mask) > 0 {
		if codec.Ptr >= len(buffer) {
			return ErrBufferInsufficient
		}
		codec.Size--
		*bits = (*bits << unit) | uint64(buffer[codec.Ptr])
		codec.Ptr++
	}

	*bits <<= gap * unit
	codec.Size = 0
	return nil
}

func (codec *VarCodec) sizeOfGap(width int) (int, int) {
	var ms = mbit.OnesCount(^uint(0)) // machine bit width for an int
	var size = ms - 8                 // bit width of effective size
	var mask = -1 ^ (-1 << size)      // mask of effective size

	var gap = 0 // gap between encoded size and decoded size
	if codec.Size > 0 {
		if width > codec.Size {
			gap = width - codec.Size
		}
		var sign = -1 << (ms - 1) // single sign bit for an int
		codec.Size = sign | (gap << size) | (codec.Size & mask)
	} else {
		gap = (codec.Size >> size) & 0x7F
	}

	return gap, mask
}
