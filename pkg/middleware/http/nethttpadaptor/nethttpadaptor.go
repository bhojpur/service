package nethttpadaptor

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
	"io/ioutil"
	"net"
	"net/http"
	"strconv"

	"github.com/valyala/fasthttp"

	"github.com/bhojpur/service/pkg/utils/logger"
)

// NewNetHTTPHandlerFunc wraps a fasthttp.RequestHandler in a http.HandlerFunc.
func NewNetHTTPHandlerFunc(logger logger.Logger, h fasthttp.RequestHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := fasthttp.RequestCtx{}
		remoteIP := net.ParseIP(r.RemoteAddr)
		remoteAddr := net.IPAddr{remoteIP, ""} //nolint
		c.Init(&fasthttp.Request{}, &remoteAddr, nil)

		if r.Body != nil {
			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logger.Errorf("error reading request body, %+v", err)

				return
			}
			c.Request.SetBody(reqBody)
		}
		c.Request.SetRequestURI(r.URL.RequestURI())
		c.Request.URI().SetScheme(r.URL.Scheme)
		c.Request.SetHost(r.Host)
		c.Request.Header.SetMethod(r.Method)
		c.Request.Header.Set("Proto", r.Proto)
		major := strconv.Itoa(r.ProtoMajor)
		minor := strconv.Itoa(r.ProtoMinor)
		c.Request.Header.Set("Protomajor", major)
		c.Request.Header.Set("Protominor", minor)
		c.Request.Header.SetContentType(r.Header.Get("Content-Type"))
		c.Request.Header.SetContentLength(int(r.ContentLength))
		c.Request.Header.SetReferer(r.Referer())
		c.Request.Header.SetUserAgent(r.UserAgent())
		for _, cookie := range r.Cookies() {
			c.Request.Header.SetCookie(cookie.Name, cookie.Value)
		}
		for k, v := range r.Header {
			for _, i := range v {
				c.Request.Header.Add(k, i)
			}
		}

		ctx := r.Context()
		reqCtx, ok := ctx.(*fasthttp.RequestCtx)
		if ok {
			reqCtx.VisitUserValues(func(k []byte, v interface{}) {
				c.SetUserValueBytes(k, v)
			})
		}

		h(&c)

		c.Response.Header.VisitAll(func(k []byte, v []byte) {
			w.Header().Add(string(k), string(v))
		})
		status := c.Response.StatusCode()
		w.WriteHeader(status)

		c.Response.BodyWriteTo(w)
	})
}
