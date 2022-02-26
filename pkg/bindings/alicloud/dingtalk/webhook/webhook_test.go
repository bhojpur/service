package webhook

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
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/utils/logger"
)

func TestPublishMsg(t *testing.T) { //nolint:paralleltest
	msg := "{\"type\": \"text\",\"text\": {\"content\": \"hello\"}}"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("{\"errcode\":0}"))
		require.NoError(t, err)
		if r.Method != "POST" {
			t.Errorf("Expected 'POST' request, got '%s'", r.Method)
		}
		if r.URL.EscapedPath() != "/test" {
			t.Errorf("Expected request to '/test', got '%s'", r.URL.EscapedPath())
		}

		body, err := ioutil.ReadAll(r.Body)
		require.Nil(t, err)
		assert.Equal(t, msg, string(body))
	}))
	defer ts.Close()

	m := bindings.Metadata{Name: "test", Properties: map[string]string{
		"url":    ts.URL + "/test",
		"secret": "",
		"id":     "x",
	}}

	d := NewDingTalkWebhook(logger.NewLogger("test"))
	err := d.Init(m)
	require.NoError(t, err)

	req := &bindings.InvokeRequest{Data: []byte(msg), Operation: bindings.CreateOperation, Metadata: map[string]string{}}
	_, err = d.Invoke(req)
	require.NoError(t, err)
}

func TestBindingReadAndInvoke(t *testing.T) { //nolint:paralleltest
	msg := "{\"type\": \"text\",\"text\": {\"content\": \"hello\"}}"

	m := bindings.Metadata{
		Name: "test",
		Properties: map[string]string{
			"url":    "/test",
			"secret": "",
			"id":     "x",
		},
	}

	d := NewDingTalkWebhook(logger.NewLogger("test"))
	err := d.Init(m)
	assert.NoError(t, err)

	var count int32
	ch := make(chan bool, 1)

	handler := func(in *bindings.ReadResponse) ([]byte, error) {
		assert.Equal(t, msg, string(in.Data))
		atomic.AddInt32(&count, 1)
		ch <- true

		return nil, nil
	}

	err = d.Read(handler)
	require.NoError(t, err)

	req := &bindings.InvokeRequest{Data: []byte(msg), Operation: bindings.GetOperation, Metadata: map[string]string{}}
	_, err = d.Invoke(req)
	require.NoError(t, err)

	select {
	case <-ch:
		require.True(t, atomic.LoadInt32(&count) > 0)
	case <-time.After(time.Second):
		require.FailNow(t, "read timeout")
	}
}
