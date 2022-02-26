package sms

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
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/utils/logger"
)

type mockTransport struct {
	response     *http.Response
	errToReturn  error
	request      *http.Request
	requestCount int32
}

func (t *mockTransport) reset() {
	atomic.StoreInt32(&t.requestCount, 0)
	t.request = nil
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	atomic.AddInt32(&t.requestCount, 1)
	t.request = req

	return t.response, t.errToReturn
}

func TestInit(t *testing.T) {
	m := bindings.Metadata{}
	m.Properties = map[string]string{"toNumber": "toNumber", "fromNumber": "fromNumber"}
	tw := NewSMS(logger.NewLogger("test"))
	err := tw.Init(m)
	assert.NotNil(t, err)
}

func TestParseDuration(t *testing.T) {
	m := bindings.Metadata{}
	m.Properties = map[string]string{
		"toNumber": "toNumber", "fromNumber": "fromNumber",
		"accountSid": "accountSid", "authToken": "authToken", "timeout": "badtimeout",
	}
	tw := NewSMS(logger.NewLogger("test"))
	err := tw.Init(m)
	assert.NotNil(t, err)
}

func TestWriteShouldSucceed(t *testing.T) {
	httpTransport := &mockTransport{
		response: &http.Response{StatusCode: 200, Body: ioutil.NopCloser(strings.NewReader(""))},
	}
	m := bindings.Metadata{}
	m.Properties = map[string]string{
		"toNumber": "toNumber", "fromNumber": "fromNumber",
		"accountSid": "accountSid", "authToken": "authToken",
	}
	tw := NewSMS(logger.NewLogger("test"))
	tw.httpClient = &http.Client{
		Transport: httpTransport,
	}
	err := tw.Init(m)
	assert.Nil(t, err)

	t.Run("Should succeed with expected url and headers", func(t *testing.T) {
		httpTransport.reset()
		_, err := tw.Invoke(&bindings.InvokeRequest{
			Data: []byte("hello world"),
			Metadata: map[string]string{
				toNumber: "toNumber",
			},
		})

		assert.Nil(t, err)
		assert.Equal(t, int32(1), httpTransport.requestCount)
		assert.Equal(t, "https://api.twilio.com/2010-04-01/Accounts/accountSid/Messages.json", httpTransport.request.URL.String())
		assert.NotNil(t, httpTransport.request)
		assert.Equal(t, "application/x-www-form-urlencoded", httpTransport.request.Header.Get("Content-Type"))
		assert.Equal(t, "application/json", httpTransport.request.Header.Get("Accept"))
		authUserName, authPassword, _ := httpTransport.request.BasicAuth()
		assert.Equal(t, "accountSid", authUserName)
		assert.Equal(t, "authToken", authPassword)
	})
}

func TestWriteShouldFail(t *testing.T) {
	httpTransport := &mockTransport{
		response: &http.Response{StatusCode: 200, Body: ioutil.NopCloser(strings.NewReader(""))},
	}
	m := bindings.Metadata{}
	m.Properties = map[string]string{
		"fromNumber": "fromNumber",
		"accountSid": "accountSid", "authToken": "authToken",
	}
	tw := NewSMS(logger.NewLogger("test"))
	tw.httpClient = &http.Client{
		Transport: httpTransport,
	}
	err := tw.Init(m)
	assert.Nil(t, err)

	t.Run("Missing 'to' should fail", func(t *testing.T) {
		httpTransport.reset()
		_, err := tw.Invoke(&bindings.InvokeRequest{
			Data:     []byte("hello world"),
			Metadata: map[string]string{},
		})

		assert.NotNil(t, err)
	})

	t.Run("Twilio call failed should be returned", func(t *testing.T) {
		httpTransport.reset()
		httpErr := errors.New("twilio fake error")
		httpTransport.errToReturn = httpErr
		_, err := tw.Invoke(&bindings.InvokeRequest{
			Data: []byte("hello world"),
			Metadata: map[string]string{
				toNumber: "toNumber",
			},
		})

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), httpErr.Error())
	})

	t.Run("Twilio call returns status not >=200 and <300", func(t *testing.T) {
		httpTransport.reset()
		httpTransport.response.StatusCode = 401
		_, err := tw.Invoke(&bindings.InvokeRequest{
			Data: []byte("hello world"),
			Metadata: map[string]string{
				toNumber: "toNumber",
			},
		})

		assert.NotNil(t, err)
	})
}
