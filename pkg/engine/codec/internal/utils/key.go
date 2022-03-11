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
	"encoding/hex"
	"strings"
)

const (
	// EmptyKey mark an empty key
	EmptyKey byte = 0
)

// IsEmptyKey determine if observe is empty
func IsEmptyKey(observe byte) bool {
	return observe == byte(EmptyKey)
}

// KeyOf parse hex string to byte
func KeyOf(hexStr string) byte {
	if strings.HasPrefix(hexStr, "0x") {
		hexStr = strings.TrimPrefix(hexStr, "0x")
	} else if strings.HasPrefix(hexStr, "0X") {
		hexStr = strings.TrimPrefix(hexStr, "0X")
	}

	data, err := hex.DecodeString(hexStr)
	if err != nil {
		DefaultLogger.Errorf("hex.DecodeString error: %v", err)
		return 0x00
	}

	if len(data) == 0 {
		DefaultLogger.Errorf("hex.DecodeString data is []")
		return 0x00
	}

	return data[0]
}

// ForbidUserKey forbid user set that key
func ForbidUserKey(key byte) bool {
	if (key >= 0x01 && key <= 0x0f) || key >= 0x40 {
		return true
	}

	return false
}

// AllowSignalKey allow set that signal key
func AllowSignalKey(key byte) bool {
	switch key {
	case 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f:
		return true
	}
	return false
}
