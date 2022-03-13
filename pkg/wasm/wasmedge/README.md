# Bhojpur Service - WasmEdge Integration

The [WasmEdge](https://github.com/WasmEdge/WasmEdge) is a high performance WebAssembly runtime optimized
for server-side applications. This library provides methods of accessing WasmEdge.

## Getting Started

The `WasmEdge` requires `Go` version >= `1.17`.

```bash
$ go version
go version go1.17.6 linux/amd64
```

The software developers must [install the WasmEdge shared library](https://github.com/WasmEdge/WasmEdge/blob/master/docs/install.md) with the same `WasmEdge` release version.

```bash
curl -sSf https://raw.githubusercontent.com/WasmEdge/WasmEdge/master/utils/install.sh | bash -s -- -v 0.9.1
```

For the software developers, who need the `TensorFlow` or `Image` extension for `WasmEdge`, please install
the `WasmEdge` with extensions:

```bash
curl -sSf https://raw.githubusercontent.com/WasmEdge/WasmEdge/master/utils/install.sh | bash -s -- -e all -v 0.9.1
```

Please note that the `TensorFlow` and `Image` extensions are only for the `Linux` platforms.

Install the `WasmEdge` package and build in your Go project directory:

```bash
go get github.com/second-state/WasmEdge-go/wasmedge@v0.9.1
go build
```

## WasmEdge Extensions

By default, the `WasmEdge` only turns on the basic runtime. The `WasmEdge` has the following extensions:

### Tensorflow

- It supports host functions in [WasmEdge-tensorflow](https://github.com/second-state/WasmEdge-tensorflow)
- The `TensorFlow` extension when installing `WasmEdge` is required. Please install `WasmEdge` with the `-e tensorflow` command
- For using this extension, the tag `tensorflow` when building is required:

    ```bash
    go build -tags tensorflow
    ```

### Image

- This extension supports the host functions in [WasmEdge-image](https://github.com/second-state/WasmEdge-image)
- The `Image` extension when installing `WasmEdge` is required. Please install `WasmEdge` with the `-e image` command
- For using this extension, the tag `image` when building is required:

    ```bash
    go build -tags image
    ```

Users can also turn on the multiple extensions when building:

```bash
go build -tags image,tensorflow
```
