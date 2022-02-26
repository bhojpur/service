package query

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

type Filter interface {
	Parse(interface{}) error
}

func parseFilter(obj interface{}) (Filter, error) {
	m, ok := obj.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("filter unit must be a map")
	}
	if len(m) != 1 {
		return nil, fmt.Errorf("filter unit must have a single element")
	}
	for k, v := range m {
		switch k {
		case "EQ":
			f := &EQ{}
			err := f.Parse(v)

			return f, err
		case "IN":
			f := &IN{}
			err := f.Parse(v)

			return f, err
		case "AND":
			f := &AND{}
			err := f.Parse(v)

			return f, err
		case "OR":
			f := &OR{}
			err := f.Parse(v)

			return f, err
		default:
			return nil, fmt.Errorf("unsupported filter %q", k)
		}
	}

	return nil, nil
}

type EQ struct {
	Key string
	Val interface{}
}

func (f *EQ) Parse(obj interface{}) error {
	m, ok := obj.(map[string]interface{})
	if !ok {
		return fmt.Errorf("EQ filter must be a map")
	}
	if len(m) != 1 {
		return fmt.Errorf("EQ filter must contain a single key/value pair")
	}
	for k, v := range m {
		f.Key = k
		f.Val = v
	}

	return nil
}

type IN struct {
	Key  string
	Vals []interface{}
}

func (f *IN) Parse(obj interface{}) error {
	m, ok := obj.(map[string]interface{})
	if !ok {
		return fmt.Errorf("IN filter must be a map")
	}
	if len(m) != 1 {
		return fmt.Errorf("IN filter must contain a single key/value pair")
	}
	for k, v := range m {
		f.Key = k
		if f.Vals, ok = v.([]interface{}); !ok {
			return fmt.Errorf("IN filter value must be an array")
		}
	}

	return nil
}

type AND struct {
	Filters []Filter
}

func (f *AND) Parse(obj interface{}) (err error) {
	f.Filters, err = parseFilters("AND", obj)

	return
}

type OR struct {
	Filters []Filter
}

func (f *OR) Parse(obj interface{}) (err error) {
	f.Filters, err = parseFilters("OR", obj)

	return
}

func parseFilters(t string, obj interface{}) ([]Filter, error) {
	arr, ok := obj.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%s filter must be an array", t)
	}
	if len(arr) < 2 {
		return nil, fmt.Errorf("%s filter must contain at least two entries", t)
	}
	filters := make([]Filter, len(arr))
	for i, entry := range arr {
		var err error
		if filters[i], err = parseFilter(entry); err != nil {
			return nil, err
		}
	}

	return filters, nil
}
