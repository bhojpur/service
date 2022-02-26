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
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/agrea/ptr"

	"github.com/bhojpur/service/pkg/state"
	"github.com/bhojpur/service/pkg/state/query"
	"github.com/bhojpur/service/pkg/utils/logger"
)

type Query struct {
	query  string
	params []interface{}
	limit  int
	skip   *int64
}

func (q *Query) VisitEQ(f *query.EQ) (string, error) {
	return q.whereFieldEqual(f.Key, f.Val), nil
}

func (q *Query) VisitIN(f *query.IN) (string, error) {
	if len(f.Vals) == 0 {
		return "", fmt.Errorf("empty IN operator for key %q", f.Key)
	}

	str := "("
	str += q.whereFieldEqual(f.Key, f.Vals[0])

	for _, v := range f.Vals[1:] {
		str += " OR "
		str += q.whereFieldEqual(f.Key, v)
	}
	str += ")"
	return str, nil
}

func (q *Query) visitFilters(op string, filters []query.Filter) (string, error) {
	var (
		arr []string
		str string
		err error
	)

	for _, fil := range filters {
		switch f := fil.(type) {
		case *query.EQ:
			if str, err = q.VisitEQ(f); err != nil {
				return "", err
			}
			arr = append(arr, str)
		case *query.IN:
			if str, err = q.VisitIN(f); err != nil {
				return "", err
			}
			arr = append(arr, str)
		case *query.OR:
			if str, err = q.VisitOR(f); err != nil {
				return "", err
			}
			arr = append(arr, str)
		case *query.AND:
			if str, err = q.VisitAND(f); err != nil {
				return "", err
			}
			arr = append(arr, str)
		default:
			return "", fmt.Errorf("unsupported filter type %#v", f)
		}
	}

	sep := fmt.Sprintf(" %s ", op)

	return fmt.Sprintf("(%s)", strings.Join(arr, sep)), nil
}

func (q *Query) VisitAND(f *query.AND) (string, error) {
	return q.visitFilters("AND", f.Filters)
}

func (q *Query) VisitOR(f *query.OR) (string, error) {
	return q.visitFilters("OR", f.Filters)
}

func (q *Query) Finalize(filters string, qq *query.Query) error {
	q.query = fmt.Sprintf("SELECT key, value, xmin as etag FROM %s", tableName)

	if filters != "" {
		q.query += fmt.Sprintf(" WHERE %s", filters)
	}

	if len(qq.Sort) > 0 {
		q.query += " ORDER BY "

		for sortIndex, sortItem := range qq.Sort {
			if sortIndex > 0 {
				q.query += ", "
			}
			q.query += translateFieldToFilter(sortItem.Key)
			if sortItem.Order != "" {
				q.query += fmt.Sprintf(" %s", sortItem.Order)
			}
		}
	}

	if qq.Page.Limit > 0 {
		q.query += fmt.Sprintf(" LIMIT %d", qq.Page.Limit)
		q.limit = qq.Page.Limit
	}

	if len(qq.Page.Token) != 0 {
		skip, err := strconv.ParseInt(qq.Page.Token, 10, 64)
		if err != nil {
			return err
		}
		q.query += fmt.Sprintf(" OFFSET %d", skip)
		q.skip = &skip
	}

	return nil
}

func (q *Query) execute(logger logger.Logger, db *sql.DB) ([]state.QueryItem, string, error) {
	rows, err := db.Query(q.query, q.params...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	ret := []state.QueryItem{}
	for rows.Next() {
		var (
			key  string
			data []byte
			etag int
		)
		if err = rows.Scan(&key, &data, &etag); err != nil {
			return nil, "", err
		}
		result := state.QueryItem{
			Key:  key,
			Data: data,
			ETag: ptr.String(strconv.Itoa(etag)),
		}
		ret = append(ret, result)
	}

	if err = rows.Err(); err != nil {
		return nil, "", err
	}

	var token string
	if q.limit != 0 {
		var skip int64
		if q.skip != nil {
			skip = *q.skip
		}
		token = strconv.FormatInt(skip+int64(len(ret)), 10)
	}

	return ret, token, nil
}

func (q *Query) addParamValueAndReturnPosition(value interface{}) int {
	q.params = append(q.params, fmt.Sprintf("%v", value))
	return len(q.params)
}

func translateFieldToFilter(key string) string {
	// add preceding "value"
	key = "value." + key

	fieldParts := strings.Split(key, ".")
	filterField := fieldParts[0]
	fieldParts = fieldParts[1:]

	for fieldIndex, fieldPart := range fieldParts {
		filterField += "->"

		if fieldIndex+1 == len(fieldParts) {
			filterField += ">"
		}

		filterField += fmt.Sprintf("'%s'", fieldPart)
	}

	return filterField
}

func (q *Query) whereFieldEqual(key string, value interface{}) string {
	position := q.addParamValueAndReturnPosition(value)
	filterField := translateFieldToFilter(key)
	query := fmt.Sprintf("%s=$%v", filterField, position)
	return query
}
