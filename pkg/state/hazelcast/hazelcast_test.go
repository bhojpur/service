package hazelcast

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

	"github.com/bhojpur/service/pkg/state"
)

func TestValidateMetadata(t *testing.T) {
	t.Run("without required configuration", func(t *testing.T) {
		properties := map[string]string{}
		m := state.Metadata{
			Properties: properties,
		}
		err := validateMetadata(m)
		assert.NotNil(t, err)
	})

	t.Run("without server configuration", func(t *testing.T) {
		properties := map[string]string{
			"hazelcastMap": "foo-map",
		}
		m := state.Metadata{
			Properties: properties,
		}
		err := validateMetadata(m)
		assert.NotNil(t, err)
	})

	t.Run("without map configuration", func(t *testing.T) {
		properties := map[string]string{
			"hazelcastServers": "hz1:5701",
		}
		m := state.Metadata{
			Properties: properties,
		}
		err := validateMetadata(m)
		assert.NotNil(t, err)
	})

	t.Run("with valid configuration", func(t *testing.T) {
		properties := map[string]string{
			"hazelcastServers": "hz1:5701",
			"hazelcastMap":     "foo-map",
		}
		m := state.Metadata{
			Properties: properties,
		}
		err := validateMetadata(m)
		assert.Nil(t, err)
	})
}
