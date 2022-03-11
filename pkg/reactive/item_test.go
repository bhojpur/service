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

func Test_SendItems_Variadic(t *testing.T) {
	defer goleak.VerifyNone(t)
	ch := make(chan Item, 3)
	go SendItems(context.Background(), ch, CloseChannel, 1, 2, 3)
	Assert(context.Background(), t, FromChannel(ch), HasItems(1, 2, 3), HasNoError())
}

func Test_SendItems_VariadicWithError(t *testing.T) {
	defer goleak.VerifyNone(t)
	ch := make(chan Item, 3)
	go SendItems(context.Background(), ch, CloseChannel, 1, errFoo, 3)
	Assert(context.Background(), t, FromChannel(ch), HasItems(1, 3), HasError(errFoo))
}

func Test_SendItems_Slice(t *testing.T) {
	defer goleak.VerifyNone(t)
	ch := make(chan Item, 3)
	go SendItems(context.Background(), ch, CloseChannel, []int{1, 2, 3})
	Assert(context.Background(), t, FromChannel(ch), HasItems(1, 2, 3), HasNoError())
}

func Test_SendItems_SliceWithError(t *testing.T) {
	defer goleak.VerifyNone(t)
	ch := make(chan Item, 3)
	go SendItems(context.Background(), ch, CloseChannel, []interface{}{1, errFoo, 3})
	Assert(context.Background(), t, FromChannel(ch), HasItems(1, 3), HasError(errFoo))
}

func Test_Item_SendBlocking(t *testing.T) {
	defer goleak.VerifyNone(t)
	ch := make(chan Item, 1)
	defer close(ch)
	Of(5).SendBlocking(ch)
	assert.Equal(t, 5, (<-ch).V)
}

func Test_Item_SendContext_True(t *testing.T) {
	defer goleak.VerifyNone(t)
	ch := make(chan Item, 1)
	defer close(ch)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	assert.True(t, Of(5).SendContext(ctx, ch))
}

func Test_Item_SendContext_False(t *testing.T) {
	defer goleak.VerifyNone(t)
	ch := make(chan Item, 1)
	defer close(ch)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	assert.False(t, Of(5).SendContext(ctx, ch))
}

func Test_Item_SendNonBlocking(t *testing.T) {
	defer goleak.VerifyNone(t)
	ch := make(chan Item, 1)
	defer close(ch)
	assert.True(t, Of(5).SendNonBlocking(ch))
	assert.False(t, Of(5).SendNonBlocking(ch))
}
