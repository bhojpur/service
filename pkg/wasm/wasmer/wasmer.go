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

// It is a complete and mature WebAssembly runtime.
//
// Key Features
//
// • Easy to use: The Wasmer API mimics the standard WebAssembly API,
// • Fast: Wasmer executes the WebAssembly modules as fast as possible, close to native speed,
// • Safe: All calls to WebAssembly will be fast, but more importantly, complete safe and sandboxed.
//
// Quick Introduction
//
// The wasmer brings the required API to execute WebAssembly modules. In a nutshell,
// wasmer compiles the WebAssembly module into compiled code, and then executes it.
// wasmer is designed to work in various environments and platforms: From nano
// single-board computers to large and powerful servers, including more exotic
// ones. To address those requirements, Wasmer provides two engines and three
// compilers.
//
// Succinctly, an engine is responsible to drive the compilation and the execution
// of a WebAssembly module. By extension, a headless engine can only execute a
// WebAssembly module, i.e. a module that has previously been compiled, or compiled,
// serialized and deserialized. By default, the wasmer package comes with 2 headless
// engines:
//
// • Universal, the compiled machine code lives in memory,
//
// • Dylib, the compiled machine code lives in a shared object file (.so, .dylib,
// or .dll), and is natively executed.
//
// The wasmer comes with three compilers:
//
// • Singlepass: Super fast compilation times, slower execution times. Not prone
// to JIT-bombs. Ideal for blockchains.
//
// • Cranelift: Fast compilation times, fast execution times. Ideal for development.
//
// • LLVM: Slow compilation times, very fast execution times (close to native).
// Ideal for production use.
//
// WebAssembly C API standard
//
// Wasmer —the runtime— is written in Rust; C and C++ bindings exist. This package
// relies on the so-called wasm_c_api, https://github.com/WebAssembly/wasm-c-api,
// which is the new standard C API, implemented inside Wasmer as the Wasmer C API,
// https://wasmerio.github.io/wasmer/crates/wasmer_c_api/. This standard is
// characterized as a living standard. The API is not yet stable, even though it
// shows maturity overtime. However, the Wasmer C API provides some extensions,
// like the wasi_* or wasmer_* types and functions, which aren't yet defined by
// the standard. It commits to keep a semantic versioning over the API, regardless
// what happens with the C API.
//
// Examples
//
// The very basic example is the following
//
//	// Let's assume we don't have WebAssembly bytes at hand. We will write
//  // WebAssembly manually.
//	wasmBytes := []byte(`
//		(module
//		  (type (func (param i32 i32) (result i32)))
//		  (func (type 0)
//		    local.get 0
//		    local.get 1
//		    i32.add)
//		  (export "sum" (func 0)))
//	`)
//
//	// Create an Engine
//	engine := wasmer.NewEngine()
//
//	// Create a Store
//	store := wasmer.NewStore(engine)
//
//	// Let's compile the module.
//	module, err := wasmer.NewModule(store, wasmBytes)
//
//	if err != nil {
//		fmt.Println("Failed to compile module:", err)
//	}
//
//	// Create an empty import object.
//	importObject := wasmer.NewImportObject()
//
//	// Let's instantiate the WebAssembly module.
//	instance, err := wasmer.NewInstance(module, importObject)
//
//	if err != nil {
//		panic(fmt.Sprintln("Failed to instantiate the module:", err))
//	}
//
//	// Now let's execute the `sum` function.
//	sum, err := instance.Exports.GetFunction("sum")
//
//	if err != nil {
//		panic(fmt.Sprintln("Failed to get the `add_one` function:", err))
//	}
//
//	result, err := sum(1, 2)
//
//	if err != nil {
//		panic(fmt.Sprintln("Failed to call the `add_one` function:", err))
//	}
//
//	fmt.Println("Results of `sum`:", result)
//
//	// Output:
//	// Results of `sum`: 3
//
// That's it. Now explore the API! Some pointers for the adventurers:
//
// • The basic elements are Module and Instance,
//
// • Exports of an instance are represented by the Exports type,
//
// • Maybe your module needs to import Function, Memory, Global or Table? Well,
// there is the ImportObject for that!
