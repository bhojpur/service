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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNInt32(t *testing.T) {
	testNVarInt32(t, -1, []byte{0xFF})
	testNVarInt32(t, -5, []byte{0xFB})
	testNVarInt32(t, 63, []byte{0x3F})
	testNVarInt32(t, -65, []byte{0xBF})
	testNVarInt32(t, 127, []byte{0x7F})
	testNVarInt32(t, 255, []byte{0x00, 0xFF})
	testNVarInt32(t, -4097, []byte{0xEF, 0xFF})
	testNVarInt32(t, -8193, []byte{0xDF, 0xFF})
	testNVarInt32(t, -2097152, []byte{0xE0, 0x00, 0x00})
	testNVarInt32(t, -134217729, []byte{0xF7, 0xFF, 0xFF, 0xFF})
	testNVarInt32(t, -2147483648, []byte{0x80, 0x00, 0x00, 0x00})
}

func TestNUInt32(t *testing.T) {
	testNVarUInt32(t, 1, []byte{0x01})
	testNVarUInt32(t, 127, []byte{0x7F})
	testNVarUInt32(t, 128, []byte{0x00, 0x80})
	testNVarUInt32(t, 130, []byte{0x00, 0x82})
	testNVarUInt32(t, 1048576, []byte{0x10, 0x00, 0x00})
	testNVarUInt32(t, 134217728, []byte{0x08, 0x00, 0x00, 0x00})
	testNVarUInt32(t, 4294967295, []byte{0xFF})
}

func TestNInt64(t *testing.T) {
	testNVarInt64(t, 0, []byte{0x00})
	testNVarInt64(t, 1, []byte{0x01})
	testNVarInt64(t, -1, []byte{0xFF})
}

func TestNUInt64(t *testing.T) {
	testNVarUInt64(t, 0, []byte{0x00})
	testNVarUInt64(t, 1, []byte{0x01})
	testNVarUInt64(t, 18446744073709551615, []byte{0xFF})
}

func testNVarInt32(t *testing.T, value int32, bytes []byte) {
	var msg = fmt.Sprintf("tester %v (%X): %X", value, uint32(value), bytes)
	var size = SizeOfNVarInt32(value)
	assert.Equal(t, len(bytes), size, msg)

	buffer := make([]byte, len(bytes))
	codec := VarCodec{Size: size}
	assert.Nil(t, codec.EncodeNVarInt32(buffer, value), msg)
	assert.Equal(t, bytes, buffer, msg)

	var val int32
	codec = VarCodec{Size: len(bytes)}
	assert.Nil(t, codec.DecodeNVarInt32(bytes, &val), msg)
	assert.Equal(t, value, val, msg)
}

func testNVarUInt32(t *testing.T, value uint32, bytes []byte) {
	var msg = fmt.Sprintf("tester %v (%X): %X", value, value, bytes)
	var size = SizeOfNVarUInt32(value)
	assert.Equal(t, len(bytes), size, msg)

	buffer := make([]byte, len(bytes))
	codec := VarCodec{Size: size}
	assert.Nil(t, codec.EncodeNVarUInt32(buffer, value), msg)
	assert.Equal(t, bytes, buffer, msg)

	var val uint32
	codec = VarCodec{Size: len(bytes)}
	assert.Nil(t, codec.DecodeNVarUInt32(bytes, &val), msg)
	assert.Equal(t, value, val, msg)
}

func testNVarInt64(t *testing.T, value int64, bytes []byte) {
	var msg = fmt.Sprintf("tester %v (%X): %X", value, uint64(value), bytes)
	var size = SizeOfNVarInt64(value)
	assert.Equal(t, len(bytes), size, msg)

	buffer := make([]byte, len(bytes))
	codec := VarCodec{Size: size}
	assert.Nil(t, codec.EncodeNVarInt64(buffer, value), msg)
	assert.Equal(t, bytes, buffer, msg)

	var val int64
	codec = VarCodec{Size: len(bytes)}
	assert.Nil(t, codec.DecodeNVarInt64(bytes, &val), msg)
	assert.Equal(t, value, val, msg)
}

func testNVarUInt64(t *testing.T, value uint64, bytes []byte) {
	var msg = fmt.Sprintf("tester %v (%X): %X", value, value, bytes)
	var size = SizeOfNVarUInt64(value)
	assert.Equal(t, len(bytes), size, msg)

	buffer := make([]byte, len(bytes))
	codec := VarCodec{Size: size}
	assert.Nil(t, codec.EncodeNVarUInt64(buffer, value), msg)
	assert.Equal(t, bytes, buffer, msg)

	var val uint64
	codec = VarCodec{Size: len(bytes)}
	assert.Nil(t, codec.DecodeNVarUInt64(bytes, &val), msg)
	assert.Equal(t, value, val)
}
