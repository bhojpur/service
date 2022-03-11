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
	"bytes"
	"io"

	"github.com/bhojpur/service/pkg/utils/encoding"
	"github.com/bhojpur/service/pkg/utils/spec"
)

// StreamPacket implement the Packet interface.
type StreamPacket struct {
	t         spec.T
	l         spec.L
	vbuf      []byte
	vr        io.Reader
	chunkMode bool
	chunkSize int
}

var _ Packet = &StreamPacket{}

// SeqID returns the sequenceID of this packet
func (p *StreamPacket) SeqID() int { return p.t.Sid() }

// Size returns the size of whole packet.
func (p *StreamPacket) Size() int {
	// T.Size + L.Size + V.Size
	return p.t.Size() + p.l.Size() + p.l.VSize()
}

// VSize returns the size of V.
func (p *StreamPacket) VSize() int { return p.l.VSize() }

// Bytes return the raw bytes of this packet. V will be absent if
// is in chunked mode
func (p *StreamPacket) Bytes() []byte {
	buf := new(bytes.Buffer)
	// the raw bytes of T and L
	p.writeTL(buf)
	// p.valbuf stores the raw bytes of V
	buf.Write(p.vbuf)

	return buf.Bytes()
}

// VReader return an io.Reader which can be read as the content of V.
func (p *StreamPacket) VReader() io.Reader {
	if !p.chunkMode {
		return bytes.NewReader(p.vbuf)
	}
	return p.vr
}

// Reader return an io.Reader which can be read as the whole bytes of
// this packet. This function only available if this V of packet is in
// chunked mode.
func (p *StreamPacket) Reader() io.Reader {
	if !p.chunkMode {
		buf := new(bytes.Buffer)
		buf.Write(p.t.Bytes())
		buf.Write(p.l.Bytes())
		buf.Write(p.vbuf)
		return buf
	}

	buf := new(bytes.Buffer)
	// T and L of this packet
	p.writeTL(buf)
	// V of this packet
	buf.Write(p.vbuf)

	return &chunkVReader{
		buf:        buf,
		src:        p.vr,
		totalSize:  p.Size(),
		ChunkVSize: p.VSize(),
	}
}

// IsStreamMode returns a bool value indicates if the V of
// this packet is in stream mode
func (p *StreamPacket) IsStreamMode() bool {
	return p.chunkMode
}

// IsNodeMode returns a bool value indicates if this packet
// is node mode
func (p *StreamPacket) IsNodeMode() bool {
	return p.t.IsNodeMode()
}

// write the raw bytes of T and L to given buf
func (p *StreamPacket) writeTL(buf *bytes.Buffer) {
	buf.Write(p.t.Bytes())
	buf.Write(p.l.Bytes())
}

// BytesV return V as bytes
func (p *StreamPacket) BytesV() []byte {
	return p.vbuf
}

// UTF8StringV return V as utf-8 string
func (p *StreamPacket) UTF8StringV() string {
	return string(p.vbuf)
}

// Int32V return V as int32
func (p *StreamPacket) Int32V() (val int32, err error) {
	codec := encoding.VarCodec{Size: len(p.vbuf)}
	err = codec.DecodeNVarInt32(p.vbuf, &val)
	return val, err
}

// UInt32V return V as uint32
func (p *StreamPacket) UInt32V() (val uint32, err error) {
	codec := encoding.VarCodec{Size: len(p.vbuf)}
	err = codec.DecodeNVarUInt32(p.vbuf, &val)
	return val, err
}

// Int64V return V as int64
func (p *StreamPacket) Int64V() (val int64, err error) {
	codec := encoding.VarCodec{Size: len(p.vbuf)}
	err = codec.DecodeNVarInt64(p.vbuf, &val)
	return val, err
}

// UInt64V return V as uint64
func (p *StreamPacket) UInt64V() (val uint64, err error) {
	codec := encoding.VarCodec{Size: len(p.vbuf)}
	err = codec.DecodeNVarUInt64(p.vbuf, &val)
	return val, err
}

// Float32V return V as float32
func (p *StreamPacket) Float32V() (val float32, err error) {
	codec := encoding.VarCodec{Size: len(p.vbuf)}
	err = codec.DecodeVarFloat32(p.vbuf, &val)
	return val, err
}

// Float64V return V as float64
func (p *StreamPacket) Float64V() (val float64, err error) {
	codec := encoding.VarCodec{Size: len(p.vbuf)}
	err = codec.DecodeVarFloat64(p.vbuf, &val)
	return val, err
}

// BoolV return V as bool
func (p *StreamPacket) BoolV() (val bool, err error) {
	codec := encoding.VarCodec{Size: len(p.vbuf)}
	err = codec.DecodePVarBool(p.vbuf, &val)
	return val, err
}
