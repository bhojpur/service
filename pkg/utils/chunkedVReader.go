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
	"io/ioutil"
)

type chunkVReader struct {
	src        io.Reader     // the reader parts of V
	buf        *bytes.Buffer // the bytes parts of V
	totalSize  int           // size of whole buffer of this packet
	off        int           // last read op
	ChunkVSize int           // the size of chunked V
}

// Read implement io.Reader interface
func (r *chunkVReader) Read(p []byte) (n int, err error) {
	if r.src == nil {
		return 0, nil
	}

	if r.off >= r.totalSize {
		return 0, io.EOF
	}

	if r.off < r.totalSize-r.ChunkVSize {
		n, err := r.buf.Read(p)
		r.off += n
		if err != nil {
			if err == io.EOF {
				return n, nil
			} else {
				return 0, err
			}
		}
		return n, nil
	}
	n, err = r.src.Read(p)
	r.off += n
	if err != nil {
		return n, err
	}
	return n, nil
}

// WriteTo implement io.WriteTo interface
func (r *chunkVReader) WriteTo(w io.Writer) (n int64, err error) {
	if r.src == nil {
		return 0, nil
	}

	// first, write existed buffer
	m, err := w.Write(r.buf.Bytes())
	if err != nil {
		return 0, err
	}
	n += int64(m)

	// last, write from reader
	buf, err := ioutil.ReadAll(r.src)
	if err != nil && err != io.EOF {
		return 0, errWriteFromReader
	}
	m, err = w.Write(buf)
	if err != nil {
		return 0, err
	}

	n += int64(m)
	return n, nil
}
