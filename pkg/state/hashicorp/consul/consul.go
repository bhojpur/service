package consul

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

	"github.com/agrea/ptr"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"

	"github.com/bhojpur/service/pkg/state"
	"github.com/bhojpur/service/pkg/utils/logger"
)

// Consul is a state store implementation for HashiCorp Consul.
type Consul struct {
	state.DefaultBulkStore
	client        *api.Client
	keyPrefixPath string
	logger        logger.Logger
}

type consulConfig struct {
	Datacenter    string `json:"datacenter"`
	HTTPAddr      string `json:"httpAddr"`
	ACLToken      string `json:"aclToken"`
	Scheme        string `json:"scheme"`
	KeyPrefixPath string `json:"keyPrefixPath"`
}

// NewConsulStateStore returns a new consul state store.
func NewConsulStateStore(logger logger.Logger) *Consul {
	s := &Consul{logger: logger}
	s.DefaultBulkStore = state.NewDefaultBulkStore(s)

	return s
}

// Init does metadata and config parsing and initializes the
// Consul client.
func (c *Consul) Init(metadata state.Metadata) error {
	consulConfig, err := metadataToConfig(metadata.Properties)
	if err != nil {
		return fmt.Errorf("couldn't convert metadata properties: %s", err)
	}

	var keyPrefixPath string
	if consulConfig.KeyPrefixPath == "" {
		keyPrefixPath = "app"
	}

	config := &api.Config{
		Datacenter: consulConfig.Datacenter,
		Address:    consulConfig.HTTPAddr,
		Token:      consulConfig.ACLToken,
		Scheme:     consulConfig.Scheme,
	}

	client, err := api.NewClient(config)
	if err != nil {
		return errors.Wrap(err, "initializing consul client")
	}

	c.client = client
	c.keyPrefixPath = keyPrefixPath

	return nil
}

// Features returns the features available in this state store.
func (c *Consul) Features() []state.Feature {
	// Etag is just returned and not handled in set or delete operations.
	return nil
}

func metadataToConfig(connInfo map[string]string) (*consulConfig, error) {
	b, err := json.Marshal(connInfo)
	if err != nil {
		return nil, err
	}

	var config consulConfig
	err = json.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// Get retrieves a Consul KV item.
func (c *Consul) Get(req *state.GetRequest) (*state.GetResponse, error) {
	queryOpts := &api.QueryOptions{}
	if req.Options.Consistency == state.Strong {
		queryOpts.RequireConsistent = true
	}

	resp, queryMeta, err := c.client.KV().Get(fmt.Sprintf("%s/%s", c.keyPrefixPath, req.Key), queryOpts)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return &state.GetResponse{}, nil
	}

	return &state.GetResponse{
		Data: resp.Value,
		ETag: ptr.String(queryMeta.LastContentHash),
	}, nil
}

// Set saves a Consul KV item.
func (c *Consul) Set(req *state.SetRequest) error {
	var reqValByte []byte
	b, ok := req.Value.([]byte)
	if ok {
		reqValByte = b
	} else {
		reqValByte, _ = json.Marshal(req.Value)
	}

	keyWithPath := fmt.Sprintf("%s/%s", c.keyPrefixPath, req.Key)

	_, err := c.client.KV().Put(&api.KVPair{
		Key:   keyWithPath,
		Value: reqValByte,
	}, nil)
	if err != nil {
		return fmt.Errorf("couldn't set key %s: %s", keyWithPath, err)
	}

	return nil
}

func (c *Consul) Ping() error {
	return nil
}

// Delete performes a Consul KV delete operation.
func (c *Consul) Delete(req *state.DeleteRequest) error {
	keyWithPath := fmt.Sprintf("%s/%s", c.keyPrefixPath, req.Key)
	_, err := c.client.KV().Delete(keyWithPath, nil)
	if err != nil {
		return fmt.Errorf("couldn't delete key %s: %s", keyWithPath, err)
	}

	return nil
}
