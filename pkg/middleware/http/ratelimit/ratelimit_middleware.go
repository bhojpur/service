package ratelimit

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
	"fmt"
	"strconv"

	"github.com/didip/tollbooth"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"

	"github.com/bhojpur/service/pkg/middleware"
	"github.com/bhojpur/service/pkg/middleware/http/nethttpadaptor"
	"github.com/bhojpur/service/pkg/utils/logger"
)

// Metadata is the ratelimit middleware config.
type rateLimitMiddlewareMetadata struct {
	MaxRequestsPerSecond float64 `json:"maxRequestsPerSecond"`
}

const (
	maxRequestsPerSecondKey = "maxRequestsPerSecond"

	// Defaults.
	defaultMaxRequestsPerSecond = 100
)

// NewRateLimitMiddleware returns a new ratelimit middleware.
func NewRateLimitMiddleware(logger logger.Logger) *Middleware {
	return &Middleware{logger: logger}
}

// Middleware is an ratelimit middleware.
type Middleware struct {
	logger logger.Logger
}

// GetHandler returns the HTTP handler provided by the middleware.
func (m *Middleware) GetHandler(metadata middleware.Metadata) (func(h fasthttp.RequestHandler) fasthttp.RequestHandler, error) {
	meta, err := m.getNativeMetadata(metadata)
	if err != nil {
		return nil, err
	}

	limiter := tollbooth.NewLimiter(meta.MaxRequestsPerSecond, nil)

	return func(h fasthttp.RequestHandler) fasthttp.RequestHandler {
		limitHandler := tollbooth.LimitFuncHandler(limiter, nethttpadaptor.NewNetHTTPHandlerFunc(m.logger, h))
		wrappedHandler := fasthttpadaptor.NewFastHTTPHandlerFunc(limitHandler.ServeHTTP)

		return func(ctx *fasthttp.RequestCtx) {
			wrappedHandler(ctx)
		}
	}, nil
}

func (m *Middleware) getNativeMetadata(metadata middleware.Metadata) (*rateLimitMiddlewareMetadata, error) {
	var middlewareMetadata rateLimitMiddlewareMetadata

	middlewareMetadata.MaxRequestsPerSecond = defaultMaxRequestsPerSecond
	if val, ok := metadata.Properties[maxRequestsPerSecondKey]; ok {
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing ratelimit middleware property %s: %+v", maxRequestsPerSecondKey, err)
		}
		if f <= 0 {
			return nil, fmt.Errorf("ratelimit middleware property %s must be a positive value", maxRequestsPerSecondKey)
		}
		middlewareMetadata.MaxRequestsPerSecond = f
	}

	return &middlewareMetadata, nil
}
