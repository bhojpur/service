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
	"fmt"
)

// Normalize converts map[interface{}]interface{} to map[string]interface{} to normalize
// for JSON and usage in component initialization.
// nolint:cyclop
func Normalize(i interface{}) (interface{}, error) {
	var err error
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			if strKey, ok := k.(string); ok {
				if m2[strKey], err = Normalize(v); err != nil {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("error parsing config field: %v", k)
			}
		}

		return m2, nil
	case map[string]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			if m2[k], err = Normalize(v); err != nil {
				return nil, err
			}
		}

		return m2, nil
	case []interface{}:
		for i, v := range x {
			if x[i], err = Normalize(v); err != nil {
				return nil, err
			}
		}
	}

	return i, nil
}
