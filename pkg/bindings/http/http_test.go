package http_test

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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bhojpur/service/pkg/utils/logger"

	"github.com/bhojpur/service/pkg/bindings"
	binding_http "github.com/bhojpur/service/pkg/bindings/http"
)

func TestOperations(t *testing.T) {
	opers := (*binding_http.HTTPSource)(nil).Operations()
	assert.Equal(t, []bindings.OperationKind{
		bindings.CreateOperation,
		"get",
		"head",
		"post",
		"put",
		"patch",
		"delete",
		"options",
		"trace",
	}, opers)
}

func TestInit(t *testing.T) {
	var path string

	s := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			path = req.URL.Path
			input := req.Method
			if req.Body != nil {
				defer req.Body.Close()
				b, _ := ioutil.ReadAll(req.Body)
				if len(b) > 0 {
					input = string(b)
				}
			}
			inputFromHeader := req.Header.Get("X-Input")
			if inputFromHeader != "" {
				input = inputFromHeader
			}
			w.Header().Set("Content-Type", "text/plain")
			if input == "internal server error" {
				w.WriteHeader(500)
			}
			w.Write([]byte(strings.ToUpper(input)))
		}),
	)
	defer s.Close()

	m := bindings.Metadata{
		Properties: map[string]string{
			"url": s.URL,
		},
	}
	hs := binding_http.NewHTTP(logger.NewLogger("test"))
	err := hs.Init(m)
	require.NoError(t, err)

	tests := map[string]struct {
		input     string
		operation string
		metadata  map[string]string
		path      string
		err       string
	}{
		"get": {
			input:     "GET",
			operation: "get",
			metadata:  nil,
			path:      "/",
			err:       "",
		},
		"request headers": {
			input:     "OVERRIDE",
			operation: "get",
			metadata:  map[string]string{"X-Input": "override"},
			path:      "/",
			err:       "",
		},
		"post": {
			input:     "expected",
			operation: "post",
			metadata:  map[string]string{"path": "/test"},
			path:      "/test",
			err:       "",
		},
		"put": {
			input:     "expected",
			operation: "put",
			metadata:  map[string]string{"path": "/test"},
			path:      "/test",
			err:       "",
		},
		"patch": {
			input:     "expected",
			operation: "patch",
			metadata:  map[string]string{"path": "/test"},
			path:      "/test",
			err:       "",
		},
		"delete": {
			input:     "DELETE",
			operation: "delete",
			metadata:  nil,
			path:      "/",
			err:       "",
		},
		"options": {
			input:     "OPTIONS",
			operation: "options",
			metadata:  nil,
			path:      "/",
			err:       "",
		},
		"trace": {
			input:     "TRACE",
			operation: "trace",
			metadata:  nil,
			path:      "/",
			err:       "",
		},
		"backward compatibility": {
			input:     "expected",
			operation: "create",
			metadata:  map[string]string{"path": "/test"},
			path:      "/test",
			err:       "",
		},
		"invalid path": {
			input:     "expected",
			operation: "POST",
			metadata:  map[string]string{"path": "/../test"},
			path:      "",
			err:       "invalid path: /../test",
		},
		"invalid operation": {
			input:     "notvalid",
			operation: "notvalid",
			metadata:  map[string]string{"path": "/test"},
			path:      "/test",
			err:       "invalid operation: notvalid",
		},
		"internal server error": {
			input:     "internal server error",
			operation: "post",
			metadata:  map[string]string{"path": "/"},
			path:      "/",
			err:       "received status code 500",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			response, err := hs.Invoke(&bindings.InvokeRequest{
				Data:      []byte(tc.input),
				Metadata:  tc.metadata,
				Operation: bindings.OperationKind(tc.operation),
			})
			if tc.err == "" {
				require.NoError(t, err)
				assert.Equal(t, tc.path, path)
				assert.Equal(t, strings.ToUpper(tc.input), string(response.Data))
				assert.Equal(t, "text/plain", response.Metadata["Content-Type"])
			} else {
				require.Error(t, err)
				assert.Equal(t, tc.err, err.Error())
			}
		})
	}
}
