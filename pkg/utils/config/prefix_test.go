package config_test

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

	"github.com/bhojpur/service/pkg/utils/config"
)

func TestPrefixedBy(t *testing.T) {
	tests := map[string]struct {
		prefix   string
		input    interface{}
		expected interface{}
		err      string
	}{
		"map of string to string": {
			prefix: "test",
			input: map[string]string{
				"":        "",
				"ignore":  "don't include me",
				"testOne": "include me",
				"testTwo": "and me",
			},
			expected: map[string]string{
				"one": "include me",
				"two": "and me",
			},
		},
		"map of string to interface{}": {
			prefix: "test",
			input: map[string]interface{}{
				"":        "",
				"ignore":  "don't include me",
				"testOne": "include me",
				"testTwo": "and me",
			},
			expected: map[string]interface{}{
				"one": "include me",
				"two": "and me",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := config.PrefixedBy(tc.input, tc.prefix)
			if tc.err != "" {
				if assert.Error(t, err) {
					assert.Equal(t, tc.err, err.Error())
				}
			} else {
				assert.Equal(t, tc.expected, actual, "unexpected output")
			}
		})
	}
}
