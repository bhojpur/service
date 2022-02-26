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

type mockCancelInstanceClient struct {
	zbc.Client
	cmd1 *mockCancelInstanceStep1
}

type mockCancelInstanceStep1 struct {
	commands.CancelInstanceStep1
	cmd2 *mockDispatchCancelProcessInstanceCommand
}

type mockDispatchCancelProcessInstanceCommand struct {
	commands.DispatchCancelProcessInstanceCommand
	processInstanceKey int64
}

func (mc *mockCancelInstanceClient) NewCancelInstanceCommand() commands.CancelInstanceStep1 {
	mc.cmd1 = &mockCancelInstanceStep1{
		cmd2: &mockDispatchCancelProcessInstanceCommand{},
	}

	return mc.cmd1
}

func (cmd1 *mockCancelInstanceStep1) ProcessInstanceKey(processInstanceKey int64) commands.DispatchCancelProcessInstanceCommand {
	cmd1.cmd2.processInstanceKey = processInstanceKey

	return cmd1.cmd2
}

func (cmd2 *mockDispatchCancelProcessInstanceCommand) Send(context.Context) (*pb.CancelProcessInstanceResponse, error) {
	return &pb.CancelProcessInstanceResponse{}, nil
}

func TestCancelInstance(t *testing.T) {
	testLogger := logger.NewLogger("test")

	t.Run("processInstanceKey is mandatory", func(t *testing.T) {
		cmd := ZeebeCommand{logger: testLogger}
		req := &bindings.InvokeRequest{Operation: CancelInstanceOperation}
		_, err := cmd.Invoke(req)
		assert.Error(t, err, ErrMissingProcessInstanceKey)
	})

	t.Run("cancel a command", func(t *testing.T) {
		payload := cancelInstancePayload{
			ProcessInstanceKey: new(int64),
		}
		data, err := json.Marshal(payload)
		assert.NoError(t, err)

		req := &bindings.InvokeRequest{Data: data, Operation: CancelInstanceOperation}

		var mc mockCancelInstanceClient

		cmd := ZeebeCommand{logger: testLogger, client: &mc}
		_, err = cmd.Invoke(req)
		assert.NoError(t, err)

		assert.Equal(t, *payload.ProcessInstanceKey, mc.cmd1.cmd2.processInstanceKey)
	})
}
