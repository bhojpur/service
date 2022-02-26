package state

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

const (
	FirstWrite = "first-write"
	LastWrite  = "last-write"
	Strong     = "strong"
	Eventual   = "eventual"
)

// CheckRequestOptions checks if request options use supported keywords.
func CheckRequestOptions(options interface{}) error {
	switch o := options.(type) {
	case SetStateOption:
		if err := validateConsistencyOption(o.Consistency); err != nil {
			return err
		}
		if err := validateConcurrencyOption(o.Concurrency); err != nil {
			return err
		}
	case DeleteStateOption:
		// no support in golang for multiple condition in type switch, so need to check explicitly
		if err := validateConsistencyOption(o.Consistency); err != nil {
			return err
		}
		if err := validateConcurrencyOption(o.Concurrency); err != nil {
			return err
		}
	case GetStateOption:
		if err := validateConsistencyOption(o.Consistency); err != nil {
			return err
		}
	}

	return nil
}

func validateConcurrencyOption(c string) error {
	if c != "" && c != FirstWrite && c != LastWrite {
		return fmt.Errorf("unrecognized concurrency model '%s'", c)
	}

	return nil
}

func validateConsistencyOption(c string) error {
	if c != "" && c != Strong && c != Eventual {
		return fmt.Errorf("unrecognized consistency model '%s'", c)
	}

	return nil
}

// SetWithOptions handles SetRequest with request options.
func SetWithOptions(method func(req *SetRequest) error, req *SetRequest) error {
	return method(req)
}

// DeleteWithOptions handles DeleteRequest with options.
func DeleteWithOptions(method func(req *DeleteRequest) error, req *DeleteRequest) error {
	return method(req)
}
