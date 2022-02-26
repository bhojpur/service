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

// GetResponse is the response object for getting state.
type GetResponse struct {
	Data        []byte            `json:"data"`
	ETag        *string           `json:"etag,omitempty"`
	Metadata    map[string]string `json:"metadata"`
	ContentType *string           `json:"contentType,omitempty"`
}

// BulkGetResponse is the response object for bulk get response.
type BulkGetResponse struct {
	Key         string            `json:"key"`
	Data        []byte            `json:"data"`
	ETag        *string           `json:"etag,omitempty"`
	Metadata    map[string]string `json:"metadata"`
	Error       string            `json:"error,omitempty"`
	ContentType *string           `json:"contentType,omitempty"`
}

// QueryResponse is the response object for querying state.
type QueryResponse struct {
	Results  []QueryItem       `json:"results"`
	Token    string            `json:"token,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// QueryItem is an object representing a single entry in query results.
type QueryItem struct {
	Key         string  `json:"key"`
	Data        []byte  `json:"data"`
	ETag        *string `json:"etag,omitempty"`
	Error       string  `json:"error,omitempty"`
	ContentType *string `json:"contentType,omitempty"`
}
