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
	"fmt"
	"strconv"
	"time"

	"github.com/bhojpur/service/pkg/utils/config"
)

type Settings struct {
	// The Redis host
	Host string `mapstructure:"redisHost"`
	// The Redis password
	Password string `mapstructure:"redisPassword"`
	// The Redis username
	Username string `mapstructure:"redisUsername"`
	// Database to be selected after connecting to the server.
	DB int `mapstructure:"redisDB"`
	// The redis type node or cluster
	RedisType string `mapstructure:"redisType"`
	// Maximum number of retries before giving up.
	// A value of -1 (not 0) disables retries
	// Default is 3 retries
	RedisMaxRetries int `mapstructure:"redisMaxRetries"`
	// Minimum backoff between each retry.
	// Default is 8 milliseconds; -1 disables backoff.
	RedisMinRetryInterval Duration `mapstructure:"redisMinRetryInterval"`
	// Maximum backoff between each retry.
	// Default is 512 milliseconds; -1 disables backoff.
	RedisMaxRetryInterval Duration `mapstructure:"redisMaxRetryInterval"`
	// Dial timeout for establishing new connections.
	DialTimeout Duration `mapstructure:"dialTimeout"`
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Use value -1 for no timeout and 0 for default.
	ReadTimeout Duration `mapstructure:"readTimeout"`
	// Timeout for socket writes. If reached, commands will fail
	WriteTimeout Duration `mapstructure:"writeTimeout"`
	// Maximum number of socket connections.
	PoolSize int `mapstructure:"poolSize"`
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	MinIdleConns int `mapstructure:"minIdleConns"`
	// Connection age at which client retires (closes) the connection.
	// Default is to not close aged connections.
	MaxConnAge Duration `mapstructure:"maxConnAge"`
	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout Duration `mapstructure:"poolTimeout"`
	// Amount of time after which client closes idle connections.
	// Should be less than server's timeout.
	// Default is 5 minutes. -1 disables idle timeout check.
	IdleTimeout Duration `mapstructure:"idleTimeout"`
	// Frequency of idle checks made by idle connections reaper.
	// Default is 1 minute. -1 disables idle connections reaper,
	// but idle connections are still discarded by the client
	// if IdleTimeout is set.
	IdleCheckFrequency Duration `mapstructure:"idleCheckFrequency"`
	// The master name
	SentinelMasterName string `mapstructure:"sentinelMasterName"`
	// Use Redis Sentinel for automatic failover.
	Failover bool `mapstructure:"failover"`

	// A flag to enables TLS by setting InsecureSkipVerify to true
	EnableTLS bool `mapstructure:"enableTLS"`
}

func (s *Settings) Decode(in interface{}) error {
	if err := config.Decode(in, s); err != nil {
		return fmt.Errorf("decode failed. %w", err)
	}

	return nil
}

type Duration time.Duration

func (r *Duration) DecodeString(value string) error {
	if val, err := strconv.Atoi(value); err == nil {
		if val < 0 {
			*r = Duration(val)

			return nil
		}
		*r = Duration(time.Duration(val) * time.Millisecond)

		return nil
	}

	// Convert it by parsing
	d, err := time.ParseDuration(value)
	if err != nil {
		return err
	}

	*r = Duration(d)

	return nil
}
