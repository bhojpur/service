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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/utils/logger"
)

func TestTableStoreMetadata(t *testing.T) {
	m := bindings.Metadata{}
	m.Properties = map[string]string{"accessKeyID": "ACCESSKEYID", "accessKey": "ACCESSKEY", "instanceName": "INSTANCENAME", "tableName": "TABLENAME", "endpoint": "ENDPOINT"}
	aliCloudTableStore := AliCloudTableStore{}

	meta, err := aliCloudTableStore.parseMetadata(m)

	assert.Nil(t, err)
	assert.Equal(t, "ACCESSKEYID", meta.AccessKeyID)
	assert.Equal(t, "ACCESSKEY", meta.AccessKey)
	assert.Equal(t, "INSTANCENAME", meta.InstanceName)
	assert.Equal(t, "TABLENAME", meta.TableName)
	assert.Equal(t, "ENDPOINT", meta.Endpoint)
}

func TestDataEncodeAndDecode(t *testing.T) {
	if !isLiveTest() {
		return
	}

	aliCloudTableStore := NewAliCloudTableStore(logger.NewLogger("test"))

	metadata := bindings.Metadata{
		Properties: getTestProperties(),
	}
	aliCloudTableStore.Init(metadata)

	// test create
	putData := map[string]interface{}{
		"pk1":     "data1",
		"column1": "the string value of column1",
		"column2": int64(2),
	}
	data, err := json.Marshal(putData)
	assert.Nil(t, err)
	putRowReq := &bindings.InvokeRequest{
		Operation: bindings.CreateOperation,
		Metadata: map[string]string{
			tableName:   "app_test_table2",
			primaryKeys: "pk1",
		},
		Data: data,
	}

	putInvokeResp, err := aliCloudTableStore.Invoke(putRowReq)

	assert.Nil(t, err)
	assert.NotNil(t, putInvokeResp)

	putRowReq.Data, _ = json.Marshal(map[string]interface{}{
		"pk1":     "data2",
		"column1": "the string value of column1",
		"column2": int64(2),
	})

	putInvokeResp, err = aliCloudTableStore.Invoke(putRowReq)

	assert.Nil(t, err)
	assert.NotNil(t, putInvokeResp)

	// test get
	getData, err := json.Marshal(map[string]interface{}{
		"pk1": "data1",
	})
	assert.Nil(t, err)
	getInvokeReq := &bindings.InvokeRequest{
		Operation: bindings.GetOperation,
		Metadata: map[string]string{
			tableName:   "app_test_table2",
			primaryKeys: "pk1",
			columnToGet: "column1,column2,column3",
		},
		Data: getData,
	}

	getInvokeResp, err := aliCloudTableStore.Invoke(getInvokeReq)

	assert.Nil(t, err)
	assert.NotNil(t, getInvokeResp)

	respData := make(map[string]interface{})
	err = json.Unmarshal(getInvokeResp.Data, &respData)

	assert.Nil(t, err)

	assert.Equal(t, putData["column1"], respData["column1"])
	assert.Equal(t, putData["column2"], int64(respData["column2"].(float64)))

	// test list
	listData, err := json.Marshal([]map[string]interface{}{
		{
			"pk1": "data1",
		},
		{
			"pk1": "data2",
		},
	})
	assert.Nil(t, err)

	listReq := &bindings.InvokeRequest{
		Operation: bindings.ListOperation,
		Metadata: map[string]string{
			tableName:   "app_test_table2",
			primaryKeys: "pk1",
			columnToGet: "column1,column2,column3",
		},
		Data: listData,
	}

	listResp, err := aliCloudTableStore.Invoke(listReq)
	assert.Nil(t, err)
	assert.NotNil(t, listResp)

	listRespData := make([]map[string]interface{}, len(listData))
	err = json.Unmarshal(listResp.Data, &listRespData)

	assert.Nil(t, err)
	assert.Len(t, listRespData, 2)

	assert.Equal(t, listRespData[0]["column1"], putData["column1"])
	assert.Equal(t, listRespData[1]["pk1"], "data2")

	// test delete
	deleteData, err := json.Marshal(map[string]interface{}{
		"pk1": "data1",
	})
	assert.Nil(t, err)

	deleteReq := &bindings.InvokeRequest{
		Operation: bindings.DeleteOperation,
		Metadata: map[string]string{
			tableName:   "app_test_table2",
			primaryKeys: "pk1",
		},
		Data: deleteData,
	}

	deleteResp, err := aliCloudTableStore.Invoke(deleteReq)

	assert.Nil(t, err)
	assert.NotNil(t, deleteResp)

	getInvokeResp, err = aliCloudTableStore.Invoke(getInvokeReq)

	assert.Nil(t, err)
	assert.Nil(t, getInvokeResp.Data)
}

func getTestProperties() map[string]string {
	return map[string]string{
		"accessKeyID":  "****",
		"accessKey":    "****",
		"instanceName": "app-test",
		"tableName":    "app_test_table2",
		"endpoint":     "https://app-test.cn-hangzhou.ots.aliyuncs.com",
	}
}

func isLiveTest() bool {
	return os.Getenv("RUN_LIVE_ROCKETMQ_TEST") == "true"
}
