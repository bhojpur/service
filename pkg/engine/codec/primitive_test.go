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
	"testing"

	"github.com/stretchr/testify/assert"
)

// Packet 2 bytes
func TestLackLengthPacket(t *testing.T) {
	buf := []byte{0x01}
	expected := "invalid Bhojpur Service packet minimal size"
	_, _, _, err := DecodePrimitivePacket(buf)
	if err.Error() != expected {
		t.Errorf("err actual = %v, and Expected = %v", err, expected)
	}
}

func TestPacketWrongLength(t *testing.T) {
	buf := []byte{0x04, 0x00, 0x02, 0x01}
	expected := "malformed, Length can not smaller than 1"
	_, _, _, err := DecodePrimitivePacket(buf)
	assert.NoError(t, err)

	if err != nil && err.Error() != expected {
		t.Errorf("err should %v, actual = %v", expected, err)
	}
}

// 0x04:-1
func TestPacketRead(t *testing.T) {
	buf := []byte{0x04, 0x01, 0x7F}
	expectedTag := byte(0x04)
	var expectedLength uint32 = 1
	expectedValue := []byte{0x7F}

	res, endPos, _, err := DecodePrimitivePacket(buf)
	assert.NoError(t, err)

	assert.Equal(t, expectedTag, res.SeqID())
	assert.Equal(t, expectedLength, res.length)

	if !_compareByteSlice(res.valBuf, expectedValue) {
		t.Errorf("res.raw actual = %v, and Expected = %v", res.valBuf, expectedValue)
	}

	assert.Equal(t, endPos, 3)
}

// 0x0A:2
func TestParseInt32(t *testing.T) {
	buf := []byte{0x0A, 0x02, 0x81, 0x7F}
	expectedValue := int32(255)

	res, _, _, err := DecodePrimitivePacket(buf)
	assert.NoError(t, err)

	target, err := res.ToInt32()
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, target)
}

// 0x0B:"C"
func TestParseString(t *testing.T) {
	buf := []byte{0x0B, 0x01, 0x43}
	expectedValue := "C"

	res, _, _, err := DecodePrimitivePacket(buf)
	assert.NoError(t, err)

	target, err := res.ToUTF8String()
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, target)
}

// 0x0B:"C"
func TestParseEmptyString(t *testing.T) {
	buf := []byte{0x0B, 0x00}
	expectedValue := ""

	res, _, _, err := DecodePrimitivePacket(buf)
	assert.NoError(t, err)

	target, err := res.ToUTF8String()
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, target)
}

// compares two slice, every element is equal
func _compareByteSlice(left []byte, right []byte) bool {
	if len(left) != len(right) {
		return false
	}

	for i, v := range left {
		if v != right[i] {
			return false
		}
	}

	return true
}
