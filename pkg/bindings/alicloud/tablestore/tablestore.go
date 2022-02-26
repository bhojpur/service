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
	"encoding/json"
	"strings"
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/utils/logger"

	"github.com/pkg/errors"
)

const (
	tableName   = "tableName"
	columnToGet = "columnToGet"
	primaryKeys = "primaryKeys"

	invokeStartTimeKey = "start-time"
	invokeEndTimeKey   = "end-time"
	invokeDurationKey  = "duration"
)

type tablestoreMetadata struct {
	Endpoint     string `json:"endpoint"`
	AccessKeyID  string `json:"accessKeyID"`
	AccessKey    string `json:"accessKey"`
	InstanceName string `json:"instanceName"`
	TableName    string `json:"tableName"`
}

type AliCloudTableStore struct {
	logger   logger.Logger
	client   *tablestore.TableStoreClient
	metadata tablestoreMetadata
}

func NewAliCloudTableStore(log logger.Logger) *AliCloudTableStore {
	return &AliCloudTableStore{
		logger: log,
		client: nil,
	}
}

func (s *AliCloudTableStore) Init(metadata bindings.Metadata) error {
	m, err := s.parseMetadata(metadata)
	if err != nil {
		return err
	}

	s.metadata = *m
	s.client = tablestore.NewClient(m.Endpoint, m.InstanceName, m.AccessKeyID, m.AccessKey)

	return nil
}

func (s *AliCloudTableStore) Invoke(req *bindings.InvokeRequest) (*bindings.InvokeResponse, error) {
	if req == nil {
		return nil, errors.Errorf("invoke request required")
	}

	startTime := time.Now().UTC()
	resp := &bindings.InvokeResponse{
		Metadata: map[string]string{
			invokeStartTimeKey: startTime.Format(time.RFC3339Nano),
		},
	}

	switch req.Operation {
	case bindings.GetOperation:
		err := s.get(req, resp)
		if err != nil {
			return nil, err
		}
	case bindings.ListOperation:
		err := s.list(req, resp)
		if err != nil {
			return nil, err
		}
	case bindings.CreateOperation:
		err := s.create(req, resp)
		if err != nil {
			return nil, err
		}
	case bindings.DeleteOperation:
		err := s.delete(req, resp)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.Errorf("invalid operation type: %s. Expected %s, %s, %s, or %s",
			req.Operation, bindings.GetOperation, bindings.ListOperation, bindings.CreateOperation, bindings.DeleteOperation)
	}

	endTime := time.Now().UTC()
	resp.Metadata[invokeEndTimeKey] = endTime.Format(time.RFC3339Nano)
	resp.Metadata[invokeDurationKey] = endTime.Sub(startTime).String()

	return resp, nil
}

func (s *AliCloudTableStore) Operations() []bindings.OperationKind {
	return []bindings.OperationKind{bindings.CreateOperation, bindings.DeleteOperation, bindings.GetOperation, bindings.ListOperation}
}

