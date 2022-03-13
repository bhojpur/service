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
import (
	"fmt"
	"unsafe"
)

// ImportObject contains all of the import data used when
// instantiating a WebAssembly module.
type ImportObject struct {
	externs map[string]map[string]IntoExtern
}

// NewImportObject instantiates a new empty ImportObject.
//
//   imports := NewImportObject()
//
func NewImportObject() *ImportObject {
	return &ImportObject{
		externs: make(map[string]map[string]IntoExtern),
	}
}

func (self *ImportObject) intoInner(module *Module) (*C.wasm_extern_vec_t, error) {
	cExterns := &C.wasm_extern_vec_t{}

	var externs []*C.wasm_extern_t
	var numberOfExterns uint

	for _, importType := range module.Imports() {
		namespace := importType.Module()
		name := importType.Name()

		if self.externs[namespace][name] == nil {
			return nil, &Error{
				message: fmt.Sprintf("Missing import: `%s`.`%s`", namespace, name),
			}
		}

		externs = append(externs, self.externs[namespace][name].IntoExtern().inner())
		numberOfExterns++
	}

	if numberOfExterns > 0 {
		C.wasm_extern_vec_new(cExterns, C.size_t(numberOfExterns), (**C.wasm_extern_t)(unsafe.Pointer(&externs[0])))
	}

	return cExterns, nil
}

// ContainsNamespace returns true if the ImportObject contains the given namespace (or module name)
//
//   imports := NewImportObject()
//   _ = imports.ContainsNamespace("env") // false
//
func (self *ImportObject) ContainsNamespace(name string) bool {
	_, exists := self.externs[name]

	return exists
}

// Register registers a namespace (or module name) in the ImportObject.
//
// It takes two arguments: the namespace name and a map with imports names as key and externs as values.
//
// Note:️ An extern is anything implementing IntoExtern: Function, Global, Memory, Table.
//
//   imports := NewImportObject()
//   importObject.Register(
//   	"env",
//   	map[string]wasmer.IntoExtern{
//   		"host_function": hostFunction,
//   		"host_global": hostGlobal,
//   	},
//  )
//
// Note:️ The namespace (or module name) may be empty:
//
//   imports := NewImportObject()
//   importObject.Register(
//   	"",
//   	map[string]wasmer.IntoExtern{
//    		"host_function": hostFunction,
//   		"host_global": hostGlobal,
//   	},
//   )
//
func (self *ImportObject) Register(namespaceName string, namespace map[string]IntoExtern) {
	_, exists := self.externs[namespaceName]

	if exists == false {
		self.externs[namespaceName] = namespace
	} else {
		for key, value := range namespace {
			self.externs[namespaceName][key] = value
		}
	}
}
