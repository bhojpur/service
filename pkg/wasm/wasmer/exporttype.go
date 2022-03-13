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
import (
	"runtime"
	"unsafe"
)

type exportTypes struct {
	_inner      C.wasm_exporttype_vec_t
	exportTypes []*ExportType
}

func newExportTypes(module *Module) *exportTypes {
	self := &exportTypes{}
	C.wasm_module_exports(module.inner(), &self._inner)

	runtime.SetFinalizer(self, func(self *exportTypes) {
		self.close()
	})

	numberOfExportTypes := int(self.inner().size)
	types := make([]*ExportType, numberOfExportTypes)
	firstExportType := unsafe.Pointer(self.inner().data)
	sizeOfExportTypePointer := unsafe.Sizeof(firstExportType)

	var currentTypePointer *C.wasm_exporttype_t

	for nth := 0; nth < numberOfExportTypes; nth++ {
		currentTypePointer = *(**C.wasm_exporttype_t)(unsafe.Pointer(uintptr(firstExportType) + uintptr(nth)*sizeOfExportTypePointer))
		exportType := newExportType(currentTypePointer, self)
		types[nth] = exportType
	}

	self.exportTypes = types

	return self
}

func (self *exportTypes) inner() *C.wasm_exporttype_vec_t {
	return &self._inner
}

func (self *exportTypes) close() {
	runtime.SetFinalizer(self, nil)
	C.wasm_exporttype_vec_delete(&self._inner)
}

// ExportType is a descriptor for an exported WebAssembly value.
type ExportType struct {
	_inner   *C.wasm_exporttype_t
	_ownedBy interface{}
}

func newExportType(pointer *C.wasm_exporttype_t, ownedBy interface{}) *ExportType {
	exportType := &ExportType{_inner: pointer, _ownedBy: ownedBy}

	if ownedBy == nil {
		runtime.SetFinalizer(exportType, func(self *ExportType) {
			self.Close()
		})
	}

	return exportType
}

// NewExportType instantiates a new ExportType with a name and an extern type.
//
// Note: An extern type is anything implementing IntoExternType: FunctionType, GlobalType, MemoryType, TableType.
//
//   valueType := NewValueType(I32)
//   globalType := NewGlobalType(valueType, CONST)
//   exportType := NewExportType("a_global", globalType)
func NewExportType(name string, ty IntoExternType) *ExportType {
	nameName := newName(name)
	externType := ty.IntoExternType().inner()
	externTypeCopy := C.wasm_externtype_copy(externType)

	runtime.KeepAlive(externType)

	exportType := C.wasm_exporttype_new(&nameName, externTypeCopy)

	return newExportType(exportType, nil)
}

func (self *ExportType) inner() *C.wasm_exporttype_t {
	return self._inner
}

func (self *ExportType) ownedBy() interface{} {
	if self._ownedBy == nil {
		return self
	}

	return self._ownedBy
}

// Name returns the name of the export type.
//
//   exportType := NewExportType("a_global", globalType)
//   exportType.Name() // "global"
func (self *ExportType) Name() string {
	byteVec := C.wasm_exporttype_name(self.inner())
	name := C.GoStringN(byteVec.data, C.int(byteVec.size))

	runtime.KeepAlive(self)

	return name
}

// Type returns the type of the export type.
//
//   exportType := NewExportType("a_global", globalType)
//   exportType.Type() // ExternType
func (self *ExportType) Type() *ExternType {
	ty := C.wasm_exporttype_type(self.inner())

	runtime.KeepAlive(self)

	return newExternType(ty, self.ownedBy())
}

// Force to close the ExportType.
//
// A runtime finalizer is registered on the ExportType, but it is
// possible to force the destruction of the ExportType by calling
// Close manually.
func (self *ExportType) Close() {
	runtime.SetFinalizer(self, nil)
	C.wasm_exporttype_delete(self.inner())
}
