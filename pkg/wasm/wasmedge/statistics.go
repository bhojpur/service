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

type Statistics struct {
	_inner *C.WasmEdge_StatisticsContext
	_own   bool
}

func NewStatistics() *Statistics {
	stat := C.WasmEdge_StatisticsCreate()
	if stat == nil {
		return nil
	}
	return &Statistics{_inner: stat, _own: true}
}

func (self *Statistics) GetInstrCount() uint {
	return uint(C.WasmEdge_StatisticsGetInstrCount(self._inner))
}

func (self *Statistics) GetInstrPerSecond() float64 {
	return float64(C.WasmEdge_StatisticsGetInstrPerSecond(self._inner))
}

func (self *Statistics) GetTotalCost() uint {
	return uint(C.WasmEdge_StatisticsGetTotalCost(self._inner))
}

func (self *Statistics) SetCostTable(table []uint64) {
	var ptr *uint64 = nil
	if len(table) > 0 {
		ptr = &(table[0])
	}
	C.WasmEdge_StatisticsSetCostTable(self._inner, (*C.uint64_t)(ptr), C.uint32_t(len(table)))
}

func (self *Statistics) SetCostLimit(limit uint) {
	C.WasmEdge_StatisticsSetCostLimit(self._inner, C.uint64_t(limit))
}

func (self *Statistics) Release() {
	if self._own {
		C.WasmEdge_StatisticsDelete(self._inner)
	}
	self._inner = nil
	self._own = false
}
