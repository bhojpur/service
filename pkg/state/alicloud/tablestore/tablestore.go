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

	"encoding/json"

	"github.com/agrea/ptr"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	jsoniter "github.com/json-iterator/go"

	"github.com/bhojpur/service/pkg/state"
	"github.com/bhojpur/service/pkg/utils/logger"
)

const (
	stateKey   = "stateKey"
	stateValue = "stateValue"
	sateEtag   = "sateEtag"
)

type AliCloudTableStore struct {
	logger   logger.Logger
	client   tablestore.TableStoreApi
	metadata tablestoreMetadata
	features []state.Feature
}

type tablestoreMetadata struct {
	Endpoint     string `json:"endpoint"`
	AccessKeyID  string `json:"accessKeyID"`
	AccessKey    string `json:"accessKey"`
	InstanceName string `json:"instanceName"`
	TableName    string `json:"tableName"`
}

func NewAliCloudTableStore(logger logger.Logger) *AliCloudTableStore {
	return &AliCloudTableStore{
		features: []state.Feature{state.FeatureETag, state.FeatureTransactional},
		logger:   logger,
	}
}

func (s *AliCloudTableStore) Init(metadata state.Metadata) error {
	m, err := s.parse(metadata)
	if err != nil {
		return err
	}

	s.metadata = *m
	s.client = tablestore.NewClient(m.Endpoint, m.InstanceName, m.AccessKeyID, m.AccessKey)

	return nil
}

func (s *AliCloudTableStore) Features() []state.Feature {
	return s.features
}

func (s *AliCloudTableStore) Get(req *state.GetRequest) (*state.GetResponse, error) {
	criteria := &tablestore.SingleRowQueryCriteria{
		PrimaryKey: s.primaryKey(req.Key),
		TableName:  s.metadata.TableName,
		MaxVersion: 1,
	}

	rowGetReq := &tablestore.GetRowRequest{
		SingleRowQueryCriteria: criteria,
	}

	resp, err := s.client.GetRow(rowGetReq)
	if err != nil {
		return nil, err
	}

	getResp := s.getResp(resp.Columns)

	return getResp, nil
}

func (s *AliCloudTableStore) getResp(columns []*tablestore.AttributeColumn) *state.GetResponse {
	getResp := &state.GetResponse{}

	for _, column := range columns {
		if column.ColumnName == stateValue {
			getResp.Data = unmarshal(column.Value)
		} else if column.ColumnName == sateEtag {
			getResp.ETag = ptr.String(column.Value.(string))
		}
	}

	return getResp
}

func (s *AliCloudTableStore) BulkGet(reqs []state.GetRequest) (bool, []state.BulkGetResponse, error) {
	// "len == 0": empty request, directly return empty response
	if len(reqs) == 0 {
		return true, []state.BulkGetResponse{}, nil
	}

	mqCriteria := &tablestore.MultiRowQueryCriteria{
		TableName:  s.metadata.TableName,
		MaxVersion: 1,
	}

	for _, req := range reqs {
		mqCriteria.AddRow(s.primaryKey(req.Key))
	}

	batchGetReq := &tablestore.BatchGetRowRequest{}
	batchGetReq.MultiRowQueryCriteria = append(batchGetReq.MultiRowQueryCriteria, mqCriteria)
	batchGetResp, err := s.client.BatchGetRow(batchGetReq)
	responseList := make([]state.BulkGetResponse, 0, 10)
	if err != nil {
		return false, nil, err
	}

	for _, row := range batchGetResp.TableToRowsResult[mqCriteria.TableName] {
		resp := s.getResp(row.Columns)

		responseList = append(responseList, state.BulkGetResponse{
			Data: resp.Data,
			ETag: resp.ETag,
			Key:  row.PrimaryKey.PrimaryKeys[0].Value.(string),
		})
	}

	return true, responseList, nil
}

func (s *AliCloudTableStore) Set(req *state.SetRequest) error {
	change := s.updateRowChange(req)

	request := &tablestore.UpdateRowRequest{
		UpdateRowChange: change,
	}

	_, err := s.client.UpdateRow(request)

	return err
}

func (s *AliCloudTableStore) updateRowChange(req *state.SetRequest) *tablestore.UpdateRowChange {
	change := &tablestore.UpdateRowChange{
		PrimaryKey: s.primaryKey(req.Key),
		TableName:  s.metadata.TableName,
	}

	value, _ := marshal(req.Value)
	change.PutColumn(stateValue, value)

	if req.ETag != nil {
		change.PutColumn(sateEtag, *req.ETag)
	}

	change.SetCondition(tablestore.RowExistenceExpectation_IGNORE)

	return change
}

func marshal(value interface{}) ([]byte, error) {
	v, _ := jsoniter.MarshalToString(value)

	return []byte(v), nil
}

func unmarshal(val interface{}) []byte {
	var output string

	jsoniter.UnmarshalFromString(string(val.([]byte)), &output)

	return []byte(output)
}

func (s *AliCloudTableStore) Delete(req *state.DeleteRequest) error {
	change := s.deleteRowChange(req)

	deleteRowReq := &tablestore.DeleteRowRequest{
		DeleteRowChange: change,
	}

	_, err := s.client.DeleteRow(deleteRowReq)

	return err
}

func (s *AliCloudTableStore) deleteRowChange(req *state.DeleteRequest) *tablestore.DeleteRowChange {
	change := &tablestore.DeleteRowChange{
		PrimaryKey: s.primaryKey(req.Key),
		TableName:  s.metadata.TableName,
	}
	change.SetCondition(tablestore.RowExistenceExpectation_EXPECT_EXIST)

	return change
}

func (s *AliCloudTableStore) BulkSet(reqs []state.SetRequest) error {
	return s.batchWrite(reqs, nil)
}

func (s *AliCloudTableStore) BulkDelete(reqs []state.DeleteRequest) error {
	return s.batchWrite(nil, reqs)
}

func (s *AliCloudTableStore) batchWrite(setReqs []state.SetRequest, deleteReqs []state.DeleteRequest) error {
	bathReq := &tablestore.BatchWriteRowRequest{
		IsAtomic: true,
	}

	for i := range setReqs {
		bathReq.AddRowChange(s.updateRowChange(&setReqs[i]))
	}

	for i := range deleteReqs {
		bathReq.AddRowChange(s.deleteRowChange(&deleteReqs[i]))
	}

	_, err := s.client.BatchWriteRow(bathReq)
	if err != nil {
		return err
	}

	return nil
}

func (s *AliCloudTableStore) Ping() error {
	return nil
}

func (s *AliCloudTableStore) parse(metadata state.Metadata) (*tablestoreMetadata, error) {
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

func (s *AliCloudTableStore) primaryKey(key string) *tablestore.PrimaryKey {
	pk := &tablestore.PrimaryKey{}
	pk.AddPrimaryKeyColumn(stateKey, key)

	return pk
}
