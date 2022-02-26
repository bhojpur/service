package redis

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
	"reflect"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/configuration"
	"github.com/bhojpur/service/pkg/utils/logger"
)

func TestConfigurationStore_Get(t *testing.T) {
	s, c := setupMiniredis()
	defer s.Close()
	assert.Nil(t, s.Set("testKey", "testValue"))
	assert.Nil(t, s.Set("testKey2", "testValue2"))

	type fields struct {
		client   *redis.Client
		json     jsoniter.API
		metadata metadata
		replicas int
		logger   logger.Logger
	}
	type args struct {
		ctx context.Context
		req *configuration.GetRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *configuration.GetResponse
		wantErr bool
	}{
		{
			name: "normal get redis value",
			fields: fields{
				client: c,
				json:   jsoniter.ConfigFastest,
				logger: logger.NewLogger("test"),
			},
			args: args{
				req: &configuration.GetRequest{
					Keys: []string{"testKey"},
				},
				ctx: context.Background(),
			},
			want: &configuration.GetResponse{
				Items: []*configuration.Item{
					{
						Key:      "testKey",
						Value:    "testValue",
						Metadata: make(map[string]string),
					},
				},
			},
		},
		{
			name: "get with no request key",
			fields: fields{
				client: c,
				json:   jsoniter.ConfigFastest,
				logger: logger.NewLogger("test"),
			},
			args: args{
				req: &configuration.GetRequest{},
				ctx: context.Background(),
			},
			want: &configuration.GetResponse{
				Items: []*configuration.Item{
					{
						Key:      "testKey",
						Value:    "testValue",
						Metadata: make(map[string]string),
					}, {
						Key:      "testKey2",
						Value:    "testValue2",
						Metadata: make(map[string]string),
					},
				},
			},
		},
		{
			name: "get with not exists key",
			fields: fields{
				client: c,
				json:   jsoniter.ConfigFastest,
				logger: logger.NewLogger("test"),
			},
			args: args{
				req: &configuration.GetRequest{
					Keys: []string{"notExistKey"},
				},
				ctx: context.Background(),
			},
			want: &configuration.GetResponse{
				Items: []*configuration.Item{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ConfigurationStore{
				client:   tt.fields.client,
				json:     tt.fields.json,
				metadata: tt.fields.metadata,
				replicas: tt.fields.replicas,
				logger:   tt.fields.logger,
			}
			got, err := r.Get(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				t.Errorf("Get() got configuration response is nil")
				return
			}

			if len(got.Items) != len(tt.want.Items) {
				t.Errorf("Get() got len = %v, want len = %v", len(got.Items), len(tt.want.Items))
				return
			}

			if len(got.Items) == 0 {
				return
			}

			for k := range got.Items {
				assert.Equal(t, tt.want.Items[k], got.Items[k])
			}
		})
	}
}

func TestParseConnectedSlaves(t *testing.T) {
	store := &ConfigurationStore{logger: logger.NewLogger("test")}

	t.Run("Empty info", func(t *testing.T) {
		slaves := store.parseConnectedSlaves("")
		assert.Equal(t, 0, slaves, "connected slaves must be 0")
	})

	t.Run("connectedSlaves property is not included", func(t *testing.T) {
		slaves := store.parseConnectedSlaves("# Replication\r\nrole:master\r\n")
		assert.Equal(t, 0, slaves, "connected slaves must be 0")
	})

	t.Run("connectedSlaves is 2", func(t *testing.T) {
		slaves := store.parseConnectedSlaves("# Replication\r\nrole:master\r\nconnected_slaves:2\r\n")
		assert.Equal(t, 2, slaves, "connected slaves must be 2")
	})

	t.Run("connectedSlaves is 1", func(t *testing.T) {
		slaves := store.parseConnectedSlaves("# Replication\r\nrole:master\r\nconnected_slaves:1")
		assert.Equal(t, 1, slaves, "connected slaves must be 1")
	})
}

func TestNewRedisConfigurationStore(t *testing.T) {
	type args struct {
		logger logger.Logger
	}
	tests := []struct {
		name string
		args args
		want configuration.Store
	}{
		{
			args: args{
				logger: logger.NewLogger("test"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRedisConfigurationStore(tt.args.logger)
			assert.NotNil(t, got)
		})
	}
}

func Test_parseRedisMetadata(t *testing.T) {
	type args struct {
		meta configuration.Metadata
	}
	testProperties := make(map[string]string)
	testProperties[host] = "testHost"
	testProperties[password] = "testPassword"
	testProperties[enableTLS] = "true"
	testProperties[maxRetries] = "10"
	testProperties[maxRetryBackoff] = "1000000000"
	testProperties[failover] = "true"
	testProperties[sentinelMasterName] = "tesSentinelMasterName"
	tests := []struct {
		name    string
		args    args
		want    metadata
		wantErr bool
	}{
		{
			args: args{
				meta: configuration.Metadata{
					Properties: testProperties,
				},
			},
			want: metadata{
				host:               "testHost",
				password:           "testPassword",
				enableTLS:          true,
				maxRetries:         10,
				maxRetryBackoff:    time.Second,
				failover:           true,
				sentinelMasterName: "tesSentinelMasterName",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRedisMetadata(tt.args.meta)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRedisMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseRedisMetadata() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func setupMiniredis() (*miniredis.Miniredis, *redis.Client) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	opts := &redis.Options{
		Addr: s.Addr(),
		DB:   defaultDB,
	}

	return s, redis.NewClient(opts)
}
