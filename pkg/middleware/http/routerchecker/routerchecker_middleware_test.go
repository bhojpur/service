package routerchecker

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
	"github.com/valyala/fasthttp"

	"github.com/bhojpur/service/pkg/middleware"
	"github.com/bhojpur/service/pkg/utils/logger"
)

type RouterOutput struct{}

func (ro *RouterOutput) handle(ctx *fasthttp.RequestCtx) {
	ctx.Error(string(ctx.RequestURI()), fasthttp.StatusOK)
}

func TestRequestHandlerWithIllegalRouterRule(t *testing.T) {
	meta := middleware.Metadata{Properties: map[string]string{
		"rule": "^[A-Za-z0-9/._-]+$",
	}}
	log := logger.NewLogger("routerchecker.test")
	rchecker := NewMiddleware(log)
	handler, err := rchecker.GetHandler(meta)
	assert.Nil(t, err)

	var ctx fasthttp.RequestCtx
	ctx.Request.SetHost("localhost:5001")
	ctx.Request.SetRequestURI("/v1/invoke/qcg.default/method/ cat password")
	ctx.Request.Header.SetMethod("GET")

	output := new(RouterOutput)
	handler(output.handle)(&ctx)
	assert.Equal(t, fasthttp.StatusBadRequest, ctx.Response.Header.StatusCode())
}

func TestRequestHandlerWithLegalRouterRule(t *testing.T) {
	meta := middleware.Metadata{Properties: map[string]string{
		"rule": "^[A-Za-z0-9/._-]+$",
	}}

	log := logger.NewLogger("routerchecker.test")
	rchecker := NewMiddleware(log)
	handler, err := rchecker.GetHandler(meta)
	assert.Nil(t, err)

	var ctx fasthttp.RequestCtx
	ctx.Request.SetHost("localhost:5001")
	ctx.Request.SetRequestURI("/v1/invoke/qcg.default/method")
	ctx.Request.Header.SetMethod("GET")

	output := new(RouterOutput)
	handler(output.handle)(&ctx)
	assert.Equal(t, fasthttp.StatusOK, ctx.Response.Header.StatusCode())
}
