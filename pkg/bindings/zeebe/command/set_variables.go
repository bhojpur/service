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

var (
	ErrMissingElementInstanceKey = errors.New("elementInstanceKey is a required attribute")
	ErrMissingVariables          = errors.New("variables is a required attribute")
)

type setVariablesPayload struct {
	ElementInstanceKey *int64      `json:"elementInstanceKey"`
	Local              bool        `json:"local"`
	Variables          interface{} `json:"variables"`
}

func (z *ZeebeCommand) setVariables(req *bindings.InvokeRequest) (*bindings.InvokeResponse, error) {
	var payload setVariablesPayload
	err := json.Unmarshal(req.Data, &payload)
	if err != nil {
		return nil, err
	}

	if payload.ElementInstanceKey == nil {
		return nil, ErrMissingElementInstanceKey
	}

	if payload.Variables == nil {
		return nil, ErrMissingVariables
	}

	cmd, err := z.client.NewSetVariablesCommand().
		ElementInstanceKey(*payload.ElementInstanceKey).
		VariablesFromObject(payload.Variables)
	if err != nil {
		return nil, err
	}

	response, err := cmd.Local(payload.Local).Send(context.Background())
	if err != nil {
		return nil, fmt.Errorf("cannot set variables for element instance key %d: %w", payload.ElementInstanceKey, err)
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal response to json: %w", err)
	}

	return &bindings.InvokeResponse{
		Data: jsonResponse,
	}, nil
}
