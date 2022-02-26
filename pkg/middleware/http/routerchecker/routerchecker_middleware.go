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
	"encoding/json"
	"regexp"

	"github.com/valyala/fasthttp"

	"github.com/bhojpur/service/pkg/middleware"
	"github.com/bhojpur/service/pkg/utils/logger"
)

// Metadata is the routerchecker middleware config.
type Metadata struct {
	Rule string `json:"rule"`
}

// NewRouterCheckerMiddleware returns a new routerchecker middleware.
func NewMiddleware(logger logger.Logger) *Middleware {
	return &Middleware{logger: logger}
}

// Middleware is an routerchecker middleware.
type Middleware struct {
	logger logger.Logger
}

// GetHandler retruns the HTTP handler provided by the middleware.
func (m *Middleware) GetHandler(metadata middleware.Metadata) (
	func(h fasthttp.RequestHandler) fasthttp.RequestHandler, error) {
	meta, err := m.getNativeMetadata(metadata)
	if err != nil {
		return nil, err
	}

	return func(h fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			isPass, err := regexp.MatchString(meta.Rule, string(ctx.RequestURI()))
			if err != nil {
				m.logger.Error("regexp match failed", err.Error())
				ctx.Error("regexp match failed", fasthttp.StatusBadRequest)

				return
			}
			if !isPass {
				ctx.Error("invalid router", fasthttp.StatusBadRequest)

				return
			}
			h(ctx)
		}
	}, nil
}

func (m *Middleware) getNativeMetadata(metadata middleware.Metadata) (*Metadata, error) {
	b, err := json.Marshal(metadata.Properties)
	if err != nil {
		return nil, err
	}

	var middlewareMetadata Metadata
	err = json.Unmarshal(b, &middlewareMetadata)
	if err != nil {
		return nil, err
	}

	return &middlewareMetadata, nil
}
