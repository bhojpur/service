package jobworker

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
	"testing"

	"github.com/camunda-cloud/zeebe/clients/go/pkg/zbc"
	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/bindings/zeebe"
	"github.com/bhojpur/service/pkg/utils/logger"
)

type mockClientFactory struct {
	zeebe.ClientFactory
	metadata bindings.Metadata
	error    error
}

type mockClient struct {
	zbc.Client
}

func (mcf mockClientFactory) Get(metadata bindings.Metadata) (zbc.Client, error) {
	mcf.metadata = metadata

	if mcf.error != nil {
		return nil, mcf.error
	}

	return mockClient{}, nil
}

func TestInit(t *testing.T) {
	testLogger := logger.NewLogger("test")

	t.Run("jobType is mandatory", func(t *testing.T) {
		metadata := bindings.Metadata{}
		var mcf mockClientFactory

		jobWorker := ZeebeJobWorker{clientFactory: &mcf, logger: testLogger}
		err := jobWorker.Init(metadata)

		assert.Error(t, err, ErrMissingJobType)
	})

	t.Run("sets client from client factory", func(t *testing.T) {
		metadata := bindings.Metadata{
			Properties: map[string]string{"jobType": "a"},
		}
		mcf := mockClientFactory{
			metadata: metadata,
		}
		jobWorker := ZeebeJobWorker{clientFactory: mcf, logger: testLogger}
		err := jobWorker.Init(metadata)

		assert.NoError(t, err)

		mc, err := mcf.Get(metadata)

		assert.NoError(t, err)
		assert.Equal(t, mc, jobWorker.client)
		assert.Equal(t, metadata, mcf.metadata)
	})

	t.Run("returns error if client could not be instantiated properly", func(t *testing.T) {
		errParsing := errors.New("error on parsing metadata")
		metadata := bindings.Metadata{}
		mcf := mockClientFactory{
			error: errParsing,
		}

		jobWorker := ZeebeJobWorker{clientFactory: mcf, logger: testLogger}
		err := jobWorker.Init(metadata)
		assert.Error(t, err, errParsing)
	})

	t.Run("sets client from client factory", func(t *testing.T) {
		metadata := bindings.Metadata{
			Properties: map[string]string{"jobType": "a"},
		}
		mcf := mockClientFactory{
			metadata: metadata,
		}

		jobWorker := ZeebeJobWorker{clientFactory: mcf, logger: testLogger}
		err := jobWorker.Init(metadata)

		assert.NoError(t, err)

		mc, err := mcf.Get(metadata)

		assert.NoError(t, err)
		assert.Equal(t, mc, jobWorker.client)
		assert.Equal(t, metadata, mcf.metadata)
	})
}
