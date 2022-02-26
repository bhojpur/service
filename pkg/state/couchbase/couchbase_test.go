package couchbase

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

	"github.com/stretchr/testify/assert"
	"gopkg.in/couchbase/gocb.v1"

	"github.com/bhojpur/service/pkg/state"
)

func TestValidateMetadata(t *testing.T) {
	t.Run("with mandatory fields", func(t *testing.T) {
		props := map[string]string{
			couchbaseURL: "foo://bar",
			username:     "pramila",
			password:     "secret",
			bucketName:   "testbucket",
		}
		metadata := state.Metadata{Properties: props}

		err := validateMetadata(metadata)
		assert.Equal(t, nil, err)
	})
	t.Run("with optional fields", func(t *testing.T) {
		props := map[string]string{
			couchbaseURL:                  "foo://bar",
			username:                      "pramila",
			password:                      "secret",
			bucketName:                    "testbucket",
			numReplicasDurablePersistence: "1",
			numReplicasDurableReplication: "2",
		}
		metadata := state.Metadata{Properties: props}

		err := validateMetadata(metadata)
		assert.Equal(t, nil, err)
	})
	t.Run("With missing couchbase URL", func(t *testing.T) {
		props := map[string]string{
			username:   "pramila",
			password:   "secret",
			bucketName: "testbucket",
		}
		metadata := state.Metadata{Properties: props}
		err := validateMetadata(metadata)
		assert.NotNil(t, err)
	})
	t.Run("With missing username", func(t *testing.T) {
		props := map[string]string{
			couchbaseURL: "foo://bar",
			password:     "secret",
			bucketName:   "testbucket",
		}
		metadata := state.Metadata{Properties: props}
		err := validateMetadata(metadata)
		assert.NotNil(t, err)
	})
	t.Run("With missing password", func(t *testing.T) {
		props := map[string]string{
			couchbaseURL: "foo://bar",
			username:     "pramila",
			bucketName:   "testbucket",
		}
		metadata := state.Metadata{Properties: props}
		err := validateMetadata(metadata)
		assert.NotNil(t, err)
	})
	t.Run("With missing bucket", func(t *testing.T) {
		props := map[string]string{
			couchbaseURL: "foo://bar",
			username:     "pramila",
			password:     "secret",
		}
		metadata := state.Metadata{Properties: props}
		err := validateMetadata(metadata)
		assert.NotNil(t, err)
	})
	t.Run("With invalid durable replication", func(t *testing.T) {
		props := map[string]string{
			couchbaseURL:                  "foo://bar",
			username:                      "pramila",
			password:                      "secret",
			numReplicasDurableReplication: "junk",
		}
		metadata := state.Metadata{Properties: props}
		err := validateMetadata(metadata)
		assert.NotNil(t, err)
	})
	t.Run("With invalid durable persistence", func(t *testing.T) {
		props := map[string]string{
			couchbaseURL:                  "foo://bar",
			username:                      "pramila",
			password:                      "secret",
			numReplicasDurablePersistence: "junk",
		}
		metadata := state.Metadata{Properties: props}
		err := validateMetadata(metadata)
		assert.NotNil(t, err)
	})
}

func TestETagToCas(t *testing.T) {
	t.Run("with valid string", func(t *testing.T) {
		casStr := "1572938024378368000"
		ver := uint64(1572938024378368000)
		var expectedCas gocb.Cas = gocb.Cas(ver)
		cas, err := eTagToCas(casStr)
		assert.Equal(t, nil, err)
		assert.Equal(t, expectedCas, cas)
	})
	t.Run("with empty string", func(t *testing.T) {
		_, err := eTagToCas("")
		assert.NotNil(t, err)
	})
}
