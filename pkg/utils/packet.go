package utils

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
	"io"
)

var (
	errBuildIncomplete = errors.New("service.Encoder: invalid structure of packet")
	errInvalidAdding   = errors.New("service.Encoder: can not add this Packet after StreamPacket has been add")
	errNonStreamPacket = errors.New("service.Packet: this packet is not in node mode")
	errWriteFromReader = errors.New("service.streamV: write from reader error")
	errNotNodeMode     = errors.New("service.Encoder: packet should be in node mode can be add other packets as child")
	errNilReader       = errors.New("service.Decoder: nil source reader")
)

// Packet decribe a Bhojpur Service codec packet
type Packet interface {
	// SeqID returns the sequence ID of this packet.
	SeqID() int
	// Size returns the size of whole packet.
	Size() int
	// VSize returns the size of V.
	VSize() int
	// Bytes returns the whole bytes of this packet.
	Bytes() []byte
	// Reader returns an io.Reader which returns whole bytes.
	Reader() io.Reader
	// GetValReader returns an io.Reader which holds V.
	VReader() io.Reader
	// IsStreamMode returns a bool value indicates if the V of
	// this packet is in stream mode
	IsStreamMode() bool
	// IsNodeMode returns a bool value indicates if this packet
	// is node mode
	IsNodeMode() bool

	// BytesV return V as bytes
	BytesV() []byte
	// StringV return V as utf-8 string
	UTF8StringV() string
	// Int32V return V as int32
	Int32V() (val int32, err error)
	// UInt32V return V as uint32
	UInt32V() (val uint32, err error)
	// Int64V return V as int64
	Int64V() (val int64, err error)
	// UInt64V return V as uint64
	UInt64V() (val uint64, err error)
	// Float32V return V as float32
	Float32V() (val float32, err error)
	// Float64V return V as float64
	Float64V() (val float64, err error)
	// BoolV return V as bool
	BoolV() (val bool, err error)
}
