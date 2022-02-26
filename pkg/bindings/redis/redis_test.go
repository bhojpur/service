package redis

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

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/utils/logger"
)

const (
	testData = `{"data":"data"}`
	testKey  = "test"
)

func TestInvoke(t *testing.T) {
	s, c := setupMiniredis()
	defer s.Close()

	bind := &Redis{
		client: c,
		logger: logger.NewLogger("test"),
	}
	bind.ctx, bind.cancel = context.WithCancel(context.Background())

	_, err := c.Do(context.Background(), "GET", testKey).Result()
	assert.Equal(t, redis.Nil, err)

	bindingRes, err := bind.Invoke(&bindings.InvokeRequest{
		Data:     []byte(testData),
		Metadata: map[string]string{"key": testKey},
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, true, bindingRes == nil)

	getRes, err := c.Do(context.Background(), "GET", testKey).Result()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, getRes == testData)
}

func setupMiniredis() (*miniredis.Miniredis, *redis.Client) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	opts := &redis.Options{
		Addr: s.Addr(),
		DB:   0,
	}

	return s, redis.NewClient(opts)
}
