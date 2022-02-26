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
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSetRequestWithOptions is used to test request options.
func TestSetRequestWithOptions(t *testing.T) {
	t.Run("set with default options", func(t *testing.T) {
		counter := 0
		SetWithOptions(func(req *SetRequest) error {
			counter++

			return nil
		}, &SetRequest{})
		assert.Equal(t, 1, counter, "should execute only once")
	})

	t.Run("set with no explicit options", func(t *testing.T) {
		counter := 0
		SetWithOptions(func(req *SetRequest) error {
			counter++

			return nil
		}, &SetRequest{
			Options: SetStateOption{},
		})
		assert.Equal(t, 1, counter, "should execute only once")
	})
}

// TestCheckRequestOptions is used to validate request options.
func TestCheckRequestOptions(t *testing.T) {
	t.Run("set state options", func(t *testing.T) {
		ro := SetStateOption{Concurrency: FirstWrite, Consistency: Eventual}
		err := CheckRequestOptions(ro)
		assert.NoError(t, err)
	})
	t.Run("delete state options", func(t *testing.T) {
		ro := DeleteStateOption{Concurrency: FirstWrite, Consistency: Eventual}
		err := CheckRequestOptions(ro)
		assert.NoError(t, err)
	})
	t.Run("get state options", func(t *testing.T) {
		ro := GetStateOption{Consistency: Eventual}
		err := CheckRequestOptions(ro)
		assert.NoError(t, err)
	})
	t.Run("invalid state options", func(t *testing.T) {
		ro := SetStateOption{Concurrency: "invalid", Consistency: Eventual}
		err := CheckRequestOptions(ro)
		assert.Error(t, err)
	})
}
