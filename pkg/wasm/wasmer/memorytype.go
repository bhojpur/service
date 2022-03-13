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
import "C"
import "runtime"

// MemoryType classifies linear memories and their size range.
//
// See also
//
// Specification: https://webassembly.github.io/spec/core/syntax/types.html#memory-types
//
type MemoryType struct {
	_inner   *C.wasm_memorytype_t
	_ownedBy interface{}
}

func newMemoryType(pointer *C.wasm_memorytype_t, ownedBy interface{}) *MemoryType {
	memoryType := &MemoryType{_inner: pointer, _ownedBy: ownedBy}

	if ownedBy == nil {
		runtime.SetFinalizer(memoryType, func(memoryType *MemoryType) {
			C.wasm_memorytype_delete(memoryType.inner())
		})
	}

	return memoryType
}

// NewMemoryType instantiates a new MemoryType given some Limits.
//
//   limits := NewLimits(1, 4)
//   memoryType := NewMemoryType(limits)
//
func NewMemoryType(limits *Limits) *MemoryType {
	pointer := C.wasm_memorytype_new(limits.inner())

	return newMemoryType(pointer, nil)
}

func (self *MemoryType) inner() *C.wasm_memorytype_t {
	return self._inner
}

func (self *MemoryType) ownedBy() interface{} {
	if self._ownedBy == nil {
		return self
	}

	return self._ownedBy
}

// Limits returns the MemoryType's Limits.
//
//   limits := NewLimits(1, 4)
//   memoryType := NewMemoryType(limits)
//   _ = memoryType.Limits()
//
func (self *MemoryType) Limits() *Limits {
	limits := newLimits(C.wasm_memorytype_limits(self.inner()), self.ownedBy())

	runtime.KeepAlive(self)

	return limits
}

// IntoExternType converts the MemoryType into an ExternType.
//
//   limits := NewLimits(1, 4)
//   memoryType := NewMemoryType(limits)
//   externType = memoryType.IntoExternType()
//
func (self *MemoryType) IntoExternType() *ExternType {
	pointer := C.wasm_memorytype_as_externtype_const(self.inner())

	return newExternType(pointer, self.ownedBy())
}
