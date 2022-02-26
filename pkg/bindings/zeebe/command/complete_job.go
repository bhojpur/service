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
	"fmt"

	"github.com/camunda-cloud/zeebe/clients/go/pkg/commands"

	"github.com/bhojpur/service/pkg/bindings"
)

type completeJobPayload struct {
	JobKey    *int64      `json:"jobKey"`
	Variables interface{} `json:"variables"`
}

func (z *ZeebeCommand) completeJob(req *bindings.InvokeRequest) (*bindings.InvokeResponse, error) {
	var payload completeJobPayload
	err := json.Unmarshal(req.Data, &payload)
	if err != nil {
		return nil, err
	}

	if payload.JobKey == nil {
		return nil, ErrMissingJobKey
	}

	cmd1 := z.client.NewCompleteJobCommand()
	cmd2 := cmd1.JobKey(*payload.JobKey)
	var cmd3 commands.DispatchCompleteJobCommand = cmd2
	if payload.Variables != nil {
		cmd3, err = cmd2.VariablesFromObject(payload.Variables)
		if err != nil {
			return nil, err
		}
	}

	_, err = cmd3.Send(context.Background())
	if err != nil {
		return nil, fmt.Errorf("cannot complete job for key %d: %w", payload.JobKey, err)
	}

	return &bindings.InvokeResponse{}, nil
}
