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
	"reflect"
	"sync"
	"unsafe"
)

type hostFunctionSignature func(data interface{}, mem *Memory, params []interface{}) ([]interface{}, Result)

type hostFunctionManager struct {
	mu sync.Mutex
	// Valid next index of map. Use and increase this index when gc is empty.
	idx uint
	// Recycled entries of map. Use entry in this slide when allocate a new host function.
	gc    []uint
	data  map[uint]interface{}
	funcs map[uint]hostFunctionSignature
}

func (self *hostFunctionManager) add(hostfunc hostFunctionSignature, hostdata interface{}) uint {
	self.mu.Lock()
	defer self.mu.Unlock()

	var realidx uint
	if len(self.gc) > 0 {
		realidx = self.gc[len(self.gc)-1]
		self.gc = self.gc[0 : len(self.gc)-1]
	} else {
		realidx = self.idx
		self.idx++
	}
	self.funcs[realidx] = hostfunc
	self.data[realidx] = hostdata
	return realidx
}

func (self *hostFunctionManager) get(i uint) (hostFunctionSignature, interface{}) {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.funcs[i], self.data[i]
}

func (self *hostFunctionManager) del(i uint) {
	self.mu.Lock()
	defer self.mu.Unlock()
	delete(self.funcs, i)
	delete(self.data, i)
	self.gc = append(self.gc, i)
}

var hostfuncMgr = hostFunctionManager{
	idx:   0,
	data:  make(map[uint]interface{}),
	funcs: make(map[uint]hostFunctionSignature),
}

//export wasmedgego_HostFuncInvokeImpl
func wasmedgego_HostFuncInvokeImpl(fn uintptr, data *C.void, mem *C.WasmEdge_MemoryInstanceContext, params *C.WasmEdge_Value, paramlen C.uint32_t, returns *C.WasmEdge_Value, returnlen C.uint32_t) C.WasmEdge_Result {
	gomem := &Memory{
		_inner: mem,
	}

	goparams := make([]interface{}, uint(paramlen))
	var cparams []C.WasmEdge_Value
	if paramlen > 0 {
		sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&cparams)))
		sliceHeader.Cap = int(paramlen)
		sliceHeader.Len = int(paramlen)
		sliceHeader.Data = uintptr(unsafe.Pointer(params))
		for i := 0; i < int(paramlen); i++ {
			goparams[i] = fromWasmEdgeValue(cparams[i])
			if cparams[i].Type == C.WasmEdge_ValType_ExternRef && !goparams[i].(ExternRef)._valid {
				panic("External reference is released")
			}
		}
	}

	gofunc, godata := hostfuncMgr.get(uint(fn))
	goreturns, err := gofunc(godata, gomem, goparams)

	var creturns []C.WasmEdge_Value
	if returnlen > 0 && goreturns != nil {
		sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&creturns)))
		sliceHeader.Cap = int(returnlen)
		sliceHeader.Len = int(returnlen)
		sliceHeader.Data = uintptr(unsafe.Pointer(returns))
		for i, val := range goreturns {
			if i < int(returnlen) {
				creturns[i] = toWasmEdgeValue(val)
			}
		}
	}

	return C.WasmEdge_Result{Code: C.uint8_t(err.code)}
}
