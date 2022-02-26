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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestETagError(t *testing.T) {
	t.Run("invalid with context error", func(t *testing.T) {
		cerr := errors.New("error1")
		err := NewETagError(ETagInvalid, cerr)

		assert.Equal(t, invalidPrefix+": error1", err.Error())
	})

	t.Run("invalid without context error", func(t *testing.T) {
		err := NewETagError(ETagInvalid, nil)

		assert.Equal(t, invalidPrefix, err.Error())
	})

	t.Run("mismatch with context error", func(t *testing.T) {
		cerr := errors.New("error1")
		err := NewETagError(ETagMismatch, cerr)

		assert.Equal(t, mismatchPrefix+": error1", err.Error())
	})

	t.Run("mismatch without context error", func(t *testing.T) {
		err := NewETagError(ETagMismatch, nil)

		assert.Equal(t, mismatchPrefix, err.Error())
	})

	t.Run("valid kind", func(t *testing.T) {
		err := NewETagError(ETagMismatch, nil)

		assert.IsType(t, ETagMismatch, err.kind)
	})
}
