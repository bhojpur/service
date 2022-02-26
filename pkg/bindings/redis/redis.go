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
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"

	"github.com/bhojpur/service/pkg/bindings"
	rediscomponent "github.com/bhojpur/service/pkg/internal/component/redis"
	"github.com/bhojpur/service/pkg/utils/logger"
)

// Redis is a redis output binding.
type Redis struct {
	client         redis.UniversalClient
	clientSettings *rediscomponent.Settings
	logger         logger.Logger

	ctx    context.Context
	cancel context.CancelFunc
}

// NewRedis returns a new redis bindings instance.
func NewRedis(logger logger.Logger) *Redis {
	return &Redis{logger: logger}
}

// Init performs metadata parsing and connection creation.
func (r *Redis) Init(meta bindings.Metadata) (err error) {
	r.client, r.clientSettings, err = rediscomponent.ParseClientFromProperties(meta.Properties, nil)
	if err != nil {
		return err
	}

	r.ctx, r.cancel = context.WithCancel(context.Background())

	_, err = r.client.Ping(r.ctx).Result()
	if err != nil {
		return fmt.Errorf("redis binding: error connecting to redis at %s: %s", r.clientSettings.Host, err)
	}

	return err
}

func (r *Redis) Operations() []bindings.OperationKind {
	return []bindings.OperationKind{bindings.CreateOperation}
}

func (r *Redis) Invoke(req *bindings.InvokeRequest) (*bindings.InvokeResponse, error) {
	if val, ok := req.Metadata["key"]; ok && val != "" {
		key := val
		_, err := r.client.Do(r.ctx, "SET", key, req.Data).Result()
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	return nil, errors.New("redis binding: missing key on write request metadata")
}

func (r *Redis) Close() error {
	r.cancel()

	return r.client.Close()
}
