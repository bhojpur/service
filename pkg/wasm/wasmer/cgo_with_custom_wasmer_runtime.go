//go:build custom_wasmer_runtime
// +build custom_wasmer_runtime

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

// // With the `customlib` tag, the user is expected to provide the
// // `CGO_LDFLAGS` and `CGO_CFLAGS` that point to appropriate Wasmer build
// // directories, e.g.:
// //
// // ```sh
// // export CGO_CFLAGS="-I/wasmer/lib/c-api/"
// //
// // export CGO_LDFLAGS="-Wl,-rpath,/wasmer/target/x86_64-unknown-linux-musl/release/ -L/wasmer/target/x86_64-unknown-linux-musl/release/ -pthread -lwasmer_c_api -lm -ldl -static"
// //
// // export CC=/usr/bin/musl-gcc
// //
// // go build -tags custom_wasmer_runtime
// // ```
//
// #include <wasmer.h>
import "C"
