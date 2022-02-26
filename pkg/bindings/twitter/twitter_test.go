package twitter

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
	"os"
	"testing"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/utils/logger"
)

const (
	testTwitterConsumerKey    = "test-consumerKey"
	testTwitterConsumerSecret = "test-consumerSecret"
	testTwitterAccessToken    = "test-accessToken"
	testTwitterAccessSecret   = "test-accessSecret"
)

func getTestMetadata() bindings.Metadata {
	m := bindings.Metadata{}
	m.Properties = map[string]string{
		"consumerKey":    testTwitterConsumerKey,
		"consumerSecret": testTwitterConsumerSecret,
		"accessToken":    testTwitterAccessToken,
		"accessSecret":   testTwitterAccessSecret,
	}

	return m
}

func getRuntimeMetadata() map[string]string {
	return map[string]string{
		"consumerKey":    os.Getenv("CONSUMER_KEY"),
		"consumerSecret": os.Getenv("CONSUMER_SECRET"),
		"accessToken":    os.Getenv("ACCESS_TOKEN"),
		"accessSecret":   os.Getenv("ACCESS_SECRET"),
	}
}

// go test -v -count=1 ./bindings/twitter/.
func TestInit(t *testing.T) {
	m := getTestMetadata()
	tw := NewTwitter(logger.NewLogger("test"))
	err := tw.Init(m)
	assert.Nilf(t, err, "error initializing valid metadata properties")
}

// TestReadError excutes the Read method and fails before the Twitter API call
// go test -v -count=1 -run TestReadError ./bindings/twitter/.
func TestReadError(t *testing.T) {
	tw := NewTwitter(logger.NewLogger("test"))
	m := getTestMetadata()
	err := tw.Init(m)
	assert.Nilf(t, err, "error initializing valid metadata properties")

	tw.Read(func(res *bindings.ReadResponse) ([]byte, error) {
		t.Logf("result: %+v", res)
		assert.NotNilf(t, err, "no error on read with invalid credentials")

		return nil, nil
	})
}

// TestRead executes the Read method which calls Twiter API
// env RUN_LIVE_TW_TEST=true go test -v -count=1 -run TestReed ./bindings/twitter/.
func TestReed(t *testing.T) {
	if os.Getenv("RUN_LIVE_TW_TEST") != "true" {
		t.SkipNow() // skip this test until able to read credentials in test infra
	}
	m := bindings.Metadata{}
	m.Properties = getRuntimeMetadata()
	// add query
	m.Properties["query"] = "microsoft"
	tw := NewTwitter(logger.NewLogger("test"))
	tw.logger.SetOutputLevel(logger.DebugLevel)
	err := tw.Init(m)
	assert.Nilf(t, err, "error initializing read")

	counter := 0
	err = tw.Read(func(res *bindings.ReadResponse) ([]byte, error) {
		counter++
		t.Logf("tweet[%d]", counter)
		var tweet twitter.Tweet
		json.Unmarshal(res.Data, &tweet)
		assert.NotEmpty(t, tweet.IDStr, "tweet should have an ID")
		os.Exit(0)

		return nil, nil
	})
	assert.Nilf(t, err, "error on read")
}

// TestInvoke executes the Invoke method which calls Twiter API
// test tokens must be set
// env RUN_LIVE_TW_TEST=true go test -v -count=1 -run TestInvoke ./bindings/twitter/.
func TestInvoke(t *testing.T) {
	if os.Getenv("RUN_LIVE_TW_TEST") != "true" {
		t.SkipNow() // skip this test until able to read credentials in test infra
	}
	m := bindings.Metadata{}
	m.Properties = getRuntimeMetadata()
	tw := NewTwitter(logger.NewLogger("test"))
	tw.logger.SetOutputLevel(logger.DebugLevel)
	err := tw.Init(m)
	assert.Nilf(t, err, "error initializing Invoke")

	req := &bindings.InvokeRequest{
		Metadata: map[string]string{
			"query": "microsoft",
		},
	}

	resp, err := tw.Invoke(req)
	assert.Nilf(t, err, "error on invoke")
	assert.NotNil(t, resp)
}
