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

// Encoder is the tool for creating a Bhojpur Service packet easily
type Encoder struct {
	tag           spec.T
	len           *spec.L
	valReader     io.Reader
	valReaderSize int
	nodes         map[int]Packet
	state         int
	size          int32 // size of value
	isStreamMode  bool
	valbuf        *bytes.Buffer
	done          bool
	seqID         int
	isNodeMode    bool
}

// SetSeqID set sequenceID of a Bhojpur Service packet, if this packet contains other
// Bhojpur Service packets, isNode should set to true
func (b *Encoder) SetSeqID(seqID int, isNode bool) {
	// init
	b.valbuf = new(bytes.Buffer)
	b.nodes = make(map[int]Packet)
	// set seqID
	b.seqID = seqID
	b.isNodeMode = isNode
}

// SetBytesV set bytes type as V
func (b *Encoder) SetBytesV(buf []byte) {
	b.size += int32(len(buf))
	b.valbuf.Write(buf)
	b.isStreamMode = false
	b.state |= 0x04
}

// SetReaderV set io.Reader type as V
func (b *Encoder) SetReaderV(r io.Reader, size int) {
	b.isStreamMode = true
	b.valReader = r
	b.state |= 0x04
	b.size += int32(size)
	b.valReaderSize = size
}

// AddPacket add a Bhojpur Service Packet child to this packet, this packet must be NodeMode
func (b *Encoder) AddPacket(child Packet) error {
	// only packet is in node mode can add other packets
	if !b.isNodeMode {
		return errNotNodeMode
	}

	if b.done {
		return errInvalidAdding
	}
	b.nodes[child.SeqID()] = child
	buf := child.Bytes()
	b.SetBytesV(buf)
	return nil
}

// AddStreamPacket will put a StreamPacket in chunked mode to current packet.
func (b *Encoder) AddStreamPacket(child Packet) (err error) {
	// if this packet is in stream mode, can not add any packets
	if b.done {
		return errInvalidAdding
	}

	// only accept packet in stream mode
	if !child.IsStreamMode() {
		return errNonStreamPacket
	}

	// set the valReader of this packet to the child's
	b.valReader = child.VReader()

	// valReaderSize will be the same as child's
	b.valReaderSize = child.VSize()
	// add this child packet
	b.nodes[child.SeqID()] = child
	// add the size of child's V to L of this packet
	b.size += int32(child.Size())
	// put the bytes of child to valbuf
	buf := child.Bytes()
	b.valbuf.Write(buf)
	// update state
	b.state |= 0x04
	b.isStreamMode = true
	b.done = true
	return nil
}

// Packet return a Bhojpur Service Packet instance.
func (b *Encoder) Packet() (Packet, error) {
	err := b.generateT()
	if err != nil {
		return nil, err
	}

	err = b.generateL()
	if err != nil {
		return nil, err
	}

	if b.state != 0x07 {
		return nil, errBuildIncomplete
	}

	if b.isStreamMode {
		return &StreamPacket{
			t:         b.tag,
			l:         *b.len,
			vr:        b.valReader,
			vbuf:      b.valbuf.Bytes(),
			chunkMode: true,
			chunkSize: b.valReaderSize,
		}, err
	}

	// not streaming mode
	return &StreamPacket{
		t:         b.tag,
		l:         *b.len,
		vbuf:      b.valbuf.Bytes(),
		chunkMode: false,
	}, err
}

// will generate T of a TLV.
func (b *Encoder) generateT() error {
	t, err := spec.NewT(b.seqID)
	t.SetNodeMode(b.isNodeMode)
	if err != nil {
		return err
	}
	b.tag = t
	b.state |= 0x01
	return nil
}

// will generate L of a TLV.
func (b *Encoder) generateL() error {
	l, err := spec.NewL(int(b.size))
	if err != nil {
		return err
	}
	b.len = &l
	b.state |= 0x02
	return nil
}
