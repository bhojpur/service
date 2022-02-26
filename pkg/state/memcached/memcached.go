package memcached

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
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	jsoniter "github.com/json-iterator/go"

	"github.com/bhojpur/service/pkg/state"
	"github.com/bhojpur/service/pkg/state/utils"
	"github.com/bhojpur/service/pkg/utils/logger"
)

const (
	hosts              = "hosts"
	maxIdleConnections = "maxIdleConnections"
	timeout            = "timeout"
	ttlInSeconds       = "ttlInSeconds"
	// These defaults are already provided by gomemcache.
	defaultMaxIdleConnections = 2
	defaultTimeout            = 1000 * time.Millisecond
)

type Memcached struct {
	state.DefaultBulkStore
	client *memcache.Client
	json   jsoniter.API
	logger logger.Logger
}

type memcachedMetadata struct {
	hosts              []string
	maxIdleConnections int
	timeout            time.Duration
}

func NewMemCacheStateStore(logger logger.Logger) *Memcached {
	s := &Memcached{
		json:   jsoniter.ConfigFastest,
		logger: logger,
	}
	s.DefaultBulkStore = state.NewDefaultBulkStore(s)

	return s
}

func (m *Memcached) Init(metadata state.Metadata) error {
	meta, err := getMemcachedMetadata(metadata)
	if err != nil {
		return err
	}

	client := memcache.New(meta.hosts...)
	client.Timeout = meta.timeout
	client.MaxIdleConns = meta.maxIdleConnections

	m.client = client

	err = client.Ping()
	if err != nil {
		return err
	}

	return nil
}

// Features returns the features available in this state store.
func (m *Memcached) Features() []state.Feature {
	return nil
}

func getMemcachedMetadata(metadata state.Metadata) (*memcachedMetadata, error) {
	meta := memcachedMetadata{
		maxIdleConnections: defaultMaxIdleConnections,
		timeout:            defaultTimeout,
	}

	if val, ok := metadata.Properties[hosts]; ok && val != "" {
		meta.hosts = strings.Split(val, ",")
	} else {
		return nil, errors.New("missing or empty hosts field from metadata")
	}

	if val, ok := metadata.Properties[maxIdleConnections]; ok && val != "" {
		p, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("error parsing maxIdleConnections")
		}
		meta.maxIdleConnections = p
	}

	if val, ok := metadata.Properties[timeout]; ok && val != "" {
		p, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("error parsing timeout")
		}
		meta.timeout = time.Duration(p) * time.Millisecond
	}

	return &meta, nil
}

func (m *Memcached) parseTTL(req *state.SetRequest) (*int32, error) {
	if val, ok := req.Metadata[ttlInSeconds]; ok && val != "" {
		parsedVal, err := strconv.ParseInt(val, 10, 0)
		if err != nil {
			return nil, err
		}
		parsedInt := int32(parsedVal)

		return &parsedInt, nil
	}

	return nil, nil
}

func (m *Memcached) setValue(req *state.SetRequest) error {
	var bt []byte
	ttl, err := m.parseTTL(req)
	if err != nil {
		return fmt.Errorf("failed to parse ttl %s: %s", req.Key, err)
	}

	bt, _ = utils.Marshal(req.Value, m.json.Marshal)
	if ttl != nil {
		err = m.client.Set(&memcache.Item{Key: req.Key, Value: bt, Expiration: *ttl})
	} else {
		err = m.client.Set(&memcache.Item{Key: req.Key, Value: bt})
	}
	if err != nil {
		return fmt.Errorf("failed to set key %s: %s", req.Key, err)
	}

	return nil
}

func (m *Memcached) Delete(req *state.DeleteRequest) error {
	err := m.client.Delete(req.Key)
	if err != nil {
		return err
	}

	return nil
}

func (m *Memcached) Get(req *state.GetRequest) (*state.GetResponse, error) {
	item, err := m.client.Get(req.Key)
	if err != nil {
		// Return nil for status 204
		if errors.Is(err, memcache.ErrCacheMiss) {
			return &state.GetResponse{}, nil
		}

		return &state.GetResponse{}, err
	}

	return &state.GetResponse{
		Data: item.Value,
	}, nil
}

func (m *Memcached) Set(req *state.SetRequest) error {
	return state.SetWithOptions(m.setValue, req)
}

func (m *Memcached) Ping() error {
	return nil
}
