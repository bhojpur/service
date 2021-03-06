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
	"fmt"
)

const (
	FILTER = "filter"
	SORT   = "sort"
	PAGE   = "page"
	ASC    = "ASC"
	DESC   = "DESC"
)

type Sorting struct {
	Key   string `json:"key"`
	Order string `json:"order,omitempty"`
}

type Pagination struct {
	Limit int    `json:"limit"`
	Token string `json:"token,omitempty"`
}

type Query struct {
	Filters map[string]interface{} `json:"filter"`
	Sort    []Sorting              `json:"sort"`
	Page    Pagination             `json:"page"`

	// derived from Filters
	Filter Filter
}

type Visitor interface {
	// returns "equal" expression
	VisitEQ(*EQ) (string, error)
	// returns "in" expression
	VisitIN(*IN) (string, error)
	// returns "and" expression
	VisitAND(*AND) (string, error)
	// returns "or" expression
	VisitOR(*OR) (string, error)
	// receives concatenated filters and finalizes the native query
	Finalize(string, *Query) error
}

type Builder struct {
	visitor Visitor
}

func NewQueryBuilder(visitor Visitor) *Builder {
	return &Builder{
		visitor: visitor,
	}
}

func (h *Builder) BuildQuery(q *Query) error {
	filters, err := h.buildFilter(q.Filter)
	if err != nil {
		return err
	}

	return h.visitor.Finalize(filters, q)
}

func (h *Builder) buildFilter(filter Filter) (string, error) {
	if filter == nil {
		return "", nil
	}
	switch f := filter.(type) {
	case *EQ:
		return h.visitor.VisitEQ(f)
	case *IN:
		return h.visitor.VisitIN(f)
	case *OR:
		return h.visitor.VisitOR(f)
	case *AND:
		return h.visitor.VisitAND(f)
	default:
		return "", fmt.Errorf("unsupported filter type %#v", filter)
	}
}

func (q *Query) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	if elem, ok := m[FILTER]; ok {
		q.Filter, err = parseFilter(elem)
		if err != nil {
			return err
		}
	}
	// setting sorting
	if elem, ok := m[SORT]; ok {
		arr, ok := elem.([]interface{})
		if !ok {
			return fmt.Errorf("%q must be an array", SORT)
		}
		jdata, err := json.Marshal(arr)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(jdata, &q.Sort); err != nil {
			return err
		}
	}
	// setting pagination
	if elem, ok := m[PAGE]; ok {
		page, ok := elem.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%q must be a map", PAGE)
		}
		jdata, err := json.Marshal(page)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(jdata, &q.Page); err != nil {
			return err
		}
	}

	return nil
}
