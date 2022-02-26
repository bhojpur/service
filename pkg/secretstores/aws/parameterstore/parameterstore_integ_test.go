//go:build integration
// +build integration

package parameterstore

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
	"os"
	"testing"

	"github.com/bhojpur/service/pkg/secretstores"
	"github.com/bhojpur/service/pkg/utils/logger"
	"github.com/stretchr/testify/assert"
)

// TestIntegrationGetSecret requires AWS specific environments for authentication AWS_DEFAULT_REGION AWS_ACCESS_KEY_ID,
// AWS_SECRET_ACCESS_KEY and AWS_SESSION_TOKEN
func TestIntegrationGetSecret(t *testing.T) {
	secretName := "/aws/secret/testing"
	sm := NewParameterStore(logger.NewLogger("test"))
	err := sm.Init(secretstores.Metadata{
		Properties: map[string]string{
			"Region":       os.Getenv("AWS_DEFAULT_REGION"),
			"AccessKey":    os.Getenv("AWS_ACCESS_KEY_ID"),
			"SecretKey":    os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"SessionToken": os.Getenv("AWS_SESSION_TOKEN"),
		},
	})
	assert.Nil(t, err)
	response, err := sm.GetSecret(secretstores.GetSecretRequest{
		Name:     secretName,
		Metadata: map[string]string{},
	})
	assert.Nil(t, err)
	assert.NotNil(t, response)
}

func TestIntegrationBulkGetSecret(t *testing.T) {
	sm := NewParameterStore(logger.NewLogger("test"))
	err := sm.Init(secretstores.Metadata{
		Properties: map[string]string{
			"Region":       os.Getenv("AWS_DEFAULT_REGION"),
			"AccessKey":    os.Getenv("AWS_ACCESS_KEY_ID"),
			"SecretKey":    os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"SessionToken": os.Getenv("AWS_SESSION_TOKEN"),
		},
	})
	assert.Nil(t, err)
	response, err := sm.BulkGetSecret(secretstores.BulkGetSecretRequest{})
	assert.Nil(t, err)
	assert.NotNil(t, response)
}
