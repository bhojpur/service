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
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloat32(t *testing.T) {
	testVarFloat32(t, 0, []byte{0x00})
	testVarFloat32(t, 1, []byte{0x3F, 0x80})
	testVarFloat32(t, 25, []byte{0x41, 0xC8})
	testVarFloat32(t, -2, []byte{0xC0})
	testVarFloat32(t, 0.25, []byte{0x3E, 0x80})
	testVarFloat32(t, 0.375, []byte{0x3E, 0xC0})
	testVarFloat32(t, 12.375, []byte{0x41, 0x46})
	testVarFloat32(t, 68.123, []byte{0x42, 0x88, 0x3E, 0xFA})
}

func TestFloat64(t *testing.T) {
	testVarFloat64(t, 0, []byte{0x00})
	testVarFloat64(t, 1, []byte{0x3F, 0xF0})
	testVarFloat64(t, 2, []byte{0x40})
	testVarFloat64(t, 23, []byte{0x40, 0x37})
	testVarFloat64(t, -2, []byte{0xC0})
	testVarFloat64(t, 0.01171875, []byte{0x3F, 0x88})
}

func testVarFloat32(t *testing.T, value float32, bytes []byte) {
	var msg = fmt.Sprintf("tester %v (%X): %X", value, math.Float32bits(value), bytes)
	var size = SizeOfVarFloat32(value)
	assert.Equal(t, len(bytes), size, msg)

	buffer := make([]byte, len(bytes))
	codec := VarCodec{Size: size}
	assert.Nil(t, codec.EncodeVarFloat32(buffer, value), msg)
	assert.Equal(t, bytes, buffer, msg)

	var val float32
	codec = VarCodec{Size: len(bytes)}
	assert.Nil(t, codec.DecodeVarFloat32(bytes, &val), msg)
	assert.Equal(t, value, val, msg)
}

func testVarFloat64(t *testing.T, value float64, bytes []byte) {
	var msg = fmt.Sprintf("tester %v (%X): %X", value, math.Float64bits(value), bytes)
	var size = SizeOfVarFloat64(value)
	assert.Equal(t, len(bytes), size, msg)

	buffer := make([]byte, len(bytes))
	codec := VarCodec{Size: size}
	assert.Nil(t, codec.EncodeVarFloat64(buffer, value), msg)
	assert.Equal(t, bytes, buffer, msg)

	var val float64
	codec = VarCodec{Size: len(bytes)}
	assert.Nil(t, codec.DecodeVarFloat64(bytes, &val), msg)
	assert.Equal(t, value, val, msg)
}
