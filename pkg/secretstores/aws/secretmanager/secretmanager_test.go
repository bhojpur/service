package secretmanager

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
	"testing"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/secretstores"
	"github.com/bhojpur/service/pkg/utils/logger"
)

const secretValue = "secret"

type mockedSM struct {
	GetSecretValueFn func(*secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error)
	secretsmanageriface.SecretsManagerAPI
}

func (m *mockedSM) GetSecretValue(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	return m.GetSecretValueFn(input)
}

func TestInit(t *testing.T) {
	m := secretstores.Metadata{}
	s := NewSecretManager(logger.NewLogger("test"))
	t.Run("Init with valid metadata", func(t *testing.T) {
		m.Properties = map[string]string{
			"AccessKey":    "a",
			"Region":       "a",
			"Endpoint":     "a",
			"SecretKey":    "a",
			"SessionToken": "a",
		}
		err := s.Init(m)
		assert.Nil(t, err)
	})
}

func TestGetSecret(t *testing.T) {
	t.Run("successfully retrieve secret", func(t *testing.T) {
		t.Run("without version id and version stage", func(t *testing.T) {
			s := smSecretStore{
				client: &mockedSM{
					GetSecretValueFn: func(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
						assert.Nil(t, input.VersionId)
						assert.Nil(t, input.VersionStage)
						secret := secretValue

						return &secretsmanager.GetSecretValueOutput{
							Name:         input.SecretId,
							SecretString: &secret,
						}, nil
					},
				},
			}

			req := secretstores.GetSecretRequest{
				Name:     "/aws/secret/testing",
				Metadata: map[string]string{},
			}
			output, e := s.GetSecret(req)
			assert.Nil(t, e)
			assert.Equal(t, "secret", output.Data[req.Name])
		})

		t.Run("with version id", func(t *testing.T) {
			s := smSecretStore{
				client: &mockedSM{
					GetSecretValueFn: func(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
						assert.NotNil(t, input.VersionId)
						secret := secretValue

						return &secretsmanager.GetSecretValueOutput{
							Name:         input.SecretId,
							SecretString: &secret,
						}, nil
					},
				},
			}

			req := secretstores.GetSecretRequest{
				Name: "/aws/secret/testing",
				Metadata: map[string]string{
					VersionID: "1",
				},
			}
			output, e := s.GetSecret(req)
			assert.Nil(t, e)
			assert.Equal(t, secretValue, output.Data[req.Name])
		})

		t.Run("with version stage", func(t *testing.T) {
			s := smSecretStore{
				client: &mockedSM{
					GetSecretValueFn: func(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
						assert.NotNil(t, input.VersionStage)
						secret := secretValue

						return &secretsmanager.GetSecretValueOutput{
							Name:         input.SecretId,
							SecretString: &secret,
						}, nil
					},
				},
			}

			req := secretstores.GetSecretRequest{
				Name: "/aws/secret/testing",
				Metadata: map[string]string{
					VersionStage: "dev",
				},
			}
			output, e := s.GetSecret(req)
			assert.Nil(t, e)
			assert.Equal(t, secretValue, output.Data[req.Name])
		})
	})

	t.Run("unsuccessfully retrieve secret", func(t *testing.T) {
		s := smSecretStore{
			client: &mockedSM{
				GetSecretValueFn: func(input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
					return nil, fmt.Errorf("failed due to any reason")
				},
			},
		}
		req := secretstores.GetSecretRequest{
			Name:     "/aws/secret/testing",
			Metadata: map[string]string{},
		}
		_, err := s.GetSecret(req)
		assert.NotNil(t, err)
	})
}
