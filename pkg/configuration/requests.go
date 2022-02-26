package configuration

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

// Item represents a configuration item with name, content and other information.
type Item struct {
	Key      string            `json:"key"`
	Value    string            `json:"value,omitempty"`
	Version  string            `json:"version,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// GetRequest is the object describing a request to get configuration.
type GetRequest struct {
	Keys     []string          `json:"keys"`
	Metadata map[string]string `json:"metadata"`
}

// SubscribeRequest is the object describing a request to subscribe configuration.
type SubscribeRequest struct {
	Keys     []string          `json:"keys"`
	Metadata map[string]string `json:"metadata"`
}

// UnsubscribeRequest is the object describing a request to unsubscribe configuration.
type UnsubscribeRequest struct {
	ID string `json:"id"`
}

// UpdateEvent is the object describing a configuration update event.
type UpdateEvent struct {
	ID    string  `json:"id"`
	Items []*Item `json:"items"`
}
