package keyvault

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

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/secretstores"
	"github.com/bhojpur/service/pkg/utils/logger"
)

func TestInit(t *testing.T) {
	m := secretstores.Metadata{}
	s := NewAzureKeyvaultSecretStore(logger.NewLogger("test"))
	t.Run("Init with valid metadata", func(t *testing.T) {
		m.Properties = map[string]string{
			"vaultName":         "foo",
			"azureTenantId":     "00000000-0000-0000-0000-000000000000",
			"azureClientId":     "00000000-0000-0000-0000-000000000000",
			"azureClientSecret": "passw0rd",
		}
		err := s.Init(m)
		assert.Nil(t, err)
		kv, ok := s.(*keyvaultSecretStore)
		assert.True(t, ok)
		assert.Equal(t, kv.vaultName, "foo")
		assert.Equal(t, kv.vaultDNSSuffix, "vault.azure.net")
		assert.NotNil(t, kv.vaultClient)
	})
	t.Run("Init with valid metadata and Azure environment", func(t *testing.T) {
		m.Properties = map[string]string{
			"vaultName":         "foo",
			"azureTenantId":     "00000000-0000-0000-0000-000000000000",
			"azureClientId":     "00000000-0000-0000-0000-000000000000",
			"azureClientSecret": "passw0rd",
			"azureEnvironment":  "AZURECHINACLOUD",
		}
		err := s.Init(m)
		assert.Nil(t, err)
		kv, ok := s.(*keyvaultSecretStore)
		assert.True(t, ok)
		assert.Equal(t, kv.vaultName, "foo")
		assert.Equal(t, kv.vaultDNSSuffix, "vault.azure.cn")
		assert.NotNil(t, kv.vaultClient)
	})
	t.Run("Init with Azure environment as part of vaultName FQDN (1) - legacy", func(t *testing.T) {
		m.Properties = map[string]string{
			"vaultName":         "foo.vault.azure.cn",
			"azureTenantId":     "00000000-0000-0000-0000-000000000000",
			"azureClientId":     "00000000-0000-0000-0000-000000000000",
			"azureClientSecret": "passw0rd",
		}
		err := s.Init(m)
		assert.Nil(t, err)
		kv, ok := s.(*keyvaultSecretStore)
		assert.True(t, ok)
		assert.Equal(t, kv.vaultName, "foo")
		assert.Equal(t, kv.vaultDNSSuffix, "vault.azure.cn")
		assert.NotNil(t, kv.vaultClient)
	})
	t.Run("Init with Azure environment as part of vaultName FQDN (2) - legacy", func(t *testing.T) {
		m.Properties = map[string]string{
			"vaultName":         "https://foo.vault.usgovcloudapi.net",
			"azureTenantId":     "00000000-0000-0000-0000-000000000000",
			"azureClientId":     "00000000-0000-0000-0000-000000000000",
			"azureClientSecret": "passw0rd",
		}
		err := s.Init(m)
		assert.Nil(t, err)
		kv, ok := s.(*keyvaultSecretStore)
		assert.True(t, ok)
		assert.Equal(t, kv.vaultName, "foo")
		assert.Equal(t, kv.vaultDNSSuffix, "vault.usgovcloudapi.net")
		assert.NotNil(t, kv.vaultClient)
	})
}
