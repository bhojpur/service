package env

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
	"os"
	"strings"

	"github.com/bhojpur/service/pkg/secretstores"
	"github.com/bhojpur/service/pkg/utils/logger"
)

type envSecretStore struct {
	logger logger.Logger
}

// NewEnvSecretStore returns a new env var secret store.
func NewEnvSecretStore(logger logger.Logger) secretstores.SecretStore {
	return &envSecretStore{
		logger: logger,
	}
}

// Init creates a Local secret store.
func (s *envSecretStore) Init(metadata secretstores.Metadata) error {
	return nil
}

// GetSecret retrieves a secret from env var using provided key.
func (s *envSecretStore) GetSecret(req secretstores.GetSecretRequest) (secretstores.GetSecretResponse, error) {
	return secretstores.GetSecretResponse{
		Data: map[string]string{
			req.Name: os.Getenv(req.Name),
		},
	}, nil
}

// BulkGetSecret retrieves all secrets in the store and returns a map of decrypted string/string values.
func (s *envSecretStore) BulkGetSecret(req secretstores.BulkGetSecretRequest) (secretstores.BulkGetSecretResponse, error) {
	r := map[string]map[string]string{}

	for _, element := range os.Environ() {
		envVariable := strings.SplitN(element, "=", 2)
		r[envVariable[0]] = map[string]string{envVariable[0]: envVariable[1]}
	}

	return secretstores.BulkGetSecretResponse{
		Data: r,
	}, nil
}
