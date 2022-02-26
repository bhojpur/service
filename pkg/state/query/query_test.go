package query

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
)

func TestQuery(t *testing.T) {
	tests := []struct {
		input string
		query Query
	}{
		{
			input: "../../tests/state/query/q1.json",
			query: Query{
				Filters: nil,
				Sort:    nil,
				Page:    Pagination{Limit: 2, Token: ""},
				Filter:  nil,
			},
		},
		{
			input: "../../tests/state/query/q2.json",
			query: Query{
				Filters: nil,
				Sort:    nil,
				Page:    Pagination{Limit: 2, Token: ""},
				Filter:  &EQ{Key: "state", Val: "CA"},
			},
		},
		{
			input: "../../tests/state/query/q3.json",
			query: Query{
				Filters: nil,
				Sort: []Sorting{
					{Key: "state", Order: "DESC"},
					{Key: "person.name", Order: ""},
				},
				Page: Pagination{Limit: 0, Token: ""},
				Filter: &AND{
					Filters: []Filter{
						&EQ{Key: "person.org", Val: "A"},
						&IN{Key: "state", Vals: []interface{}{"CA", "WA"}},
					},
				},
			},
		},
		{
			input: "../../tests/state/query/q4.json",
			query: Query{
				Filters: nil,
				Sort: []Sorting{
					{Key: "state", Order: "DESC"},
					{Key: "person.name", Order: ""},
				},
				Page: Pagination{Limit: 2, Token: ""},
				Filter: &OR{
					Filters: []Filter{
						&EQ{Key: "person.org", Val: "A"},
						&AND{
							Filters: []Filter{
								&EQ{Key: "person.org", Val: "B"},
								&IN{Key: "state", Vals: []interface{}{"CA", "WA"}},
							},
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		data, err := ioutil.ReadFile(test.input)
		assert.NoError(t, err)
		var q Query
		err = json.Unmarshal(data, &q)
		assert.NoError(t, err)
		assert.Equal(t, test.query, q)
	}
}
