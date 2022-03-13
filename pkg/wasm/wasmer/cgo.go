//go:build !custom_wasmer_runtime
// +build !custom_wasmer_runtime

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

// #cgo CFLAGS: -I${SRCDIR}/include
// #cgo LDFLAGS: -lwasmer
//
// #cgo linux,amd64 LDFLAGS: -Wl,-rpath,${SRCDIR}/lib/linux-amd64 -L${SRCDIR}/lib/linux-amd64
// //#cgo linux,arm64 LDFLAGS: -Wl,-rpath,${SRCDIR}/lib/linux-aarch64 -L${SRCDIR}/lib/linux-aarch64
// #cgo darwin,amd64 LDFLAGS: -Wl,-rpath,${SRCDIR}/lib/darwin-amd64 -L${SRCDIR}/lib/darwin-amd64
// #cgo darwin,arm64 LDFLAGS: -Wl,-rpath,${SRCDIR}/lib/darwin-aarch64 -L${SRCDIR}/lib/darwin-aarch64
//
// #include <wasmer.h>
import "C"

import (
	_ "github.com/bhojpur/service/pkg/wasm/wasmer/include"
	_ "github.com/bhojpur/service/pkg/wasm/wasmer/lib"
)
