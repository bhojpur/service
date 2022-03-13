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
import "runtime"

// Engine is used by the Store to drive the compilation and the
// execution of a WebAssembly module.
type Engine struct {
	_inner *C.wasm_engine_t
}

func newEngine(engine *C.wasm_engine_t) *Engine {
	self := &Engine{
		_inner: engine,
	}

	runtime.SetFinalizer(self, func(self *Engine) {
		C.wasm_engine_delete(self.inner())
	})

	return self
}

// NewEngine instantiates and returns a new Engine with the default configuration.
//
//   engine := NewEngine()
//
func NewEngine() *Engine {
	return newEngine(C.wasm_engine_new())
}

// NewEngineWithConfig instantiates and returns a new Engine with the given configuration.
//
//   config := NewConfig()
//   engine := NewEngineWithConfig(config)
//
func NewEngineWithConfig(config *Config) *Engine {
	return newEngine(C.wasm_engine_new_with_config(config.inner()))
}

// NewUniversalEngine instantiates and returns a new Universal engine.
//
//   engine := NewUniversalEngine()
//
func NewUniversalEngine() *Engine {
	config := NewConfig()
	config.UseUniversalEngine()

	return NewEngineWithConfig(config)
}

// NewDylibEngine instantiates and returns a new Dylib engine.
//
//   engine := NewDylibEngine()
//
func NewDylibEngine() *Engine {
	config := NewConfig()
	config.UseDylibEngine()

	return NewEngineWithConfig(config)
}

func (self *Engine) inner() *C.wasm_engine_t {
	return self._inner
}

// NewJITEngine is a deprecated function. Please use NewUniversalEngine instead.
func NewJITEngine() *Engine {
	return NewUniversalEngine()
}

// NewNativeEngine is a deprecated function. Please use NewDylibEngine instead.
func NewNativeEngine() *Engine {
	return NewDylibEngine()
}
