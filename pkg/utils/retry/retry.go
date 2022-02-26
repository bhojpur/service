package retry

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
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"

	"github.com/bhojpur/service/pkg/utils/config"
)

// PolicyType denotes if the back off delay should be constant or exponential.
type PolicyType int

const (
	// PolicyConstant is a backoff policy that always returns the same backoff delay.
	PolicyConstant PolicyType = iota
	// PolicyExponential is a backoff implementation that increases the backoff period
	// for each retry attempt using a randomization function that grows exponentially.
	PolicyExponential
)

// Config encapsulates the back off policy configuration.
type Config struct {
	Policy PolicyType `mapstructure:"policy"`

	// Constant back off
	Duration time.Duration `mapstructure:"duration"`

	// Exponential back off
	InitialInterval     time.Duration `mapstructure:"initialInterval"`
	RandomizationFactor float32       `mapstructure:"randomizationFactor"`
	Multiplier          float32       `mapstructure:"multiplier"`
	MaxInterval         time.Duration `mapstructure:"maxInterval"`
	MaxElapsedTime      time.Duration `mapstructure:"maxElapsedTime"`

	// Additional options
	MaxRetries int64 `mapstructure:"maxRetries"`
}

// DefaultConfig represents the default configuration for a
// `Config`.
func DefaultConfig() Config {
	return Config{
		Policy:              PolicyConstant,
		Duration:            5 * time.Second,
		InitialInterval:     backoff.DefaultInitialInterval,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          backoff.DefaultMultiplier,
		MaxInterval:         backoff.DefaultMaxInterval,
		MaxElapsedTime:      backoff.DefaultMaxElapsedTime,
		MaxRetries:          -1,
	}
}

// DefaultConfigWithNoRetry represents the default configuration with `MaxRetries` set to 0.
// This may be useful for those brokers which can handles retries on its own.
func DefaultConfigWithNoRetry() Config {
	c := DefaultConfig()
	c.MaxRetries = 0

	return c
}

// DecodeConfig decodes a Go struct into a `Config`.
func DecodeConfig(c *Config, input interface{}) error {
	// Use the deefault config if `c` is empty/zero value.
	var emptyConfig Config
	if *c == emptyConfig {
		*c = DefaultConfig()
	}

	return config.Decode(input, c)
}

// DecodeConfigWithPrefix decodes a Go struct into a `Config`.
func DecodeConfigWithPrefix(c *Config, input interface{}, prefix string) error {
	input, err := config.PrefixedBy(input, prefix)
	if err != nil {
		return err
	}

	return DecodeConfig(c, input)
}

// NewBackOff returns a BackOff instance for use with `NotifyRecover`
// or `backoff.RetryNotify` directly. The instance will not stop due to
// context cancellation. To support cancellation (recommended), use
// `NewBackOffWithContext`.
//
// Since the underlying backoff implementations are not always thread safe,
// `NewBackOff` or `NewBackOffWithContext` should be called each time
// `RetryNotifyRecover` or `backoff.RetryNotify` is used.
func (c *Config) NewBackOff() backoff.BackOff {
	var b backoff.BackOff
	switch c.Policy {
	case PolicyConstant:
		b = backoff.NewConstantBackOff(c.Duration)
	case PolicyExponential:
		eb := backoff.NewExponentialBackOff()
		eb.InitialInterval = c.InitialInterval
		eb.RandomizationFactor = float64(c.RandomizationFactor)
		eb.Multiplier = float64(c.Multiplier)
		eb.MaxInterval = c.MaxInterval
		eb.MaxElapsedTime = c.MaxElapsedTime
		b = eb
	}

	if c.MaxRetries >= 0 {
		b = backoff.WithMaxRetries(b, uint64(c.MaxRetries))
	}

	return b
}

// NewBackOffWithContext returns a BackOff instance for use with `RetryNotifyRecover`
// or `backoff.RetryNotify` directly. The provided context is used to cancel retries
// if it is canceled.
//
// Since the underlying backoff implementations are not always thread safe,
// `NewBackOff` or `NewBackOffWithContext` should be called each time
// `RetryNotifyRecover` or `backoff.RetryNotify` is used.
func (c *Config) NewBackOffWithContext(ctx context.Context) backoff.BackOff {
	b := c.NewBackOff()

	return backoff.WithContext(b, ctx)
}

// NotifyRecover is a wrapper around backoff.RetryNotify that adds another callback for when an operation
// previously failed but has since recovered. The main purpose of this wrapper is to call `notify` only when
// the operations fails the first time and `recovered` when it finally succeeds. This can be helpful in limiting
// log messages to only the events that operators need to be alerted on.
func NotifyRecover(operation backoff.Operation, b backoff.BackOff, notify backoff.Notify, recovered func()) error {
	var notified bool

	return backoff.RetryNotify(func() error {
		err := operation()

		if err == nil && notified {
			notified = false
			recovered()
		}

		return err
	}, b, func(err error, d time.Duration) {
		if !notified {
			notify(err, d)
			notified = true
		}
	})
}

// DecodeString handles converting a string value to `p`.
func (p *PolicyType) DecodeString(value string) error {
	switch strings.ToLower(value) {
	case "constant":
		*p = PolicyConstant
	case "exponential":
		*p = PolicyExponential
	default:
		return errors.Errorf("unexpected back off policy type: %s", value)
	}

	return nil
}
