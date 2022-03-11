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
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToStringSlice(t *testing.T) {
	value := reflect.ValueOf([]string{"a", "b"})
	out, ok := ToStringSlice(value.Interface())
	assert.True(t, ok, "must be successfully converted")
	assert.Equal(t, "a", out[0], fmt.Sprintf("value does not match(%v): %v", "a", out[0]))
	assert.Equal(t, "b", out[1], fmt.Sprintf("value does not match(%v): %v", "a", out[1]))
}

func TestToInt64Slice(t *testing.T) {
	value := reflect.ValueOf([]int64{1, 2})
	out, ok := ToInt64Slice(value.Interface())
	assert.True(t, ok, "must be successfully converted")
	assert.Equal(t, int64(1), out[0], fmt.Sprintf("value does not match(%v): %v", int64(1), out[0]))
	assert.Equal(t, int64(2), out[1], fmt.Sprintf("value does not match(%v): %v", int64(2), out[1]))
}

func TestToUInt64Slice(t *testing.T) {
	value := reflect.ValueOf([]uint64{1, 2})
	out, ok := ToUInt64Slice(value.Interface())
	assert.True(t, ok, "must be successfully converted")
	assert.Equal(t, uint64(1), out[0], fmt.Sprintf("value does not match(%v): %v", uint64(1), out[0]))
	assert.Equal(t, uint64(2), out[1], fmt.Sprintf("value does not match(%v): %v", uint64(2), out[1]))
}

func TestToUFloat64Slice(t *testing.T) {
	value := reflect.ValueOf([]float64{1, 2})
	out, ok := ToUFloat64Slice(value.Interface())
	assert.True(t, ok, "must be successfully converted")
	assert.Equal(t, float64(1), out[0], fmt.Sprintf("value does not match(%v): %v", float64(1), out[0]))
	assert.Equal(t, float64(2), out[1], fmt.Sprintf("value does not match(%v): %v", float64(2), out[1]))
}

func TestToBoolSlice(t *testing.T) {
	value := reflect.ValueOf([]bool{true, false})
	out, ok := ToBoolSlice(value.Interface())
	assert.True(t, ok, "must be successfully converted")
	assert.Equal(t, true, out[0], fmt.Sprintf("value does not match(%v): %v", true, out[0]))
	assert.Equal(t, false, out[1], fmt.Sprintf("value does not match(%v): %v", false, out[1]))
}
