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
	"time"

	"github.com/camunda-cloud/zeebe/clients/go/pkg/commands"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/pb"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/zbc"
	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/bindings"
	contrib_metadata "github.com/bhojpur/service/pkg/metadata"
	"github.com/bhojpur/service/pkg/utils/logger"
)

type mockPublishMessageClient struct {
	zbc.Client
	cmd1 *mockPublishMessageCommandStep1
}

type mockPublishMessageCommandStep1 struct {
	commands.PublishMessageCommandStep1
	cmd2        *mockPublishMessageCommandStep2
	messageName string
}

type mockPublishMessageCommandStep2 struct {
	commands.PublishMessageCommandStep2
	cmd3           *mockPublishMessageCommandStep3
	correlationKey string
}

type mockPublishMessageCommandStep3 struct {
	commands.PublishMessageCommandStep3
	messageID  string
	timeToLive time.Duration
	variables  interface{}
}

func (mc *mockPublishMessageClient) NewPublishMessageCommand() commands.PublishMessageCommandStep1 {
	mc.cmd1 = &mockPublishMessageCommandStep1{
		cmd2: &mockPublishMessageCommandStep2{
			cmd3: &mockPublishMessageCommandStep3{},
		},
	}

	return mc.cmd1
}

func (cmd1 *mockPublishMessageCommandStep1) MessageName(messageName string) commands.PublishMessageCommandStep2 {
	cmd1.messageName = messageName

	return cmd1.cmd2
}

func (cmd2 *mockPublishMessageCommandStep2) CorrelationKey(correlationKey string) commands.PublishMessageCommandStep3 {
	cmd2.correlationKey = correlationKey

	return cmd2.cmd3
}

//nolint // MessageId comes from the Zeebe client API and cannot be written as MessageID
func (cmd3 *mockPublishMessageCommandStep3) MessageId(messageID string) commands.PublishMessageCommandStep3 {
	cmd3.messageID = messageID

	return cmd3
}

func (cmd3 *mockPublishMessageCommandStep3) TimeToLive(timeToLive time.Duration) commands.PublishMessageCommandStep3 {
	cmd3.timeToLive = timeToLive

	return cmd3
}

func (cmd3 *mockPublishMessageCommandStep3) VariablesFromObject(variables interface{}) (commands.PublishMessageCommandStep3, error) {
	cmd3.variables = variables

	return cmd3, nil
}

func (cmd3 *mockPublishMessageCommandStep3) Send(context.Context) (*pb.PublishMessageResponse, error) {
	return &pb.PublishMessageResponse{}, nil
}

func TestPublishMessage(t *testing.T) {
	testLogger := logger.NewLogger("test")

	t.Run("messageName is mandatory", func(t *testing.T) {
		cmd := ZeebeCommand{logger: testLogger}
		req := &bindings.InvokeRequest{Operation: PublishMessageOperation}
		_, err := cmd.Invoke(req)
		assert.Error(t, err, ErrMissingMessageName)
	})

	t.Run("send message with mandatory fields", func(t *testing.T) {
		payload := publishMessagePayload{
			MessageName: "a",
		}
		data, err := json.Marshal(payload)
		assert.NoError(t, err)

		req := &bindings.InvokeRequest{Data: data, Operation: PublishMessageOperation}

		var mc mockPublishMessageClient

		cmd := ZeebeCommand{logger: testLogger, client: &mc}
		_, err = cmd.Invoke(req)
		assert.NoError(t, err)

		assert.Equal(t, payload.MessageName, mc.cmd1.messageName)
		assert.Equal(t, payload.CorrelationKey, mc.cmd1.cmd2.correlationKey)
	})

	t.Run("send message with optional fields", func(t *testing.T) {
		payload := publishMessagePayload{
			MessageName:    "a",
			CorrelationKey: "b",
			MessageID:      "c",
			TimeToLive:     contrib_metadata.Duration{Duration: 1 * time.Second},
			Variables: map[string]interface{}{
				"key": "value",
			},
		}
		data, err := json.Marshal(payload)
		assert.NoError(t, err)

		req := &bindings.InvokeRequest{Data: data, Operation: PublishMessageOperation}

		var mc mockPublishMessageClient

		cmd := ZeebeCommand{logger: testLogger, client: &mc}
		_, err = cmd.Invoke(req)
		assert.NoError(t, err)

		assert.Equal(t, payload.MessageName, mc.cmd1.messageName)
		assert.Equal(t, payload.CorrelationKey, mc.cmd1.cmd2.correlationKey)
		assert.Equal(t, payload.MessageID, mc.cmd1.cmd2.cmd3.messageID)
		assert.Equal(t, payload.TimeToLive.Duration, mc.cmd1.cmd2.cmd3.timeToLive)
		assert.Equal(t, payload.Variables, mc.cmd1.cmd2.cmd3.variables)
	})
}
