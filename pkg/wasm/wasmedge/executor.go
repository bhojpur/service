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
import (
	"unsafe"
)

type Executor struct {
	_inner *C.WasmEdge_ExecutorContext
	_own   bool
}

func NewExecutor() *Executor {
	executor := C.WasmEdge_ExecutorCreate(nil, nil)
	if executor == nil {
		return nil
	}
	return &Executor{_inner: executor, _own: true}
}

func NewExecutorWithConfig(conf *Configure) *Executor {
	executor := C.WasmEdge_ExecutorCreate(conf._inner, nil)
	if executor == nil {
		return nil
	}
	return &Executor{_inner: executor, _own: true}
}

func NewExecutorWithStatistics(stat *Statistics) *Executor {
	executor := C.WasmEdge_ExecutorCreate(nil, stat._inner)
	if executor == nil {
		return nil
	}
	return &Executor{_inner: executor, _own: true}
}

func NewExecutorWithConfigAndStatistics(conf *Configure, stat *Statistics) *Executor {
	executor := C.WasmEdge_ExecutorCreate(conf._inner, stat._inner)
	if executor == nil {
		return nil
	}
	return &Executor{_inner: executor, _own: true}
}

func (self *Executor) Instantiate(store *Store, ast *AST) error {
	res := C.WasmEdge_ExecutorInstantiate(self._inner, store._inner, ast._inner)
	if !C.WasmEdge_ResultOK(res) {
		return newError(res)
	}
	return nil
}

func (self *Executor) RegisterImport(store *Store, imp *ImportObject) error {
	res := C.WasmEdge_ExecutorRegisterImport(self._inner, store._inner, imp._inner)
	if !C.WasmEdge_ResultOK(res) {
		return newError(res)
	}
	return nil
}

func (self *Executor) RegisterModule(store *Store, ast *AST, modname string) error {
	modstr := toWasmEdgeStringWrap(modname)
	res := C.WasmEdge_ExecutorRegisterModule(self._inner, store._inner, ast._inner, modstr)
	if !C.WasmEdge_ResultOK(res) {
		return newError(res)
	}
	return nil
}

func (self *Executor) Invoke(store *Store, funcname string, params ...interface{}) ([]interface{}, error) {
	funcstr := toWasmEdgeStringWrap(funcname)
	funccxt := store.FindFunction(funcname)
	var ftype *FunctionType
	if funccxt == nil {
		// If find function failed, set function type as NULL and keep running to let the Executor to handle the error.
		ftype = &FunctionType{_inner: nil, _own: false}
	} else {
		ftype = funccxt.GetFunctionType()
	}
	cparams := toWasmEdgeValueSlide(params...)
	creturns := make([]C.WasmEdge_Value, ftype.GetReturnsLength())
	var ptrparams *C.WasmEdge_Value = nil
	var ptrreturns *C.WasmEdge_Value = nil
	if len(cparams) > 0 {
		ptrparams = (*C.WasmEdge_Value)(unsafe.Pointer(&cparams[0]))
	}
	if len(creturns) > 0 {
		ptrreturns = (*C.WasmEdge_Value)(unsafe.Pointer(&creturns[0]))
	}
	res := C.WasmEdge_ExecutorInvoke(
		self._inner, store._inner, funcstr,
		ptrparams, C.uint32_t(len(cparams)),
		ptrreturns, C.uint32_t(len(creturns)))
	if !C.WasmEdge_ResultOK(res) {
		return nil, newError(res)
	}
	return fromWasmEdgeValueSlide(creturns), nil
}

func (self *Executor) InvokeRegistered(store *Store, modname string, funcname string, params ...interface{}) ([]interface{}, error) {
	modstr := toWasmEdgeStringWrap(modname)
	funcstr := toWasmEdgeStringWrap(funcname)
	funccxt := store.FindFunctionRegistered(modname, funcname)
	var ftype *FunctionType
	if funccxt == nil {
		// If find function failed, set function type as NULL and keep running to let the Executor to handle the error.
		ftype = &FunctionType{_inner: nil, _own: false}
	} else {
		ftype = funccxt.GetFunctionType()
	}
	cparams := toWasmEdgeValueSlide(params...)
	creturns := make([]C.WasmEdge_Value, ftype.GetReturnsLength())
	var ptrparams *C.WasmEdge_Value = nil
	var ptrreturns *C.WasmEdge_Value = nil
	if len(cparams) > 0 {
		ptrparams = (*C.WasmEdge_Value)(unsafe.Pointer(&cparams[0]))
	}
	if len(creturns) > 0 {
		ptrreturns = (*C.WasmEdge_Value)(unsafe.Pointer(&creturns[0]))
	}
	res := C.WasmEdge_ExecutorInvokeRegistered(
		self._inner, store._inner, modstr, funcstr,
		ptrparams, C.uint32_t(len(cparams)),
		ptrreturns, C.uint32_t(len(creturns)))
	if !C.WasmEdge_ResultOK(res) {
		return nil, newError(res)
	}
	return fromWasmEdgeValueSlide(creturns), nil
}

func (self *Executor) Release() {
	if self._own {
		C.WasmEdge_ExecutorDelete(self._inner)
	}
	self._inner = nil
	self._own = false
}
