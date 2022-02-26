package azure

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
	"fmt"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Azure/go-autorest/autorest/azure"

	"github.com/bhojpur/service/pkg/utils/logger"
)

const (
	storageAccountKeyKey = "accountKey"
)

// GetAzureStorageCredentials returns a azblob.Credential object that can be used to authenticate an Azure Blob Storage SDK pipeline.
// First it tries to authenticate using shared key credentials (using an account key) if present. It falls back to attempting to use Azure AD (via a service principal or MSI).
func GetAzureStorageCredentials(log logger.Logger, accountName string, metadata map[string]string) (azblob.Credential, *azure.Environment, error) {
	settings, err := NewEnvironmentSettings("storage", metadata)
	if err != nil {
		return nil, nil, err
	}

	// Try using shared key credentials first
	accountKey, ok := metadata[storageAccountKeyKey]
	if ok && accountKey != "" {
		credential, newSharedKeyErr := azblob.NewSharedKeyCredential(accountName, accountKey)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid credentials with error: %s", newSharedKeyErr.Error())
		}

		return credential, settings.AzureEnvironment, nil
	}

	// Fallback to using Azure AD
	spt, err := settings.GetServicePrincipalToken()
	if err != nil {
		return nil, nil, err
	}
	var tokenRefresher azblob.TokenRefresher = func(credential azblob.TokenCredential) time.Duration {
		log.Debug("Refreshing Azure Storage auth token")
		err := spt.Refresh()
		if err != nil {
			panic(err)
		}
		token := spt.Token()
		credential.SetToken(token.AccessToken)

		// Make the token expire 2 minutes earlier to get some extra buffer
		exp := token.Expires().Sub(time.Now().Add(2 * time.Minute))
		log.Debug("Received new token, valid for", exp)

		return exp
	}
	credential := azblob.NewTokenCredential("", tokenRefresher)

	return credential, settings.AzureEnvironment, nil
}
