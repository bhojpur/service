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
	"fmt"
)

type ETagErrorKind string

const (
	mismatchPrefix = "possible etag mismatch. error from state store"
	invalidPrefix  = "invalid etag value"

	ETagInvalid  ETagErrorKind = "invalid"
	ETagMismatch ETagErrorKind = "mismatch"
)

// ETagError is a custom error type for etag exceptions.
type ETagError struct {
	err  error
	kind ETagErrorKind
}

func (e *ETagError) Kind() ETagErrorKind {
	return e.kind
}

func (e *ETagError) Error() string {
	var prefix string

	switch e.kind {
	case ETagInvalid:
		prefix = invalidPrefix
	case ETagMismatch:
		prefix = mismatchPrefix
	}

	if e.err != nil {
		return fmt.Sprintf("%s: %s", prefix, e.err)
	}

	return errors.New(prefix).Error()
}

// NewETagError returns an ETagError wrapping an existing context error.
func NewETagError(kind ETagErrorKind, err error) *ETagError {
	return &ETagError{
		err:  err,
		kind: kind,
	}
}

// BulkDeleteRowMismatchError represents mismatch in rowcount while deleting rows.
type BulkDeleteRowMismatchError struct {
	expected uint64
	affected uint64
}

func (e *BulkDeleteRowMismatchError) Error() string {
	return fmt.Sprintf("delete affected only %d rows, expected %d", e.affected, e.expected)
}

// BulkDeleteRowMismatchError returns a BulkDeleteRowMismatchError.
func NewBulkDeleteRowMismatchError(expected, affected uint64) *BulkDeleteRowMismatchError {
	return &BulkDeleteRowMismatchError{
		expected: expected,
		affected: affected,
	}
}
