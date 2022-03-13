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

func TestImportTypeForFunctionType(t *testing.T) {
	params := NewValueTypes(I32, I64)
	results := NewValueTypes(F32)
	functionType := NewFunctionType(params, results)

	module := "foo"
	name := "bar"
	importType := NewImportType(module, name, functionType)
	assert.Equal(t, importType.Module(), module)
	assert.Equal(t, importType.Name(), name)

	externType := importType.Type()
	assert.Equal(t, externType.Kind(), FUNCTION)

	functionTypeAgain := externType.IntoFunctionType()
	assert.Equal(t, len(functionTypeAgain.Params()), len(params))
	assert.Equal(t, len(functionTypeAgain.Results()), len(results))
}

func TestImportTypeForGlobalType(t *testing.T) {
	valueType := NewValueType(I32)
	globalType := NewGlobalType(valueType, MUTABLE)

	module := "foo"
	name := "bar"
	importType := NewImportType(module, name, globalType)
	assert.Equal(t, importType.Module(), module)
	assert.Equal(t, importType.Name(), name)

	externType := importType.Type()
	assert.Equal(t, externType.Kind(), GLOBAL)

	globalTypeAgain := externType.IntoGlobalType()
	assert.Equal(t, globalTypeAgain.ValueType().Kind(), I32)
	assert.Equal(t, globalTypeAgain.Mutability(), MUTABLE)
}

func TestImportTypeForTableType(t *testing.T) {
	valueType := NewValueType(I32)

	var minimum uint32 = 1
	var maximum uint32 = 7
	limits, err := NewLimits(minimum, maximum)
	assert.NoError(t, err)

	tableType := NewTableType(valueType, limits)

	module := "foo"
	name := "bar"
	importType := NewImportType(module, name, tableType)
	assert.Equal(t, importType.Module(), module)
	assert.Equal(t, importType.Name(), name)

	externType := importType.Type()
	assert.Equal(t, externType.Kind(), TABLE)

	tableTypeAgain := externType.IntoTableType()
	valueTypeAgain := tableTypeAgain.ValueType()
	assert.Equal(t, valueTypeAgain.Kind(), I32)

	limitsAgain := tableTypeAgain.Limits()
	assert.Equal(t, limitsAgain.Minimum(), minimum)
	assert.Equal(t, limitsAgain.Maximum(), maximum)
}

func TestImportTypeForMemoryType(t *testing.T) {
	var minimum uint32 = 1
	var maximum uint32 = 7
	limits, err := NewLimits(minimum, maximum)
	assert.NoError(t, err)

	memoryType := NewMemoryType(limits)

	module := "foo"
	name := "bar"
	importType := NewImportType(module, name, memoryType)
	assert.Equal(t, importType.Module(), module)
	assert.Equal(t, importType.Name(), name)

	externType := importType.Type()
	assert.Equal(t, externType.Kind(), MEMORY)

	memoryTypeAgain := externType.IntoMemoryType()
	limitsAgain := memoryTypeAgain.Limits()
	assert.Equal(t, limitsAgain.Minimum(), minimum)
	assert.Equal(t, limitsAgain.Maximum(), maximum)
}
