package command

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
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/metadata"
)

var ErrMissingMessageName = errors.New("messageName is a required attribute")

type publishMessagePayload struct {
	MessageName    string            `json:"messageName"`
	CorrelationKey string            `json:"correlationKey"`
	MessageID      string            `json:"messageId"`
	TimeToLive     metadata.Duration `json:"timeToLive"`
	Variables      interface{}       `json:"variables"`
}

func (z *ZeebeCommand) publishMessage(req *bindings.InvokeRequest) (*bindings.InvokeResponse, error) {
	var payload publishMessagePayload
	err := json.Unmarshal(req.Data, &payload)
	if err != nil {
		return nil, err
	}

	if payload.MessageName == "" {
		return nil, ErrMissingMessageName
	}

	cmd := z.client.NewPublishMessageCommand().
		MessageName(payload.MessageName).
		CorrelationKey(payload.CorrelationKey)

	if payload.MessageID != "" {
		cmd = cmd.MessageId(payload.MessageID)
	}

	if payload.TimeToLive.Duration != time.Duration(0) {
		cmd = cmd.TimeToLive(payload.TimeToLive.Duration)
	}

	if payload.Variables != nil {
		cmd, err = cmd.VariablesFromObject(payload.Variables)
		if err != nil {
			return nil, err
		}
	}

	ctx := context.Background()
	response, err := cmd.Send(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot publish message with name %s: %w", payload.MessageName, err)
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal response to json: %w", err)
	}

	return &bindings.InvokeResponse{
		Data: jsonResponse,
	}, nil
}
