package hazelcast

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
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/hazelcast/hazelcast-go-client"
	jsoniter "github.com/json-iterator/go"

	"github.com/bhojpur/service/pkg/state"
	"github.com/bhojpur/service/pkg/utils/logger"
)

const (
	hazelcastServers = "hazelcastServers"
	hazelcastMap     = "hazelcastMap"
)

// Hazelcast state store.
type Hazelcast struct {
	state.DefaultBulkStore
	hxCtx  *context.Context
	hzMap  *hazelcast.Map
	json   jsoniter.API
	logger logger.Logger
}

var defaultStore Hazelcast

// NewHazelcastStore returns a new hazelcast backed state store.
func NewHazelcastStore(logger logger.Logger) *Hazelcast {
	s := &Hazelcast{
		json:   jsoniter.ConfigFastest,
		logger: logger,
	}
	s.DefaultBulkStore = state.NewDefaultBulkStore(s)

	return s
}

func validateMetadata(metadata state.Metadata) error {
	if metadata.Properties[hazelcastServers] == "" {
		return errors.New("hazelcast error: missing hazelcast servers")
	}
	if metadata.Properties[hazelcastMap] == "" {
		return errors.New("hazelcast error: missing hazelcast map name")
	}

	return nil
}

// Init does metadata and connection parsing.
func (store *Hazelcast) Init(metadata state.Metadata) error {
	err := validateMetadata(metadata)
	if err != nil {
		return err
	}
	servers := metadata.Properties[hazelcastServers]

	hzConfig := hazelcast.NewConfig()
	hzConfig.Cluster.Network.SetAddresses(strings.Split(servers, ",")...)

	ctx := context.TODO()
	if ctx == nil {
		return fmt.Errorf("hazelcast error: context creation")
	} else {
		store.hxCtx = &ctx
	}
	client, err := hazelcast.StartNewClientWithConfig(*store.hxCtx, hzConfig)
	mapName := fmt.Sprintf("app-%d", rand.Int())
	if err != nil {
		return fmt.Errorf("hazelcast error: %v", err)
	}
	store.hzMap, err = client.GetMap(*store.hxCtx, mapName)
	//store.hzMap, err = client.GetMap(metadata.Properties[hazelcastMap])

	if err != nil {
		return fmt.Errorf("hazelcast error: %v", err)
	}

	return nil
}

// Features returns the features available in this state store.
func (store *Hazelcast) Features() []state.Feature {
	return nil
}

// Set stores value for a key to Hazelcast.
func (store *Hazelcast) Set(req *state.SetRequest) error {
	err := state.CheckRequestOptions(req)
	if err != nil {
		return err
	}

	var value string
	b, ok := req.Value.([]byte)
	if ok {
		value = string(b)
	} else {
		value, err = store.json.MarshalToString(req.Value)
		if err != nil {
			return fmt.Errorf("hazelcast error: failed to set key %s: %s", req.Key, err)
		}
	}
	_, err = store.hzMap.Put(*store.hxCtx, req.Key, value)

	if err != nil {
		return fmt.Errorf("hazelcast error: failed to set key %s: %s", req.Key, err)
	}

	return nil
}

// Get retrieves state from Hazelcast with a key.
func (store *Hazelcast) Get(req *state.GetRequest) (*state.GetResponse, error) {
	resp, err := store.hzMap.Get(*store.hxCtx, req.Key)
	if err != nil {
		return nil, fmt.Errorf("hazelcast error: failed to get value for %s: %s", req.Key, err)
	}

	// HZ Get API returns nil response if key does not exist in the map
	if resp == nil {
		return &state.GetResponse{}, nil
	}
	value, err := store.json.Marshal(&resp)
	if err != nil {
		return nil, fmt.Errorf("hazelcast error: %v", err)
	}

	return &state.GetResponse{
		Data: value,
	}, nil
}

func (store *Hazelcast) Ping() error {
	return nil
}

// Delete performs a delete operation.
func (store *Hazelcast) Delete(req *state.DeleteRequest) error {
	err := state.CheckRequestOptions(req.Options)
	if err != nil {
		return err
	}
	err = store.hzMap.Delete(*store.hxCtx, req.Key)
	if err != nil {
		return fmt.Errorf("hazelcast error: failed to delete key - %s", req.Key)
	}

	return nil
}
