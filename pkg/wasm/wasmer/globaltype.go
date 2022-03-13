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

type GlobalMutability C.wasm_mutability_t

const (
	// Represents a global that is immutable.
	IMMUTABLE = GlobalMutability(C.WASM_CONST)

	// Represents a global that is mutable.
	MUTABLE = GlobalMutability(C.WASM_VAR)
)

// String returns the GlobalMutability as a string.
//
//   IMMUTABLE.String() // "const"
//   MUTABLE.String()   // "var"
//
func (self GlobalMutability) String() string {
	switch self {
	case IMMUTABLE:
		return "const"
	case MUTABLE:
		return "var"
	}
	panic("Unknown mutability") // unreachable
}

// GlobalType classifies global variables, which hold a value and can either be mutable or immutable.
//
// See also
//
// Specification: https://webassembly.github.io/spec/core/syntax/types.html#global-types
//
type GlobalType struct {
	_inner   *C.wasm_globaltype_t
	_ownedBy interface{}
}

func newGlobalType(pointer *C.wasm_globaltype_t, ownedBy interface{}) *GlobalType {
	globalType := &GlobalType{_inner: pointer, _ownedBy: ownedBy}

	if ownedBy == nil {
		runtime.SetFinalizer(globalType, func(globalType *GlobalType) {
			C.wasm_globaltype_delete(globalType.inner())
		})
	}

	return globalType
}

// NewGlobalType instantiates a new GlobalType from a ValueType and a GlobalMutability
//
//   valueType := NewValueType(I32)
//   globalType := NewGlobalType(valueType, IMMUTABLE)
//
func NewGlobalType(valueType *ValueType, mutability GlobalMutability) *GlobalType {
	pointer := C.wasm_globaltype_new(valueType.inner(), C.wasm_mutability_t(mutability))

	return newGlobalType(pointer, nil)
}

func (self *GlobalType) inner() *C.wasm_globaltype_t {
	return self._inner
}

func (self *GlobalType) ownedBy() interface{} {
	if self._ownedBy == nil {
		return self
	}

	return self._ownedBy
}

// ValueType returns the GlobalType's ValueType
//
//   valueType := NewValueType(I32)
//   globalType := NewGlobalType(valueType, IMMUTABLE)
//   globalType.ValueType().Kind().String() // "i32"
//
func (self *GlobalType) ValueType() *ValueType {
	pointer := C.wasm_globaltype_content(self.inner())

	runtime.KeepAlive(self)

	return newValueType(pointer, self.ownedBy())
}

// Mutability returns the GlobalType's GlobalMutability
//
//   valueType := NewValueType(I32)
//   globalType := NewGlobalType(valueType, IMMUTABLE)
//   globalType.Mutability().String() // "const"
//
func (self *GlobalType) Mutability() GlobalMutability {
	mutability := GlobalMutability(C.wasm_globaltype_mutability(self.inner()))

	runtime.KeepAlive(self)

	return mutability
}

// IntoExternType converts the GlobalType into an ExternType.
//
//   valueType := NewValueType(I32)
//   globalType := NewGlobalType(valueType, IMMUTABLE)
//   externType = globalType.IntoExternType()
//
func (self *GlobalType) IntoExternType() *ExternType {
	pointer := C.wasm_globaltype_as_externtype_const(self.inner())

	return newExternType(pointer, self.ownedBy())
}
