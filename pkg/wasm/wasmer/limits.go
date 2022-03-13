package wasmer

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

// #include <wasmer.h>
//
// uint32_t to_limit_max_unbound() {
//     return wasm_limits_max_default;
// }
import "C"
import "runtime"

// LimitMaxUnbound returns the value used to represent an unbound
// limit, i.e. when a limit only has a min but not a max. See Limit.
func LimitMaxUnbound() uint32 {
	return uint32(C.to_limit_max_unbound())
}

// Limits classify the size range of resizeable storage associated
// with memory types and table types.
//
// See also
//
// Specification: https://webassembly.github.io/spec/core/syntax/types.html#limits
type Limits struct {
	_inner C.wasm_limits_t
}

func newLimits(pointer *C.wasm_limits_t, ownedBy interface{}) *Limits {
	limits, err := NewLimits(uint32(pointer.min), uint32(pointer.max))

	if err != nil {
		return nil
	}

	if ownedBy != nil {
		runtime.KeepAlive(ownedBy)
	}

	return limits
}

// NewLimits instantiates a new Limits which describes the Memory used.
// The minimum and maximum parameters are "number of memory pages".
//
// ️Note: Each page is 64 KiB in size.
//
// Note: You cannot Memory.Grow the Memory beyond the maximum defined here.
func NewLimits(minimum uint32, maximum uint32) (*Limits, error) {
	if minimum > maximum {
		return nil, newErrorWith("The minimum limit is greater than the maximum one")
	}

	return &Limits{
		_inner: C.wasm_limits_t{
			min: C.uint32_t(minimum),
			max: C.uint32_t(maximum),
		},
	}, nil
}

func (self *Limits) inner() *C.wasm_limits_t {
	return &self._inner
}

// Minimum returns the minimum size of the Memory allocated in "number of pages".
//
// Note:️ Each page is 64 KiB in size.
func (self *Limits) Minimum() uint32 {
	return uint32(self.inner().min)
}

// Maximum returns the maximum size of the Memory allocated in "number of pages".
//
// Each page is 64 KiB in size.
//
// Note: You cannot Memory.Grow beyond this defined maximum size.
func (self *Limits) Maximum() uint32 {
	return uint32(self.inner().max)
}
