package blobstorage

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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/state"
	"github.com/bhojpur/service/pkg/utils/logger"
)

func TestInit(t *testing.T) {
	m := state.Metadata{}
	s := NewAzureBlobStorageStore(logger.NewLogger("logger"))
	t.Run("Init with valid metadata", func(t *testing.T) {
		m.Properties = map[string]string{
			"accountName":   "acc",
			"accountKey":    "e+Dnvl8EOxYxV94nurVaRQ==",
			"containerName": "app",
		}
		err := s.Init(m)
		assert.Nil(t, err)
		assert.Equal(t, "acc.blob.core.windows.net", s.containerURL.URL().Host)
		assert.Equal(t, "/app", s.containerURL.URL().Path)
	})

	t.Run("Init with missing metadata", func(t *testing.T) {
		m.Properties = map[string]string{
			"invalidValue": "a",
		}
		err := s.Init(m)
		assert.NotNil(t, err)
		assert.Equal(t, err, fmt.Errorf("missing or empty accountName field from metadata"))
	})

	t.Run("Init with invalid account name", func(t *testing.T) {
		m.Properties = map[string]string{
			"accountName":   "invalid-account",
			"accountKey":    "e+Dnvl8EOxYxV94nurVaRQ==",
			"containerName": "app",
		}
		err := s.Init(m)
		assert.NotNil(t, err)
	})
}

func TestGetBlobStorageMetaData(t *testing.T) {
	t.Run("Nothing at all passed", func(t *testing.T) {
		m := make(map[string]string)
		_, err := getBlobStorageMetadata(m)

		assert.NotNil(t, err)
	})

	t.Run("All parameters passed and parsed", func(t *testing.T) {
		m := make(map[string]string)
		m["accountName"] = "acc"
		m["containerName"] = "app"
		meta, err := getBlobStorageMetadata(m)

		assert.Nil(t, err)
		assert.Equal(t, "acc", meta.accountName)
		assert.Equal(t, "app", meta.containerName)
	})
}

func TestFileName(t *testing.T) {
	t.Run("Valid composite key", func(t *testing.T) {
		key := getFileName("app_id||key")
		assert.Equal(t, "key", key)
	})

	t.Run("No delimiter present", func(t *testing.T) {
		key := getFileName("key")
		assert.Equal(t, "key", key)
	})
}

func TestBlobHTTPHeaderGeneration(t *testing.T) {
	s := NewAzureBlobStorageStore(logger.NewLogger("logger"))
	t.Run("Content type is set from request, forward compatibility", func(t *testing.T) {
		contentType := "application/json"
		req := &state.SetRequest{
			ContentType: &contentType,
		}

		blobHeaders, err := s.createBlobHTTPHeadersFromRequest(req)
		assert.Nil(t, err)
		assert.Equal(t, "application/json", blobHeaders.ContentType)
	})
	t.Run("Content type and metadata provided (conflict), content type chosen", func(t *testing.T) {
		contentType := "application/json"
		req := &state.SetRequest{
			ContentType: &contentType,
			Metadata: map[string]string{
				contentType: "text/plain",
			},
		}

		blobHeaders, err := s.createBlobHTTPHeadersFromRequest(req)
		assert.Nil(t, err)
		assert.Equal(t, "application/json", blobHeaders.ContentType)
	})
	t.Run("ContentType not provided, metadata provided set backward compatibility", func(t *testing.T) {
		req := &state.SetRequest{
			Metadata: map[string]string{
				contentType: "text/plain",
			},
		}

		blobHeaders, err := s.createBlobHTTPHeadersFromRequest(req)
		assert.Nil(t, err)
		assert.Equal(t, "text/plain", blobHeaders.ContentType)
	})
}
