package tester

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

	engine "github.com/bhojpur/service/pkg/engine/codec"

	"github.com/stretchr/testify/assert"
)

func TestBasicTestData(t *testing.T) {
	input := BasicTestData{
		Vstring:  "foo",
		Vint32:   int32(127),
		Vint64:   int64(-1),
		Vuint32:  uint32(130),
		Vuint64:  uint64(18446744073709551615),
		Vfloat32: float32(0.25),
		Vfloat64: float64(23),
		Vbool:    true,
	}
	assert.NotEmpty(t, input, "Should not equal empty")
	assert.Equal(t, "foo", input.Vstring, fmt.Sprintf("value does not match(%v): %v", "foo", input.Vstring))
}

func TestObservableTestData(t *testing.T) {
	type ObservableTestData struct {
		A float32 `bhojpur:"0x10"`
		B string  `bhojpur:"0x11"`
	}

	codec := engine.NewCodec(0x20)
	obj := ObservableTestData{A: float32(456), B: "bhojpur"}
	buf, _ := codec.Marshal(obj)
	fmt.Printf("%#v\n", buf)
	target := []byte{0x81, 0xf, 0xa0, 0xd, 0x10, 0x2, 0x43, 0xe4, 0x11, 0x07, 0x62, 0x68, 0x6f, 0x6a, 0x70, 0x75, 0x72}
	for i, v := range target {
		assert.Equal(t, v, buf[i], fmt.Sprintf("should be: [%#x], but is [%#x]", v, buf[i]))
	}
}
