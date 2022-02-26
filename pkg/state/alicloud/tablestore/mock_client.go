package tablestore

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
	"bytes"
	"encoding/binary"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
)

type mockClient struct {
	tablestore.TableStoreClient

	data map[string][]byte
}

func (m *mockClient) DeleteRow(request *tablestore.DeleteRowRequest) (*tablestore.DeleteRowResponse, error) {
	var key string
	for _, col := range request.DeleteRowChange.PrimaryKey.PrimaryKeys {
		if col.ColumnName == stateKey {
			key = col.Value.(string)

			break
		}
	}

	delete(m.data, key)

	return nil, nil
}

func (m *mockClient) GetRow(request *tablestore.GetRowRequest) (*tablestore.GetRowResponse, error) {
	var key string
	for _, col := range request.SingleRowQueryCriteria.PrimaryKey.PrimaryKeys {
		if col.ColumnName == stateKey {
			key = col.Value.(string)

			break
		}
	}

	val := m.data[key]

	resp := &tablestore.GetRowResponse{
		Columns: []*tablestore.AttributeColumn{{
			ColumnName: stateValue,
			Value:      val,
		}},
	}

	return resp, nil
}

func (m *mockClient) UpdateRow(req *tablestore.UpdateRowRequest) (*tablestore.UpdateRowResponse, error) {
	change := req.UpdateRowChange

	var val []byte
	var key string

	for _, col := range change.PrimaryKey.PrimaryKeys {
		if col.ColumnName == stateKey {
			key = col.Value.(string)

			break
		}
	}

	for _, col := range change.Columns {
		if col.ColumnName == stateValue {
			buf := &bytes.Buffer{}
			binary.Write(buf, binary.BigEndian, col.Value)
			val = buf.Bytes()

			break
		}
	}

	m.data[key] = val

	return nil, nil
}

func (m *mockClient) BatchGetRow(request *tablestore.BatchGetRowRequest) (*tablestore.BatchGetRowResponse, error) {
	resp := &tablestore.BatchGetRowResponse{
		TableToRowsResult: map[string][]tablestore.RowResult{},
	}

	for _, criteria := range request.MultiRowQueryCriteria {
		tableRes := resp.TableToRowsResult[criteria.TableName]
		if tableRes == nil {
			tableRes = []tablestore.RowResult{}
		}
		for _, keys := range criteria.PrimaryKey {
			for _, key := range keys.PrimaryKeys {
				if key.ColumnName == stateKey {
					pk := key.Value.(string)

					if m.data[pk] == nil {
						continue
					}

					value := m.data[key.Value.(string)]
					tableRes = append(tableRes, tablestore.RowResult{
						TableName: criteria.TableName,
						Columns: []*tablestore.AttributeColumn{
							{
								ColumnName: stateValue,
								Value:      value,
							},
						},
						PrimaryKey: tablestore.PrimaryKey{
							PrimaryKeys: []*tablestore.PrimaryKeyColumn{
								{
									ColumnName: stateKey,
									Value:      key.Value,
								},
							},
						},
					})
					resp.TableToRowsResult[criteria.TableName] = tableRes

					break
				}
			}
		}
	}

	return resp, nil
}

func (m *mockClient) BatchWriteRow(request *tablestore.BatchWriteRowRequest) (*tablestore.BatchWriteRowResponse, error) {
	resp := &tablestore.BatchWriteRowResponse{}
	for _, changes := range request.RowChangesGroupByTable {
		for _, change := range changes {
			switch inst := change.(type) {
			case *tablestore.UpdateRowChange:
				var pk string
				for _, col := range inst.PrimaryKey.PrimaryKeys {
					if col.ColumnName == stateKey {
						pk = col.Value.(string)

						break
					}
				}

				for _, col := range inst.Columns {
					if col.ColumnName == stateValue {
						buf := &bytes.Buffer{}
						binary.Write(buf, binary.BigEndian, col.Value)
						m.data[pk] = buf.Bytes()
					}
				}

			case *tablestore.DeleteRowChange:
				for _, col := range inst.PrimaryKey.PrimaryKeys {
					if col.ColumnName == stateKey {
						delete(m.data, col.Value.(string))

						break
					}
				}
			}
		}
	}

	return resp, nil
}
