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
	"github.com/stretchr/testify/require"

	"github.com/bhojpur/service/pkg/utils/config"
)

func TestNormalize(t *testing.T) {
	tests := map[string]struct {
		input    interface{}
		expected interface{}
		err      string
	}{
		"simple": {input: "test", expected: "test"},
		"map of string to interface{}": {
			input: map[string]interface{}{
				"test": "1234",
				"nested": map[string]interface{}{
					"value": "5678",
				},
			}, expected: map[string]interface{}{
				"test": "1234",
				"nested": map[string]interface{}{
					"value": "5678",
				},
			},
		},
		"map of string to interface{} with error": {
			input: map[string]interface{}{
				"test": "1234",
				"nested": map[interface{}]interface{}{
					5678: "5678",
				},
			}, err: "error parsing config field: 5678",
		},
		"map of interface{} to interface{}": {
			input: map[string]interface{}{
				"test": "1234",
				"nested": map[interface{}]interface{}{
					"value": "5678",
				},
			}, expected: map[string]interface{}{
				"test": "1234",
				"nested": map[string]interface{}{
					"value": "5678",
				},
			},
		},
		"map of interface{} to interface{} with error": {
			input: map[interface{}]interface{}{
				"test": "1234",
				"nested": map[interface{}]interface{}{
					5678: "5678",
				},
			}, err: "error parsing config field: 5678",
		},
		"slice of interface{}": {
			input: []interface{}{
				map[interface{}]interface{}{
					"value": "5678",
				},
			}, expected: []interface{}{
				map[string]interface{}{
					"value": "5678",
				},
			},
		},
		"slice of interface{} with error": {
			input: []interface{}{
				map[interface{}]interface{}{
					1234: "1234",
				},
			}, err: "error parsing config field: 1234",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := config.Normalize(tc.input)
			if tc.err != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.expected, actual)
		})
	}
}
