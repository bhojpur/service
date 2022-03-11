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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataFrameEncode(t *testing.T) {
	var userDataTag byte = 0x15
	d := NewDataFrame()
	d.SetCarriage(userDataTag, []byte("bhojpur"))

	tidbuf := []byte(d.TransactionID())
	result := []byte{
		0x80 | byte(TagOfDataFrame), byte(len(tidbuf) + 4 + 8),
		0x80 | byte(TagOfMetaFrame), byte(len(tidbuf) + 2),
		byte(TagOfTransactionID), byte(len(tidbuf))}
	result = append(result, tidbuf...)
	result = append(result, 0x80|byte(TagOfPayloadFrame), 0x06,
		userDataTag, 0x04, 0x79, 0x6F, 0x6D, 0x6F)
	assert.Equal(t, result, d.Encode())
}

func TestDataFrameDecode(t *testing.T) {
	var userDataTag byte = 0x15
	buf := []byte{
		0x80 | byte(TagOfDataFrame), 0x10,
		0x80 | byte(TagOfMetaFrame), 0x06,
		byte(TagOfTransactionID), 0x04, 0x31, 0x32, 0x33, 0x34,
		0x80 | byte(TagOfPayloadFrame), 0x06,
		userDataTag, 0x04, 0x79, 0x6F, 0x6D, 0x6F}
	data, err := DecodeToDataFrame(buf)
	assert.NoError(t, err)
	assert.EqualValues(t, "1234", data.TransactionID())
	assert.EqualValues(t, userDataTag, data.GetDataTag())
	assert.EqualValues(t, []byte("bhojpur"), data.GetCarriage())
}
