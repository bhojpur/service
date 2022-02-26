package file

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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bhojpur/service/pkg/secretstores"
	"github.com/bhojpur/service/pkg/utils/logger"
)

const secretValue = "secret"

func TestInit(t *testing.T) {
	m := secretstores.Metadata{}
	s := localSecretStore{
		logger: logger.NewLogger("test"),
		readLocalFileFn: func(secretsFile string) (map[string]interface{}, error) {
			return nil, nil
		},
	}
	t.Run("Init with valid metadata", func(t *testing.T) {
		m.Properties = map[string]string{
			"secretsFile":     "a",
			"nestedSeparator": "a",
		}
		err := s.Init(m)
		assert.Nil(t, err)
	})

	t.Run("Init with missing metadata", func(t *testing.T) {
		m.Properties = map[string]string{
			"dummy": "a",
		}
		err := s.Init(m)
		assert.NotNil(t, err)
		assert.Equal(t, err, fmt.Errorf("missing local secrets file in metadata"))
	})
}

func TestSeparator(t *testing.T) {
	m := secretstores.Metadata{}
	s := localSecretStore{
		logger: logger.NewLogger("test"),
		readLocalFileFn: func(secretsFile string) (map[string]interface{}, error) {
			return map[string]interface{}{
				"root": map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
			}, nil
		},
	}
	t.Run("Init with custom separator", func(t *testing.T) {
		m.Properties = map[string]string{
			"secretsFile":     "a",
			"nestedSeparator": ".",
		}
		err := s.Init(m)
		assert.Nil(t, err)

		req := secretstores.GetSecretRequest{
			Name:     "root.key1",
			Metadata: map[string]string{},
		}
		output, err := s.GetSecret(req)
		assert.Nil(t, err)
		assert.Equal(t, "value1", output.Data[req.Name])
	})

	t.Run("Init with default separator", func(t *testing.T) {
		m.Properties = map[string]string{
			"secretsFile": "a",
		}
		err := s.Init(m)
		assert.Nil(t, err)

		req := secretstores.GetSecretRequest{
			Name:     "root:key2",
			Metadata: map[string]string{},
		}
		output, err := s.GetSecret(req)
		assert.Nil(t, err)
		assert.Equal(t, "value2", output.Data[req.Name])
	})
}

func TestGetSecret(t *testing.T) {
	m := secretstores.Metadata{}
	m.Properties = map[string]string{
		"secretsFile":     "a",
		"nestedSeparator": "a",
	}
	s := localSecretStore{
		logger: logger.NewLogger("test"),
		readLocalFileFn: func(secretsFile string) (map[string]interface{}, error) {
			secrets := make(map[string]interface{})
			secrets["secret"] = secretValue

			return secrets, nil
		},
	}
	s.Init(m)

	t.Run("successfully retrieve secrets", func(t *testing.T) {
		req := secretstores.GetSecretRequest{
			Name:     "secret",
			Metadata: map[string]string{},
		}
		output, e := s.GetSecret(req)
		assert.Nil(t, e)
		assert.Equal(t, "secret", output.Data[req.Name])
	})

	t.Run("unsuccessfully retrieve secret", func(t *testing.T) {
		req := secretstores.GetSecretRequest{
			Name:     "NoSecret",
			Metadata: map[string]string{},
		}
		_, err := s.GetSecret(req)
		assert.NotNil(t, err)
		assert.Equal(t, err, fmt.Errorf("secret %s not found", req.Name))
	})
}

func TestBulkGetSecret(t *testing.T) {
	m := secretstores.Metadata{}
	m.Properties = map[string]string{
		"secretsFile":     "a",
		"nestedSeparator": "a",
	}
	s := localSecretStore{
		logger: logger.NewLogger("test"),
		readLocalFileFn: func(secretsFile string) (map[string]interface{}, error) {
			secrets := make(map[string]interface{})
			secrets["secret"] = secretValue

			return secrets, nil
		},
	}
	s.Init(m)

	t.Run("successfully retrieve secrets", func(t *testing.T) {
		req := secretstores.BulkGetSecretRequest{}
		output, e := s.BulkGetSecret(req)
		assert.Nil(t, e)
		assert.Equal(t, "secret", output.Data["secret"]["secret"])
	})
}

func TestMultiValuedSecrets(t *testing.T) {
	m := secretstores.Metadata{}
	m.Properties = map[string]string{
		"secretsFile": "a",
		"multiValued": "true",
	}
	s := localSecretStore{
		logger: logger.NewLogger("test"),
		readLocalFileFn: func(secretsFile string) (map[string]interface{}, error) {
			//nolint:gosec
			secretsJSON := `
			{
				"parent": {
					"child1": "12345",
					"child2": {
						"child3": "67890",
						"child4": "00000"
					}
				}
			}
			`
			var secrets map[string]interface{}
			err := json.Unmarshal([]byte(secretsJSON), &secrets)

			return secrets, err
		},
	}
	err := s.Init(m)
	require.NoError(t, err)

	t.Run("successfully retrieve a single multi-valued secret", func(t *testing.T) {
		req := secretstores.GetSecretRequest{
			Name: "parent",
		}
		resp, err := s.GetSecret(req)
		require.NoError(t, err)
		assert.Equal(t, map[string]string{
			"child1":        "12345",
			"child2:child3": "67890",
			"child2:child4": "00000",
		}, resp.Data)
	})

	t.Run("successfully retrieve multi-valued secrets", func(t *testing.T) {
		req := secretstores.BulkGetSecretRequest{}
		resp, err := s.BulkGetSecret(req)
		require.NoError(t, err)
		assert.Equal(t, map[string]map[string]string{
			"parent": {
				"child1":        "12345",
				"child2:child3": "67890",
				"child2:child4": "00000",
			},
		}, resp.Data)
	})
}