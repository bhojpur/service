package kubernetes

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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/utils/logger"
)

func TestParseMetadata(t *testing.T) {
	nsName := "fooNamespace"
	t.Run("parse metadata", func(t *testing.T) {
		resyncPeriod := time.Second * 15
		m := bindings.Metadata{}
		m.Properties = map[string]string{"namespace": nsName, "resyncPeriodInSec": "15"}

		i := kubernetesInput{logger: logger.NewLogger("test")}
		i.parseMetadata(m)

		assert.Equal(t, nsName, i.namespace, "The namespaces should be the same.")
		assert.Equal(t, resyncPeriod, i.resyncPeriodInSec, "The resyncPeriod should be the same.")
	})
	t.Run("parse metadata no namespace", func(t *testing.T) {
		m := bindings.Metadata{}
		m.Properties = map[string]string{"resyncPeriodInSec": "15"}

		i := kubernetesInput{logger: logger.NewLogger("test")}
		err := i.parseMetadata(m)

		assert.NotNil(t, err, "Expected err to be returned.")
		assert.Equal(t, "namespace is missing in metadata", err.Error(), "Error message not same.")
	})
	t.Run("parse metadata invalid resync period", func(t *testing.T) {
		m := bindings.Metadata{}
		m.Properties = map[string]string{"namespace": nsName, "resyncPeriodInSec": "invalid"}

		i := kubernetesInput{logger: logger.NewLogger("test")}
		err := i.parseMetadata(m)

		assert.Nil(t, err, "Expected err to be nil.")
		assert.Equal(t, nsName, i.namespace, "The namespaces should be the same.")
		assert.Equal(t, time.Second*10, i.resyncPeriodInSec, "The resyncPeriod should be the same.")
	})
}
