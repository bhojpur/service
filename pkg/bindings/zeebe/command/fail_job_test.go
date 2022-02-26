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
	"testing"

	"github.com/camunda-cloud/zeebe/clients/go/pkg/commands"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/pb"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/zbc"
	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/utils/logger"
)

type mockFailJobClient struct {
	zbc.Client
	cmd1 *mockFailJobCommandStep1
}

type mockFailJobCommandStep1 struct {
	commands.FailJobCommandStep1
	cmd2   *mockFailJobCommandStep2
	jobKey int64
}

type mockFailJobCommandStep2 struct {
	commands.FailJobCommandStep2
	cmd3    *mockFailJobCommandStep3
	retries int32
}

type mockFailJobCommandStep3 struct {
	commands.FailJobCommandStep3
	errorMessage string
}

func (mc *mockFailJobClient) NewFailJobCommand() commands.FailJobCommandStep1 {
	mc.cmd1 = &mockFailJobCommandStep1{
		cmd2: &mockFailJobCommandStep2{
			cmd3: &mockFailJobCommandStep3{},
		},
	}

	return mc.cmd1
}

func (cmd1 *mockFailJobCommandStep1) JobKey(jobKey int64) commands.FailJobCommandStep2 {
	cmd1.jobKey = jobKey

	return cmd1.cmd2
}

func (cmd2 *mockFailJobCommandStep2) Retries(retries int32) commands.FailJobCommandStep3 {
	cmd2.retries = retries

	return cmd2.cmd3
}

func (cmd3 *mockFailJobCommandStep3) ErrorMessage(errorMessage string) commands.FailJobCommandStep3 {
	cmd3.errorMessage = errorMessage

	return cmd3
}

func (cmd3 *mockFailJobCommandStep3) Send(context.Context) (*pb.FailJobResponse, error) {
	return &pb.FailJobResponse{}, nil
}

func TestFailJob(t *testing.T) {
	testLogger := logger.NewLogger("test")

	t.Run("jobKey is mandatory", func(t *testing.T) {
		cmd := ZeebeCommand{logger: testLogger}
		req := &bindings.InvokeRequest{Operation: FailJobOperation}
		_, err := cmd.Invoke(req)
		assert.Error(t, err, ErrMissingJobKey)
	})

	t.Run("retries is mandatory", func(t *testing.T) {
		payload := failJobPayload{
			JobKey: new(int64),
		}
		data, err := json.Marshal(payload)
		assert.NoError(t, err)

		cmd := ZeebeCommand{logger: testLogger}
		req := &bindings.InvokeRequest{Data: data, Operation: FailJobOperation}
		_, err = cmd.Invoke(req)
		assert.Error(t, err, ErrMissingRetries)
	})

	t.Run("fail a job", func(t *testing.T) {
		payload := failJobPayload{
			JobKey:       new(int64),
			Retries:      new(int32),
			ErrorMessage: "a",
		}
		data, err := json.Marshal(payload)
		assert.NoError(t, err)

		req := &bindings.InvokeRequest{Data: data, Operation: FailJobOperation}

		var mc mockFailJobClient

		cmd := ZeebeCommand{logger: testLogger, client: &mc}
		_, err = cmd.Invoke(req)
		assert.NoError(t, err)

		assert.Equal(t, *payload.JobKey, mc.cmd1.jobKey)
		assert.Equal(t, *payload.Retries, mc.cmd1.cmd2.retries)
		assert.Equal(t, payload.ErrorMessage, mc.cmd1.cmd2.cmd3.errorMessage)
	})
}
