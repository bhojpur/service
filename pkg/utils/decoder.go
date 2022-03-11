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

	"github.com/bhojpur/service/pkg/utils/spec"
)

// Decoder is the tool for decoding Bhojpur Service packet from stream
type Decoder struct {
	tag spec.T
	len *spec.L
	rd  io.Reader
}

// NewDecoder returns a Decoder from an io.Reader
func NewDecoder(reader io.Reader) *Decoder {
	return &Decoder{
		rd: reader,
	}
}

// SeqID return the SequenceID of the decoding packet
func (d *Decoder) SeqID() int {
	return d.tag.Sid()
}

// UnderlyingReader returns the reader this decoder using
func (d *Decoder) UnderlyingReader() io.Reader {
	return d.rd
}

// ReadHeader will block until io.EOF or recieve T and L of a packet.
func (d *Decoder) ReadHeader() error {
	// only read T and L
	return d.readTL()
}

// GetChunkedPacket will block until io.EOF or recieve V of a packet in chunked mode.
func (d *Decoder) GetChunkedPacket() Packet {
	return &StreamPacket{
		t:         d.tag,
		l:         *d.len,
		vr:        d.rd,
		chunkMode: true,
		chunkSize: d.len.VSize(),
	}
}

// GetFullfilledPacket read full Packet from given io.Reader
func (d *Decoder) GetFullfilledPacket() (packet Packet, err error) {
	// read V
	buf := new(bytes.Buffer)
	total := 0
	for {
		valbuf := make([]byte, d.len.VSize())
		n, err := d.rd.Read(valbuf)
		if n > 0 {
			total += n
			buf.Write(valbuf[:n])
		}
		if total >= d.len.VSize() || err != nil {
			break
		}
	}

	packet = &StreamPacket{
		t:         d.tag,
		l:         *d.len,
		vbuf:      buf.Bytes(),
		chunkMode: false,
	}

	return packet, nil
}

func (d *Decoder) readTL() (err error) {
	if d.rd == nil {
		return errNilReader
	}

	// read T
	d.tag, err = spec.ReadT(d.rd)
	if err != nil {
		return err
	}

	// read L
	d.len, err = spec.ReadL(d.rd)

	return err
}
