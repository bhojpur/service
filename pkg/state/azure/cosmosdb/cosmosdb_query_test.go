package cosmosdb

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

	"github.com/a8m/documentdb"
	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/state/query"
)

func TestCosmosDbKeyReplace(t *testing.T) {
	tests := []struct{ input, expected string }{
		{
			input:    "c.a",
			expected: "c.a",
		},
		{
			input:    "c.a.b",
			expected: "c.a.b",
		},
		{
			input:    "c.value",
			expected: "c['value']",
		},
		{
			input:    "c.value.a",
			expected: "c['value'].a",
		},
		{
			input:    "c.value.value",
			expected: "c['value']['value']",
		},
		{
			input:    "c.value.a.value",
			expected: "c['value'].a['value']",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, replaceKeywords(test.input))
	}
}

func TestCosmosDbQuery(t *testing.T) {
	tests := []struct {
		input string
		query documentdb.Query
	}{
		{
			input: "../../../tests/state/query/q1.json",
			query: documentdb.Query{
				Query:      "SELECT * FROM c",
				Parameters: nil,
			},
		},
		{
			input: "../../../tests/state/query/q2.json",
			query: documentdb.Query{
				Query: "SELECT * FROM c WHERE c['value'].state = @__param__0__",
				Parameters: []documentdb.Parameter{
					{
						Name:  "@__param__0__",
						Value: "CA",
					},
				},
			},
		},
		{
			input: "../../../tests/state/query/q3.json",
			query: documentdb.Query{
				Query: "SELECT * FROM c WHERE c['value'].person.org = @__param__0__ AND c['value'].state IN (@__param__1__, @__param__2__) ORDER BY c['value'].state DESC, c['value'].person.name ASC",
				Parameters: []documentdb.Parameter{
					{
						Name:  "@__param__0__",
						Value: "A",
					},
					{
						Name:  "@__param__1__",
						Value: "CA",
					},
					{
						Name:  "@__param__2__",
						Value: "WA",
					},
				},
			},
		},
		{
			input: "../../../tests/state/query/q4.json",
			query: documentdb.Query{
				Query: "SELECT * FROM c WHERE c['value'].person.org = @__param__0__ OR (c['value'].person.org = @__param__1__ AND c['value'].state IN (@__param__2__, @__param__3__)) ORDER BY c['value'].state DESC, c['value'].person.name ASC",
				Parameters: []documentdb.Parameter{
					{
						Name:  "@__param__0__",
						Value: "A",
					},
					{
						Name:  "@__param__1__",
						Value: "B",
					},
					{
						Name:  "@__param__2__",
						Value: "CA",
					},
					{
						Name:  "@__param__3__",
						Value: "WA",
					},
				},
			},
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
