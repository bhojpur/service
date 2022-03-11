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
	"net"

	"github.com/bhojpur/service/pkg/engine/core/auth"
	"github.com/bhojpur/service/pkg/engine/core/store"
	"github.com/lucas-clemente/quic-go"
)

type ServerOptions struct {
	QuicConfig *quic.Config
	TLSConfig  *tls.Config
	Addr       string
	Auths      []auth.Authentication
	Store      store.Store
	Conn       net.PacketConn
}

func WithAddr(addr string) ServerOption {
	return func(o *ServerOptions) {
		o.Addr = addr
	}
}

func WithAuth(auth auth.Authentication) ServerOption {
	return func(o *ServerOptions) {
		o.Auths = append(o.Auths, auth)
	}
}

func WithStore(store store.Store) ServerOption {
	return func(o *ServerOptions) {
		o.Store = store
	}
}

func WithServerTLSConfig(tc *tls.Config) ServerOption {
	return func(o *ServerOptions) {
		o.TLSConfig = tc
	}
}

func WithServerQuicConfig(qc *quic.Config) ServerOption {
	return func(o *ServerOptions) {
		o.QuicConfig = qc
	}
}

func WithConn(conn net.PacketConn) ServerOption {
	return func(o *ServerOptions) {
		o.Conn = conn
	}
}
