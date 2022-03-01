package firestore

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

	"github.com/bhojpur/service/pkg/state"
)

func TestGetFirestoreMetadata(t *testing.T) {
	t.Run("With correct properties", func(t *testing.T) {
		properties := map[string]string{
			"type":                        "service_account",
			"project_id":                  "myprojectid",
			"private_key_id":              "123",
			"private_key":                 "mykey",
			"client_email":                "me@123.iam.gserviceaccount.com",
			"client_id":                   "456",
			"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
			"token_uri":                   "https://oauth2.googleapis.com/token",
			"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
			"client_x509_cert_url":        "https://www.googleapis.com/robot/v1/metadata/x509/x",
		}
		m := state.Metadata{
			Properties: properties,
		}
		metadata, err := getFirestoreMetadata(m)
		assert.Nil(t, err)
		assert.Equal(t, "service_account", metadata.Type)
		assert.Equal(t, "myprojectid", metadata.ProjectID)
		assert.Equal(t, "123", metadata.PrivateKeyID)
		assert.Equal(t, "mykey", metadata.PrivateKey)
		assert.Equal(t, defaultEntityKind, metadata.EntityKind)
	})

	t.Run("With incorrect properties", func(t *testing.T) {
		properties := map[string]string{
			"type":           "service_account",
			"project_id":     "myprojectid",
			"private_key_id": "123",
			"private_key":    "mykey",
		}
		m := state.Metadata{
			Properties: properties,
		}
		_, err := getFirestoreMetadata(m)
		assert.NotNil(t, err)
	})
}
