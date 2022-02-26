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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/utils/logger"
)

func TestParseMetadata(t *testing.T) {
	m := bindings.Metadata{}
	m.Properties = map[string]string{"Collection": "a", "Database": "a", "MasterKey": "a", "PartitionKey": "a", "URL": "a"}
	cosmosDB := CosmosDB{logger: logger.NewLogger("test")}
	meta, err := cosmosDB.parseMetadata(m)
	assert.Nil(t, err)
	assert.Equal(t, "a", meta.Collection)
	assert.Equal(t, "a", meta.Database)
	assert.Equal(t, "a", meta.MasterKey)
	assert.Equal(t, "a", meta.PartitionKey)
	assert.Equal(t, "a", meta.URL)
}

func TestPartitionKeyValue(t *testing.T) {
	m := bindings.Metadata{}
	m.Properties = map[string]string{"Collection": "a", "Database": "a", "MasterKey": "a", "PartitionKey": "a", "URL": "a"}
	cosmosDB := CosmosDB{logger: logger.NewLogger("test")}
	var obj interface{}
	jsonStr := `{"name": "name", "empty" : "", "address": { "planet" : { "name": "earth" }, "zip" : "zipcode" }}`
	json.Unmarshal([]byte(jsonStr), &obj)

	// Valid single partition key
	val, err := cosmosDB.getPartitionKeyValue("name", obj)
	assert.Nil(t, err)
	assert.Equal(t, "name", val)

	// Not existing key
	_, err = cosmosDB.getPartitionKeyValue("notexists", obj)
	assert.NotNil(t, err)

	// // Empty value for the key
	_, err = cosmosDB.getPartitionKeyValue("empty", obj)
	assert.NotNil(t, err)

	// Valid nested partition key
	val, err = cosmosDB.getPartitionKeyValue("address.zip", obj)
	assert.Nil(t, err)
	assert.Equal(t, "zipcode", val)

	// Valid nested three level partition key
	val, err = cosmosDB.getPartitionKeyValue("address.planet.name", obj)
	assert.Nil(t, err)
	assert.Equal(t, "earth", val)

	// Invalid nested partition key
	_, err = cosmosDB.getPartitionKeyValue("address.notexists", obj)
	assert.NotNil(t, err)

	// Empty key is passed
	_, err = cosmosDB.getPartitionKeyValue("", obj)
	assert.NotNil(t, err)
}
