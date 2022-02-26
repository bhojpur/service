package sms

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
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/utils/logger"
)

const (
	toNumber      = "toNumber"
	fromNumber    = "fromNumber"
	accountSid    = "accountSid"
	authToken     = "authToken"
	timeout       = "timeout"
	twilioURLBase = "https://api.twilio.com/2010-04-01/Accounts/"
)

type SMS struct {
	metadata   twilioMetadata
	logger     logger.Logger
	httpClient *http.Client
}

type twilioMetadata struct {
	toNumber   string
	fromNumber string
	accountSid string
	authToken  string
	timeout    time.Duration
}

func NewSMS(logger logger.Logger) *SMS {
	return &SMS{
		logger:     logger,
		httpClient: &http.Client{},
	}
}

func (t *SMS) Init(metadata bindings.Metadata) error {
	twilioM := twilioMetadata{
		timeout: time.Minute * 5,
	}

	if metadata.Properties[fromNumber] == "" {
		return errors.New("\"fromNumber\" is a required field")
	}
	if metadata.Properties[accountSid] == "" {
		return errors.New("\"accountSid\" is a required field")
	}
	if metadata.Properties[authToken] == "" {
		return errors.New("\"authToken\" is a required field")
	}

	twilioM.toNumber = metadata.Properties[toNumber]
	twilioM.fromNumber = metadata.Properties[fromNumber]
	twilioM.accountSid = metadata.Properties[accountSid]
	twilioM.authToken = metadata.Properties[authToken]
	if metadata.Properties[timeout] != "" {
		t, err := time.ParseDuration(metadata.Properties[timeout])
		if err != nil {
			return fmt.Errorf("error parsing timeout: %s", err)
		}
		twilioM.timeout = t
	}

	t.metadata = twilioM
	t.httpClient.Timeout = twilioM.timeout

	return nil
}

func (t *SMS) Operations() []bindings.OperationKind {
	return []bindings.OperationKind{bindings.CreateOperation}
}

func (t *SMS) Invoke(req *bindings.InvokeRequest) (*bindings.InvokeResponse, error) {
	toNumberValue := t.metadata.toNumber
	if toNumberValue == "" {
		toNumberFromRequest, ok := req.Metadata[toNumber]
		if !ok || toNumberFromRequest == "" {
			return nil, errors.New("twilio missing \"toNumber\" field")
		}
		toNumberValue = toNumberFromRequest
	}

	v := url.Values{}
	v.Set("To", toNumberValue)
	v.Set("From", t.metadata.fromNumber)
	v.Set("Body", string(req.Data))
	vDr := *strings.NewReader(v.Encode())

	twilioURL := fmt.Sprintf("%s%s/Messages.json", twilioURLBase, t.metadata.accountSid)
	httpReq, err := http.NewRequestWithContext(context.Background(), "POST", twilioURL, &vDr)
	if err != nil {
		return nil, err
	}
	httpReq.SetBasicAuth(t.metadata.accountSid, t.metadata.authToken)
	httpReq.Header.Add("Accept", "application/json")
	httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := t.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return nil, fmt.Errorf("error from Twilio: %s", resp.Status)
	}

	return nil, nil
}
