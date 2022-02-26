package env

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

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/secretstores"
	"github.com/bhojpur/service/pkg/utils/logger"
)

func TestInit(t *testing.T) {
	secret := "secret1"
	key := "TEST_SECRET"

	s := envSecretStore{logger: logger.NewLogger("test")}

	os.Setenv(key, secret)
	assert.Equal(t, secret, os.Getenv(key))

	t.Run("Test init", func(t *testing.T) {
		err := s.Init(secretstores.Metadata{})
		assert.Nil(t, err)
	})

	t.Run("Test set and get", func(t *testing.T) {
		err := s.Init(secretstores.Metadata{})
		assert.Nil(t, err)
		resp, err := s.GetSecret(secretstores.GetSecretRequest{Name: key})
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, secret, resp.Data[key])
	})

	t.Run("Test bulk get", func(t *testing.T) {
		err := s.Init(secretstores.Metadata{})
		assert.Nil(t, err)
		resp, err := s.BulkGetSecret(secretstores.BulkGetSecretRequest{})
		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, secret, resp.Data[key][key])
	})
}
