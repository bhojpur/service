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
	"errors"
	"fmt"

	"github.com/camunda-cloud/zeebe/clients/go/pkg/zbc"

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/bindings/zeebe"
	"github.com/bhojpur/service/pkg/utils/logger"
)

const (
	// operations.
	TopologyOperation         bindings.OperationKind = "topology"
	DeployProcessOperation    bindings.OperationKind = "deploy-process"
	CreateInstanceOperation   bindings.OperationKind = "create-instance"
	CancelInstanceOperation   bindings.OperationKind = "cancel-instance"
	SetVariablesOperation     bindings.OperationKind = "set-variables"
	ResolveIncidentOperation  bindings.OperationKind = "resolve-incident"
	PublishMessageOperation   bindings.OperationKind = "publish-message"
	ActivateJobsOperation     bindings.OperationKind = "activate-jobs"
	CompleteJobOperation      bindings.OperationKind = "complete-job"
	FailJobOperation          bindings.OperationKind = "fail-job"
	UpdateJobRetriesOperation bindings.OperationKind = "update-job-retries"
	ThrowErrorOperation       bindings.OperationKind = "throw-error"
)

var (
	ErrMissingJobKey        = errors.New("jobKey is a required attribute")
	ErrUnsupportedOperation = func(operation bindings.OperationKind) error {
		return fmt.Errorf("unsupported operation: %v", operation)
	}
)

// ZeebeCommand executes Zeebe commands.
type ZeebeCommand struct {
	clientFactory zeebe.ClientFactory
	client        zbc.Client
	logger        logger.Logger
}

// NewZeebeCommand returns a new ZeebeCommand instance.
func NewZeebeCommand(logger logger.Logger) *ZeebeCommand {
	return &ZeebeCommand{clientFactory: zeebe.NewClientFactoryImpl(logger), logger: logger}
}

// Init does metadata parsing and connection creation.
func (z *ZeebeCommand) Init(metadata bindings.Metadata) error {
	client, err := z.clientFactory.Get(metadata)
	if err != nil {
		return err
	}

	z.client = client

	return nil
}

func (z *ZeebeCommand) Operations() []bindings.OperationKind {
	return []bindings.OperationKind{
		TopologyOperation,
		DeployProcessOperation,
		CreateInstanceOperation,
		CancelInstanceOperation,
		SetVariablesOperation,
		ResolveIncidentOperation,
		PublishMessageOperation,
		ActivateJobsOperation,
		CompleteJobOperation,
		FailJobOperation,
		UpdateJobRetriesOperation,
		ThrowErrorOperation,
	}
}

func (z *ZeebeCommand) Invoke(req *bindings.InvokeRequest) (*bindings.InvokeResponse, error) {
	switch req.Operation {
	case TopologyOperation:
		return z.topology()
	case DeployProcessOperation:
		return z.deployProcess(req)
	case CreateInstanceOperation:
		return z.createInstance(req)
	case CancelInstanceOperation:
		return z.cancelInstance(req)
	case SetVariablesOperation:
		return z.setVariables(req)
	case ResolveIncidentOperation:
		return z.resolveIncident(req)
	case PublishMessageOperation:
		return z.publishMessage(req)
	case ActivateJobsOperation:
		return z.activateJobs(req)
	case CompleteJobOperation:
		return z.completeJob(req)
	case FailJobOperation:
		return z.failJob(req)
	case UpdateJobRetriesOperation:
		return z.updateJobRetries(req)
	case ThrowErrorOperation:
		return z.throwError(req)
	case bindings.GetOperation:
		fallthrough
	case bindings.CreateOperation:
		fallthrough
	case bindings.DeleteOperation:
		fallthrough
	case bindings.ListOperation:
		fallthrough
	default:
		return nil, ErrUnsupportedOperation(req.Operation)
	}
}
