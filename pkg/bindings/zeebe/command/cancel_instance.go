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

	"github.com/bhojpur/service/pkg/bindings"
)

var ErrMissingProcessInstanceKey = errors.New("processInstanceKey is a required attribute")

type cancelInstancePayload struct {
	ProcessInstanceKey *int64 `json:"processInstanceKey"`
}

func (z *ZeebeCommand) cancelInstance(req *bindings.InvokeRequest) (*bindings.InvokeResponse, error) {
	var payload cancelInstancePayload
	err := json.Unmarshal(req.Data, &payload)
	if err != nil {
		return nil, err
	}

	if payload.ProcessInstanceKey == nil {
		return nil, ErrMissingProcessInstanceKey
	}

	_, err = z.client.NewCancelInstanceCommand().
		ProcessInstanceKey(*payload.ProcessInstanceKey).
		Send(context.Background())
	if err != nil {
		return nil, fmt.Errorf("cannot cancel instance for process instance key %d: %w", payload.ProcessInstanceKey, err)
	}

	return &bindings.InvokeResponse{}, nil
}
