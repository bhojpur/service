package aerospike

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

func TestValidateMetadataForValidInputs(t *testing.T) {
	type testCase struct {
		name       string
		properties map[string]string
	}
	tests := []testCase{
		{"with mandatory fields", map[string]string{
			hosts:     "host1:1234",
			namespace: "foobarnamespace",
		}},
		{"with multiple hosts", map[string]string{
			hosts:     "host1:7777,host2:8888,host3:9999",
			namespace: "foobarnamespace",
		}},
		{"with optional fields", map[string]string{
			hosts:     "host1:1234",
			namespace: "foobarnamespace",
			set:       "fooset",
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			metadata := state.Metadata{Properties: test.properties}
			err := validateMetadata(metadata)
			assert.Nil(t, err)
		})
	}
}

func TestValidateMetadataForInvalidInputs(t *testing.T) {
	type testCase struct {
		name       string
		properties map[string]string
	}
	tests := []testCase{
		{"With missing hosts", map[string]string{
			namespace: "foobarnamespace",
			set:       "fooset",
		}},
		{"With invalid hosts 1", map[string]string{
			hosts:     "host1",
			namespace: "foobarnamespace",
			set:       "fooset",
		}},
		{"With invalid hosts 2", map[string]string{
			hosts:     "host1:8080,host2",
			namespace: "foobarnamespace",
			set:       "fooset",
		}},
		{"With missing namspace", map[string]string{
			hosts: "host1:1234",
			set:   "fooset",
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			metadata := state.Metadata{Properties: test.properties}
			err := validateMetadata(metadata)
			assert.NotNil(t, err)
		})
	}
}

func TestParseHostsForValidInputs(t *testing.T) {
	type testCase struct {
		name      string
		hostPorts string
	}
	tests := []testCase{
		{"valid host ports", "host1:1234"},
		{"valid multiple host ports", "host1:7777,host2:8888,host3:9999"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := parseHosts(test.hostPorts)
			assert.Nil(t, err)
			assert.NotNil(t, result)
			assert.True(t, len(result) >= 1)
		})
	}
}

func TestParseHostsForInvalidInputs(t *testing.T) {
	type testCase struct {
		name      string
		hostPorts string
	}
	tests := []testCase{
		{"missing port", "host1"},
		{"multiple entries missing port", "host1:1234,host2"},
		{"invalid port", "host1:foo"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := parseHosts(test.hostPorts)
			assert.NotNil(t, err)
		})
	}
}

func TestConvertETag(t *testing.T) {
	t.Run("valid conversion", func(t *testing.T) {
		result, err := convertETag("42")
		assert.Nil(t, err)
		assert.Equal(t, uint32(42), result)
	})

	t.Run("invalid conversion", func(t *testing.T) {
		_, err := convertETag("junk")
		assert.NotNil(t, err)
	})
}
