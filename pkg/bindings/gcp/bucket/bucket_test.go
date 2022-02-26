package bucket

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

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/utils/logger"
)

func TestParseMetadata(t *testing.T) {
	t.Run("Has correct metadata", func(t *testing.T) {
		m := bindings.Metadata{}
		m.Properties = map[string]string{
			"auth_provider_x509_cert_url": "my_auth_provider_x509",
			"auth_uri":                    "my_auth_uri",
			"Bucket":                      "my_bucket",
			"client_x509_cert_url":        "my_client_x509",
			"client_email":                "my_email@mail.app",
			"client_id":                   "my_client_id",
			"private_key":                 "my_private_key",
			"private_key_id":              "my_private_key_id",
			"project_id":                  "my_project_id",
			"token_uri":                   "my_token_uri",
			"type":                        "my_type",
		}
		gs := GCPStorage{logger: logger.NewLogger("test")}
		meta, _, err := gs.parseMetadata(m)
		assert.Nil(t, err)

		assert.Equal(t, "my_auth_provider_x509", meta.AuthProviderCertURL)
		assert.Equal(t, "my_auth_uri", meta.AuthURI)
		assert.Equal(t, "my_bucket", meta.Bucket)
		assert.Equal(t, "my_client_x509", meta.ClientCertURL)
		assert.Equal(t, "my_email@mail.app", meta.ClientEmail)
		assert.Equal(t, "my_client_id", meta.ClientID)
		assert.Equal(t, "my_private_key", meta.PrivateKey)
		assert.Equal(t, "my_private_key_id", meta.PrivateKeyID)
		assert.Equal(t, "my_project_id", meta.ProjectID)
		assert.Equal(t, "my_token_uri", meta.TokenURI)
		assert.Equal(t, "my_type", meta.Type)
	})

	t.Run("check backward compatibility", func(t *testing.T) {
		gs := GCPStorage{logger: logger.NewLogger("test")}

		request := bindings.InvokeRequest{}
		request.Operation = bindings.CreateOperation
		request.Metadata = map[string]string{
			"name": "my_file.txt",
		}
		result := gs.handleBackwardCompatibilityForMetadata(request.Metadata)
		assert.NotEmpty(t, result["key"])
	})
}

