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
	"time"

	"github.com/bhojpur/service/pkg/engine/logger"
	pkgtls "github.com/bhojpur/service/pkg/engine/tls"
	"github.com/lucas-clemente/quic-go"
)

var _ Listener = (*defaultListener)(nil)

type defaultListener struct {
	c *quic.Config
	quic.Listener
}

func newListener() *defaultListener {
	return &defaultListener{}
}

func (l *defaultListener) Name() string {
	return "QUIC-Server"
}

func (l *defaultListener) Listen(conn net.PacketConn, tlsConfig *tls.Config, quicConfig *quic.Config) error {
	var err error
	// tls config
	var tc *tls.Config = tlsConfig
	if tc == nil {
		tc, err = pkgtls.CreateServerTLSConfig(conn.LocalAddr().String())
		if err != nil {
			logger.Errorf("%sCreateServerTLSConfig: %v", ServerLogPrefix, err)
			return err
		}
	}
	// quic config
	var c *quic.Config = quicConfig
	if c == nil {
		c = &quic.Config{
			Versions:                       []quic.VersionNumber{quic.Version1, quic.VersionDraft29},
			MaxIdleTimeout:                 time.Second * 5,
			KeepAlive:                      true,
			MaxIncomingStreams:             1000,
			MaxIncomingUniStreams:          1000,
			HandshakeIdleTimeout:           time.Second * 3,
			InitialStreamReceiveWindow:     1024 * 1024 * 2,
			InitialConnectionReceiveWindow: 1024 * 1024 * 2,
			DisablePathMTUDiscovery:        true,
			// Tracer:                      getQlogConfig("server"),
		}
	}
	l.c = c

	listener, err := quic.Listen(conn, tc, l.c)
	if err != nil {
		return err
	}
	l.Listener = listener
	return nil
}

func (l *defaultListener) Versions() []string {
	vers := make([]string, 0)
	for _, v := range l.c.Versions {
		vers = append(vers, v.String())
	}
	return vers
}
