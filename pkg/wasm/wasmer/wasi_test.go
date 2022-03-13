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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWasiVersion(t *testing.T) {
	assert.Equal(t, WASI_VERSION_LATEST.String(), "__latest__")
	assert.Equal(t, WASI_VERSION_SNAPSHOT0.String(), "wasi_unstable")
	assert.Equal(t, WASI_VERSION_SNAPSHOT1.String(), "wasi_snapshot_preview1")
	assert.Equal(t, WASI_VERSION_INVALID.String(), "__unknown__")
}

func TestWasiGetVersion(t *testing.T) {
	engine := NewEngine()
	store := NewStore(engine)
	module, err := NewModule(store, testGetBytes("wasi.wasm"))
	assert.NoError(t, err)

	assert.Equal(t, GetWasiVersion(module), WASI_VERSION_SNAPSHOT1)
}

func TestWasiWithCapturedStdout(t *testing.T) {
	engine := NewEngine()
	store := NewStore(engine)
	module, err := NewModule(store, testGetBytes("wasi.wasm"))
	assert.NoError(t, err)

	wasiEnv, err := NewWasiStateBuilder("test-program").
		Argument("--foo").
		Environment("ABC", "DEF").
		Environment("X", "ZY").
		MapDirectory("the_host_current_directory", ".").
		CaptureStdout().
		Finalize()
	assert.NoError(t, err)

	importObject, err := wasiEnv.GenerateImportObject(store, module)

	instance, err := NewInstance(module, importObject)
	assert.NoError(t, err)

	start, err := instance.Exports.GetWasiStartFunction()
	assert.NoError(t, err)

	start()

	stdout := string(wasiEnv.ReadStdout())

	assert.Equal(
		t,
		stdout,
		"Found program name: `test-program`\n"+
			"Found 1 arguments: --foo\n"+
			"Found 2 environment variables: ABC=DEF, X=ZY\n"+
			"Found 1 preopened directories: DirEntry(\"/the_host_current_directory\")\n",
	)
}
