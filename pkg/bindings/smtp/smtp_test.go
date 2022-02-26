package smtp

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

	t.Run("Has correct metadata (default priority)", func(t *testing.T) {
		m := bindings.Metadata{}
		m.Properties = map[string]string{
			"host":          "mailserver.bhojpur.net",
			"port":          "25",
			"user":          "user@bhojpur.net",
			"password":      "P@$$w0rd!",
			"skipTLSVerify": "true",
			"emailFrom":     "from@bhojpur.net",
			"emailTo":       "to@bhojpur.net",
			"emailCC":       "cc@bhojpur.net",
			"emailBCC":      "bcc@bhojpur.net",
			"subject":       "Test email",
		}
		r := Mailer{logger: logger}
		smtpMeta, err := r.parseMetadata(m)
		assert.Nil(t, err)
		assert.Equal(t, "mailserver.bhojpur.net", smtpMeta.Host)
		assert.Equal(t, 25, smtpMeta.Port)
		assert.Equal(t, "user@bhojpur.net", smtpMeta.User)
		assert.Equal(t, "P@$$w0rd!", smtpMeta.Password)
		assert.Equal(t, true, smtpMeta.SkipTLSVerify)
		assert.Equal(t, "from@bhojpur.net", smtpMeta.EmailFrom)
		assert.Equal(t, "to@bhojpur.net", smtpMeta.EmailTo)
		assert.Equal(t, "cc@bhojpur.net", smtpMeta.EmailCC)
		assert.Equal(t, "bcc@bhojpur.net", smtpMeta.EmailBCC)
		assert.Equal(t, "Test email", smtpMeta.Subject)
		assert.Equal(t, 3, smtpMeta.Priority)
	})
	t.Run("Has correct metadata (no default value for priority)", func(t *testing.T) {
		m := bindings.Metadata{}
		m.Properties = map[string]string{
			"host":          "mailserver.bhojpur.net",
			"port":          "25",
			"user":          "user@bhojpur.net",
			"password":      "P@$$w0rd!",
			"skipTLSVerify": "true",
			"emailFrom":     "from@bhojpur.net",
			"emailTo":       "to@bhojpur.net",
			"emailCC":       "cc@bhojpur.net",
			"emailBCC":      "bcc@bhojpur.net",
			"subject":       "Test email",
			"priority":      "1",
		}
		r := Mailer{logger: logger}
		smtpMeta, err := r.parseMetadata(m)
		assert.Nil(t, err)
		assert.Equal(t, "mailserver.bhojpur.net", smtpMeta.Host)
		assert.Equal(t, 25, smtpMeta.Port)
		assert.Equal(t, "user@bhojpur.net", smtpMeta.User)
		assert.Equal(t, "P@$$w0rd!", smtpMeta.Password)
		assert.Equal(t, true, smtpMeta.SkipTLSVerify)
		assert.Equal(t, "from@bhojpur.net", smtpMeta.EmailFrom)
		assert.Equal(t, "to@bhojpur.net", smtpMeta.EmailTo)
		assert.Equal(t, "cc@bhojpur.net", smtpMeta.EmailCC)
		assert.Equal(t, "bcc@bhojpur.net", smtpMeta.EmailBCC)
		assert.Equal(t, "Test email", smtpMeta.Subject)
		assert.Equal(t, 1, smtpMeta.Priority)
	})
	t.Run("Incorrrect  metadata (invalid priority)", func(t *testing.T) {
		m := bindings.Metadata{}
		m.Properties = map[string]string{
			"host":          "mailserver.bhojpur.net",
			"port":          "25",
			"user":          "user@bhojpur.net",
			"password":      "P@$$w0rd!",
			"skipTLSVerify": "true",
			"emailFrom":     "from@bhojpur.net",
			"emailTo":       "to@bhojpur.net",
			"emailCC":       "cc@bhojpur.net",
			"emailBCC":      "bcc@bhojpur.net",
			"subject":       "Test email",
			"priority":      "0",
		}
		r := Mailer{logger: logger}
		smtpMeta, err := r.parseMetadata(m)
		assert.NotNil(t, smtpMeta)
		assert.NotNil(t, err)
	})
	t.Run("Incorrrect  metadata (user, no password)", func(t *testing.T) {
		m := bindings.Metadata{}
		m.Properties = map[string]string{
			"host":          "mailserver.bhojpur.net",
			"port":          "25",
			"user":          "user@bhojpur.net",
			"skipTLSVerify": "true",
			"emailFrom":     "from@bhojpur.net",
			"emailTo":       "to@bhojpur.net",
			"emailCC":       "cc@bhojpur.net",
			"emailBCC":      "bcc@bhojpur.net",
			"subject":       "Test email",
			"priority":      "0",
		}
		r := Mailer{logger: logger}
		smtpMeta, err := r.parseMetadata(m)
		assert.NotNil(t, smtpMeta)
		assert.NotNil(t, err)
	})
	t.Run("Incorrrect  metadata (no user, password)", func(t *testing.T) {
		m := bindings.Metadata{}
		m.Properties = map[string]string{
			"host":          "mailserver.bhojpur.net",
			"port":          "25",
			"password":      "P@$$w0rd!",
			"skipTLSVerify": "true",
			"emailFrom":     "from@bhojpur.net",
			"emailTo":       "to@bhojpur.net",
			"emailCC":       "cc@bhojpur.net",
			"emailBCC":      "bcc@bhojpur.net",
			"subject":       "Test email",
			"priority":      "0",
		}
		r := Mailer{logger: logger}
		smtpMeta, err := r.parseMetadata(m)
		assert.NotNil(t, smtpMeta)
		assert.NotNil(t, err)
	})
}

