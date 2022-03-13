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
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"sync/atomic"

	svcsvr "github.com/bhojpur/service/pkg/engine"
	bindgen "github.com/bhojpur/service/pkg/wasm/bindgen"
	edge "github.com/bhojpur/service/pkg/wasm/wasmedge"
)

var (
	counter uint64
)

const ImageDataKey = 0x10

func main() {
	// Connect to the Bhojpur Service-Processor function
	sfn := svcsvr.NewStreamFunction("image-recognition", svcsvr.WithProcessorAddr("localhost:9900"))
	defer sfn.Close()

	// set only monitoring data
	sfn.SetObserveDataTags(ImageDataKey)

	// set handler
	sfn.SetHandler(Handler)

	// start
	err := sfn.Connect()
	if err != nil {
		log.Print("❌ Connect to the Bhojpur Service-Processor function failure: ", err)
		os.Exit(1)
	}

	select {}
}

// Handler process the data in the stream
func Handler(img []byte) (byte, []byte) {
	// Initialize WasmEdge's VM
	vmConf, vm := initVM()
	bg := bindgen.Instantiate(vm)
	defer bg.Release()
	defer vm.Release()
	defer vmConf.Release()

	// recognize the image
	res, err := bg.Execute("infer", img)
	if err == nil {
		fmt.Println("Go: run bindgen -- infer:", res)
	} else {
		fmt.Println("Go: run bindgen -- infer FAILED")
	}

	// print logs
	hash := genSha1(img)
	log.Printf("✅ received image-%d hash %v, img_size=%d \n", atomic.AddUint64(&counter, 1), hash, len(img))

	return 0x11, nil
}

// genSha1 generate the hash value of the image
func genSha1(buf []byte) string {
	h := sha1.New()
	h.Write(buf)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// initVM initialize WasmEdge's VM
func initVM() (*edge.Configure, *edge.VM) {
	edge.SetLogErrorLevel()
	/// Set Tensorflow not to print debug info
	os.Setenv("TF_CPP_MIN_LOG_LEVEL", "3")
	os.Setenv("TF_CPP_MIN_VLOG_LEVEL", "3")

	/// Create configure
	vmConf := edge.NewConfigure(edge.WASI)

	/// Create VM with configure
	vm := edge.NewVMWithConfig(vmConf)

	/// Init WASI
	var wasi = vm.GetImportObject(edge.WASI)
	wasi.InitWasi(
		os.Args[1:],     /// The args
		os.Environ(),    /// The envs
		[]string{".:."}, /// The mapping directories
	)

	/// Register WasmEdge-tensorflow and WasmEdge-image
	var tfobj = edge.NewTensorflowImportObject()
	var tfliteobj = edge.NewTensorflowLiteImportObject()
	vm.RegisterImport(tfobj)
	vm.RegisterImport(tfliteobj)
	var imgobj = edge.NewImageImportObject()
	vm.RegisterImport(imgobj)

	/// Instantiate wasm
	vm.LoadWasmFile("rust_mobilenet_food_lib.so")
	vm.Validate()

	return vmConf, vm
}
