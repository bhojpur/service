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
	"crypto/tls"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"

	"github.com/bhojpur/service/pkg/configuration"
	"github.com/bhojpur/service/pkg/configuration/redis/internal"
	"github.com/bhojpur/service/pkg/utils/logger"
)

const (
	connectedSlavesReplicas  = "connected_slaves:"
	infoReplicationDelimiter = "\r\n"
	host                     = "redisHost"
	password                 = "redisPassword"
	enableTLS                = "enableTLS"
	maxRetries               = "maxRetries"
	maxRetryBackoff          = "maxRetryBackoff"
	failover                 = "failover"
	sentinelMasterName       = "sentinelMasterName"
	defaultBase              = 10
	defaultBitSize           = 0
	defaultDB                = 0
	defaultMaxRetries        = 3
	defaultMaxRetryBackoff   = time.Second * 2
	defaultEnableTLS         = false
	keySpacePrefix           = "__keyspace@0__:"
	keySpaceAny              = "__keyspace@0__:*"
)

// ConfigurationStore is a Redis configuration store.
type ConfigurationStore struct {
	client               redis.UniversalClient
	json                 jsoniter.API
	metadata             metadata
	replicas             int
	subscribeStopChanMap sync.Map

	logger logger.Logger
}

// NewRedisConfigurationStore returns a new redis state store.
func NewRedisConfigurationStore(logger logger.Logger) configuration.Store {
	s := &ConfigurationStore{
		json:   jsoniter.ConfigFastest,
		logger: logger,
	}

	return s
}

func parseRedisMetadata(meta configuration.Metadata) (metadata, error) {
	m := metadata{}

	if val, ok := meta.Properties[host]; ok && val != "" {
		m.host = val
	} else {
		return m, errors.New("redis store error: missing host address")
	}

	if val, ok := meta.Properties[password]; ok && val != "" {
		m.password = val
	}

	m.enableTLS = defaultEnableTLS
	if val, ok := meta.Properties[enableTLS]; ok && val != "" {
		tls, err := strconv.ParseBool(val)
		if err != nil {
			return m, fmt.Errorf("redis store error: can't parse enableTLS field: %s", err)
		}
		m.enableTLS = tls
	}

	m.maxRetries = defaultMaxRetries
	if val, ok := meta.Properties[maxRetries]; ok && val != "" {
		parsedVal, err := strconv.ParseInt(val, defaultBase, defaultBitSize)
		if err != nil {
			return m, fmt.Errorf("redis store error: can't parse maxRetries field: %s", err)
		}
		m.maxRetries = int(parsedVal)
	}

	m.maxRetryBackoff = defaultMaxRetryBackoff
	if val, ok := meta.Properties[maxRetryBackoff]; ok && val != "" {
		parsedVal, err := strconv.ParseInt(val, defaultBase, defaultBitSize)
		if err != nil {
			return m, fmt.Errorf("redis store error: can't parse maxRetryBackoff field: %s", err)
		}
		m.maxRetryBackoff = time.Duration(parsedVal)
	}

	if val, ok := meta.Properties[failover]; ok && val != "" {
		failover, err := strconv.ParseBool(val)
		if err != nil {
			return m, fmt.Errorf("redis store error: can't parse failover field: %s", err)
		}
		m.failover = failover
	}

	// set the sentinelMasterName only with failover == true.
	if m.failover {
		if val, ok := meta.Properties[sentinelMasterName]; ok && val != "" {
			m.sentinelMasterName = val
		} else {
			return m, errors.New("redis store error: missing sentinelMasterName")
		}
	}

	return m, nil
}

// Init does metadata and connection parsing.
func (r *ConfigurationStore) Init(metadata configuration.Metadata) error {
	m, err := parseRedisMetadata(metadata)
	if err != nil {
		return err
	}
	r.metadata = m

	if r.metadata.failover {
		r.client = r.newFailoverClient(m)
	} else {
		r.client = r.newClient(m)
	}

	if _, err = r.client.Ping(context.TODO()).Result(); err != nil {
		return fmt.Errorf("redis store: error connecting to redis at %s: %s", m.host, err)
	}

	r.replicas, err = r.getConnectedSlaves()

	return err
}

func (r *ConfigurationStore) newClient(m metadata) *redis.Client {
	opts := &redis.Options{
		Addr:            m.host,
		Password:        m.password,
		DB:              defaultDB,
		MaxRetries:      m.maxRetries,
		MaxRetryBackoff: m.maxRetryBackoff,
	}

	// tell the linter to skip a check here.
	/* #nosec */
	if m.enableTLS {
		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify: m.enableTLS,
		}
	}

	return redis.NewClient(opts)
}

func (r *ConfigurationStore) newFailoverClient(m metadata) *redis.Client {
	opts := &redis.FailoverOptions{
		MasterName:      r.metadata.sentinelMasterName,
		SentinelAddrs:   []string{r.metadata.host},
		DB:              defaultDB,
		MaxRetries:      m.maxRetries,
		MaxRetryBackoff: m.maxRetryBackoff,
	}

	/* #nosec */
	if m.enableTLS {
		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify: m.enableTLS,
		}
	}

	return redis.NewFailoverClient(opts)
}