func TestMergeWithRequestMetadata(t *testing.T) {
	t.Run("Has merged metadata", func(t *testing.T) {
		smtpMeta := Metadata{
			Host:          "mailserver.bhojpur.net",
			Port:          25,
			User:          "user@bhojpur.net",
			SkipTLSVerify: true,
			Password:      "P@$$w0rd!",
			EmailFrom:     "from@bhojpur.net",
			EmailTo:       "to@bhojpur.net",
			EmailCC:       "cc@bhojpur.net",
			EmailBCC:      "bcc@bhojpur.net",
			Subject:       "Test email",
		}

		request := bindings.InvokeRequest{}
		request.Metadata = map[string]string{
			"emailFrom": "req-from@bhojpur.net",
			"emailTo":   "req-to@bhojpur.net",
			"emailCC":   "req-cc@bhojpur.net",
			"emailBCC":  "req-bcc@bhojpur.net",
			"subject":   "req-Test email",
			"priority":  "1",
		}

		mergedMeta, err := smtpMeta.mergeWithRequestMetadata(&request)

		assert.Nil(t, err)

		assert.Equal(t, "mailserver.bhojpur.net", mergedMeta.Host)
		assert.Equal(t, 25, mergedMeta.Port)
		assert.Equal(t, "user@bhojpur.net", mergedMeta.User)
		assert.Equal(t, "P@$$w0rd!", mergedMeta.Password)
		assert.Equal(t, true, mergedMeta.SkipTLSVerify)
		assert.Equal(t, "req-from@bhojpur.net", mergedMeta.EmailFrom)
		assert.Equal(t, "req-to@bhojpur.net", mergedMeta.EmailTo)
		assert.Equal(t, "req-cc@bhojpur.net", mergedMeta.EmailCC)
		assert.Equal(t, "req-bcc@bhojpur.net", mergedMeta.EmailBCC)
		assert.Equal(t, "req-Test email", mergedMeta.Subject)
		assert.Equal(t, 1, mergedMeta.Priority)
	})
}

func TestMergeWithNoRequestMetadata(t *testing.T) {
	t.Run("Has no merged metadata", func(t *testing.T) {
		smtpMeta := Metadata{
			Host:          "mailserver.bhojpur.net",
			Port:          25,
			User:          "user@bhojpur.net",
			SkipTLSVerify: true,
			Password:      "P@$$w0rd!",
			EmailFrom:     "from@bhojpur.net",
			EmailTo:       "to@bhojpur.net",
			EmailCC:       "cc@bhojpur.net",
			EmailBCC:      "bcc@bhojpur.net",
			Subject:       "Test email",
			Priority:      1,
		}

		request := bindings.InvokeRequest{}
		request.Metadata = map[string]string{}

		mergedMeta, err := smtpMeta.mergeWithRequestMetadata(&request)

		assert.Nil(t, err)
		assert.Equal(t, "mailserver.bhojpur.net", mergedMeta.Host)
		assert.Equal(t, 25, mergedMeta.Port)
		assert.Equal(t, "user@bhojpur.net", mergedMeta.User)
		assert.Equal(t, "P@$$w0rd!", mergedMeta.Password)
		assert.Equal(t, true, mergedMeta.SkipTLSVerify)
		assert.Equal(t, "from@bhojpur.net", mergedMeta.EmailFrom)
		assert.Equal(t, "to@bhojpur.net", mergedMeta.EmailTo)
		assert.Equal(t, "cc@bhojpur.net", mergedMeta.EmailCC)
		assert.Equal(t, "bcc@bhojpur.net", mergedMeta.EmailBCC)
		assert.Equal(t, "Test email", mergedMeta.Subject)
		assert.Equal(t, 1, mergedMeta.Priority)
	})
}

