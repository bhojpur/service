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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/utils/logger"
)

func TestParseMetadata(t *testing.T) {
	m := bindings.Metadata{Properties: map[string]string{
		"gatewayAddr":            "172.0.0.1:1234",
		"gatewayKeepAlive":       "5s",
		"caCertificatePath":      "/cert/path",
		"usePlaintextConnection": "true",
	}}
	client := ClientFactoryImpl{logger: logger.NewLogger("test")}
	meta, err := client.parseMetadata(m)
	assert.NoError(t, err)
	assert.Equal(t, "172.0.0.1:1234", meta.GatewayAddr)
	assert.Equal(t, 5*time.Second, meta.GatewayKeepAlive.Duration)
	assert.Equal(t, "/cert/path", meta.CaCertificatePath)
	assert.Equal(t, true, meta.UsePlaintextConnection)
}

func TestGatewayAddrMetadataIsMandatory(t *testing.T) {
	m := bindings.Metadata{}
	client := ClientFactoryImpl{logger: logger.NewLogger("test")}
	meta, err := client.parseMetadata(m)
	assert.Nil(t, meta)
	assert.Error(t, err)
	assert.Equal(t, err, ErrMissingGatewayAddr)
}

func TestParseMetadataDefaultValues(t *testing.T) {
	m := bindings.Metadata{Properties: map[string]string{"gatewayAddr": "172.0.0.1:1234"}}
	client := ClientFactoryImpl{logger: logger.NewLogger("test")}
	meta, err := client.parseMetadata(m)
	assert.NoError(t, err)
	assert.Equal(t, time.Duration(0), meta.GatewayKeepAlive.Duration)
	assert.Equal(t, "", meta.CaCertificatePath)
	assert.Equal(t, false, meta.UsePlaintextConnection)
}
