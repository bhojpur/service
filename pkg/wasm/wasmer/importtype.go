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

type importTypes struct {
	_inner      C.wasm_importtype_vec_t
	importTypes []*ImportType
}

func newImportTypes(module *Module) *importTypes {
	self := &importTypes{}
	C.wasm_module_imports(module.inner(), &self._inner)

	runtime.SetFinalizer(self, func(self *importTypes) {
		self.close()
	})

	numberOfImportTypes := int(self.inner().size)
	types := make([]*ImportType, numberOfImportTypes)
	firstImportType := unsafe.Pointer(self.inner().data)
	sizeOfImportTypePointer := unsafe.Sizeof(firstImportType)

	var currentTypePointer *C.wasm_importtype_t

	for nth := 0; nth < numberOfImportTypes; nth++ {
		currentTypePointer = *(**C.wasm_importtype_t)(unsafe.Pointer(uintptr(firstImportType) + uintptr(nth)*sizeOfImportTypePointer))
		importType := newImportType(currentTypePointer, self)
		types[nth] = importType
	}

	self.importTypes = types

	return self
}

func (self *importTypes) inner() *C.wasm_importtype_vec_t {
	return &self._inner
}

func (self *importTypes) close() {
	runtime.SetFinalizer(self, nil)
	C.wasm_importtype_vec_delete(&self._inner)
}

// ImportType is a descriptor for an imported value into a WebAssembly
// module.
type ImportType struct {
	_inner   *C.wasm_importtype_t
	_ownedBy interface{}
}

func newImportType(pointer *C.wasm_importtype_t, ownedBy interface{}) *ImportType {
	importType := &ImportType{_inner: pointer, _ownedBy: ownedBy}

	if ownedBy == nil {
		runtime.SetFinalizer(importType, func(self *ImportType) {
			self.Close()
		})
	}

	return importType
}

// NewImportType instantiates a new ImportType with a module name (or
// namespace), a name and an extern type.
//
// Note:️ An extern type is anything implementing IntoExternType:
// FunctionType, GlobalType, MemoryType, TableType.
//
//   valueType := NewValueType(I32)
//   globalType := NewGlobalType(valueType, CONST)
//   importType := NewImportType("ns", "host_global", globalType)
//
func NewImportType(module string, name string, ty IntoExternType) *ImportType {
	moduleName := newName(module)
	nameName := newName(name)
	externType := ty.IntoExternType().inner()
	externTypeCopy := C.wasm_externtype_copy(externType)

	runtime.KeepAlive(externType)

	importType := C.wasm_importtype_new(&moduleName, &nameName, externTypeCopy)

	return newImportType(importType, nil)
}

func (self *ImportType) inner() *C.wasm_importtype_t {
	return self._inner
}

func (self *ImportType) ownedBy() interface{} {
	if self._ownedBy == nil {
		return self
	}

	return self._ownedBy
}

// Module returns the ImportType's module name (or namespace).
//
//   valueType := NewValueType(I32)
//   globalType := NewGlobalType(valueType, CONST)
//   importType := NewImportType("ns", "host_global", globalType)
//   _ = importType.Module()
//
func (self *ImportType) Module() string {
	byteVec := C.wasm_importtype_module(self.inner())
	module := C.GoStringN(byteVec.data, C.int(byteVec.size))

	runtime.KeepAlive(self)

	return module
}

// Name returns the ImportType's name.
//
//   valueType := NewValueType(I32)
//   globalType := NewGlobalType(valueType, CONST)
//   importType := NewImportType("ns", "host_global", globalType)
//   _ = importType.Name()
//
func (self *ImportType) Name() string {
	byteVec := C.wasm_importtype_name(self.inner())
	name := C.GoStringN(byteVec.data, C.int(byteVec.size))

	runtime.KeepAlive(self)

	return name
}

// Type returns the ImportType's type as an ExternType.
//
//   valueType := NewValueType(I32)
//   globalType := NewGlobalType(valueType, CONST)
//   importType := NewImportType("ns", "host_global", globalType)
//   _ = importType.Type()
//
func (self *ImportType) Type() *ExternType {
	ty := C.wasm_importtype_type(self.inner())

	runtime.KeepAlive(self)

	return newExternType(ty, self.ownedBy())
}

// Force to close the ImportType.
//
// A runtime finalizer is registered on the ImportType, but it is
// possible to force the destruction of the ImportType by calling
// Close manually.
func (self *ImportType) Close() {
	runtime.SetFinalizer(self, nil)
	C.wasm_importtype_delete(self.inner())
}
