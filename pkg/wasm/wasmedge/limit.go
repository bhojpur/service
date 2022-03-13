package wasmedge

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

// #include <wasmedge/wasmedge.h>
import "C"

type Limit struct {
	min    uint
	max    uint
	hasmax bool
}

func NewLimit(minVal uint) *Limit {
	l := &Limit{
		min:    minVal,
		max:    minVal,
		hasmax: false,
	}
	return l
}

func NewLimitWithMax(minVal uint, maxVal uint) *Limit {
	if maxVal >= minVal {
		return &Limit{
			min:    minVal,
			max:    maxVal,
			hasmax: true,
		}
	}
	return nil
}

func (l *Limit) HasMax() bool {
	return l.hasmax
}

func (l *Limit) GetMin() uint {
	return l.min
}

func (l *Limit) GetMax() uint {
	return l.max
}
