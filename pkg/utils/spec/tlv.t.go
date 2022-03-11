package spec

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
	"io"
)

// T is the Tag in a TLV structure
type T byte

// NewT returns a T with sequenceID. If this packet contains other
// packets, this packet will be a "node packet", the T of this packet
// will set MSB to T.
func NewT(seqID int) (T, error) {
	if seqID < 0 || seqID > maxSeqID {
		return 0, errInvalidSeqID
	}

	return T(seqID), nil
}

// Sid returns the sequenceID of this packet.
func (t T) Sid() int {
	return int(t & wipeFlagBits)
}

// Bytes returns raw bytes of T.
func (t T) Bytes() []byte {
	return []byte{byte(t)}
}

// IsNodeMode will return true if this packet contains other packets.
// Otherwise return flase.
func (t T) IsNodeMode() bool {
	return t&flagBitNode == flagBitNode
}

// SetNodeMode will set T to indicates this packet contains
// other packets.
func (t *T) SetNodeMode(flag bool) {
	if flag {
		*t |= flagBitNode
	}
}

// Size return the size of T raw bytes.
func (t T) Size() int {
	return 1
}

// ReadT read T from a bufio.Reader
func ReadT(rd io.Reader) (T, error) {
	b, err := readByte(rd)
	return T(b), err
}
