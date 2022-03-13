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

func TestTableType(t *testing.T) {
	valueType := NewValueType(I32)

	var minimum uint32 = 1
	var maximum uint32 = 7
	limits, err := NewLimits(minimum, maximum)
	assert.NoError(t, err)

	tableType := NewTableType(valueType, limits)

	valueTypeAgain := tableType.ValueType()
	assert.Equal(t, valueTypeAgain.Kind(), I32)

	limitsAgain := tableType.Limits()
	assert.Equal(t, limitsAgain.Minimum(), minimum)
	assert.Equal(t, limitsAgain.Maximum(), maximum)
}

func TestTableTypeIntoExternTypeAndBack(t *testing.T) {
	valueType := NewValueType(I32)

	var minimum uint32 = 1
	var maximum uint32 = 7
	limits, err := NewLimits(minimum, maximum)
	assert.NoError(t, err)

	tableType := NewTableType(valueType, limits)
	externType := tableType.IntoExternType()
	assert.Equal(t, externType.Kind(), TABLE)

	tableTypeAgain := externType.IntoTableType()

	valueTypeAgain := tableTypeAgain.ValueType()
	assert.Equal(t, valueTypeAgain.Kind(), I32)

	limitsAgain := tableTypeAgain.Limits()
	assert.Equal(t, limitsAgain.Minimum(), minimum)
	assert.Equal(t, limitsAgain.Maximum(), maximum)
}