func (r *ConfigurationStore) getConnectedSlaves() (int, error) {
	res, err := r.client.Do(context.Background(), "INFO", "replication").Result()
	if err != nil {
		return 0, err
	}

	// Response example: https://redis.io/commands/info#return-value
	// # Replication\r\nrole:master\r\nconnected_slaves:1\r\n
	s, _ := strconv.Unquote(fmt.Sprintf("%q", res))
	if len(s) == 0 {
		return 0, nil
	}

	return r.parseConnectedSlaves(s), nil
}

func (r *ConfigurationStore) parseConnectedSlaves(res string) int {
	infos := strings.Split(res, infoReplicationDelimiter)
	for _, info := range infos {
		if strings.Contains(info, connectedSlavesReplicas) {
			parsedReplicas, _ := strconv.ParseUint(info[len(connectedSlavesReplicas):], 10, 32)
			return int(parsedReplicas)
		}
	}

	return 0
}

func (r *ConfigurationStore) Get(ctx context.Context, req *configuration.GetRequest) (*configuration.GetResponse, error) {
	keys := req.Keys
	var err error
	if len(keys) == 0 {
		if keys, err = r.client.Keys(ctx, "*").Result(); err != nil {
			r.logger.Errorf("failed to all keys, error is %s", err)
		}
	}

	items := make([]*configuration.Item, 0, 16)

	// query by keys
	for _, redisKey := range keys {
		item := &configuration.Item{
			Metadata: map[string]string{},
		}

		redisValue, err := r.client.Get(ctx, redisKey).Result()
		if err != nil {
			return &configuration.GetResponse{}, fmt.Errorf("fail to get configuration for redis key=%s, error is %s", redisKey, err)
		}
		val, version := internal.GetRedisValueAndVersion(redisValue)
		item.Key = redisKey
		item.Version = version
		item.Value = val

		if item.Value != "" {
			items = append(items, item)
		}
	}

	return &configuration.GetResponse{
		Items: items,
	}, nil
}

func (r *ConfigurationStore) Subscribe(ctx context.Context, req *configuration.SubscribeRequest, handler configuration.UpdateHandler) (string, error) {
	subscribeID := uuid.New().String()
	if len(req.Keys) == 0 {
		// subscribe all keys
		stop := make(chan struct{})
		r.subscribeStopChanMap.Store(subscribeID, stop)
		go r.doSubscribe(ctx, req, handler, keySpaceAny, subscribeID, stop)
		return subscribeID, nil
	}
	for _, k := range req.Keys {
		// subscribe single key
		stop := make(chan struct{})
		keySpacePrefixAndKey := keySpacePrefix + k
		if oldStopChan, ok := r.subscribeStopChanMap.Load(keySpacePrefixAndKey); ok {
			// already exist subscription
			close(oldStopChan.(chan struct{}))
		}
		r.subscribeStopChanMap.Store(subscribeID, stop)
		go r.doSubscribe(ctx, req, handler, keySpacePrefixAndKey, subscribeID, stop)
	}
	return subscribeID, nil
}

func (r *ConfigurationStore) Unsubscribe(ctx context.Context, req *configuration.UnsubscribeRequest) error {
	if oldStopChan, ok := r.subscribeStopChanMap.Load(req.ID); ok {
		// already exist subscription
		r.subscribeStopChanMap.Delete(req.ID)
		close(oldStopChan.(chan struct{}))
	}
	return nil
}

func (r *ConfigurationStore) doSubscribe(ctx context.Context, req *configuration.SubscribeRequest, handler configuration.UpdateHandler, redisChannel4revision string, id string, stop chan struct{}) {
	// enable notify-keyspace-events by redis Set command
	r.client.ConfigSet(ctx, "notify-keyspace-events", "KA")
	p := r.client.Subscribe(ctx, redisChannel4revision)
	for {
		select {
		case <-stop:
			return
		case <-ctx.Done():
			return
		case msg := <-p.Channel():
			r.handleSubscribedChange(ctx, req, handler, msg, id)
		}
	}
}

func (r *ConfigurationStore) handleSubscribedChange(ctx context.Context, req *configuration.SubscribeRequest, handler configuration.UpdateHandler, msg *redis.Message, id string) {
	defer func() {
		if err := recover(); err != nil {
			r.logger.Errorf("panic in handleSubscribedChange(）method and recovered: %s", err)
		}
	}()
	targetKey, err := internal.ParseRedisKeyFromEvent(msg.Channel)
	if err != nil {
		r.logger.Errorf("parse redis key failed: %s", err)
		return
	}

	// get all keys if only one is changed
	getResponse, err := r.Get(ctx, &configuration.GetRequest{
		Metadata: req.Metadata,
		Keys:     []string{targetKey},
	})
	if err != nil {
		r.logger.Errorf("get response from redis failed: %s", err)
		return
	}

	e := &configuration.UpdateEvent{
		Items: getResponse.Items,
		ID:    id,
	}
	err = handler(ctx, e)
	if err != nil {
		r.logger.Errorf("fail to call handler to notify event for configuration update subscribe: %s", err)
	}
}
