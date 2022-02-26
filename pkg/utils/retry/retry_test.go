package retry_test

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
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/utils/retry"
)

var errRetry = errors.New("Testing")

func TestDecode(t *testing.T) {
	tests := map[string]struct {
		config    interface{}
		overrides func(config *retry.Config)
		err       string
	}{
		"invalid policy type": {
			config: map[string]interface{}{
				"backOffPolicy": "invalid",
			},
			overrides: nil,
			err:       "1 error(s) decoding:\n\n* error decoding 'policy': invalid PolicyType \"invalid\": unexpected back off policy type: invalid",
		},
		"default": {
			config:    map[string]interface{}{},
			overrides: nil,
			err:       "",
		},
		"constant default": {
			config: map[string]interface{}{
				"backOffPolicy": "constant",
			},
			overrides: nil,
			err:       "",
		},
		"constant with duraction": {
			config: map[string]interface{}{
				"backOffPolicy":   "constant",
				"backOffDuration": "10s",
			},
			overrides: func(config *retry.Config) {
				config.Duration = 10 * time.Second
			},
			err: "",
		},
		"exponential default": {
			config: map[string]interface{}{
				"backOffPolicy": "exponential",
			},
			overrides: func(config *retry.Config) {
				config.Policy = retry.PolicyExponential
			},
			err: "",
		},
		"exponential with string settings": {
			config: map[string]interface{}{
				"backOffPolicy":              "exponential",
				"backOffInitialInterval":     "1000", // 1s
				"backOffRandomizationFactor": "1.0",
				"backOffMultiplier":          "2.0",
				"backOffMaxInterval":         "120000",  // 2m
				"backOffMaxElapsedTime":      "1800000", // 30m
			},
			overrides: func(config *retry.Config) {
				config.Policy = retry.PolicyExponential
				config.InitialInterval = 1 * time.Second
				config.RandomizationFactor = 1.0
				config.Multiplier = 2.0
				config.MaxInterval = 2 * time.Minute
				config.MaxElapsedTime = 30 * time.Minute
			},
			err: "",
		},
		"exponential with typed settings": {
			config: map[string]interface{}{
				"backOffPolicy":              "exponential",
				"backOffInitialInterval":     "1000ms", // 1s
				"backOffRandomizationFactor": 1.0,
				"backOffMultiplier":          2.0,
				"backOffMaxInterval":         "120s", // 2m
				"backOffMaxElapsedTime":      "30m",  // 30m
			},
			overrides: func(config *retry.Config) {
				config.Policy = retry.PolicyExponential
				config.InitialInterval = 1 * time.Second
				config.RandomizationFactor = 1.0
				config.Multiplier = 2.0
				config.MaxInterval = 2 * time.Minute
				config.MaxElapsedTime = 30 * time.Minute
			},
			err: "",
		},
		"map[string]string settings": {
			config: map[string]string{
				"backOffPolicy":              "exponential",
				"backOffInitialInterval":     "1000ms", // 1s
				"backOffRandomizationFactor": "1.0",
				"backOffMultiplier":          "2.0",
				"backOffMaxInterval":         "120s", // 2m
				"backOffMaxElapsedTime":      "30m",  // 30m
			},
			overrides: func(config *retry.Config) {
				config.Policy = retry.PolicyExponential
				config.InitialInterval = 1 * time.Second
				config.RandomizationFactor = 1.0
				config.Multiplier = 2.0
				config.MaxInterval = 2 * time.Minute
				config.MaxElapsedTime = 30 * time.Minute
			},
			err: "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var actual retry.Config
			err := retry.DecodeConfigWithPrefix(&actual, tc.config, "backOff")
			if tc.err != "" {
				if assert.Error(t, err) {
					assert.Equal(t, tc.err, err.Error())
				}
			} else {
				b := actual.NewBackOff()
				config := retry.DefaultConfig()
				if tc.overrides != nil {
					tc.overrides(&config)
				}
				assert.Equal(t, config, actual, "unexpected decoded configuration")
				if actual.Policy == retry.PolicyConstant {
					_, ok := b.(*backoff.ConstantBackOff)
					assert.True(t, ok)
				} else if actual.Policy == retry.PolicyExponential {
					_, ok := b.(*backoff.ExponentialBackOff)
					assert.True(t, ok)
				}
			}
		})
	}
}

func TestRetryNotifyRecoverNoetries(t *testing.T) {
	config := retry.DefaultConfigWithNoRetry()
	config.Duration = 1

	var operationCalls, notifyCalls, recoveryCalls int

	b := config.NewBackOff()
	err := retry.NotifyRecover(func() error {
		operationCalls++

		return errRetry
	}, b, func(err error, d time.Duration) {
		notifyCalls++
	}, func() {
		recoveryCalls++
	})

	assert.Error(t, err)
	assert.Equal(t, errRetry, err)
	assert.Equal(t, 1, operationCalls)
	assert.Equal(t, 0, notifyCalls)
	assert.Equal(t, 0, recoveryCalls)
}

func TestRetryNotifyRecoverMaxRetries(t *testing.T) {
	config := retry.DefaultConfig()
	config.MaxRetries = 3
	config.Duration = 1

	var operationCalls, notifyCalls, recoveryCalls int

	b := config.NewBackOff()
	err := retry.NotifyRecover(func() error {
		operationCalls++

		return errRetry
	}, b, func(err error, d time.Duration) {
		notifyCalls++
	}, func() {
		recoveryCalls++
	})

	assert.Error(t, err)
	assert.Equal(t, errRetry, err)
	assert.Equal(t, 4, operationCalls)
	assert.Equal(t, 1, notifyCalls)
	assert.Equal(t, 0, recoveryCalls)
}

func TestRetryNotifyRecoverRecovery(t *testing.T) {
	config := retry.DefaultConfig()
	config.MaxRetries = 3
	config.Duration = 1

	var operationCalls, notifyCalls, recoveryCalls int

	b := config.NewBackOff()
	err := retry.NotifyRecover(func() error {
		operationCalls++

		if operationCalls >= 2 {
			return nil
		}

		return errRetry
	}, b, func(err error, d time.Duration) {
		notifyCalls++
	}, func() {
		recoveryCalls++
	})

	assert.NoError(t, err)
	assert.Equal(t, 2, operationCalls)
	assert.Equal(t, 1, notifyCalls)
	assert.Equal(t, 1, recoveryCalls)
}

func TestRetryNotifyRecoverCancel(t *testing.T) {
	config := retry.DefaultConfig()
	config.Policy = retry.PolicyConstant
	config.Duration = 1 * time.Minute

	var notifyCalls, recoveryCalls int

	ctx, cancel := context.WithCancel(context.Background())
	b := config.NewBackOffWithContext(ctx)
	errC := make(chan error, 1)
	startedC := make(chan struct{}, 100)

	go func() {
		errC <- retry.NotifyRecover(func() error {
			return errRetry
		}, b, func(err error, d time.Duration) {
			notifyCalls++
			startedC <- struct{}{}
		}, func() {
			recoveryCalls++
		})
	}()

	<-startedC
	cancel()

	err := <-errC
	assert.Error(t, err)
	assert.True(t, errors.Is(err, context.Canceled))
	assert.Equal(t, 1, notifyCalls)
	assert.Equal(t, 0, recoveryCalls)
}

func TestCheckEmptyConfig(t *testing.T) {
	var config retry.Config
	err := retry.DecodeConfig(&config, map[string]interface{}{})
	assert.NoError(t, err)
	defaultConfig := retry.DefaultConfig()
	assert.Equal(t, config, defaultConfig)
}
