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

type Compiler struct {
	_inner *C.WasmEdge_CompilerContext
	_own   bool
}

func NewCompiler() *Compiler {
	compiler := C.WasmEdge_CompilerCreate(nil)
	if compiler == nil {
		return nil
	}
	return &Compiler{_inner: compiler, _own: true}
}

func NewCompilerWithConfig(conf *Configure) *Compiler {
	compiler := C.WasmEdge_CompilerCreate(conf._inner)
	if compiler == nil {
		return nil
	}
	return &Compiler{_inner: compiler, _own: true}
}

func (self *Compiler) Compile(inpath string, outpath string) error {
	cinpath := C.CString(inpath)
	coutpath := C.CString(outpath)
	defer C.free(unsafe.Pointer(cinpath))
	defer C.free(unsafe.Pointer(coutpath))
	res := C.WasmEdge_CompilerCompile(self._inner, cinpath, coutpath)
	if !C.WasmEdge_ResultOK(res) {
		return newError(res)
	}
	return nil
}

func (self *Compiler) Release() {
	if self._own {
		C.WasmEdge_CompilerDelete(self._inner)
	}
	self._inner = nil
	self._own = false
}
