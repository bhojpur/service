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
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"

	azauth "github.com/bhojpur/service/pkg/authentication/azure"
	"github.com/bhojpur/service/pkg/secretstores"
	"github.com/bhojpur/service/pkg/utils/logger"
)

// Keyvault secret store component metadata properties
// This is in addition to what's defined in authentication/azure.
const (
	componentVaultName = "vaultName"
	VersionID          = "version_id"
	secretItemIDPrefix = "/secrets/"
)

type keyvaultSecretStore struct {
	vaultName      string
	vaultClient    *azsecrets.Client
	vaultDNSSuffix string

	logger logger.Logger
}

// NewAzureKeyvaultSecretStore returns a new Azure Key Vault secret store.
func NewAzureKeyvaultSecretStore(logger logger.Logger) secretstores.SecretStore {
	return &keyvaultSecretStore{
		vaultName:   "",
		vaultClient: nil,
		logger:      logger,
	}
}

// Init creates a Azure Key Vault client.
func (k *keyvaultSecretStore) Init(metadata secretstores.Metadata) error {
	// Fix for maintaining backwards compatibility with a change introduced in 1.3 that allowed specifying an Azure environment by setting a FQDN for vault name
	// This should be considered deprecated and users should rely the "azureEnvironment" metadata instead, but it's maintained here for backwards-compatibility
	if vaultName, ok := metadata.Properties[componentVaultName]; ok {
		keyVaultSuffixToEnvironment := map[string]string{
			".vault.azure.net":         "AZUREPUBLICCLOUD",
			".vault.azure.cn":          "AZURECHINACLOUD",
			".vault.usgovcloudapi.net": "AZUREUSGOVERNMENTCLOUD",
			".vault.microsoftazure.de": "AZUREGERMANCLOUD",
		}
		for suffix, environment := range keyVaultSuffixToEnvironment {
			if strings.HasSuffix(vaultName, suffix) {
				metadata.Properties["azureEnvironment"] = environment
				vaultName = strings.TrimSuffix(vaultName, suffix)
				if strings.HasPrefix(vaultName, "https://") {
					vaultName = strings.TrimPrefix(vaultName, "https://")
				}
				metadata.Properties[componentVaultName] = vaultName

				break
			}
		}
	}

	// Initialization code
	settings, err := azauth.NewEnvironmentSettings("keyvault", metadata.Properties)
	if err != nil {
		return err
	}

	k.vaultName = settings.Values[componentVaultName]
	k.vaultDNSSuffix = settings.AzureEnvironment.KeyVaultDNSSuffix

	cred, err := settings.GetTokenCredential()
	if err != nil {
		return err
	}
	coreClientOpts := azcore.ClientOptions{
		Telemetry: policy.TelemetryOptions{
			ApplicationID: "app-" + logger.AppVersion,
		},
	}
	k.vaultClient, err = azsecrets.NewClient(k.getVaultURI(), cred, &azsecrets.ClientOptions{
		ClientOptions: coreClientOpts,
	})
	if err != nil {
		return err
	}

	return nil
}

// GetSecret retrieves a secret using a key and returns a map of decrypted string/string values.
func (k *keyvaultSecretStore) GetSecret(req secretstores.GetSecretRequest) (secretstores.GetSecretResponse, error) {
	opts := &azsecrets.GetSecretOptions{}
	if value, ok := req.Metadata[VersionID]; ok {
		opts.Version = value
	}

	secretResp, err := k.vaultClient.GetSecret(context.TODO(), req.Name, opts)
	if err != nil {
		return secretstores.GetSecretResponse{}, err
	}

	secretValue := ""
	if secretResp.Value != nil {
		secretValue = *secretResp.Value
	}

	return secretstores.GetSecretResponse{
		Data: map[string]string{
			req.Name: secretValue,
		},
	}, nil
}

// BulkGetSecret retrieves all secrets in the store and returns a map of decrypted string/string values.
func (k *keyvaultSecretStore) BulkGetSecret(req secretstores.BulkGetSecretRequest) (secretstores.BulkGetSecretResponse, error) {
	maxResults, err := k.getMaxResultsFromMetadata(req.Metadata)
	if err != nil {
		return secretstores.BulkGetSecretResponse{}, err
	}

	resp := secretstores.BulkGetSecretResponse{
		Data: map[string]map[string]string{},
	}

	secretIDPrefix := k.getVaultURI() + secretItemIDPrefix

	pager := k.vaultClient.ListSecrets(&azsecrets.ListSecretsOptions{
		MaxResults: maxResults,
	})

	for pager.NextPage(context.TODO()) {
		pr := pager.PageResponse()
		for _, secret := range pr.Secrets {
			if secret.Attributes == nil || secret.Attributes.Enabled == nil || !*secret.Attributes.Enabled {
				continue
			}

			secretName := strings.TrimPrefix(*secret.ID, secretIDPrefix)
			secretResp, err := k.vaultClient.GetSecret(context.TODO(), secretName, nil)
			if err != nil {
				return secretstores.BulkGetSecretResponse{}, err
			}

			secretValue := ""
			if secretResp.Value != nil {
				secretValue = *secretResp.Value
			}

			resp.Data[secretName] = map[string]string{secretName: secretValue}
		}
	}

	if pager.Err() != nil {
		return secretstores.BulkGetSecretResponse{}, pager.Err()
	}

	return resp, nil
}

// getVaultURI returns Azure Key Vault URI.
func (k *keyvaultSecretStore) getVaultURI() string {
	return fmt.Sprintf("https://%s.%s", k.vaultName, k.vaultDNSSuffix)
}

func (k *keyvaultSecretStore) getMaxResultsFromMetadata(metadata map[string]string) (*int32, error) {
	if s, ok := metadata["maxresults"]; ok && s != "" {
		/* #nosec */
		val, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		converted := int32(val)

		return &converted, nil
	}

	return nil, nil
}
