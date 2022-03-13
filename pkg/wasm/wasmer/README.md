# Bhojpur Service - Wasmer Integration

A complete and mature WebAssembly runtime for based on [Wasmer](https://github.com/wasmerio/wasmer).

## Key Features

- **Easy to use**: The `wasmer` API mimics the standard WebAssembly API,
- **Fast**: `wasmer` executes the WebAssembly modules as fast as possible, close to **native speed**,
- **Safe**: All calls to WebAssembly will be fast, but more importantly, completely safe and sandboxed.

## Installation

To install the library, follow the classical:

```sh
$ go get github.com/bhojpur/service/pkg/wasm/wasmer
```

And, you're ready to get fun!

## Supported platforms

This library embeds the `Wasmer` runtime compiled as `shared library` objects, and so uses
[`cgo`](https://golang.org/cmd/cgo/) to consume it. A set of precompiled shared library
objects are provided. Thus, this library works (and is tested) on the following platforms:

<table>
  <thead>
    <tr>
      <th>Platform</th>
      <th>Architecture</th>
      <th>Triple</th>
      <th>Status</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td rowspan="2">Linux</td>
      <td><code>amd64</code></td>
      <td><code>x86_64-unknown-linux-gnu</code></td>
      <td>✅</td>
    </tr>
    <tr>
      <td><code>aarch64</code></td>
      <td><code>aarch64-unknown-linux-gnu</code></td>
      <td>✅</td>
    </tr>
    <tr>
      <td rowspan="2">Darwin</td>
      <td><code>amd64</code></td>
      <td><code>x86_64-apple-darwin</code></td>
      <td>✅</td>
    </tr>
    <tr>
      <td><code>aarch64</code></td>
      <td><code>aarch64-apple-darwin</code></td>
      <td>⏳</td>
    </tr>
    <tr>
      <td>Windows</td>
      <td><code>amd64</code></td>
      <td><code>x86_64-pc-windows-msvc</code></td>
      <td>⏳</td>
    </tr>
  </tbody>
</table>

<details>
  <summary>What to do, if your platform is missing?</summary>
  
  Up to now, there is no script to automate that process. We are working on it.

  Here are the steps to do that manually:
  
  ```sh
  $ # Build the new Wasmer C API shared object library.
  $ cargo build --release
  $
  $ # Configure CGO
  $ export CGO_CFLAGS="-I$(pwd)/pkg/wasm/wasmer/include/"
  $ export CGO_LDFLAGS="-Wl,-rpath,$(pwd)/target/release/ -L$(pwd)/target/release/ -lwasmer"
  $
  $ # Run the tests.
  $ just test -tags custom_wasmer_runtime
  ```
</details>

## Examples

There is a quick program in `internal/simple.go`, written in Rust:

```rust
#[no_mangle]
pub extern "C" fn sum(x: i32, y: i32) -> i32 {
    x + y
}
```

A compiled WebAssembly binary is included in `internal/simple/simple.wasm`.

Then, we can execute it as follows:

```go
package main

import (
	"fmt"
	"io/ioutil"
	wasmer "github.com/bhojpur/service/pkg/wasm/wasmer"
)

func main() {
    wasmBytes, _ := ioutil.ReadFile("simple.wasm")

    engine := wasmer.NewEngine()
    store := wasmer.NewStore(engine)

    // Compiles the module
    module, _ := wasmer.NewModule(store, wasmBytes)

    // Instantiates the module
    importObject := wasmer.NewImportObject()
    instance, _ := wasmer.NewInstance(module, importObject)

    // Gets the `sum` exported function from the WebAssembly instance.
    sum, _ := instance.Exports.GetFunction("sum")

    // Calls that exported function with Go standard values. The WebAssembly
    // types are inferred and values are casted automatically.
    result, _ := sum(5, 37)

    fmt.Println(result) // 42!
}
```

And then, finally, enjoy by running:

```bash
$ cd internal/simple
$ go run simple.go
42
```

## Testing

Run the tests with the following command:

```bash
$ just test
```

## What is WebAssembly?

Quoting [the WebAssembly site](https://webassembly.org/):

> WebAssembly (abbreviated Wasm) is a binary instruction format for a
> stack-based virtual machine. Wasm is designed as a portable target
> for compilation of high-level languages like C/C++/Rust, enabling
> deployment on the web for client and server applications.

About speed:

> WebAssembly aims to execute at native speed by taking advantage of
> [common hardware
> capabilities](https://webassembly.org/docs/portability/#assumptions-for-efficient-execution)
> available on a wide range of platforms.

About safety:

> WebAssembly describes a memory-safe, sandboxed [execution
> environment](https://webassembly.org/docs/semantics/#linear-memory) […].

## FAQ

## How to run Go programs compiled to WebAssembly modules with `wasmer`?

Let's start by emphasing that `wasmer` is a WebAssembly runtime. It allows to _run_
WebAssembly inside [Bhojpur Service](https://gihub.com/bhojpur/service) engine. It's
not a tool to _compile_ a Go program into WebAssembly. While the Go compiler supports
WebAssembly, [it does not support WASI](https://github.com/golang/go/issues/31105)
(WebAssembly System Interface). It generates an ABI that is deeply tied to JavaScript,
and one needs to use the `wasm_exec.js` file provided by the Go toolchain, which doesn't
work outside a JavaScript host.

Fortunately, there are two solutions to this problem:

1. Use [TinyGo](https://tinygo.org) to compile your Go program to WebAssembly with the
  `-target wasi` option, e.g.:

   ```bash
   $ tinygo build -o module.wasm -target wasi .
   ```

   The generated WebAssembly module will be portable across all WebAssembly runtimes that
   support WASI.

2. Use the Go compiler with adapters. Let's see how to compile:

   ```bash
   $ GOOS=js GOARCH=wasm go build -o module.wasm .
   ```

   (the `GOOS=js` is the sign that JavaScript is targeted, not a surprise).

   Then, pick one adapter (they are written by the community):

   - [`mattn/gowasmer`](https://github.com/mattn/gowasmer/),
   - [`go-wasm-adapter/go-wasm`](https://github.com/go-wasm-adapter/go-wasm/),

   and, follow their documentation.

We highly recommend the first solution (with TinyGo) if it works for you as the WebAssembly
module will be portable across all WebAssembly runtimes. It's not a hacky solution based on
adapters; it's the right way to… Go.
