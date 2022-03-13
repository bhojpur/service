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

// Units of WebAssembly pages (as specified to be 65,536 bytes).
type Pages C.wasm_memory_pages_t

// Represents a memory page size.
const WasmPageSize = uint(0x10000)

// Represents the maximum number of pages.
const WasmMaxPages = uint(0x10000)

// Represents the minimum number of pages.
const WasmMinPages = uint(0x100)

// ToUint32 converts a Pages to a native Go uint32 which is the Pages' size.
//
//   memory, _ := instance.Exports.GetMemory("exported_memory")
//   size := memory.Size().ToUint32()
//
func (self *Pages) ToUint32() uint32 {
	return uint32(C.wasm_memory_pages_t(*self))
}

// ToBytes converts a Pages to a native Go uint which is the Pages' size in bytes.
//
//   memory, _ := instance.Exports.GetMemory("exported_memory")
//   size := memory.Size().ToBytes()
//
func (self *Pages) ToBytes() uint {
	return uint(self.ToUint32()) * WasmPageSize
}
