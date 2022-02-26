package postgresql

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
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/state/query"
)

func TestPostgresqlQueryBuildQuery(t *testing.T) {
	tests := []struct {
		input string
		query string
	}{
		{
			input: "../../tests/state/query/q1.json",
			query: "SELECT key, value, xmin as etag FROM state LIMIT 2",
		},
		{
			input: "../../tests/state/query/q2.json",
			query: "SELECT key, value, xmin as etag FROM state WHERE value->>'state'=$1 LIMIT 2",
		},
		{
			input: "../../tests/state/query/q2-token.json",
			query: "SELECT key, value, xmin as etag FROM state WHERE value->>'state'=$1 LIMIT 2 OFFSET 2",
		},
		{
			input: "../../tests/state/query/q3.json",
			query: "SELECT key, value, xmin as etag FROM state WHERE (value->'person'->>'org'=$1 AND (value->>'state'=$2 OR value->>'state'=$3)) ORDER BY value->>'state' DESC, value->'person'->>'name'",
		},
		{
			input: "../../tests/state/query/q4.json",
			query: "SELECT key, value, xmin as etag FROM state WHERE (value->'person'->>'org'=$1 OR (value->'person'->>'org'=$2 AND (value->>'state'=$3 OR value->>'state'=$4))) ORDER BY value->>'state' DESC, value->'person'->>'name' LIMIT 2",
		},
		{
			input: "../../tests/state/query/q5.json",
			query: "SELECT key, value, xmin as etag FROM state WHERE (value->'person'->>'org'=$1 AND (value->'person'->>'name'=$2 OR (value->>'state'=$3 OR value->>'state'=$4))) ORDER BY value->>'state' DESC, value->'person'->>'name' LIMIT 2",
		},
	}
	for _, test := range tests {
		data, err := ioutil.ReadFile(test.input)
		assert.NoError(t, err)
		var qq query.Query
		err = json.Unmarshal(data, &qq)
		assert.NoError(t, err)

		q := &Query{}
		qbuilder := query.NewQueryBuilder(q)
		err = qbuilder.BuildQuery(&qq)
		assert.NoError(t, err)
		assert.Equal(t, test.query, q.query)
	}
}
