package zeebe

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
	"errors"

	"github.com/camunda-cloud/zeebe/clients/go/pkg/zbc"

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/metadata"
	"github.com/bhojpur/service/pkg/utils/logger"
)

var ErrMissingGatewayAddr = errors.New("gatewayAddr is a required attribute")

// ClientFactory enables injection for testing.
type ClientFactory interface {
	Get(metadata bindings.Metadata) (zbc.Client, error)
}

type ClientFactoryImpl struct {
	logger logger.Logger
}

// https://docs.zeebe.io/operations/authentication.html
type clientMetadata struct {
	GatewayAddr            string            `json:"gatewayAddr"`
	GatewayKeepAlive       metadata.Duration `json:"gatewayKeepAlive"`
	CaCertificatePath      string            `json:"caCertificatePath"`
	UsePlaintextConnection bool              `json:"usePlainTextConnection,string"`
}

// NewClientFactoryImpl returns a new ClientFactory instance.
func NewClientFactoryImpl(logger logger.Logger) *ClientFactoryImpl {
	return &ClientFactoryImpl{logger: logger}
}

func (c *ClientFactoryImpl) Get(metadata bindings.Metadata) (zbc.Client, error) {
	meta, err := c.parseMetadata(metadata)
	if err != nil {
		return nil, err
	}

	client, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         meta.GatewayAddr,
		UsePlaintextConnection: meta.UsePlaintextConnection,
		CaCertificatePath:      meta.CaCertificatePath,
		KeepAlive:              meta.GatewayKeepAlive.Duration,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *ClientFactoryImpl) parseMetadata(metadata bindings.Metadata) (*clientMetadata, error) {
	b, err := json.Marshal(metadata.Properties)
	if err != nil {
		return nil, err
	}

	var m clientMetadata
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}

	if m.GatewayAddr == "" {
		return nil, ErrMissingGatewayAddr
	}

	return &m, nil
}
