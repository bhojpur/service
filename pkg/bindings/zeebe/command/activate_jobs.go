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

var (
	ErrMissingJobType           = errors.New("jobType is a required attribute")
	ErrMissingMaxJobsToActivate = errors.New("maxJobsToActivate is a required attribute")
)

type activateJobsPayload struct {
	JobType           string            `json:"jobType"`
	MaxJobsToActivate *int32            `json:"maxJobsToActivate"`
	Timeout           metadata.Duration `json:"timeout"`
	WorkerName        string            `json:"workerName"`
	FetchVariables    []string          `json:"fetchVariables"`
}

func (z *ZeebeCommand) activateJobs(req *bindings.InvokeRequest) (*bindings.InvokeResponse, error) {
	var payload activateJobsPayload
	err := json.Unmarshal(req.Data, &payload)
	if err != nil {
		return nil, err
	}

	if payload.JobType == "" {
		return nil, ErrMissingJobType
	}

	if payload.MaxJobsToActivate == nil {
		return nil, ErrMissingMaxJobsToActivate
	}

	cmd := z.client.NewActivateJobsCommand().
		JobType(payload.JobType).
		MaxJobsToActivate(*payload.MaxJobsToActivate)

	if payload.Timeout.Duration != time.Duration(0) {
		cmd = cmd.Timeout(payload.Timeout.Duration)
	}

	if payload.WorkerName != "" {
		cmd = cmd.WorkerName(payload.WorkerName)
	}

	if payload.FetchVariables != nil {
		cmd = cmd.FetchVariables(payload.FetchVariables...)
	}

	ctx := context.Background()
	response, err := cmd.Send(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot activate jobs for type %s: %w", payload.JobType, err)
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal response to json: %w", err)
	}

	return &bindings.InvokeResponse{
		Data: jsonResponse,
	}, nil
}