func TestMergeWithRequestMetadata_invalidPriorityTooHigh(t *testing.T) {
	t.Run("Has merged metadata", func(t *testing.T) {
		smtpMeta := Metadata{
			Host:          "mailserver.bhojpur.net",
			Port:          25,
			User:          "user@bhojpur.net",
			SkipTLSVerify: true,
			Password:      "P@$$w0rd!",
			EmailFrom:     "from@bhojpur.net",
			EmailTo:       "to@bhojpur.net",
			EmailCC:       "cc@bhojpur.net",
			EmailBCC:      "bcc@bhojpur.net",
			Subject:       "Test email",
			Priority:      2,
		}

		request := bindings.InvokeRequest{}
		request.Metadata = map[string]string{
			"emailFrom": "req-from@bhojpur.net",
			"emailTo":   "req-to@bhojpur.net",
			"emailCC":   "req-cc@bhojpur.net",
			"emailBCC":  "req-bcc@bhojpur.net",
			"subject":   "req-Test email",
			"priority":  "6",
		}

		mergedMeta, err := smtpMeta.mergeWithRequestMetadata(&request)

		assert.NotNil(t, mergedMeta)
		assert.NotNil(t, err)
	})
}

func TestMergeWithRequestMetadata_invalidPriorityTooLow(t *testing.T) {
	t.Run("Has merged metadata", func(t *testing.T) {
		smtpMeta := Metadata{
			Host:          "mailserver.bhojpur.net",
			Port:          25,
			User:          "user@bhojpur.net",
			SkipTLSVerify: true,
			Password:      "P@$$w0rd!",
			EmailFrom:     "from@bhojpur.net",
			EmailTo:       "to@bhojpur.net",
			EmailCC:       "cc@bhojpur.net",
			EmailBCC:      "bcc@bhojpur.net",
			Subject:       "Test email",
			Priority:      2,
		}

		request := bindings.InvokeRequest{}
		request.Metadata = map[string]string{
			"emailFrom": "req-from@bhojpur.net",
			"emailTo":   "req-to@bhojpur.net",
			"emailCC":   "req-cc@bhojpur.net",
			"emailBCC":  "req-bcc@bhojpur.net",
			"subject":   "req-Test email",
			"priority":  "0",
		}

		mergedMeta, err := smtpMeta.mergeWithRequestMetadata(&request)

		assert.NotNil(t, mergedMeta)
		assert.NotNil(t, err)
	})
}

func TestMergeWithRequestMetadata_invalidPriorityNotNumber(t *testing.T) {
	t.Run("Has merged metadata", func(t *testing.T) {
		smtpMeta := Metadata{
			Host:          "mailserver.bhojpur.net",
			Port:          25,
			User:          "user@bhojpur.net",
			SkipTLSVerify: true,
			Password:      "P@$$w0rd!",
			EmailFrom:     "from@bhojpur.net",
			EmailTo:       "to@bhojpur.net",
			EmailCC:       "cc@bhojpur.net",
			EmailBCC:      "bcc@bhojpur.net",
			Subject:       "Test email",
		}

		request := bindings.InvokeRequest{}
		request.Metadata = map[string]string{
			"emailFrom": "req-from@bhojpur.net",
			"emailTo":   "req-to@bhojpur.net",
			"emailCC":   "req-cc@bhojpur.net",
			"emailBCC":  "req-bcc@bhojpur.net",
			"subject":   "req-Test email",
			"priority":  "NoNumber",
		}

		mergedMeta, err := smtpMeta.mergeWithRequestMetadata(&request)

		assert.NotNil(t, mergedMeta)
		assert.NotNil(t, err)
	})
}
