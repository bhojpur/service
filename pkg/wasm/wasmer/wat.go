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
	"unsafe"
)

// Wat2Wasm parses a string as either WAT code or a binary Wasm module.
//
// See https://webassembly.github.io/spec/core/text/index.html.
//
// Note: This is not part of the standard Wasm C API. It is Wasmer specific.
//
//   wat := "(module)"
//   wasm, _ := Wat2Wasm(wat)
//   engine := wasmer.NewEngine()
//   store := wasmer.NewStore(engine)
//   module, _ := wasmer.NewModule(store, wasmBytes)
func Wat2Wasm(wat string) ([]byte, error) {
	var watBytes C.wasm_byte_vec_t
	var watLength = len(wat)

	C.wasm_byte_vec_new(&watBytes, C.size_t(watLength), C.CString(wat))
	defer C.wasm_byte_vec_delete(&watBytes)

	var wasm C.wasm_byte_vec_t

	err := maybeNewErrorFromWasmer(func() bool {
		C.wat2wasm(&watBytes, &wasm)

		return wasm.data == nil
	})

	if err != nil {
		return nil, err
	}

	defer C.wasm_byte_vec_delete(&wasm)

	wasmBytes := C.GoBytes(unsafe.Pointer(wasm.data), C.int(wasm.size))

	return wasmBytes, nil
}
