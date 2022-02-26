package mongodb

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
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/bhojpur/service/pkg/state"
	"github.com/bhojpur/service/pkg/state/query"
)

type Query struct {
	query  string
	filter interface{}
	opts   *options.FindOptions
}

func (q *Query) VisitEQ(f *query.EQ) (string, error) {
	// { <key>: <val> }
	return fmt.Sprintf(`{ "value.%s": %q }`, f.Key, f.Val), nil
}

func (q *Query) VisitIN(f *query.IN) (string, error) {
	// { $in: [ <val1>, <val2>, ... , <valN> ] }
	if len(f.Vals) == 0 {
		return "", fmt.Errorf("empty IN operator for key %q", f.Key)
	}
	str := fmt.Sprintf(`{ "value.%s": { "$in": [ %q`, f.Key, f.Vals[0])
	for _, v := range f.Vals[1:] {
		str += fmt.Sprintf(", %q", v)
	}
	str += " ] } }"

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

	return fmt.Sprintf(`{ "%s": [ %s ] }`, op, strings.Join(arr, ", ")), nil
}

func (q *Query) VisitAND(f *query.AND) (string, error) {
	// { $and: [ { <expression1> }, { <expression2> } , ... , { <expressionN> } ] }
	return q.visitFilters("$and", f.Filters)
}

func (q *Query) VisitOR(f *query.OR) (string, error) {
	// { $or: [ { <expression1> }, { <expression2> } , ... , { <expressionN> } ] }
	return q.visitFilters("$or", f.Filters)
}

func (q *Query) Finalize(filters string, qq *query.Query) error {
	q.query = filters
	if len(filters) == 0 {
		q.filter = bson.D{}
	} else if err := bson.UnmarshalExtJSON([]byte(filters), false, &q.filter); err != nil {
		return err
	}
	q.opts = options.Find()

	// sorting
	if len(qq.Sort) > 0 {
		sort := bson.D{}
		for _, s := range qq.Sort {
			order := 1 // ascending
			if s.Order == query.DESC {
				order = -1
			}
			sort = append(sort, bson.E{Key: "value." + s.Key, Value: order})
		}
		q.opts.SetSort(sort)
	}
	// pagination
	if qq.Page.Limit > 0 {
		q.opts.SetLimit(int64(qq.Page.Limit))
	}
	if len(qq.Page.Token) != 0 {
		skip, err := strconv.ParseInt(qq.Page.Token, 10, 64)
		if err != nil {
			return err
		}
		q.opts.SetSkip(skip)
	}

	return nil
}

func (q *Query) execute(ctx context.Context, collection *mongo.Collection) ([]state.QueryItem, string, error) {
	cur, err := collection.Find(ctx, q.filter, []*options.FindOptions{q.opts}...)
	if err != nil {
		return nil, "", err
	}
	defer cur.Close(ctx)
	ret := []state.QueryItem{}
	for cur.Next(ctx) {
		var item Item
		if err = cur.Decode(&item); err != nil {
			return nil, "", err
		}
		result := state.QueryItem{
			Key:  item.Key,
			ETag: &item.Etag,
		}

		switch obj := item.Value.(type) {
		case string:
			result.Data = []byte(obj)
		case primitive.D:
			if result.Data, err = bson.MarshalExtJSON(obj, true, true); err != nil {
				result.Error = err.Error()
			}
		default:
			if result.Data, err = json.Marshal(item.Value); err != nil {
				result.Error = err.Error()
			}
		}
		ret = append(ret, result)
	}
	if err = cur.Err(); err != nil {
		return nil, "", err
	}
	// set next query token only if limit is specified
	var token string
	if q.opts.Limit != nil && *q.opts.Limit != 0 {
		var skip int64
		if q.opts.Skip != nil {
			skip = *q.opts.Skip
		}
		token = strconv.FormatInt(skip+int64(len(ret)), 10)
	}

	return ret, token, nil
}
