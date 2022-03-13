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
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testEngine(t *testing.T, engine *Engine) {
	store := NewStore(engine)
	module, err := NewModule(store, testGetBytes("tests.wasm"))
	assert.NoError(t, err)

	instance, err := NewInstance(module, NewImportObject())
	assert.NoError(t, err)

	sum, err := instance.Exports.GetFunction("sum")
	assert.NoError(t, err)

	result, err := sum(37, 5)
	assert.NoError(t, err)
	assert.Equal(t, result, int32(42))
}

func TestEngine(t *testing.T) {
	testEngine(t, NewEngine())
}

func TestJITEngine(t *testing.T) {
	testEngine(t, NewJITEngine())
}

func TestNativeEngine(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		return
	}

	testEngine(t, NewNativeEngine())
}

func TestEngineWithTarget(t *testing.T) {
	triple, err := NewTriple("aarch64-unknown-linux-gnu")
	assert.NoError(t, err)

	cpuFeatures := NewCpuFeatures()
	assert.NoError(t, err)

	target := NewTarget(triple, cpuFeatures)

	config := NewConfig()
	config.UseTarget(target)

	engine := NewEngineWithConfig(config)
	store := NewStore(engine)

	module, err := NewModule(store, testGetBytes("tests.wasm"))
	assert.NoError(t, err)

	_ = module
}
