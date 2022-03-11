package utils

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
	"github.com/bhojpur/service/pkg/utils/encoding"
)

// SetUTF8StringV set utf-8 string type value as V
func (b *Encoder) SetUTF8StringV(v string) {
	buf := []byte(v)
	b.SetBytesV(buf)
}

// SetInt32V set an int32 type value as V
func (b *Encoder) SetInt32V(v int32) error {
	size := encoding.SizeOfNVarInt32(v)
	codec := encoding.VarCodec{Size: size}
	buf := make([]byte, size)
	err := codec.EncodeNVarInt32(buf, v)
	if err != nil {
		return err
	}
	b.SetBytesV(buf)
	return nil
}

// SetUInt32V set an uint32 type value as V
func (b *Encoder) SetUInt32V(v uint32) error {
	size := encoding.SizeOfNVarUInt32(v)
	codec := encoding.VarCodec{Size: size}
	buf := make([]byte, size)
	err := codec.EncodeNVarUInt32(buf, v)
	if err != nil {
		return err
	}
	b.SetBytesV(buf)
	return nil
}

// SetInt64V set an int64 type value as V
func (b *Encoder) SetInt64V(v int64) error {
	size := encoding.SizeOfNVarInt64(v)
	codec := encoding.VarCodec{Size: size}
	buf := make([]byte, size)
	err := codec.EncodeNVarInt64(buf, v)
	if err != nil {
		return err
	}
	b.SetBytesV(buf)
	return nil
}

// SetUInt64V set an uint64 type value as V
func (b *Encoder) SetUInt64V(v uint64) error {
	size := encoding.SizeOfNVarUInt64(v)
	codec := encoding.VarCodec{Size: size}
	buf := make([]byte, size)
	err := codec.EncodeNVarUInt64(buf, v)
	if err != nil {
		return err
	}
	b.SetBytesV(buf)
	return nil
}

// SetFloat32V set an float32 type value as V
func (b *Encoder) SetFloat32V(v float32) error {
	size := encoding.SizeOfVarFloat32(v)
	codec := encoding.VarCodec{Size: size}
	buf := make([]byte, size)
	err := codec.EncodeVarFloat32(buf, v)
	if err != nil {
		return err
	}
	b.SetBytesV(buf)
	return nil
}

// SetFloat64V set an float64 type value as V
func (b *Encoder) SetFloat64V(v float64) error {
	size := encoding.SizeOfVarFloat64(v)
	codec := encoding.VarCodec{Size: size}
	buf := make([]byte, size)
	err := codec.EncodeVarFloat64(buf, v)
	if err != nil {
		return err
	}
	b.SetBytesV(buf)
	return nil
}

// SetBoolV set bool type value as V
func (b *Encoder) SetBoolV(v bool) {
	var size = encoding.SizeOfPVarUInt32(uint32(1))
	codec := encoding.VarCodec{Size: size}
	buf := make([]byte, size)
	codec.EncodePVarBool(buf, v)
	b.SetBytesV(buf)
}
