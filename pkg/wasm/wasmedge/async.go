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

type Async struct {
	_inner *C.WasmEdge_Async
	_own   bool
}

func (self *Async) WaitFor(millisec int) bool {
	return bool(C.WasmEdge_AsyncWaitFor(self._inner, C.uint64_t(millisec)))
}

func (self *Async) Cancel() {
	C.WasmEdge_AsyncCancel(self._inner)
}

func (self *Async) GetResult() ([]interface{}, error) {
	arity := C.WasmEdge_AsyncGetReturnsLength(self._inner)
	creturns := make([]C.WasmEdge_Value, arity)
	var ptrreturns *C.WasmEdge_Value = nil
	if len(creturns) > 0 {
		ptrreturns = (*C.WasmEdge_Value)(unsafe.Pointer(&creturns[0]))
	}
	res := C.WasmEdge_AsyncGet(self._inner, ptrreturns, arity)
	if !C.WasmEdge_ResultOK(res) {
		return nil, newError(res)
	}
	return fromWasmEdgeValueSlide(creturns), nil
}

func (self *Async) Release() {
	if self._own {
		C.WasmEdge_AsyncDelete(self._inner)
	}
	self._inner = nil
	self._own = false
}
