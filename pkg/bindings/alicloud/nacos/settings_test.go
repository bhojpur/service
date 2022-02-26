package nacos_test

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
	"github.com/stretchr/testify/require"

	"github.com/bhojpur/service/pkg/bindings/alicloud/nacos"
)

func TestParseMetadata(t *testing.T) { //nolint:paralleltest
	props := map[string]string{
		"endpoint":        "a",
		"region":          "b",
		"namespace":       "c",
		"accessKey":       "d",
		"secretKey":       "e",
		"updateThreadNum": "3",
	}

	var settings nacos.Settings
	err := settings.Decode(props)
	require.NoError(t, err)
	assert.Equal(t, "a", settings.Endpoint)
	assert.Equal(t, "b", settings.RegionID)
	assert.Equal(t, "c", settings.NamespaceID)
	assert.Equal(t, "d", settings.AccessKey)
	assert.Equal(t, "e", settings.SecretKey)
	assert.Equal(t, 3, settings.UpdateThreadNum)
}
