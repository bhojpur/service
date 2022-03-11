package core

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
	"crypto/tls"

	"github.com/bhojpur/service/pkg/engine/core/auth"
	"github.com/bhojpur/service/pkg/engine/core/log"
	"github.com/lucas-clemente/quic-go"
)

type ClientOptions struct {
	ObserveDataTags []byte
	QuicConfig      *quic.Config
	TLSConfig       *tls.Config
	Credential      auth.Credential
	Logger          log.Logger
}

// WithObserveDataTags sets data tag list for the client.
func WithObserveDataTags(tags ...byte) ClientOption {
	return func(o *ClientOptions) {
		o.ObserveDataTags = tags
	}
}

// WithCredential sets app auth for the client.
func WithCredential(cred auth.Credential) ClientOption {
	return func(o *ClientOptions) {
		o.Credential = cred
	}
}

// WithClientTLSConfig sets tls config for the client.
func WithClientTLSConfig(tc *tls.Config) ClientOption {
	return func(o *ClientOptions) {
		o.TLSConfig = tc
	}
}

// WithClientQuicConfig sets quic config for the client.
func WithClientQuicConfig(qc *quic.Config) ClientOption {
	return func(o *ClientOptions) {
		o.QuicConfig = qc
	}
}

// WithLogger sets logger for the client.
func WithLogger(logger log.Logger) ClientOption {
	return func(o *ClientOptions) {
		o.Logger = logger
	}
}
