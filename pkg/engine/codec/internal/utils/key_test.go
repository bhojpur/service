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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyOf(t *testing.T) {
	value := KeyOf("0x01")
	assert.Equal(t, byte(0x01), value, fmt.Sprintf("value does not match(%v): %v", byte(0x01), value))

	value = KeyOf("0X01")
	assert.Equal(t, byte(0x01), value, fmt.Sprintf("value does not match(%v): %v", byte(0x01), value))

	value = KeyOf("01")
	assert.Equal(t, byte(0x01), value, fmt.Sprintf("value does not match(%v): %v", byte(0x01), value))
}

func TestIsEmptyKey(t *testing.T) {
	assert.True(t, IsEmptyKey(0x00), "0x00 is a empty key")
}

func TestForbiddenCustomizedKey(t *testing.T) {
	assert.True(t, ForbidUserKey(0x01), "0x01 is disabled")
	assert.False(t, ForbidUserKey(0x10), "0x10 is allowed")
}

func TestAllowableSignalKey(t *testing.T) {
	assert.True(t, AllowSignalKey(0x02), "0x01 is allowed")
	assert.False(t, AllowSignalKey(0x01), "0x10 is disabled")
}
