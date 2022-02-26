package tablestorage

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
)

func TestGetTableStorageMetadata(t *testing.T) {
	t.Run("Nothing at all passed", func(t *testing.T) {
		m := make(map[string]string)
		_, err := getTablesMetadata(m)

		assert.NotNil(t, err)
	})

	t.Run("All parameters passed and parsed", func(t *testing.T) {
		m := make(map[string]string)
		m["accountName"] = "acc"
		m["accountKey"] = "key"
		m["tableName"] = "app"
		meta, err := getTablesMetadata(m)

		assert.Nil(t, err)
		assert.Equal(t, "acc", meta.accountName)
		assert.Equal(t, "key", meta.accountKey)
		assert.Equal(t, "app", meta.tableName)
	})
}

func TestPartitionAndRowKey(t *testing.T) {
	t.Run("Valid composite key", func(t *testing.T) {
		pk, rk := getPartitionAndRowKey("pk||rk")
		assert.Equal(t, "pk", pk)
		assert.Equal(t, "rk", rk)
	})

	t.Run("No delimiter present", func(t *testing.T) {
		pk, rk := getPartitionAndRowKey("pk_rk")
		assert.Equal(t, "pk_rk", pk)
		assert.Equal(t, "", rk)
	})
}
