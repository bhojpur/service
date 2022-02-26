package postgres

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
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/utils/logger"
)

const (
	testTableDDL = `CREATE TABLE IF NOT EXISTS foo (
		id bigint NOT NULL,
		v1 character varying(50) NOT NULL,
		ts TIMESTAMP)`
	testInsert = "INSERT INTO foo (id, v1, ts) VALUES (%d, 'test-%d', '%v')"
	testDelete = "DELETE FROM foo"
	testUpdate = "UPDATE foo SET ts = '%v' WHERE id = %d"
	testSelect = "SELECT * FROM foo WHERE id < 3"
)

func TestOperations(t *testing.T) {
	t.Parallel()
	t.Run("Get operation list", func(t *testing.T) {
		b := NewPostgres(nil)
		assert.NotNil(t, b)
		l := b.Operations()
		assert.Equal(t, 3, len(l))
	})
}

// SETUP TESTS
// 1. `createdb apptest`
// 2. `createuser apptest`
// 3. `psql=# grant all privileges on database apptest to apptest;``
// 4. `export POSTGRES_TEST_CONN_URL="postgres://apptest@localhost:5432/apptest"``
// 5. `go test -v -count=1 ./bindings/postgres -run ^TestPostgresIntegration`

func TestPostgresIntegration(t *testing.T) {
	url := os.Getenv("POSTGRES_TEST_CONN_URL")
	if url == "" {
		t.SkipNow()
	}

	// live DB test
	b := NewPostgres(logger.NewLogger("test"))
	m := bindings.Metadata{Properties: map[string]string{connectionURLKey: url}}
	if err := b.Init(m); err != nil {
		t.Fatal(err)
	}

	// create table
	req := &bindings.InvokeRequest{
		Operation: execOperation,
		Metadata:  map[string]string{commandSQLKey: testTableDDL},
	}

	t.Run("Invoke create table", func(t *testing.T) {
		res, err := b.Invoke(req)
		assertResponse(t, res, err)
	})

	t.Run("Invoke delete", func(t *testing.T) {
		req.Metadata[commandSQLKey] = testDelete
		res, err := b.Invoke(req)
		assertResponse(t, res, err)
	})

	t.Run("Invoke insert", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			req.Metadata[commandSQLKey] = fmt.Sprintf(testInsert, i, i, time.Now().Format(time.RFC3339))
			res, err := b.Invoke(req)
			assertResponse(t, res, err)
		}
	})

	t.Run("Invoke update", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			req.Metadata[commandSQLKey] = fmt.Sprintf(testUpdate, time.Now().Format(time.RFC3339), i)
			res, err := b.Invoke(req)
			assertResponse(t, res, err)
		}
	})

	t.Run("Invoke select", func(t *testing.T) {
		req.Operation = queryOperation
		req.Metadata[commandSQLKey] = testSelect
		res, err := b.Invoke(req)
		assertResponse(t, res, err)
	})

	t.Run("Invoke delete", func(t *testing.T) {
		req.Operation = execOperation
		req.Metadata[commandSQLKey] = testDelete
		req.Data = nil
		res, err := b.Invoke(req)
		assertResponse(t, res, err)
	})

	t.Run("Invoke close", func(t *testing.T) {
		req.Operation = closeOperation
		req.Metadata = nil
		req.Data = nil
		_, err := b.Invoke(req)
		assert.NoError(t, err)
	})

	t.Run("Close", func(t *testing.T) {
		err := b.Close()
		assert.NoError(t, err, "expected no error closing output binding")
	})
}

func assertResponse(t *testing.T, res *bindings.InvokeResponse, err error) {
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Metadata)
}
