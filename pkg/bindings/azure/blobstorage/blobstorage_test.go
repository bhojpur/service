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
	"testing"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/utils/logger"
)

func TestParseMetadata(t *testing.T) {
	m := bindings.Metadata{}
	blobStorage := NewAzureBlobStorage(logger.NewLogger("test"))

	t.Run("parse all metadata", func(t *testing.T) {
		m.Properties = map[string]string{
			"storageAccount":    "account",
			"storageAccessKey":  "key",
			"container":         "test",
			"getBlobRetryCount": "5",
			"decodeBase64":      "true",
		}
		meta, err := blobStorage.parseMetadata(m)
		assert.Nil(t, err)
		assert.Equal(t, "test", meta.Container)
		assert.Equal(t, "account", meta.StorageAccount)
		assert.Equal(t, "key", meta.StorageAccessKey)
		assert.Equal(t, true, meta.DecodeBase64)
		assert.Equal(t, 5, meta.GetBlobRetryCount)
		assert.Equal(t, azblob.PublicAccessNone, meta.PublicAccessLevel)
	})

	t.Run("parse metadata with publicAccessLevel = blob", func(t *testing.T) {
		m.Properties = map[string]string{
			"publicAccessLevel": "blob",
		}
		meta, err := blobStorage.parseMetadata(m)
		assert.Nil(t, err)
		assert.Equal(t, azblob.PublicAccessBlob, meta.PublicAccessLevel)
	})

	t.Run("parse metadata with publicAccessLevel = container", func(t *testing.T) {
		m.Properties = map[string]string{
			"publicAccessLevel": "container",
		}
		meta, err := blobStorage.parseMetadata(m)
		assert.Nil(t, err)
		assert.Equal(t, azblob.PublicAccessContainer, meta.PublicAccessLevel)
	})

	t.Run("parse metadata with invalid publicAccessLevel", func(t *testing.T) {
		m.Properties = map[string]string{
			"publicAccessLevel": "invalid",
		}
		_, err := blobStorage.parseMetadata(m)
		assert.Error(t, err)
	})
}

func TestGetOption(t *testing.T) {
	blobStorage := NewAzureBlobStorage(logger.NewLogger("test"))

	t.Run("return error if blobName is missing", func(t *testing.T) {
		r := bindings.InvokeRequest{}
		_, err := blobStorage.get(&r)
		if assert.Error(t, err) {
			assert.Equal(t, ErrMissingBlobName, err)
		}
	})
}

func TestDeleteOption(t *testing.T) {
	blobStorage := NewAzureBlobStorage(logger.NewLogger("test"))

	t.Run("return error if blobName is missing", func(t *testing.T) {
		r := bindings.InvokeRequest{}
		_, err := blobStorage.delete(&r)
		if assert.Error(t, err) {
			assert.Equal(t, ErrMissingBlobName, err)
		}
	})

	t.Run("return error for invalid deleteSnapshots", func(t *testing.T) {
		r := bindings.InvokeRequest{}
		r.Metadata = map[string]string{
			"blobName":        "foo",
			"deleteSnapshots": "invalid",
		}
		_, err := blobStorage.delete(&r)
		assert.Error(t, err)
	})
}
