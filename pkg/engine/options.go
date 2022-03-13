package engine

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

	pkgauth "github.com/bhojpur/service/pkg/engine/auth"
	engine "github.com/bhojpur/service/pkg/engine/core"
	"github.com/bhojpur/service/pkg/engine/core/auth"
	"github.com/bhojpur/service/pkg/engine/core/log"
	"github.com/lucas-clemente/quic-go"
)

const (
	// DefaultProcessorAddr is the default address of downstream processor.
	DefaultProcessorAddr = "localhost:9000"
)

// Option is a function that applies a Bhojpur Service-Client option.
type Option func(o *Options)

// Options are the options for Bhojpur Service
type Options struct {
	ProcessorAddr string // target Processor endpoint address
	// ProcessorListenAddr     string // Processor endpoint address
	ProcessorWorkflowConfig string // Processor workflow file
	MeshConfigURL           string // meshConfigURL is the URL of EdgeMesh config
	ServerOptions           []engine.ServerOption
	ClientOptions           []engine.ClientOption
	QuicConfig              *quic.Config
	TLSConfig               *tls.Config
	Logger                  log.Logger
}

// WithProcessorAddr return a new options with ProcessorAddr set to addr.
func WithProcessorAddr(addr string) Option {
	return func(o *Options) {
		o.ProcessorAddr = addr
	}
}

// // WithProcessorListenAddr return a new options with ProcessorListenAddr set to addr.
// func WithProcessorListenAddr(addr string) Option {
// 	return func(o *options) {
// 		o.ProcessorListenAddr = addr
// 	}
// }

// TODO: WithWorkflowConfig

// WithMeshConfigURL sets the initial EdgeMesh config URL for the Bhojpur Service-Processor.
func WithMeshConfigURL(url string) Option {
	return func(o *Options) {
		o.MeshConfigURL = url
	}
}

func WithTLSConfig(tc *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = tc
	}
}

func WithQuicConfig(qc *quic.Config) Option {
	return func(o *Options) {
		o.QuicConfig = qc
	}
}

func WithClientOptions(opts ...engine.ClientOption) Option {
	return func(o *Options) {
		o.ClientOptions = opts
	}
}

func WithServerOptions(opts ...engine.ServerOption) Option {
	return func(o *Options) {
		o.ServerOptions = opts
	}
}

// WithAuth sets the server authentication method (used by server)
func WithAuth(auth auth.Authentication) Option {
	return func(o *Options) {
		o.ServerOptions = append(
			o.ServerOptions,
			engine.WithAuth(auth),
		)
	}
}

// WithAppKeyCredential sets the client credential (used by client): AppKey
func WithAppKeyCredential(appID string, appSecret string) Option {
	return WithCredential(pkgauth.NewAppKeyCredential(appID, appSecret))
}

// WithCredential sets the client credential
func WithCredential(cred auth.Credential) Option {
	return func(o *Options) {
		o.ClientOptions = append(
			o.ClientOptions,
			engine.WithCredential(cred),
		)
	}
}

// WithObserveDataTags sets client data tag list.
func WithObserveDataTags(tags ...byte) Option {
	return func(o *Options) {
		o.ClientOptions = append(
			o.ClientOptions,
			engine.WithObserveDataTags(tags...),
		)
	}
}

// WithLogger sets the client logger
func WithLogger(logger log.Logger) Option {
	return func(o *Options) {
		o.ClientOptions = append(
			o.ClientOptions,
			engine.WithLogger(logger),
		)
	}
}

// NewOptions creates a new options for Bhojpur Service-Client.
func NewOptions(opts ...Option) *Options {
	options := &Options{}

	for _, o := range opts {
		o(options)
	}

	if options.ProcessorAddr == "" {
		options.ProcessorAddr = DefaultProcessorAddr
	}

	return options
}
