package engine

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
	"fmt"
	"log"
	"net/http"

	"github.com/bhojpur/service/pkg/engine/config"
	"github.com/bhojpur/service/pkg/engine/core"
	"github.com/bhojpur/service/pkg/engine/logger"
)

const (
	processorLogPrefix = "\033[33m[bhojpur:processor]\033[0m "
)

// Processor is the orchestrator of Bhojpur Service. There are two types of processor:
// 1. Upstream Processor, which is used to connect to multiple downstream processors,
// 2. Downstream Processor (will call it as Processor directly), which is used
// to connected by `Upstream Processor`, `Source` and `Stream Function`
type Processor interface {
	// ConfigWorkflow will register workflows from config files to processor.
	ConfigWorkflow(conf string) error

	// ConfigMesh will register edge-mesh config URL
	ConfigMesh(url string) error

	// ListenAndServe start processor as server.
	ListenAndServe() error

	// AddDownstreamProcessor will add downstream processor.
	AddDownstreamProcessor(downstream Processor) error

	// Addr returns the listen address of processor.
	Addr() string

	// Stats return insight data
	Stats() int

	// Close will close the processor.
	Close() error

	// ReadConfigFile(conf string) error
	// AddWorkflow(wf ...core.Workflow) error
	// ConfigDownstream(opts ...interface{}) error
	// Connect() error
	// RemoveDownstreamProcessor(downstream Processor) error
	// ListenAddr() string
}

// processor is the implementation of Processor interface.
type processor struct {
	name                 string
	addr                 string
	hasDownstreams       bool
	server               *core.Server
	client               *core.Client
	downstreamProcessors []Processor
}

var _ Processor = &processor{}

// NewProcessorWithOptions create a processor instance.
func NewProcessorWithOptions(name string, opts ...Option) Processor {
	options := NewOptions(opts...)
	processor := createProcessorServer(name, options)
	processor.ConfigMesh(options.MeshConfigURL)

	return processor
}

// NewProcessor create a processor instance from config files.
func NewProcessor(conf string) (Processor, error) {
	config, err := config.ParseWorkflowConfig(conf)
	if err != nil {
		logger.Errorf("%s[ERR] %v", processorLogPrefix, err)
		return nil, err
	}
	// listening address
	listenAddr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	options := NewOptions()
	options.ProcessorAddr = listenAddr
	processor := createProcessorServer(config.Name, options)
	// processor workflow
	err = processor.configWorkflow(config)

	return processor, err
}

// NewDownstreamProcessor create a processor descriptor for downstream processor.
func NewDownstreamProcessor(name string, opts ...Option) Processor {
	options := NewOptions(opts...)
	client := core.NewClient(name, core.ClientTypeUpstreamProcessor, options.ClientOptions...)

	return &processor{
		name:   name,
		addr:   options.ProcessorAddr,
		client: client,
	}
}

/*************** Server ONLY ***************/
// createProcessorServer create a processor instance as server.
func createProcessorServer(name string, options *Options) *processor {
	// create underlying QUIC server
	srv := core.NewServer(name, options.ServerOptions...)
	z := &processor{
		server: srv,
		name:   name,
		addr:   options.ProcessorAddr,
	}
	// initialize
	z.init()
	return z
}

// ConfigWorkflow will read workflows from config files and register them to processor.
func (z *processor) ConfigWorkflow(conf string) error {
	config, err := config.ParseWorkflowConfig(conf)
	if err != nil {
		logger.Errorf("%s[ERR] %v", processorLogPrefix, err)
		return err
	}
	logger.Debugf("%sConfigWorkflow config=%+v", processorLogPrefix, config)
	return z.configWorkflow(config)
}

func (z *processor) configWorkflow(config *config.WorkflowConfig) error {
	// router
	return z.server.ConfigRouter(newRouter(config))
}

func (z *processor) ConfigMesh(url string) error {
	if url == "" {
		return nil
	}

	logger.Printf("%sDownloading mesh config...", processorLogPrefix)
	// download mesh conf
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var configs []config.MeshProcessor
	err = decoder.Decode(&configs)
	if err != nil {
		logger.Errorf("%s✅ downloaded the Mesh config with err=%v", processorLogPrefix, err)
		return err
	}

	logger.Printf("%s✅ Successfully downloaded the Service Mesh config. ", processorLogPrefix)

	if len(configs) == 0 {
		return nil
	}

	for _, downstream := range configs {
		if downstream.Name == z.name {
			continue
		}
		addr := fmt.Sprintf("%s:%d", downstream.Host, downstream.Port)
		z.AddDownstreamProcessor(NewDownstreamProcessor(downstream.Name, WithProcessorAddr(addr)))
	}

	return nil
}

// ListenAndServe will start processor service.
func (z *processor) ListenAndServe() error {
	logger.Debugf("%sCreating Processor Server ...", processorLogPrefix)
	// check downstream processors
	for _, ds := range z.downstreamProcessors {
		if dsProcessor, ok := ds.(*processor); ok {
			go func(dsProcessor *processor) {
				dsProcessor.client.Connect(context.Background(), dsProcessor.addr)
				z.server.AddDownstreamServer(dsProcessor.addr, dsProcessor.client)
			}(dsProcessor)
		}
	}
	return z.server.ListenAndServe(context.Background(), z.addr)
}

// AddDownstreamProcessor will add downstream processor.
func (z *processor) AddDownstreamProcessor(downstream Processor) error {
	logger.Debugf("%sAddDownstreamProcessor: %v", processorLogPrefix, downstream)
	z.downstreamProcessors = append(z.downstreamProcessors, downstream)
	z.hasDownstreams = true
	logger.Debugf("%scurrent downstreams: %d", processorLogPrefix, len(z.downstreamProcessors))
	return nil
}

// RemoveDownstreamProcessor remove downstream processor.
func (z *processor) RemoveDownstreamProcessor(downstream Processor) error {
	index := -1
	for i, v := range z.downstreamProcessors {
		if v.Addr() == downstream.Addr() {
			index = i
			break
		}
	}

	// remove from slice
	z.downstreamProcessors = append(z.downstreamProcessors[:index], z.downstreamProcessors[index+1:]...)
	return nil
}

// Addr returns listen address of processor.
func (z *processor) Addr() string {
	return z.addr
}

// Close will close a connection. If processor is Server, close the server. If processor is Client, close the client.
func (z *processor) Close() error {
	if z.server != nil {
		if err := z.server.Close(); err != nil {
			logger.Errorf("%s Close(): %v", processorLogPrefix, err)
			return err
		}
	}
	if z.client != nil {
		if err := z.client.Close(); err != nil {
			logger.Errorf("%s Close(): %v", processorLogPrefix, err)
			return err
		}
	}
	return nil
}

// Stats inspects current server.
func (z *processor) Stats() int {
	log.Printf("[%s] all sfn connected: %d", z.name, len(z.server.StatsFunctions()))
	for k := range z.server.StatsFunctions() {
		log.Printf("[%s] -> ConnID=%v", z.name, k)
	}

	log.Printf("[%s] all downstream processors connected: %d", z.name, len(z.server.Downstreams()))
	for k, v := range z.server.Downstreams() {
		log.Printf("[%s] |> [%s] %s", z.name, k, v.ServerAddr())
	}

	log.Printf("[%s] total DataFrames received: %d", z.name, z.server.StatsCounter())

	return len(z.server.StatsFunctions())
}
