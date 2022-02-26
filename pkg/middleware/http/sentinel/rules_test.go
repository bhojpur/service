package sentinel

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
	"testing"

	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/core/isolation"
	"github.com/alibaba/sentinel-golang/core/system"
	"github.com/stretchr/testify/assert"
)

func TestFlowRules(t *testing.T) {
	rules := []*flow.Rule{
		{
			Resource:               "some-test",
			Threshold:              100,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
		},
	}

	b, _ := json.Marshal(rules)
	t.Logf("%s", b)
	err := loadRules(string(b), newFlowRuleDataSource)
	assert.Nil(t, err)
}

func TestCircuitBreakerRules(t *testing.T) {
	rules := []*circuitbreaker.Rule{
		{
			Resource:         "abc",
			Strategy:         circuitbreaker.ErrorCount,
			RetryTimeoutMs:   3000,
			MinRequestAmount: 10,
			StatIntervalMs:   5000,
			Threshold:        50,
		},
	}

	b, _ := json.Marshal(rules)
	t.Logf("%s", b)
	err := loadRules(string(b), newCircuitBreakerRuleDataSource)
	assert.Nil(t, err)
}

func TestHotspotParamRules(t *testing.T) {
	rules := `
[
	{
		"resource": "abc",
		"metricType": 1,
		"controlBehavior": 0,
		"paramIndex": 1,
		"threshold": 50,
		"burstCount": 0,
		"durationInSec": 1
	}
]
`
	err := loadRules(rules, newHotSpotParamRuleDataSource)
	assert.Nil(t, err)
}

func TestIsolationRules(t *testing.T) {
	rules := []*isolation.Rule{
		{
			Resource:   "abc",
			MetricType: isolation.Concurrency,
			Threshold:  12,
		},
	}

	b, _ := json.Marshal(rules)
	t.Logf("%s", b)
	err := loadRules(string(b), newIsolationRuleDataSource)
	assert.Nil(t, err)
}

func TestSystemRules(t *testing.T) {
	rules := []*system.Rule{
		{
			ID:           "test-id",
			MetricType:   system.InboundQPS,
			TriggerCount: 1000,
			Strategy:     system.BBR,
		},
	}

	b, _ := json.Marshal(rules)
	t.Logf("%s", b)
	err := loadRules(string(b), newSystemRuleDataSource)
	assert.Nil(t, err)
}
