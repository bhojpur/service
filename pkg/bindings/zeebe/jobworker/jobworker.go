package jobworker

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
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/camunda-cloud/zeebe/clients/go/pkg/entities"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/worker"
	"github.com/camunda-cloud/zeebe/clients/go/pkg/zbc"

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/bindings/zeebe"
	"github.com/bhojpur/service/pkg/metadata"
	"github.com/bhojpur/service/pkg/utils/logger"
)

var ErrMissingJobType = errors.New("jobType is a required attribute")

// ZeebeJobWorker allows handling jobs from the Zeebe command engine.
type ZeebeJobWorker struct {
	clientFactory zeebe.ClientFactory
	client        zbc.Client
	metadata      *jobWorkerMetadata
	logger        logger.Logger
}

// https://docs.zeebe.io/basics/job-workers.html
type jobWorkerMetadata struct {
	WorkerName     string            `json:"workerName"`
	WorkerTimeout  metadata.Duration `json:"workerTimeout"`
	RequestTimeout metadata.Duration `json:"requestTimeout"`
	JobType        string            `json:"jobType"`
	MaxJobsActive  int               `json:"maxJobsActive,string"`
	Concurrency    int               `json:"concurrency,string"`
	PollInterval   metadata.Duration `json:"pollInterval"`
	PollThreshold  float64           `json:"pollThreshold,string"`
	FetchVariables string            `json:"fetchVariables"`
}

type jobHandler struct {
	callback func(*bindings.ReadResponse) ([]byte, error)
	logger   logger.Logger
}

// NewZeebeJobWorker returns a new ZeebeJobWorker instance.
func NewZeebeJobWorker(logger logger.Logger) *ZeebeJobWorker {
	return &ZeebeJobWorker{clientFactory: zeebe.NewClientFactoryImpl(logger), logger: logger}
}

// Init does metadata parsing and connection creation.
func (z *ZeebeJobWorker) Init(metadata bindings.Metadata) error {
	meta, err := z.parseMetadata(metadata)
	if err != nil {
		return err
	}

	if meta.JobType == "" {
		return ErrMissingJobType
	}

	client, err := z.clientFactory.Get(metadata)
	if err != nil {
		return err
	}

	z.metadata = meta
	z.client = client

	return nil
}

func (z *ZeebeJobWorker) Read(handler func(*bindings.ReadResponse) ([]byte, error)) error {
	h := jobHandler{
		callback: handler,
		logger:   z.logger,
	}

	jobWorker := z.getJobWorker(h)

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, syscall.SIGTERM)
	<-exitChan

	jobWorker.Close()
	jobWorker.AwaitClose()

	return z.client.Close()
}

func (z *ZeebeJobWorker) parseMetadata(metadata bindings.Metadata) (*jobWorkerMetadata, error) {
	b, err := json.Marshal(metadata.Properties)
	if err != nil {
		return nil, err
	}

	var m jobWorkerMetadata
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (z *ZeebeJobWorker) getJobWorker(handler jobHandler) worker.JobWorker {
	cmd1 := z.client.NewJobWorker()
	cmd2 := cmd1.JobType(z.metadata.JobType)
	cmd3 := cmd2.Handler(handler.handleJob)
	if z.metadata.WorkerName != "" {
		cmd3 = cmd3.Name(z.metadata.WorkerName)
	}
	if z.metadata.WorkerTimeout.Duration != time.Duration(0) {
		cmd3 = cmd3.Timeout(z.metadata.WorkerTimeout.Duration)
	}
	if z.metadata.RequestTimeout.Duration != time.Duration(0) {
		cmd3 = cmd3.RequestTimeout(z.metadata.RequestTimeout.Duration)
	}
	if z.metadata.MaxJobsActive != 0 {
		cmd3 = cmd3.MaxJobsActive(z.metadata.MaxJobsActive)
	}
	if z.metadata.Concurrency != 0 {
		cmd3 = cmd3.Concurrency(z.metadata.Concurrency)
	}
	if z.metadata.PollInterval.Duration != time.Duration(0) {
		cmd3 = cmd3.PollInterval(z.metadata.PollInterval.Duration)
	}
	if z.metadata.PollThreshold != 0 {
		cmd3 = cmd3.PollThreshold(z.metadata.PollThreshold)
	}
	if z.metadata.FetchVariables != "" {
		cmd3 = cmd3.FetchVariables(zeebe.VariableStringToArray(z.metadata.FetchVariables)...)
	}

	return cmd3.Open()
}

func (h *jobHandler) handleJob(client worker.JobClient, job entities.Job) {
	headers, err := job.GetCustomHeadersAsMap()
	if err != nil {
		h.failJob(client, job, err)

		return
	}

	headers["X-Zeebe-Job-Key"] = strconv.FormatInt(job.Key, 10)
	headers["X-Zeebe-Job-Type"] = job.Type
	headers["X-Zeebe-Process-Instance-Key"] = strconv.FormatInt(job.ProcessInstanceKey, 10)
	headers["X-Zeebe-Bpmn-Process-Id"] = job.BpmnProcessId
	headers["X-Zeebe-Process-Definition-Version"] = strconv.FormatInt(int64(job.ProcessDefinitionVersion), 10)
	headers["X-Zeebe-Process-Definition-Key"] = strconv.FormatInt(job.ProcessDefinitionKey, 10)
	headers["X-Zeebe-Element-Id"] = job.ElementId
	headers["X-Zeebe-Element-Instance-Key"] = strconv.FormatInt(job.ElementInstanceKey, 10)
	headers["X-Zeebe-Worker"] = job.Worker
	headers["X-Zeebe-Retries"] = strconv.FormatInt(int64(job.Retries), 10)
	headers["X-Zeebe-Deadline"] = strconv.FormatInt(job.Deadline, 10)

	resultVariables, err := h.callback(&bindings.ReadResponse{
		Data:     []byte(job.Variables),
		Metadata: headers,
	})
	if err != nil {
		h.failJob(client, job, err)

		return
	}

	variablesMap := make(map[string]interface{})
	if resultVariables != nil {
		err = json.Unmarshal(resultVariables, &variablesMap)
		if err != nil {
			h.failJob(client, job, fmt.Errorf("cannot parse variables from binding result %s; got error %w", string(resultVariables), err))

			return
		}
	}

	jobKey := job.GetKey()
	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variablesMap)
	if err != nil {
		h.failJob(client, job, err)

		return
	}

	h.logger.Debugf("Complete job `%d` of type `%s`", jobKey, job.Type)

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		h.logger.Errorf("Cannot complete job `%d` of type `%s`; got error: %s", jobKey, job.Type, err.Error())

		return
	}

	h.logger.Debug("Successfully completed job")
}

func (h *jobHandler) failJob(client worker.JobClient, job entities.Job, reason error) {
	reasonMsg := reason.Error()
	h.logger.Errorf("Failed to complete job `%d` reason: %s", job.GetKey(), reasonMsg)

	ctx := context.Background()
	_, err := client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).ErrorMessage(reasonMsg).Send(ctx)
	if err != nil {
		h.logger.Errorf("Cannot fail job `%d` of type `%s`; got error: %s", job.GetKey(), job.Type, err.Error())

		return
	}
}
