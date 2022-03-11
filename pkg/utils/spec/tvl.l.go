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
	"bytes"
	"errors"
	"io"

	"github.com/bhojpur/service/pkg/utils/encoding"
)

// L is the Length in a TLV structure
type L struct {
	buf  []byte
	size int
	len  int
}

// NewL will take an int type len as parameter and return L to
// represent the sieze of V in a TLV. an integer will be encode as
// a PVarInt32 type to represent the value.
func NewL(len int) (L, error) {
	var l = L{}
	if len < -1 {
		return l, errors.New("service.L: len can't less than -1")
	}

	vallen := int32(len)
	l.size = encoding.SizeOfPVarInt32(vallen)
	codec := encoding.VarCodec{Size: l.size}
	tmp := make([]byte, l.size)
	err := codec.EncodePVarInt32(tmp, vallen)
	if err != nil {
		panic(err)
	}
	l.buf = make([]byte, l.size)
	copy(l.buf, tmp)
	l.len = len
	return l, nil
}

// Bytes will return the raw bytes of L.
func (l L) Bytes() []byte {
	return l.buf
}

// Size returns how many bytes used to represent this L.
func (l L) Size() int {
	return l.size
}

// VSize returns the size of V.
func (l L) VSize() int {
	return int(l.len)
}

// ReadL read L from bufio.Reader
func ReadL(r io.Reader) (*L, error) {
	lenbuf := bytes.Buffer{}
	for {
		b, err := readByte(r)
		if err != nil {
			return nil, err
		}
		lenbuf.WriteByte(b)
		if b&msb != msb {
			break
		}
	}

	buf := lenbuf.Bytes()

	// decode to L
	length, err := decodeL(buf)
	if err != nil {
		return nil, err
	}

	return &L{
		buf:  buf,
		len:  int(length),
		size: len(buf),
	}, nil
}

func decodeL(buf []byte) (length int32, err error) {
	codec := encoding.VarCodec{}
	err = codec.DecodePVarInt32(buf, &length)
	return length, err
}
