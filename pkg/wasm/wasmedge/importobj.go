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

type ImportObject struct {
	_inner     *C.WasmEdge_ImportObjectContext
	_hostfuncs []uint
	_own       bool
}

func NewImportObject(modname string) *ImportObject {
	obj := C.WasmEdge_ImportObjectCreate(toWasmEdgeStringWrap(modname))
	if obj == nil {
		return nil
	}
	return &ImportObject{_inner: obj, _own: true}
}

func NewWasiImportObject(args []string, envs []string, preopens []string) *ImportObject {
	cargs := toCStringArray(args)
	cenvs := toCStringArray(envs)
	cpreopens := toCStringArray(preopens)
	var ptrargs *(*C.char) = nil
	var ptrenvs *(*C.char) = nil
	var ptrpreopens *(*C.char) = nil
	if len(cargs) > 0 {
		ptrargs = &cargs[0]
	}
	if len(cenvs) > 0 {
		ptrenvs = &cenvs[0]
	}
	if len(cpreopens) > 0 {
		ptrpreopens = &cpreopens[0]
	}

	obj := C.WasmEdge_ImportObjectCreateWASI(
		ptrargs, C.uint32_t(len(cargs)),
		ptrenvs, C.uint32_t(len(cenvs)),
		ptrpreopens, C.uint32_t(len(cpreopens)))

	freeCStringArray(cargs)
	freeCStringArray(cenvs)
	freeCStringArray(cpreopens)

	if obj == nil {
		return nil
	}
	return &ImportObject{_inner: obj, _own: true}
}

func (self *ImportObject) InitWasi(args []string, envs []string, preopens []string) {
	cargs := toCStringArray(args)
	cenvs := toCStringArray(envs)
	cpreopens := toCStringArray(preopens)
	var ptrargs *(*C.char) = nil
	var ptrenvs *(*C.char) = nil
	var ptrpreopens *(*C.char) = nil
	if len(cargs) > 0 {
		ptrargs = &cargs[0]
	}
	if len(cenvs) > 0 {
		ptrenvs = &cenvs[0]
	}
	if len(cpreopens) > 0 {
		ptrpreopens = &cpreopens[0]
	}

	C.WasmEdge_ImportObjectInitWASI(self._inner,
		ptrargs, C.uint32_t(len(cargs)),
		ptrenvs, C.uint32_t(len(cenvs)),
		ptrpreopens, C.uint32_t(len(cpreopens)))

	freeCStringArray(cargs)
	freeCStringArray(cenvs)
	freeCStringArray(cpreopens)
}

func (self *ImportObject) WasiGetExitCode() uint {
	return uint(C.WasmEdge_ImportObjectWASIGetExitCode(self._inner))
}

func NewWasmEdgeProcessImportObject(allowedcmds []string, allowall bool) *ImportObject {
	ccmds := toCStringArray(allowedcmds)
	var ptrcmds *(*C.char) = nil
	if len(ccmds) > 0 {
		ptrcmds = &ccmds[0]
	}

	obj := C.WasmEdge_ImportObjectCreateWasmEdgeProcess(ptrcmds, C.uint32_t(len(ccmds)), C.bool(allowall))

	freeCStringArray(ccmds)

	if obj == nil {
		return nil
	}
	return &ImportObject{_inner: obj, _own: true}
}

func (self *ImportObject) InitWasmEdgeProcess(allowedcmds []string, allowall bool) {
	ccmds := toCStringArray(allowedcmds)
	var ptrcmds *(*C.char) = nil
	if len(ccmds) > 0 {
		ptrcmds = &ccmds[0]
	}

	C.WasmEdge_ImportObjectInitWasmEdgeProcess(self._inner, ptrcmds, C.uint32_t(len(ccmds)), C.bool(allowall))

	freeCStringArray(ccmds)
}

func (self *ImportObject) AddFunction(name string, inst *Function) {
	C.WasmEdge_ImportObjectAddFunction(self._inner, toWasmEdgeStringWrap(name), inst._inner)
	self._hostfuncs = append(self._hostfuncs, inst._index)
	inst._inner = nil
	inst._own = false
}

func (self *ImportObject) AddTable(name string, inst *Table) {
	C.WasmEdge_ImportObjectAddTable(self._inner, toWasmEdgeStringWrap(name), inst._inner)
	inst._inner = nil
	inst._own = false
}

func (self *ImportObject) AddMemory(name string, inst *Memory) {
	C.WasmEdge_ImportObjectAddMemory(self._inner, toWasmEdgeStringWrap(name), inst._inner)
	inst._inner = nil
	inst._own = false
}

func (self *ImportObject) AddGlobal(name string, inst *Global) {
	C.WasmEdge_ImportObjectAddGlobal(self._inner, toWasmEdgeStringWrap(name), inst._inner)
	inst._inner = nil
	inst._own = false
}

func (self *ImportObject) Release() {
	if self._own {
		for _, idx := range self._hostfuncs {
			hostfuncMgr.del(idx)
		}
		self._hostfuncs = []uint{}
		C.WasmEdge_ImportObjectDelete(self._inner)
	}
	self._inner = nil
	self._own = false
}
