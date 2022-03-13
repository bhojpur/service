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

func TestValueKindToString(t *testing.T) {
	assert.Equal(t, I32.String(), "i32")
	assert.Equal(t, I64.String(), "i64")
	assert.Equal(t, F32.String(), "f32")
	assert.Equal(t, F64.String(), "f64")
	assert.Equal(t, AnyRef.String(), "anyref")
	assert.Equal(t, FuncRef.String(), "funcref")
}

func TestValueKindIsNumber(t *testing.T) {
	assert.Equal(t, I32.IsNumber(), true)
	assert.Equal(t, I64.IsNumber(), true)
	assert.Equal(t, F32.IsNumber(), true)
	assert.Equal(t, F64.IsNumber(), true)
	assert.Equal(t, AnyRef.IsNumber(), false)
	assert.Equal(t, FuncRef.IsNumber(), false)
}

func TestValueKindIsReference(t *testing.T) {
	assert.Equal(t, I32.IsReference(), false)
	assert.Equal(t, I64.IsReference(), false)
	assert.Equal(t, F32.IsReference(), false)
	assert.Equal(t, F64.IsReference(), false)
	assert.Equal(t, AnyRef.IsReference(), true)
	assert.Equal(t, FuncRef.IsReference(), true)
}

func TestValueTypeKind(t *testing.T) {
	assert.Equal(t, NewValueType(I32).Kind(), I32)
	assert.Equal(t, NewValueType(I64).Kind(), I64)
	assert.Equal(t, NewValueType(F32).Kind(), F32)
	assert.Equal(t, NewValueType(F64).Kind(), F64)
	assert.Equal(t, NewValueType(AnyRef).Kind(), AnyRef)
	assert.Equal(t, NewValueType(FuncRef).Kind(), FuncRef)
}

func TestValueTypeToVecToList(t *testing.T) {
	valueTypeList := []*ValueType{
		NewValueType(I32),
		NewValueType(I64),
		NewValueType(F32),
	}
	valueTypeVec := toValueTypeVec(valueTypeList)
	assert.Equal(t, int(valueTypeVec.size), 3)

	actualValueTypeList := toValueTypeList(&valueTypeVec, nil)
	assert.Equal(t, len(valueTypeList), len(actualValueTypeList))
	for nth, value := range valueTypeList {
		assert.Equal(t, value.Kind(), actualValueTypeList[nth].Kind())
	}
}

func TestNewValueTypes(t *testing.T) {
	valueTypes := NewValueTypes(I32, I64, F32, F64)
	assert.Equal(t, len(valueTypes), 4)
	assert.Equal(t, valueTypes[0].Kind(), I32)
	assert.Equal(t, valueTypes[1].Kind(), I64)
	assert.Equal(t, valueTypes[2].Kind(), F32)
	assert.Equal(t, valueTypes[3].Kind(), F64)
}
