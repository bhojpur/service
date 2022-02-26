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

type mockCompleteJobClient struct {
	zbc.Client
	cmd1 *mockCompleteJobCommandStep1
}

type mockCompleteJobCommandStep1 struct {
	commands.CompleteJobCommandStep1
	cmd2   *mockCompleteJobCommandStep2
	jobKey int64
}

type mockCompleteJobCommandStep2 struct {
	commands.CompleteJobCommandStep2
	cmd3      *mockDispatchCompleteJobCommand
	variables interface{}
}

type mockDispatchCompleteJobCommand struct {
	commands.DispatchCompleteJobCommand
}

func (mc *mockCompleteJobClient) NewCompleteJobCommand() commands.CompleteJobCommandStep1 {
	mc.cmd1 = &mockCompleteJobCommandStep1{
		cmd2: &mockCompleteJobCommandStep2{
			cmd3: &mockDispatchCompleteJobCommand{},
		},
	}

	return mc.cmd1
}

func (cmd1 *mockCompleteJobCommandStep1) JobKey(jobKey int64) commands.CompleteJobCommandStep2 {
	cmd1.jobKey = jobKey

	return cmd1.cmd2
}

func (cmd2 *mockCompleteJobCommandStep2) VariablesFromObject(variables interface{}) (commands.DispatchCompleteJobCommand, error) {
	cmd2.variables = variables

	return cmd2.cmd3, nil
}

func (cmd3 *mockDispatchCompleteJobCommand) Send(context.Context) (*pb.CompleteJobResponse, error) {
	return &pb.CompleteJobResponse{}, nil
}

func TestCompleteJob(t *testing.T) {
	testLogger := logger.NewLogger("test")

	t.Run("elementInstanceKey is mandatory", func(t *testing.T) {
		cmd := ZeebeCommand{logger: testLogger}
		req := &bindings.InvokeRequest{Operation: CompleteJobOperation}
		_, err := cmd.Invoke(req)
		assert.Error(t, err, ErrMissingJobKey)
	})

	t.Run("complete a job", func(t *testing.T) {
		payload := completeJobPayload{
			JobKey: new(int64),
			Variables: map[string]interface{}{
				"key": "value",
			},
		}
		data, err := json.Marshal(payload)
		assert.NoError(t, err)

		req := &bindings.InvokeRequest{Data: data, Operation: CompleteJobOperation}

		var mc mockCompleteJobClient

		cmd := ZeebeCommand{logger: testLogger, client: &mc}
		_, err = cmd.Invoke(req)
		assert.NoError(t, err)

		assert.Equal(t, *payload.JobKey, mc.cmd1.jobKey)
		assert.Equal(t, payload.Variables, mc.cmd1.cmd2.variables)
	})
}
