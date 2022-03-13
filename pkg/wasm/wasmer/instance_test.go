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

func TestInstance(t *testing.T) {
	engine := NewEngine()
	store := NewStore(engine)
	module, err := NewModule(store, []byte("(module)"))
	assert.NoError(t, err)

	_, err = NewInstance(module, NewImportObject())
	assert.NoError(t, err)
}

func TestInstanceExports(t *testing.T) {
	engine := NewEngine()
	store := NewStore(engine)
	module, err := NewModule(
		store,
		[]byte(`
			(module
			  (func (export "function") (param i32 i64))
			  (global (export "global") i32 (i32.const 7))
			  (table (export "table") 0 funcref)
			  (memory (export "memory") 1))
		`),
	)
	assert.NoError(t, err)

	instance, err := NewInstance(module, NewImportObject())
	assert.NoError(t, err)

	extern, err := instance.Exports.Get("function")
	assert.NoError(t, err)
	assert.Equal(t, extern.Kind(), FUNCTION)

	function, err := instance.Exports.GetFunction("function")
	assert.NoError(t, err)
	assert.NotNil(t, function)

	global, err := instance.Exports.GetGlobal("global")
	assert.NoError(t, err)
	assert.NotNil(t, global)

	table, err := instance.Exports.GetTable("table")
	assert.NoError(t, err)
	assert.NotNil(t, table)

	memory, err := instance.Exports.GetMemory("memory")
	assert.NoError(t, err)
	assert.NotNil(t, memory)
}

func TestInstanceMissingImports(t *testing.T) {
	engine := NewEngine()
	store := NewStore(engine)
	module, err := NewModule(
		store,
		[]byte(`
			(module
			  (func (import "missing" "function"))
			  (func (import "exists" "function")))
		`),
	)
	assert.NoError(t, err)

	function := NewFunction(
		store,
		NewFunctionType(NewValueTypes(), NewValueTypes()),
		func(args []Value) ([]Value, error) {
			return []Value{}, nil
		},
	)

	importObject := NewImportObject()
	importObject.Register(
		"exists",
		map[string]IntoExtern{
			"function": function,
		},
	)

	_, err = NewInstance(module, importObject)
	assert.Error(t, err)
}

func TestInstanceTraps(t *testing.T) {
	engine := NewEngine()
	store := NewStore(engine)
	module, err := NewModule(
		store,
		[]byte(`
			(module
			  (start $start_f)
			  (type $start_t (func))
			  (func $start_f (type $start_t)
			    unreachable))
		`),
	)
	assert.NoError(t, err)

	_, err = NewInstance(module, NewImportObject())
	assert.Error(t, err)
	assert.Equal(t, "unreachable", err.Error())
}
