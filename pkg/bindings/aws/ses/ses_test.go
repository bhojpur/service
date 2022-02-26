package ses

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
	logger := logger.NewLogger("test")

	t.Run("Has correct metadata", func(t *testing.T) {
		m := bindings.Metadata{}
		m.Properties = map[string]string{
			"region":       "myRegionForSES",
			"accessKey":    "myAccessKeyForSES",
			"secretKey":    "mySecretKeyForSES",
			"sessionToken": "mySessionToken",
			"emailFrom":    "from@bhojpur.net",
			"emailTo":      "to@bhojpur.net",
			"emailCc":      "cc@bhojpur.net",
			"emailBcc":     "bcc@bhojpur.net",
			"subject":      "Test email",
		}
		r := AWSSES{logger: logger}
		smtpMeta, err := r.parseMetadata(m)
		assert.Nil(t, err)
		assert.Equal(t, "myRegionForSES", smtpMeta.Region)
		assert.Equal(t, "myAccessKeyForSES", smtpMeta.AccessKey)
		assert.Equal(t, "mySecretKeyForSES", smtpMeta.SecretKey)
		assert.Equal(t, "mySessionToken", smtpMeta.SessionToken)
		assert.Equal(t, "from@bhojpur.net", smtpMeta.EmailFrom)
		assert.Equal(t, "to@bhojpur.net", smtpMeta.EmailTo)
		assert.Equal(t, "cc@bhojpur.net", smtpMeta.EmailCc)
		assert.Equal(t, "bcc@bhojpur.net", smtpMeta.EmailBcc)
		assert.Equal(t, "Test email", smtpMeta.Subject)
	})

	t.Run("region is required", func(t *testing.T) {
		m := bindings.Metadata{}
		m.Properties = map[string]string{
			"accessKey": "myAccessKeyForSES",
			"secretKey": "mySecretKeyForSES",
			"emailFrom": "from@bhojpur.net",
			"emailTo":   "to@bhojpur.net",
			"emailCc":   "cc@bhojpur.net",
			"emailBcc":  "bcc@bhojpur.net",
			"subject":   "Test email",
		}
		r := AWSSES{logger: logger}
		_, err := r.parseMetadata(m)
		assert.Error(t, err)
	})

	t.Run("accessKey is required", func(t *testing.T) {
		m := bindings.Metadata{}
		m.Properties = map[string]string{
			"region":    "myRegionForSES",
			"secretKey": "mySecretKeyForSES",
			"emailFrom": "from@bhojpur.net",
			"emailTo":   "to@bhojpur.net",
			"emailCc":   "cc@bhojpur.net",
			"emailBcc":  "bcc@bhojpur.net",
			"subject":   "Test email",
		}
		r := AWSSES{logger: logger}
		_, err := r.parseMetadata(m)
		assert.Error(t, err)
	})

	t.Run("secretKey is required", func(t *testing.T) {
		m := bindings.Metadata{}
		m.Properties = map[string]string{
			"region":    "myRegionForSES",
			"accessKey": "myAccessKeyForSES",
			"emailFrom": "from@bhojpur.net",
			"emailTo":   "to@bhojpur.net",
			"emailCc":   "cc@bhojpur.net",
			"emailBcc":  "bcc@bhojpur.net",
			"subject":   "Test email",
		}
		r := AWSSES{logger: logger}
		_, err := r.parseMetadata(m)
		assert.Error(t, err)
	})
}

func TestMergeWithRequestMetadata(t *testing.T) {
	t.Run("Has merged metadata", func(t *testing.T) {
		sesMeta := sesMetadata{
			Region:    "myRegionForSES",
			AccessKey: "myAccessKeyForSES",
			SecretKey: "mySecretKeyForSES",
			EmailFrom: "from@bhojpur.net",
			EmailTo:   "to@bhojpur.net",
			EmailCc:   "cc@bhojpur.net",
			EmailBcc:  "bcc@bhojpur.net",
			Subject:   "Test email",
		}

		request := bindings.InvokeRequest{}
		request.Metadata = map[string]string{
			"emailFrom": "req-from@bhojpur.net",
			"emailTo":   "req-to@bhojpur.net",
			"emailCc":   "req-cc@bhojpur.net",
			"emailBcc":  "req-bcc@bhojpur.net",
			"subject":   "req-Test email",
		}

		mergedMeta := sesMeta.mergeWithRequestMetadata(&request)

		assert.Equal(t, "myRegionForSES", mergedMeta.Region)
		assert.Equal(t, "myAccessKeyForSES", mergedMeta.AccessKey)
		assert.Equal(t, "mySecretKeyForSES", mergedMeta.SecretKey)
		assert.Equal(t, "req-from@bhojpur.net", mergedMeta.EmailFrom)
		assert.Equal(t, "req-to@bhojpur.net", mergedMeta.EmailTo)
		assert.Equal(t, "req-cc@bhojpur.net", mergedMeta.EmailCc)
		assert.Equal(t, "req-bcc@bhojpur.net", mergedMeta.EmailBcc)
		assert.Equal(t, "req-Test email", mergedMeta.Subject)
	})

	t.Run("Has no merged metadata", func(t *testing.T) {
		sesMeta := sesMetadata{
			Region:    "myRegionForSES",
			AccessKey: "myAccessKeyForSES",
			SecretKey: "mySecretKeyForSES",
			EmailFrom: "from@bhojpur.net",
			EmailTo:   "to@bhojpur.net",
			EmailCc:   "cc@bhojpur.net",
			EmailBcc:  "bcc@bhojpur.net",
			Subject:   "Test email",
		}

		request := bindings.InvokeRequest{}
		request.Metadata = map[string]string{}

		mergedMeta := sesMeta.mergeWithRequestMetadata(&request)

		assert.Equal(t, "myRegionForSES", mergedMeta.Region)
		assert.Equal(t, "myAccessKeyForSES", mergedMeta.AccessKey)
		assert.Equal(t, "mySecretKeyForSES", mergedMeta.SecretKey)
		assert.Equal(t, "from@bhojpur.net", mergedMeta.EmailFrom)
		assert.Equal(t, "to@bhojpur.net", mergedMeta.EmailTo)
		assert.Equal(t, "cc@bhojpur.net", mergedMeta.EmailCc)
		assert.Equal(t, "bcc@bhojpur.net", mergedMeta.EmailBcc)
		assert.Equal(t, "Test email", mergedMeta.Subject)
	})
}
