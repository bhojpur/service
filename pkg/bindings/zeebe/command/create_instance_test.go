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

type mockCreateInstanceClient struct {
	zbc.Client
	cmd1 *mockCreateInstanceCommandStep1
}

type mockCreateInstanceCommandStep1 struct {
	commands.CreateInstanceCommandStep1
	cmd2                 *mockCreateInstanceCommandStep2
	bpmnProcessID        string
	processDefinitionKey int64
}

type mockCreateInstanceCommandStep2 struct {
	commands.CreateInstanceCommandStep2
	cmd3          *mockCreateInstanceCommandStep3
	version       int32
	latestVersion bool
}

type mockCreateInstanceCommandStep3 struct {
	commands.CreateInstanceCommandStep3
	variables interface{}
}

func (mc *mockCreateInstanceClient) NewCreateInstanceCommand() commands.CreateInstanceCommandStep1 {
	mc.cmd1 = &mockCreateInstanceCommandStep1{
		cmd2: &mockCreateInstanceCommandStep2{
			cmd3: &mockCreateInstanceCommandStep3{},
		},
	}

	return mc.cmd1
}

//nolint // BPMNProcessId comes from the Zeebe client API and cannot be written as BPMNProcessID
func (cmd1 *mockCreateInstanceCommandStep1) BPMNProcessId(bpmnProcessID string) commands.CreateInstanceCommandStep2 {
	cmd1.bpmnProcessID = bpmnProcessID

	return cmd1.cmd2
}

func (cmd1 *mockCreateInstanceCommandStep1) ProcessDefinitionKey(processDefinitionKey int64) commands.CreateInstanceCommandStep3 {
	cmd1.processDefinitionKey = processDefinitionKey

	return cmd1.cmd2.cmd3
}

func (cmd2 *mockCreateInstanceCommandStep2) Version(version int32) commands.CreateInstanceCommandStep3 {
	cmd2.version = version

	return cmd2.cmd3
}

func (cmd2 *mockCreateInstanceCommandStep2) LatestVersion() commands.CreateInstanceCommandStep3 {
	cmd2.latestVersion = true

	return cmd2.cmd3
}

func (cmd3 *mockCreateInstanceCommandStep3) VariablesFromObject(variables interface{}) (commands.CreateInstanceCommandStep3, error) {
	cmd3.variables = variables

	return cmd3, nil
}

func (cmd3 *mockCreateInstanceCommandStep3) Send(context.Context) (*pb.CreateProcessInstanceResponse, error) {
	return &pb.CreateProcessInstanceResponse{}, nil
}

func TestCreateInstance(t *testing.T) {
	testLogger := logger.NewLogger("test")

	t.Run("bpmnProcessId and processDefinitionKey are not allowed at the same time", func(t *testing.T) {
		payload := createInstancePayload{
			BpmnProcessID:        "some-id",
			ProcessDefinitionKey: new(int64),
		}
		data, err := json.Marshal(payload)
		assert.Nil(t, err)

		req := &bindings.InvokeRequest{Data: data, Operation: CreateInstanceOperation}

		var mc mockCreateInstanceClient

		cmd := ZeebeCommand{logger: testLogger, client: &mc}
		_, err = cmd.Invoke(req)
		assert.Error(t, err, ErrAmbiguousCreationVars)
	})

	t.Run("either bpmnProcessId or processDefinitionKey must be given", func(t *testing.T) {
		payload := createInstancePayload{}
		data, err := json.Marshal(payload)
		assert.NoError(t, err)

		req := &bindings.InvokeRequest{Data: data, Operation: CreateInstanceOperation}

		var mc mockCreateInstanceClient

		cmd := ZeebeCommand{logger: testLogger, client: &mc}
		_, err = cmd.Invoke(req)
		assert.Error(t, err, ErrMissingCreationVars)
	})

	t.Run("create command with bpmnProcessId and specific version", func(t *testing.T) {
		payload := createInstancePayload{
			BpmnProcessID: "some-id",
			Version:       new(int32),
		}
		data, err := json.Marshal(payload)
		assert.NoError(t, err)

		req := &bindings.InvokeRequest{Data: data, Operation: CreateInstanceOperation}

		var mc mockCreateInstanceClient

		cmd := ZeebeCommand{logger: testLogger, client: &mc}
		_, err = cmd.Invoke(req)
		assert.NoError(t, err)

		assert.Equal(t, payload.BpmnProcessID, mc.cmd1.bpmnProcessID)
		assert.Equal(t, *payload.Version, mc.cmd1.cmd2.version)
	})

	t.Run("create command with bpmnProcessId and latest version", func(t *testing.T) {
		payload := createInstancePayload{
			BpmnProcessID: "some-id",
		}
		data, err := json.Marshal(payload)
		assert.NoError(t, err)

		req := &bindings.InvokeRequest{Data: data, Operation: CreateInstanceOperation}

		var mc mockCreateInstanceClient

		cmd := ZeebeCommand{logger: testLogger, client: &mc}
		_, err = cmd.Invoke(req)
		assert.NoError(t, err)

		assert.Equal(t, payload.BpmnProcessID, mc.cmd1.bpmnProcessID)
		assert.Equal(t, true, mc.cmd1.cmd2.latestVersion)
	})

	t.Run("create command with processDefinitionKey", func(t *testing.T) {
		payload := createInstancePayload{
			ProcessDefinitionKey: new(int64),
		}
		data, err := json.Marshal(payload)
		assert.NoError(t, err)

		req := &bindings.InvokeRequest{Data: data, Operation: CreateInstanceOperation}

		var mc mockCreateInstanceClient

		cmd := ZeebeCommand{logger: testLogger, client: &mc}
		_, err = cmd.Invoke(req)
		assert.NoError(t, err)

		assert.Equal(t, *payload.ProcessDefinitionKey, mc.cmd1.processDefinitionKey)
	})

	t.Run("create command with variables", func(t *testing.T) {
		payload := createInstancePayload{
			ProcessDefinitionKey: new(int64),
			Variables: map[string]interface{}{
				"key": "value",
			},
		}
		data, err := json.Marshal(payload)
		assert.NoError(t, err)

		req := &bindings.InvokeRequest{Data: data, Operation: CreateInstanceOperation}

		var mc mockCreateInstanceClient

		cmd := ZeebeCommand{logger: testLogger, client: &mc}
		_, err = cmd.Invoke(req)
		assert.NoError(t, err)

		assert.Equal(t, *payload.ProcessDefinitionKey, mc.cmd1.processDefinitionKey)
		assert.Equal(t, payload.Variables, mc.cmd1.cmd2.cmd3.variables)
	})
}
