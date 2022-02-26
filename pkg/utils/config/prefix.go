package config

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
	"strings"
	"unicode"
)

func PrefixedBy(input interface{}, prefix string) (interface{}, error) {
	normalized, err := Normalize(input)
	if err != nil {
		// The only error that can come from normalize is if
		// input is a map[interface{}]interface{} and contains
		// a key that is not a string.
		return input, err
	}
	input = normalized

	if inputMap, ok := input.(map[string]interface{}); ok {
		converted := make(map[string]interface{}, len(inputMap))
		for k, v := range inputMap {
			if strings.HasPrefix(k, prefix) {
				key := uncapitalize(strings.TrimPrefix(k, prefix))
				converted[key] = v
			}
		}

		return converted, nil
	} else if inputMap, ok := input.(map[string]string); ok {
		converted := make(map[string]string, len(inputMap))
		for k, v := range inputMap {
			if strings.HasPrefix(k, prefix) {
				key := uncapitalize(strings.TrimPrefix(k, prefix))
				converted[key] = v
			}
		}

		return converted, nil
	}

	return input, nil
}

// uncapitalize initial capital letters in `str`.
func uncapitalize(str string) string {
	if len(str) == 0 {
		return str
	}

	vv := []rune(str) // Introduced later
	vv[0] = unicode.ToLower(vv[0])

	return string(vv)
}
