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

type Instance struct {
	_inner  *C.wasm_instance_t
	Exports *Exports

	// without this, imported functions may be freed before execution of an exported function is complete.
	imports *ImportObject
}

// NewInstance instantiates a new Instance.
//
// It takes two arguments, the Module and an ImportObject.
//
// Note:️ Instantiating a module may return TrapError if the module's
// start function traps.
//
//   wasmBytes := []byte(`...`)
//   engine := wasmer.NewEngine()
//   store := wasmer.NewStore(engine)
//   module, err := wasmer.NewModule(store, wasmBytes)
//   importObject := wasmer.NewImportObject()
//   instance, err := wasmer.NewInstance(module, importObject)
//
func NewInstance(module *Module, imports *ImportObject) (*Instance, error) {
	var traps *C.wasm_trap_t
	externs, err := imports.intoInner(module)

	if err != nil {
		return nil, err
	}

	var instance *C.wasm_instance_t

	err2 := maybeNewErrorFromWasmer(func() bool {
		instance = C.wasm_instance_new(
			module.store.inner(),
			module.inner(),
			externs,
			&traps,
		)

		return traps == nil && instance == nil
	})

	if err2 != nil {
		return nil, err2
	}

	if traps != nil {
		return nil, newErrorFromTrap(traps)
	}

	self := &Instance{
		_inner:  instance,
		Exports: newExports(instance, module),
		imports: imports,
	}

	runtime.SetFinalizer(self, func(self *Instance) {
		self.Close()
	})

	return self, nil
}

func (self *Instance) inner() *C.wasm_instance_t {
	return self._inner
}

// Force to close the Instance.
//
// A runtime finalizer is registered on the Instance, but it is
// possible to force the destruction of the Instance by calling Close
// manually.
func (self *Instance) Close() {
	runtime.SetFinalizer(self, nil)
	C.wasm_instance_delete(self.inner())
	self.Exports.Close()
}
