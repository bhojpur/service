package bindings

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
	"strconv"
)

// InvokeRequest is the object given to a Bhojpur Service output binding.
type InvokeRequest struct {
	Data      []byte            `json:"data"`
	Metadata  map[string]string `json:"metadata"`
	Operation OperationKind     `json:"operation"`
}

// OperationKind defines an output binding operation.
type OperationKind string

// Non exhaustive list of operations. A binding can add operations that are not in this list.
const (
	GetOperation    OperationKind = "get"
	CreateOperation OperationKind = "create"
	DeleteOperation OperationKind = "delete"
	ListOperation   OperationKind = "list"
)

// GetMetadataAsBool parses metadata as bool.
func (r *InvokeRequest) GetMetadataAsBool(key string) (bool, error) {
	if val, ok := r.Metadata[key]; ok {
		boolVal, err := strconv.ParseBool(val)
		if err != nil {
			return false, fmt.Errorf("error parsing metadata `%s` with value `%s` as bool: %w", key, val, err)
		}

		return boolVal, nil
	}

	return false, nil
}

// GetMetadataAsInt64 parses metadata as int64.
func (r *InvokeRequest) GetMetadataAsInt64(key string, bitSize int) (int64, error) {
	if val, ok := r.Metadata[key]; ok {
		intVal, err := strconv.ParseInt(val, 10, bitSize)
		if err != nil {
			return 0, fmt.Errorf("error parsing metadata `%s` with value `%s` as int%d: %w", key, val, bitSize, err)
		}

		return intVal, nil
	}

	return 0, nil
}
