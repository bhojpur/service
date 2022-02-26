package oauth2

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

	"github.com/fasthttp-contrib/sessions"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"golang.org/x/oauth2"

	"github.com/bhojpur/service/pkg/middleware"
)

// Metadata is the oAuth middleware config.
type oAuth2MiddlewareMetadata struct {
	ClientID       string `json:"clientID"`
	ClientSecret   string `json:"clientSecret"`
	Scopes         string `json:"scopes"`
	AuthURL        string `json:"authURL"`
	TokenURL       string `json:"tokenURL"`
	AuthHeaderName string `json:"authHeaderName"`
	RedirectURL    string `json:"redirectURL"`
	ForceHTTPS     string `json:"forceHTTPS"`
}

// NewOAuth2Middleware returns a new oAuth2 middleware.
func NewOAuth2Middleware() *Middleware {
	return &Middleware{}
}

// Middleware is an oAuth2 authentication middleware.
type Middleware struct{}

const (
	stateParam   = "state"
	savedState   = "auth-state"
	redirectPath = "redirect-url"
	codeParam    = "code"
	https        = "https://"
)

// GetHandler retruns the HTTP handler provided by the middleware.
func (m *Middleware) GetHandler(metadata middleware.Metadata) (func(h fasthttp.RequestHandler) fasthttp.RequestHandler, error) {
	meta, err := m.getNativeMetadata(metadata)
	if err != nil {
		return nil, err
	}

	return func(h fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			conf := &oauth2.Config{
				ClientID:     meta.ClientID,
				ClientSecret: meta.ClientSecret,
				Scopes:       strings.Split(meta.Scopes, ","),
				RedirectURL:  meta.RedirectURL,
				Endpoint: oauth2.Endpoint{
					AuthURL:  meta.AuthURL,
					TokenURL: meta.TokenURL,
				},
			}
			session := sessions.StartFasthttp(ctx)
			if session.GetString(meta.AuthHeaderName) != "" {
				ctx.Request.Header.Add(meta.AuthHeaderName, session.GetString(meta.AuthHeaderName))
				h(ctx)

				return
			}
			state := string(ctx.FormValue(stateParam))
			//nolint:nestif
			if state == "" {
				id, _ := uuid.NewUUID()
				session.Set(savedState, id.String())
				session.Set(redirectPath, string(ctx.RequestURI()))
				url := conf.AuthCodeURL(id.String(), oauth2.AccessTypeOffline)
				ctx.Redirect(url, 302)
			} else {
				authState := session.GetString(savedState)
				redirectURL := session.GetString(redirectPath)
				if strings.EqualFold(meta.ForceHTTPS, "true") {
					redirectURL = https + string(ctx.Request.Host()) + redirectURL
				}
				if state != authState {
					ctx.Error("invalid state", fasthttp.StatusBadRequest)
				} else {
					code := string(ctx.FormValue(codeParam))
					if code == "" {
						ctx.Error("code not found", fasthttp.StatusBadRequest)
					} else {
						token, err := conf.Exchange(context.Background(), code)
						if err != nil {
							ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
						}
						session.Set(meta.AuthHeaderName, token.Type()+" "+token.AccessToken)
						ctx.Request.Header.Add(meta.AuthHeaderName, token.Type()+" "+token.AccessToken)
						ctx.Redirect(redirectURL, 302)
					}
				}
			}
		}
	}, nil
}

func (m *Middleware) getNativeMetadata(metadata middleware.Metadata) (*oAuth2MiddlewareMetadata, error) {
	b, err := json.Marshal(metadata.Properties)
	if err != nil {
		return nil, err
	}

	var middlewareMetadata oAuth2MiddlewareMetadata
	err = json.Unmarshal(b, &middlewareMetadata)
	if err != nil {
		return nil, err
	}

	return &middlewareMetadata, nil
}
