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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/engine/codec/internal/utils"
)

func TestStructEncoderWithSignals(t *testing.T) {
	input := exampleData{
		Name:  "bhojpur",
		Noise: float32(456),
		Therm: thermometer{Temperature: float32(30), Humidity: float32(40)},
	}

	encoder := newStructEncoder(0x30, structEncoderOptionRoot(utils.RootToken),
		structEncoderOptionConfig(&structEncoderConfig{
			ZeroFields: true,
			TagName:    "bhojpur",
		}))
	inputBuf, _ := encoder.Encode(input,
		createSignal(0x02).SetString("a"),
		createSignal(0x03).SetString("b"))
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

	var mold exampleData
	err := ToObject(inputBuf[2+3+3:], &mold)
	assert.NoError(t, err, fmt.Sprintf("decode error:%v", err))
	assert.Equal(t, input.Name, mold.Name, fmt.Sprintf("value does not match(%v): %v", input.Name, mold.Name))
	assert.Equal(t, input.Noise, mold.Noise, fmt.Sprintf("value does not match(%v): %v", input.Noise, mold.Noise))
	assert.Equal(t, input.Therm.Temperature, mold.Therm.Temperature, fmt.Sprintf("value does not match(%v): %v", input.Therm.Temperature, mold.Therm.Temperature))
	assert.Equal(t, input.Therm.Humidity, mold.Therm.Humidity, fmt.Sprintf("value does not match(%v): %v", input.Therm.Humidity, mold.Therm.Humidity))
}

func TestStructEncoderWithSignalsNoRoot(t *testing.T) {
	input := exampleData{
		Name:  "bhojpur",
		Noise: float32(456),
		Therm: thermometer{Temperature: float32(30), Humidity: float32(40)},
	}

	encoder := newStructEncoder(0x30)
	inputBuf, _ := encoder.Encode(input,
		createSignal(0x02).SetString("a"),
		createSignal(0x03).SetString("b"))
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

	var mold exampleData
	err := ToObject(inputBuf[3+3:], &mold)
	assert.NoError(t, err, fmt.Sprintf("decode error:%v", err))
	assert.Equal(t, input.Name, mold.Name, fmt.Sprintf("value does not match(%v): %v", input.Name, mold.Name))
	assert.Equal(t, input.Noise, mold.Noise, fmt.Sprintf("value does not match(%v): %v", input.Noise, mold.Noise))
	assert.Equal(t, input.Therm.Temperature, mold.Therm.Temperature, fmt.Sprintf("value does not match(%v): %v", input.Therm.Temperature, mold.Therm.Temperature))
	assert.Equal(t, input.Therm.Humidity, mold.Therm.Humidity, fmt.Sprintf("value does not match(%v): %v", input.Therm.Humidity, mold.Therm.Humidity))
}

func TestStructSliceEncoderWithSignals(t *testing.T) {
	input := exampleSlice{
		Therms: []thermometer{
			{Temperature: float32(30), Humidity: float32(40)},
			{Temperature: float32(50), Humidity: float32(60)},
		},
	}

	encoder := newStructEncoder(0x30, structEncoderOptionRoot(utils.RootToken))
	inputBuf, _ := encoder.Encode(input,
		createSignal(0x02).SetString("a"),
		createSignal(0x03).SetString("b"))
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

	var mold exampleSlice
	err := ToObject(inputBuf[2+3+3:], &mold)
	assert.NoError(t, err, fmt.Sprintf("decode error:%v", err))

	assert.Equal(t, float32(30), mold.Therms[0].Temperature, fmt.Sprintf("value does not match(%v): %v", float32(30), mold.Therms[0].Temperature))
	assert.Equal(t, float32(40), mold.Therms[0].Humidity, fmt.Sprintf("value does not match(%v): %v", float32(40), mold.Therms[0].Humidity))
	assert.Equal(t, float32(50), mold.Therms[1].Temperature, fmt.Sprintf("value does not match(%v): %v", float32(50), mold.Therms[1].Temperature))
	assert.Equal(t, float32(60), mold.Therms[1].Humidity, fmt.Sprintf("value does not match(%v): %v", float32(60), mold.Therms[1].Humidity))
}

func TestStructSliceEncoderWithSignalsNoRoot(t *testing.T) {
	input := exampleSlice{
		Therms: []thermometer{
			{Temperature: float32(30), Humidity: float32(40)},
			{Temperature: float32(50), Humidity: float32(60)},
		},
	}

	encoder := newStructEncoder(0x30)
	inputBuf, _ := encoder.Encode(input,
		createSignal(0x02).SetString("a"),
		createSignal(0x03).SetString("b"))
	testPrintf("inputBuf=%v\n", utils.FormatBytes(inputBuf))

	var mold exampleSlice
	err := ToObject(inputBuf[3+3:], &mold)
	assert.NoError(t, err, fmt.Sprintf("decode error:%v", err))

	assert.Equal(t, float32(30), mold.Therms[0].Temperature, fmt.Sprintf("value does not match(%v): %v", float32(30), mold.Therms[0].Temperature))
	assert.Equal(t, float32(40), mold.Therms[0].Humidity, fmt.Sprintf("value does not match(%v): %v", float32(40), mold.Therms[0].Humidity))
	assert.Equal(t, float32(50), mold.Therms[1].Temperature, fmt.Sprintf("value does not match(%v): %v", float32(50), mold.Therms[1].Temperature))
	assert.Equal(t, float32(60), mold.Therms[1].Humidity, fmt.Sprintf("value does not match(%v): %v", float32(60), mold.Therms[1].Humidity))
}

func TestStructForbidUserKey(t *testing.T) {
	input := exampleData{
		Name:  "bhojpur",
		Noise: float32(456),
		Therm: thermometer{Temperature: float32(30), Humidity: float32(40)},
	}

	var key byte = 0x02
	assert.Panics(t, func() {
		newStructEncoder(key,
			structEncoderOptionRoot(utils.RootToken),
			structEncoderOptionConfig(&structEncoderConfig{
				ZeroFields: true,
				TagName:    "bhojpur",
			}),
			structEncoderOptionForbidUserKey(utils.ForbidUserKey)).
			Encode(input)
	}, "should forbid this Key: %#x", key)

	key = 0x0f
	assert.Panics(t, func() {
		newStructEncoder(key,
			structEncoderOptionRoot(utils.RootToken),
			structEncoderOptionConfig(&structEncoderConfig{
				ZeroFields: true,
				TagName:    "bhojpur",
			}),
			structEncoderOptionForbidUserKey(utils.ForbidUserKey)).
			Encode(input)
	}, "should forbid this Key: %#x", key)

}

func TestStructAllowSignalKey(t *testing.T) {
	input := exampleData{
		Name:  "bhojpur",
		Noise: float32(456),
		Therm: thermometer{Temperature: float32(30), Humidity: float32(40)},
	}

	var signalKey byte = 0x02
	assert.NotPanics(t, func() {
		newStructEncoder(0x30,
			structEncoderOptionRoot(utils.RootToken),
			structEncoderOptionConfig(&structEncoderConfig{
				ZeroFields: true,
				TagName:    "bhojpur",
			}),
			structEncoderOptionAllowSignalKey(utils.AllowSignalKey)).
			Encode(input, createSignal(signalKey).SetString("a"))
	}, "should allow this Signal Key: %#x", signalKey)
}