func (s *AliCloudTableStore) parseMetadata(metadata bindings.Metadata) (*tablestoreMetadata, error) {
	b, err := json.Marshal(metadata.Properties)
	if err != nil {
		return nil, err
	}

	var m tablestoreMetadata
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (s *AliCloudTableStore) get(req *bindings.InvokeRequest, resp *bindings.InvokeResponse) error {
	columns := strings.Split(req.Metadata[columnToGet], ",")
	pkNames := strings.Split(req.Metadata[primaryKeys], ",")
	pks := make([]*tablestore.PrimaryKeyColumn, len(pkNames))

	data := make(map[string]interface{})
	err := json.Unmarshal(req.Data, &data)
	if err != nil {
		return err
	}

	for idx, pkName := range pkNames {
		pks[idx] = &tablestore.PrimaryKeyColumn{
			ColumnName: pkName,
			Value:      data[pkName],
		}
	}

	criteria := &tablestore.SingleRowQueryCriteria{
		TableName:    s.getTableName(req.Metadata),
		PrimaryKey:   &tablestore.PrimaryKey{PrimaryKeys: pks},
		ColumnsToGet: columns,
		MaxVersion:   1,
	}
	getRowReq := &tablestore.GetRowRequest{
		SingleRowQueryCriteria: criteria,
	}
	getRowResp, err := s.client.GetRow(getRowReq)
	if err != nil {
		return err
	}

	ret, err := s.unmarshal(getRowResp.PrimaryKey.PrimaryKeys, getRowResp.Columns)
	if err != nil {
		return err
	}

	if ret == nil {
		resp.Data = nil

		return nil
	}

	resp.Data, err = json.Marshal(ret)

	return err
}

func (s *AliCloudTableStore) list(req *bindings.InvokeRequest, resp *bindings.InvokeResponse) error {
	columns := strings.Split(req.Metadata[columnToGet], ",")
	pkNames := strings.Split(req.Metadata[primaryKeys], ",")

	var data []map[string]interface{}
	err := json.Unmarshal(req.Data, &data)
	if err != nil {
		return err
	}

	criteria := &tablestore.MultiRowQueryCriteria{
		TableName:    s.getTableName(req.Metadata),
		ColumnsToGet: columns,
		MaxVersion:   1,
	}

	for _, item := range data {
		pk := &tablestore.PrimaryKey{}
		for _, pkName := range pkNames {
			pk.AddPrimaryKeyColumn(pkName, item[pkName])
		}
		criteria.AddRow(pk)
	}

	getRowRequest := &tablestore.BatchGetRowRequest{}
	getRowRequest.MultiRowQueryCriteria = append(getRowRequest.MultiRowQueryCriteria, criteria)
	getRowResp, err := s.client.BatchGetRow(getRowRequest)
	if err != nil {
		return err
	}

	var ret []interface{}

	for _, criteria := range getRowRequest.MultiRowQueryCriteria {
		for _, row := range getRowResp.TableToRowsResult[criteria.TableName] {
			rowData, rowErr := s.unmarshal(row.PrimaryKey.PrimaryKeys, row.Columns)
			if rowErr != nil {
				return rowErr
			}
			ret = append(ret, rowData)
		}
	}

	resp.Data, err = json.Marshal(ret)

	return err
}

func (s *AliCloudTableStore) create(req *bindings.InvokeRequest, resp *bindings.InvokeResponse) error {
	data := make(map[string]interface{})
	err := json.Unmarshal(req.Data, &data)
	if err != nil {
		return err
	}
	pkNames := strings.Split(req.Metadata[primaryKeys], ",")
	pks := make([]*tablestore.PrimaryKeyColumn, len(pkNames))
	columns := make([]tablestore.AttributeColumn, len(data)-len(pkNames))

	for idx, pk := range pkNames {
		pks[idx] = &tablestore.PrimaryKeyColumn{
			ColumnName: pk,
			Value:      data[pk],
		}
	}

	idx := 0
	for key, val := range data {
		if !contains(pkNames, key) {
			columns[idx] = tablestore.AttributeColumn{
				ColumnName: key,
				Value:      val,
			}
			idx++
		}
	}

	change := tablestore.PutRowChange{
		TableName:     s.getTableName(req.Metadata),
		PrimaryKey:    &tablestore.PrimaryKey{PrimaryKeys: pks},
		Columns:       columns,
		ReturnType:    tablestore.ReturnType_RT_NONE,
		TransactionId: nil,
	}

	change.SetCondition(tablestore.RowExistenceExpectation_IGNORE)

	putRequest := &tablestore.PutRowRequest{
		PutRowChange: &change,
	}

	_, err = s.client.PutRow(putRequest)

	if err != nil {
		return err
	}

	return nil
}

func (s *AliCloudTableStore) delete(req *bindings.InvokeRequest, resp *bindings.InvokeResponse) error {
	pkNams := strings.Split(req.Metadata[primaryKeys], ",")
	pks := make([]*tablestore.PrimaryKeyColumn, len(pkNams))
	data := make(map[string]interface{})
	err := json.Unmarshal(req.Data, &data)
	if err != nil {
		return err
	}

	for idx, pkName := range pkNams {
		pks[idx] = &tablestore.PrimaryKeyColumn{
			ColumnName: pkName,
			Value:      data[pkName],
		}
	}

	change := &tablestore.DeleteRowChange{
		TableName:  s.getTableName(req.Metadata),
		PrimaryKey: &tablestore.PrimaryKey{PrimaryKeys: pks},
	}
	change.SetCondition(tablestore.RowExistenceExpectation_IGNORE)
	deleteReq := &tablestore.DeleteRowRequest{DeleteRowChange: change}
	_, err = s.client.DeleteRow(deleteReq)

	if err != nil {
		return err
	}

	return nil
}

func (s *AliCloudTableStore) unmarshal(pks []*tablestore.PrimaryKeyColumn, columns []*tablestore.AttributeColumn) (map[string]interface{}, error) {
	if pks == nil && columns == nil {
		return nil, nil
	}

	data := make(map[string]interface{})

	for _, pk := range pks {
		data[pk.ColumnName] = pk.Value
	}

	for _, column := range columns {
		data[column.ColumnName] = column.Value
	}

	return data, nil
}

func (s *AliCloudTableStore) getTableName(metadata map[string]string) string {
	name := metadata[tableName]
	if name == "" {
		name = s.metadata.TableName
	}

	return name
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}

	return false
}
