package bearer

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
	"encoding/json"
	"strings"

	oidc "github.com/coreos/go-oidc"
	"github.com/valyala/fasthttp"

	"github.com/bhojpur/service/pkg/middleware"
	"github.com/bhojpur/service/pkg/utils/logger"
)

type bearerMiddlewareMetadata struct {
	IssuerURL string `json:"issuerURL"`
	ClientID  string `json:"clientID"`
}

// NewBearerMiddleware returns a new oAuth2 middleware.
func NewBearerMiddleware(logger logger.Logger) *Middleware {
	return &Middleware{logger: logger}
}

// Middleware is an oAuth2 authentication middleware.
type Middleware struct {
	logger logger.Logger
}

const (
	bearerPrefix       = "bearer "
	bearerPrefixLength = len(bearerPrefix)
)

// GetHandler retruns the HTTP handler provided by the middleware.
func (m *Middleware) GetHandler(metadata middleware.Metadata) (func(h fasthttp.RequestHandler) fasthttp.RequestHandler, error) {
	meta, err := m.getNativeMetadata(metadata)
	if err != nil {
		return nil, err
	}

	provider, err := oidc.NewProvider(context.Background(), meta.IssuerURL)
	if err != nil {
		return nil, err
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: meta.ClientID,
	})

	return func(h fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			authHeader := string(ctx.Request.Header.Peek(fasthttp.HeaderAuthorization))
			if !strings.HasPrefix(strings.ToLower(authHeader), bearerPrefix) {
				ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)

				return
			}
			rawToken := authHeader[bearerPrefixLength:]
			_, err := verifier.Verify(ctx, rawToken)
			if err != nil {
				ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)

				return
			}

			h(ctx)
		}
	}, nil
}

func (m *Middleware) getNativeMetadata(metadata middleware.Metadata) (*bearerMiddlewareMetadata, error) {
	b, err := json.Marshal(metadata.Properties)
	if err != nil {
		return nil, err
	}

	var middlewareMetadata bearerMiddlewareMetadata
	err = json.Unmarshal(b, &middlewareMetadata)
	if err != nil {
		return nil, err
	}

	return &middlewareMetadata, nil
}
