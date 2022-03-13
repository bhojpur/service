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

// TableType classifies tables over elements of element types within a size range.
//
// See also
//
// Specification: https://webassembly.github.io/spec/core/syntax/types.html#table-types
//
type TableType struct {
	_inner   *C.wasm_tabletype_t
	_ownedBy interface{}
}

func newTableType(pointer *C.wasm_tabletype_t, ownedBy interface{}) *TableType {
	tableType := &TableType{_inner: pointer, _ownedBy: ownedBy}

	if ownedBy == nil {
		runtime.SetFinalizer(tableType, func(tableType *TableType) {
			C.wasm_tabletype_delete(tableType.inner())
		})
	}

	return tableType
}

// NewTableType instantiates a new TableType given a ValueType and some Limits.
//
//   valueType := NewValueType(I32)
//   limits := NewLimits(1, 4)
//   tableType := NewTableType(valueType, limits)
//   _ = tableType.IntoExternType()
//
func NewTableType(valueType *ValueType, limits *Limits) *TableType {
	pointer := C.wasm_tabletype_new(valueType.inner(), limits.inner())

	return newTableType(pointer, nil)
}

func (self *TableType) inner() *C.wasm_tabletype_t {
	return self._inner
}

func (self *TableType) ownedBy() interface{} {
	if self._ownedBy == nil {
		return self
	}

	return self._ownedBy
}

// ValueType returns the TableType's ValueType.
//
//   valueType := NewValueType(I32)
//   limits := NewLimits(1, 4)
//   tableType := NewTableType(valueType, limits)
//   _ = tableType.ValueType()
//
func (self *TableType) ValueType() *ValueType {
	pointer := C.wasm_tabletype_element(self.inner())

	runtime.KeepAlive(self)

	return newValueType(pointer, self.ownedBy())
}

// Limits returns the TableType's Limits.
//
//   valueType := NewValueType(I32)
//   limits := NewLimits(1, 4)
//   tableType := NewTableType(valueType, limits)
//   _ = tableType.Limits()
//
func (self *TableType) Limits() *Limits {
	limits := newLimits(C.wasm_tabletype_limits(self.inner()), self.ownedBy())

	runtime.KeepAlive(self)

	return limits
}

// IntoExternType converts the TableType into an ExternType.
//
//   valueType := NewValueType(I32)
//   limits := NewLimits(1, 4)
//   tableType := NewTableType(valueType, limits)
//   _ = tableType.IntoExternType()
//
func (self *TableType) IntoExternType() *ExternType {
	pointer := C.wasm_tabletype_as_externtype_const(self.inner())

	return newExternType(pointer, self.ownedBy())
}
