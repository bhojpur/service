package reactive

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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func Test_Single_Get_Item(t *testing.T) {
	defer goleak.VerifyNone(t)
	var s Single = &SingleImpl{iterable: Just(1)()}
	get, err := s.Get()
	assert.NoError(t, err)
	assert.Equal(t, 1, get.V)
}

func Test_Single_Get_Error(t *testing.T) {
	defer goleak.VerifyNone(t)
	var s Single = &SingleImpl{iterable: Just(errFoo)()}
	get, err := s.Get()
	assert.NoError(t, err)
	assert.Equal(t, errFoo, get.E)
}

func Test_Single_Get_ContextCanceled(t *testing.T) {
	defer goleak.VerifyNone(t)
	ch := make(chan Item)
	defer close(ch)
	ctx, cancel := context.WithCancel(context.Background())
	var s Single = &SingleImpl{iterable: FromChannel(ch)}
	cancel()
	_, err := s.Get(WithContext(ctx))
	assert.Equal(t, ctx.Err(), err)
}

func Test_Single_Filter_True(t *testing.T) {
	defer goleak.VerifyNone(t)
	os := JustItem(1).Filter(func(i interface{}) bool {
		return i == 1
	})
	Assert(context.Background(), t, os, HasItem(1), HasNoError())
}

func Test_Single_Filter_False(t *testing.T) {
	defer goleak.VerifyNone(t)
	os := JustItem(1).Filter(func(i interface{}) bool {
		return i == 0
	})
	Assert(context.Background(), t, os, IsEmpty(), HasNoError())
}

func Test_Single_Map(t *testing.T) {
	defer goleak.VerifyNone(t)
	single := JustItem(1).Map(func(_ context.Context, i interface{}) (interface{}, error) {
		return i.(int) + 1, nil
	})
	Assert(context.Background(), t, single, HasItem(2), HasNoError())
}
