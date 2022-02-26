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

import "github.com/bhojpur/service/pkg/state/query"

// GetRequest is the object describing a state fetch request.
type GetRequest struct {
	Key      string            `json:"key"`
	Metadata map[string]string `json:"metadata"`
	Options  GetStateOption    `json:"options,omitempty"`
}

// GetStateOption controls how a state store reacts to a get request.
type GetStateOption struct {
	Consistency string `json:"consistency"` // "eventual, strong"
}

// DeleteRequest is the object describing a delete state request.
type DeleteRequest struct {
	Key      string            `json:"key"`
	ETag     *string           `json:"etag,omitempty"`
	Metadata map[string]string `json:"metadata"`
	Options  DeleteStateOption `json:"options,omitempty"`
}

// Key gets the Key on a DeleteRequest.
func (r DeleteRequest) GetKey() string {
	return r.Key
}

// Metadata gets the Metadata on a DeleteRequest.
func (r DeleteRequest) GetMetadata() map[string]string {
	return r.Metadata
}

// DeleteStateOption controls how a state store reacts to a delete request.
type DeleteStateOption struct {
	Concurrency string `json:"concurrency,omitempty"` // "concurrency"
	Consistency string `json:"consistency"`           // "eventual, strong"
}

// SetRequest is the object describing an upsert request.
type SetRequest struct {
	Key         string            `json:"key"`
	Value       interface{}       `json:"value"`
	ETag        *string           `json:"etag,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Options     SetStateOption    `json:"options,omitempty"`
	ContentType *string           `json:"contentType,omitempty"`
}

// GetKey gets the Key on a SetRequest.
func (r SetRequest) GetKey() string {
	return r.Key
}

// GetMetadata gets the Key on a SetRequest.
func (r SetRequest) GetMetadata() map[string]string {
	return r.Metadata
}

// SetStateOption controls how a state store reacts to a set request.
type SetStateOption struct {
	Concurrency string `json:"concurrency,omitempty"` // first-write, last-write
	Consistency string `json:"consistency"`           // "eventual, strong"
}

// OperationType describes a CRUD operation performed against a state store.
type OperationType string

// Upsert is an update or create operation.
const Upsert OperationType = "upsert"

// Delete is a delete operation.
const Delete OperationType = "delete"

// TransactionalStateRequest describes a transactional operation against a state store that comprises multiple types of operations
// The Request field is either a DeleteRequest or SetRequest.
type TransactionalStateRequest struct {
	Operations []TransactionalStateOperation `json:"operations"`
	Metadata   map[string]string             `json:"metadata,omitempty"`
}

// TransactionalStateOperation describes operation type, key, and value for transactional operation.
type TransactionalStateOperation struct {
	Operation OperationType `json:"operation"`
	Request   interface{}   `json:"request"`
}

// KeyInt is an interface that allows gets of the Key and Metadata inside requests.
type KeyInt interface {
	GetKey() string
	GetMetadata() map[string]string
}

type QueryRequest struct {
	Query    query.Query       `json:"query"`
	Metadata map[string]string `json:"metadata,omitempty"`
}
