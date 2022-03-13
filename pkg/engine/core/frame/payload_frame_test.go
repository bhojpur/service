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

func TestPayloadFrameEncode(t *testing.T) {
	f := NewPayloadFrame(0x13).SetCarriage([]byte("bhojpur"))
	assert.Equal(t, []byte{0x80 | byte(TagOfPayloadFrame),
		0x09, 0x13, 0x07, 0x62, 0x68, 0x6f, 0x6a, 0x70, 0x75, 0x72}, f.Encode())
}

func TestPayloadFrameDecode(t *testing.T) {
	buf := []byte{0x80 | byte(TagOfPayloadFrame),
		0x09, 0x13, 0x07, 0x62, 0x68, 0x6f, 0x6a, 0x70, 0x75, 0x72}
	payload, err := DecodeToPayloadFrame(buf)
	assert.NoError(t, err)
	assert.EqualValues(t, 0x13, payload.Tag)
	assert.Equal(t, []byte{0x62, 0x68, 0x6f, 0x6a, 0x70, 0x75, 0x72}, payload.Carriage)
}