func TestMergeWithRequestMetadata(t *testing.T) {
	t.Run("Has merged metadata", func(t *testing.T) {
		m := bindings.Metadata{}
		m.Properties = map[string]string{
			"auth_provider_x509_cert_url": "my_auth_provider_x509",
			"auth_uri":                    "my_auth_uri",
			"Bucket":                      "my_bucket",
			"client_x509_cert_url":        "my_client_x509",
			"client_email":                "my_email@mail.app",
			"client_id":                   "my_client_id",
			"private_key":                 "my_private_key",
			"private_key_id":              "my_private_key_id",
			"project_id":                  "my_project_id",
			"token_uri":                   "my_token_uri",
			"type":                        "my_type",
			"decodeBase64":                "false",
		}
		gs := GCPStorage{logger: logger.NewLogger("test")}
		meta, _, err := gs.parseMetadata(m)
		assert.Nil(t, err)

		assert.Equal(t, "my_auth_provider_x509", meta.AuthProviderCertURL)
		assert.Equal(t, "my_auth_uri", meta.AuthURI)
		assert.Equal(t, "my_bucket", meta.Bucket)
		assert.Equal(t, "my_client_x509", meta.ClientCertURL)
		assert.Equal(t, "my_email@mail.app", meta.ClientEmail)
		assert.Equal(t, "my_client_id", meta.ClientID)
		assert.Equal(t, "my_private_key", meta.PrivateKey)
		assert.Equal(t, "my_private_key_id", meta.PrivateKeyID)
		assert.Equal(t, "my_project_id", meta.ProjectID)
		assert.Equal(t, "my_token_uri", meta.TokenURI)
		assert.Equal(t, "my_type", meta.Type)
		assert.Equal(t, false, meta.DecodeBase64)

		request := bindings.InvokeRequest{}
		request.Metadata = map[string]string{
			"decodeBase64": "true",
		}

		mergedMeta, err := meta.mergeWithRequestMetadata(&request)

		assert.Nil(t, err)

		assert.Equal(t, "my_auth_provider_x509", mergedMeta.AuthProviderCertURL)
		assert.Equal(t, "my_auth_uri", mergedMeta.AuthURI)
		assert.Equal(t, "my_bucket", mergedMeta.Bucket)
		assert.Equal(t, "my_client_x509", mergedMeta.ClientCertURL)
		assert.Equal(t, "my_email@mail.app", mergedMeta.ClientEmail)
		assert.Equal(t, "my_client_id", mergedMeta.ClientID)
		assert.Equal(t, "my_private_key", mergedMeta.PrivateKey)
		assert.Equal(t, "my_private_key_id", mergedMeta.PrivateKeyID)
		assert.Equal(t, "my_project_id", mergedMeta.ProjectID)
		assert.Equal(t, "my_token_uri", mergedMeta.TokenURI)
		assert.Equal(t, "my_type", mergedMeta.Type)
		assert.Equal(t, true, mergedMeta.DecodeBase64)
	})

	t.Run("Has invalid merged metadata decodeBase64", func(t *testing.T) {
		m := bindings.Metadata{}
		m.Properties = map[string]string{
			"auth_provider_x509_cert_url": "my_auth_provider_x509",
			"auth_uri":                    "my_auth_uri",
			"Bucket":                      "my_bucket",
			"client_x509_cert_url":        "my_client_x509",
			"client_email":                "my_email@mail.app",
			"client_id":                   "my_client_id",
			"private_key":                 "my_private_key",
			"private_key_id":              "my_private_key_id",
			"project_id":                  "my_project_id",
			"token_uri":                   "my_token_uri",
			"type":                        "my_type",
			"decodeBase64":                "false",
		}
		gs := GCPStorage{logger: logger.NewLogger("test")}
		meta, _, err := gs.parseMetadata(m)
		assert.Nil(t, err)

		assert.Equal(t, "my_auth_provider_x509", meta.AuthProviderCertURL)
		assert.Equal(t, "my_auth_uri", meta.AuthURI)
		assert.Equal(t, "my_bucket", meta.Bucket)
		assert.Equal(t, "my_client_x509", meta.ClientCertURL)
		assert.Equal(t, "my_email@mail.app", meta.ClientEmail)
		assert.Equal(t, "my_client_id", meta.ClientID)
		assert.Equal(t, "my_private_key", meta.PrivateKey)
		assert.Equal(t, "my_private_key_id", meta.PrivateKeyID)
		assert.Equal(t, "my_project_id", meta.ProjectID)
		assert.Equal(t, "my_token_uri", meta.TokenURI)
		assert.Equal(t, "my_type", meta.Type)
		assert.Equal(t, false, meta.DecodeBase64)

		request := bindings.InvokeRequest{}
		request.Metadata = map[string]string{
			"decodeBase64": "hello",
		}

		mergedMeta, err := meta.mergeWithRequestMetadata(&request)

		assert.NotNil(t, err)
		assert.NotNil(t, mergedMeta)
	})
	t.Run("Has invalid merged metadata encodeBase64", func(t *testing.T) {
		m := bindings.Metadata{}
		m.Properties = map[string]string{
			"auth_provider_x509_cert_url": "my_auth_provider_x509",
			"auth_uri":                    "my_auth_uri",
			"Bucket":                      "my_bucket",
			"client_x509_cert_url":        "my_client_x509",
			"client_email":                "my_email@mail.app",
			"client_id":                   "my_client_id",
			"private_key":                 "my_private_key",
			"private_key_id":              "my_private_key_id",
			"project_id":                  "my_project_id",
			"token_uri":                   "my_token_uri",
			"type":                        "my_type",
			"decodeBase64":                "false",
			"encodeBase64":                "true",
		}
		gs := GCPStorage{logger: logger.NewLogger("test")}
		meta, _, err := gs.parseMetadata(m)
		assert.Nil(t, err)

		assert.Equal(t, "my_auth_provider_x509", meta.AuthProviderCertURL)
		assert.Equal(t, "my_auth_uri", meta.AuthURI)
		assert.Equal(t, "my_bucket", meta.Bucket)
		assert.Equal(t, "my_client_x509", meta.ClientCertURL)
		assert.Equal(t, "my_email@mail.app", meta.ClientEmail)
		assert.Equal(t, "my_client_id", meta.ClientID)
		assert.Equal(t, "my_private_key", meta.PrivateKey)
		assert.Equal(t, "my_private_key_id", meta.PrivateKeyID)
		assert.Equal(t, "my_project_id", meta.ProjectID)
		assert.Equal(t, "my_token_uri", meta.TokenURI)
		assert.Equal(t, "my_type", meta.Type)
		assert.Equal(t, false, meta.DecodeBase64)
		assert.Equal(t, true, meta.EncodeBase64)

		request := bindings.InvokeRequest{}
		request.Metadata = map[string]string{
			"encodeBase64": "hello",
		}

		mergedMeta, err := meta.mergeWithRequestMetadata(&request)

		assert.NotNil(t, err)
		assert.NotNil(t, mergedMeta)
	})
}

func TestGetOption(t *testing.T) {
	gs := GCPStorage{logger: logger.NewLogger("test")}
	gs.metadata = &gcpMetadata{}
	t.Run("return error if key is missing", func(t *testing.T) {
		r := bindings.InvokeRequest{}
		_, err := gs.get(&r)
		assert.Error(t, err)
	})
}

func TestDeleteOption(t *testing.T) {
	gs := GCPStorage{logger: logger.NewLogger("test")}
	gs.metadata = &gcpMetadata{}

	t.Run("return error if key is missing", func(t *testing.T) {
		r := bindings.InvokeRequest{}
		_, err := gs.delete(&r)
		assert.Error(t, err)
	})
}
