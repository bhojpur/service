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
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/engine/codec/internal/utils"
)

func TestBasicEncoderWithSignals(t *testing.T) {
	input := int32(456)

	encoder := newBasicEncoder(0x10, basicEncoderOptionRoot(utils.RootToken))
	inputBuf, _ := encoder.Encode(input,
		createSignal(0x02).SetString("a"),
		createSignal(0x03).SetString("b"))
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

	value, err := ToInt32(inputBuf[2+3+3:])
	assert.NoError(t, err, fmt.Sprintf("decode error:%v", err))
	assert.Equal(t, input, value, fmt.Sprintf("value does not match(%v): %v", input, value))
}

func TestBasicEncoderWithSignalsNoRoot(t *testing.T) {
	input := int32(456)

	encoder := newBasicEncoder(0x10)
	inputBuf, _ := encoder.Encode(input,
		createSignal(0x02).SetString("a"),
		createSignal(0x03).SetString("b"))
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

	value, err := ToInt32(inputBuf[3+3:])
	assert.NoError(t, err, fmt.Sprintf("decode error:%v", err))
	assert.Equal(t, input, value, fmt.Sprintf("value does not match(%v): %v", input, value))
}

func TestBasicSliceEncoderWithSignals(t *testing.T) {
	input := []int32{123, 456}

	encoder := newBasicEncoder(0x10, basicEncoderOptionRoot(utils.RootToken))
	inputBuf, _ := encoder.Encode(input,
		createSignal(0x02).SetString("a"),
		createSignal(0x03).SetString("b"))
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

	value, err := ToInt32Slice(inputBuf[2+3+3:])
	assert.NoError(t, err, fmt.Sprintf("decode error:%v", err))

	expectedValue := reflect.ValueOf(input)
	resultValue := reflect.ValueOf(value)
	for i := 0; i < expectedValue.Len(); i++ {
		assert.Equal(t, expectedValue.Index(i).Interface(), resultValue.Index(i).Interface(),
			fmt.Sprintf("Item values are not equal %v: %v",
				expectedValue.Index(i).Interface(), resultValue.Index(i).Interface()))
	}
}

func TestBasicSliceEncoderWithSignalsNoRoot(t *testing.T) {
	input := []int32{123, 456}

	encoder := newBasicEncoder(0x10)
	inputBuf, _ := encoder.Encode(input,
		createSignal(0x02).SetString("a"),
		createSignal(0x03).SetString("b"))
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

	value, err := ToInt32Slice(inputBuf[3+3:])
	assert.NoError(t, err, fmt.Sprintf("decode error:%v", err))

	expectedValue := reflect.ValueOf(input)
	resultValue := reflect.ValueOf(value)
	for i := 0; i < expectedValue.Len(); i++ {
		assert.Equal(t, expectedValue.Index(i).Interface(), resultValue.Index(i).Interface(),
			fmt.Sprintf("Item values are not equal %v: %v",
				expectedValue.Index(i).Interface(), resultValue.Index(i).Interface()))
	}
}

func TestBasicForbidUserKey(t *testing.T) {
	input := int32(456)

	var key byte = 0x02
	assert.Panics(t, func() {
		newBasicEncoder(key,
			basicEncoderOptionRoot(utils.RootToken),
			basicEncoderOptionForbidUserKey(utils.ForbidUserKey)).
			Encode(input)
	}, "should forbid this Key: %#x", key)

	key = 0x0f
	assert.Panics(t, func() {
		newBasicEncoder(key,
			basicEncoderOptionRoot(utils.RootToken),
			basicEncoderOptionForbidUserKey(utils.ForbidUserKey)).
			Encode(input)
	}, "should forbid this Key: %#x", key)
}

func TestBasicAllowSignalKey(t *testing.T) {
	input := int32(456)

	var signalKey byte = 0x02
	assert.NotPanics(t, func() {
		newBasicEncoder(0x10,
			basicEncoderOptionRoot(utils.RootToken),
			basicEncoderOptionAllowSignalKey(utils.AllowSignalKey)).
			Encode(input, createSignal(signalKey).SetString("a"))
	}, "should allow this Signal Key: %#x", signalKey)

}
