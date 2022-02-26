package cosmosdbgremlinapi

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
	"errors"
	"fmt"
	"time"

	gremcos "github.com/supplyon/gremcos"

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/utils/logger"
)

const (
	queryOperation bindings.OperationKind = "query"

	// keys from request's Data.
	commandGremlinKey = "gremlin"

	// keys from response's Data.
	respGremlinKey   = "gremlin"
	respOpKey        = "operation"
	respStartTimeKey = "start-time"
	respEndTimeKey   = "end-time"
	respDurationKey  = "duration"
)

// CosmosDBGremlinAPI allows performing state operations on collections.
type CosmosDBGremlinAPI struct {
	metadata *cosmosDBGremlinAPICredentials
	client   *gremcos.Cosmos
	logger   logger.Logger
}

type cosmosDBGremlinAPICredentials struct {
	URL       string `json:"url"`
	MasterKey string `json:"masterKey"`
	Username  string `json:"username"`
}

// NewCosmosDBGremlinAPI returns a new CosmosDBGremlinAPI instance.
func NewCosmosDBGremlinAPI(logger logger.Logger) *CosmosDBGremlinAPI {
	return &CosmosDBGremlinAPI{logger: logger}
}

// Init performs CosmosDBGremlinAPI connection parsing and connecting.
func (c *CosmosDBGremlinAPI) Init(metadata bindings.Metadata) error {
	c.logger.Debug("Initializing Cosmos Graph DB binding")

	m, err := c.parseMetadata(metadata)
	if err != nil {
		return err
	}
	c.metadata = m
	client, err := gremcos.New(c.metadata.URL,
		gremcos.WithAuth(c.metadata.Username, c.metadata.MasterKey),
	)
	if err != nil {
		return errors.New("CosmosDBGremlinAPI Error: failed to create the Cosmos Graph DB connector")
	}

	c.client = client

	return nil
}

func (c *CosmosDBGremlinAPI) parseMetadata(metadata bindings.Metadata) (*cosmosDBGremlinAPICredentials, error) {
	b, err := json.Marshal(metadata.Properties)
	if err != nil {
		return nil, err
	}

	var creds cosmosDBGremlinAPICredentials
	err = json.Unmarshal(b, &creds)
	if err != nil {
		return nil, err
	}

	return &creds, nil
}

func (c *CosmosDBGremlinAPI) Operations() []bindings.OperationKind {
	return []bindings.OperationKind{queryOperation}
}

func (c *CosmosDBGremlinAPI) Invoke(req *bindings.InvokeRequest) (*bindings.InvokeResponse, error) {
	var jsonPoint map[string]interface{}
	err := json.Unmarshal(req.Data, &jsonPoint)
	if err != nil {
		return nil, errors.New("CosmosDBGremlinAPI Error: Cannot convert request data")
	}

	gq := fmt.Sprintf("%s", jsonPoint[commandGremlinKey])

	if gq == "" {
		return nil, errors.New("CosmosDBGremlinAPI Error: missing data - gremlin query not set")
	}
	startTime := time.Now().UTC()
	resp := &bindings.InvokeResponse{
		Metadata: map[string]string{
			respOpKey:        string(req.Operation),
			respGremlinKey:   gq,
			respStartTimeKey: startTime.Format(time.RFC3339Nano),
		},
	}
	d, err := c.client.Execute(gq)
	if err != nil {
		return nil, errors.New("CosmosDBGremlinAPI Error:error excuting gremlin")
	}
	if len(d) > 0 {
		resp.Data = d[0].Result.Data
	}
	endTime := time.Now().UTC()
	resp.Metadata[respEndTimeKey] = endTime.Format(time.RFC3339Nano)
	resp.Metadata[respDurationKey] = endTime.Sub(startTime).String()

	return resp, nil
}
