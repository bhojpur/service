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

// TableSize represents the size of a table.
type TableSize C.wasm_table_size_t

// ToUint32 converts a TableSize to a native Go uint32.
//
//   table, _ := instance.Exports.GetTable("exported_table")
//   size := table.Size().ToUint32()
func (self *TableSize) ToUint32() uint32 {
	return uint32(C.wasm_table_size_t(*self))
}

// A table instance is the runtime representation of a table. It holds
// a vector of function elements and an optional maximum size, if one
// was specified in the table type at the table’s definition site.
//
// See also
//
// Specification: https://webassembly.github.io/spec/core/exec/runtime.html#table-instances
type Table struct {
	_inner   *C.wasm_table_t
	_ownedBy interface{}
}

func newTable(pointer *C.wasm_table_t, ownedBy interface{}) *Table {
	table := &Table{_inner: pointer, _ownedBy: ownedBy}

	if ownedBy == nil {
		runtime.SetFinalizer(table, func(table *Table) {
			C.wasm_table_delete(table.inner())
		})
	}

	return table
}

func (self *Table) inner() *C.wasm_table_t {
	return self._inner
}

func (self *Table) ownedBy() interface{} {
	if self._ownedBy == nil {
		return self
	}

	return self._ownedBy
}

// Size returns the Table's size.
//
//   table, _ := instance.Exports.GetTable("exported_table")
//   size := table.Size()
//
func (self *Table) Size() TableSize {
	return TableSize(C.wasm_table_size(self.inner()))
}

// IntoExtern converts the Table into an Extern.
//
//   table, _ := instance.Exports.GetTable("exported_table")
//   extern := table.IntoExtern()
//
func (self *Table) IntoExtern() *Extern {
	pointer := C.wasm_table_as_extern(self.inner())

	return newExtern(pointer, self.ownedBy())
}
