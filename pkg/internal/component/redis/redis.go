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
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	ClusterType = "cluster"
	NodeType    = "node"
)

func ParseClientFromProperties(properties map[string]string, defaultSettings *Settings) (client redis.UniversalClient, settings *Settings, err error) {
	if defaultSettings == nil {
		settings = &Settings{}
	} else {
		settings = defaultSettings
	}
	err = settings.Decode(properties)
	if err != nil {
		return nil, nil, fmt.Errorf("redis client configuration error: %w", err)
	}
	if settings.Failover {
		return newFailoverClient(settings), settings, nil
	}

	return newClient(settings), settings, nil
}

func newFailoverClient(s *Settings) redis.UniversalClient {
	if s == nil {
		return nil
	}
	opts := &redis.FailoverOptions{
		DB:                 s.DB,
		MasterName:         s.SentinelMasterName,
		SentinelAddrs:      []string{s.Host},
		Password:           s.Password,
		Username:           s.Username,
		MaxRetries:         s.RedisMaxRetries,
		MaxRetryBackoff:    time.Duration(s.RedisMaxRetryInterval),
		MinRetryBackoff:    time.Duration(s.RedisMinRetryInterval),
		DialTimeout:        time.Duration(s.DialTimeout),
		ReadTimeout:        time.Duration(s.ReadTimeout),
		WriteTimeout:       time.Duration(s.WriteTimeout),
		PoolSize:           s.PoolSize,
		MaxConnAge:         time.Duration(s.MaxConnAge),
		MinIdleConns:       s.MinIdleConns,
		PoolTimeout:        time.Duration(s.PoolTimeout),
		IdleCheckFrequency: time.Duration(s.IdleCheckFrequency),
		IdleTimeout:        time.Duration(s.IdleTimeout),
	}

	/* #nosec */
	if s.EnableTLS {
		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify: s.EnableTLS,
		}
	}

	if s.RedisType == ClusterType {
		opts.SentinelAddrs = strings.Split(s.Host, ",")

		return redis.NewFailoverClusterClient(opts)
	}

	return redis.NewFailoverClient(opts)
}

func newClient(s *Settings) redis.UniversalClient {
	if s == nil {
		return nil
	}
	if s.RedisType == ClusterType {
		options := &redis.ClusterOptions{
			Addrs:              strings.Split(s.Host, ","),
			Password:           s.Password,
			Username:           s.Username,
			MaxRetries:         s.RedisMaxRetries,
			MaxRetryBackoff:    time.Duration(s.RedisMaxRetryInterval),
			MinRetryBackoff:    time.Duration(s.RedisMinRetryInterval),
			DialTimeout:        time.Duration(s.DialTimeout),
			ReadTimeout:        time.Duration(s.ReadTimeout),
			WriteTimeout:       time.Duration(s.WriteTimeout),
			PoolSize:           s.PoolSize,
			MaxConnAge:         time.Duration(s.MaxConnAge),
			MinIdleConns:       s.MinIdleConns,
			PoolTimeout:        time.Duration(s.PoolTimeout),
			IdleCheckFrequency: time.Duration(s.IdleCheckFrequency),
			IdleTimeout:        time.Duration(s.IdleTimeout),
		}
		/* #nosec */
		if s.EnableTLS {
			options.TLSConfig = &tls.Config{
				InsecureSkipVerify: s.EnableTLS,
			}
		}

		return redis.NewClusterClient(options)
	}

	options := &redis.Options{
		Addr:               s.Host,
		Password:           s.Password,
		Username:           s.Username,
		DB:                 s.DB,
		MaxRetries:         s.RedisMaxRetries,
		MaxRetryBackoff:    time.Duration(s.RedisMaxRetryInterval),
		MinRetryBackoff:    time.Duration(s.RedisMinRetryInterval),
		DialTimeout:        time.Duration(s.DialTimeout),
		ReadTimeout:        time.Duration(s.ReadTimeout),
		WriteTimeout:       time.Duration(s.WriteTimeout),
		PoolSize:           s.PoolSize,
		MaxConnAge:         time.Duration(s.MaxConnAge),
		MinIdleConns:       s.MinIdleConns,
		PoolTimeout:        time.Duration(s.PoolTimeout),
		IdleCheckFrequency: time.Duration(s.IdleCheckFrequency),
		IdleTimeout:        time.Duration(s.IdleTimeout),
	}

	/* #nosec */
	if s.EnableTLS {
		options.TLSConfig = &tls.Config{
			InsecureSkipVerify: s.EnableTLS,
		}
	}

	return redis.NewClient(options)
}
