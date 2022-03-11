package codec

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
	"fmt"
)

// signal is builder for PrimitivePacketEncoder
type signal struct {
	encoder *PrimitivePacketEncoder
}

// createSignal create a signal
func createSignal(key byte) *signal {
	return &signal{encoder: NewPrimitivePacketEncoder(int(key))}
}

// SetString set a string Value for the signal
func (s *signal) SetString(v string) *signal {
	s.encoder.SetStringValue(v)
	return s
}

// SetString set a int64 Value for the signal
func (s *signal) SetInt64(v int64) *signal {
	s.encoder.SetInt64Value(v)
	return s
}

// SetString set a float64 Value for the signal
func (s *signal) SetFloat64(v float64) *signal {
	s.encoder.SetFloat64Value(v)
	return s
}

// ToEncoder return current PrimitivePacketEncoder, and checking legality
func (s *signal) ToEncoder(allow func(key byte) bool) *PrimitivePacketEncoder {
	if allow != nil && !allow(byte(s.encoder.seqID)) {
		panic(fmt.Errorf("it is not allowed to use this key to create a signal: %#x", byte(s.encoder.seqID)))
	}

	return s.encoder
}
