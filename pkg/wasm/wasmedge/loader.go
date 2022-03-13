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
// #include <stdlib.h>
import "C"
import (
	"unsafe"
)

type Loader struct {
	_inner *C.WasmEdge_LoaderContext
	_own   bool
}

func NewLoader() *Loader {
	loader := C.WasmEdge_LoaderCreate(nil)
	if loader == nil {
		return nil
	}
	return &Loader{_inner: loader, _own: true}
}

func NewLoaderWithConfig(conf *Configure) *Loader {
	loader := C.WasmEdge_LoaderCreate(conf._inner)
	if loader == nil {
		return nil
	}
	return &Loader{_inner: loader, _own: true}
}

func (self *Loader) LoadFile(path string) (*AST, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	var module *C.WasmEdge_ASTModuleContext = nil
	result := C.WasmEdge_LoaderParseFromFile(self._inner, &module, cpath)
	if !C.WasmEdge_ResultOK(result) {
		return nil, newError(result)
	}
	return &AST{_inner: module, _own: true}, nil
}

func (self *Loader) LoadBuffer(buf []byte) (*AST, error) {
	var module *C.WasmEdge_ASTModuleContext = nil
	result := C.WasmEdge_LoaderParseFromBuffer(self._inner, &module, (*C.uint8_t)(unsafe.Pointer(&buf[0])), C.uint32_t(len(buf)))
	if !C.WasmEdge_ResultOK(result) {
		return nil, newError(result)
	}
	return &AST{_inner: module, _own: true}, nil
}

func (self *Loader) Release() {
	if self._own {
		C.WasmEdge_LoaderDelete(self._inner)
	}
	self._inner = nil
	self._own = false
}
