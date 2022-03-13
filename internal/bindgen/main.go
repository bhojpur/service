package main

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

import (
	"fmt"
	"os"

	bindgen "github.com/bhojpur/service/pkg/wasm/bindgen"
	edge "github.com/bhojpur/service/pkg/wasm/wasmedge"
)

func main() {
	/// Expected Args[0]: program name (./bindgen_funcs)
	/// Expected Args[1]: wasm file (rust_bindgen_funcs_lib.wasm))

	if len(os.Args) <= 2 {
		fmt.Println("Bhojpur Service: bindgen requires .wasm stream function library filename")
		return
	}
	/// Set not to print debug info
	edge.SetLogErrorLevel()

	/// Create configure
	var conf = edge.NewConfigure(edge.WASI)

	/// Create VM with configure
	var vm = edge.NewVMWithConfig(conf)

	/// Init WASI
	var wasi = vm.GetImportObject(edge.WASI)
	wasi.InitWasi(
		os.Args[1:],     /// The args
		os.Environ(),    /// The envs
		[]string{".:."}, /// The mapping preopens
	)

	/// Load and validate the wasm
	vm.LoadWasmFile(os.Args[1])
	vm.Validate()

	// Instantiate the bindgen and vm
	bg := bindgen.Instantiate(vm)

	/// create_line: string, string, string -> string (inputs are JSON stringified)
	res, err := bg.Execute("create_line", "{\"x\":2.5,\"y\":7.8}", "{\"x\":2.5,\"y\":5.8}", "A thin red line")
	if err == nil {
		fmt.Println("run bindgen -- create_line:", res[0].(string))
	} else {
		fmt.Println("run bindgen -- create_line FAILED", err)
	}

	/// say: string -> string
	res, err = bg.Execute("say", "bindgen funcs test")
	if err == nil {
		fmt.Println("run bindgen -- say:", res[0].(string))
	} else {
		fmt.Println("run bindgen -- say FAILED")
	}

	/// obfusticate: string -> string
	res, err = bg.Execute("obfusticate", "A quick brown fox jumps over the lazy dog")
	if err == nil {
		fmt.Println("run bindgen -- obfusticate:", res[0].(string))
	} else {
		fmt.Println("run bindgen -- obfusticate FAILED")
	}

	/// lowest_common_multiple: i32, i32 -> i32
	res, err = bg.Execute("lowest_common_multiple", int32(123), int32(2))
	if err == nil {
		fmt.Println("run bindgen -- lowest_common_multiple:", res[0].(int32))
	} else {
		fmt.Println("run bindgen -- lowest_common_multiple FAILED")
	}

	/// sha3_digest: array -> array
	res, err = bg.Execute("sha3_digest", []byte("This is an important message"))
	if err == nil {
		fmt.Println("run bindgen -- sha3_digest:", res[0].([]byte))
	} else {
		fmt.Println("run bindgen -- sha3_digest FAILED")
	}

	/// keccak_digest: array -> array
	res, err = bg.Execute("keccak_digest", []byte("This is an important message"))
	if err == nil {
		fmt.Println("run bindgen -- keccak_digest:", res[0].([]byte))
	} else {
		fmt.Println("run bindgen -- keccak_digest FAILED")
	}

	bg.Release()
	vm.Release()
	conf.Release()
}
