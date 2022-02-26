package storagequeues

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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/Azure/azure-storage-queue-go/azqueue"

	"github.com/bhojpur/service/pkg/bindings"
	contrib_metadata "github.com/bhojpur/service/pkg/metadata"
	"github.com/bhojpur/service/pkg/utils/logger"
)

const (
	defaultTTL = time.Minute * 10
)

type consumer struct {
	callback func(*bindings.ReadResponse) ([]byte, error)
}

// QueueHelper enables injection for testnig.
type QueueHelper interface {
	Init(accountName string, accountKey string, queueName string, decodeBase64 bool) error
	Write(data []byte, ttl *time.Duration) error
	Read(ctx context.Context, consumer *consumer) error
}

// AzureQueueHelper concrete impl of queue helper.
type AzureQueueHelper struct {
	credential   *azqueue.SharedKeyCredential
	queueURL     azqueue.QueueURL
	reqURI       string
	logger       logger.Logger
	decodeBase64 bool
}

// Init sets up this helper.
func (d *AzureQueueHelper) Init(accountName string, accountKey string, queueName string, decodeBase64 bool) error {
	credential, err := azqueue.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return err
	}
	d.credential = credential
	d.decodeBase64 = decodeBase64
	u, _ := url.Parse(fmt.Sprintf(d.reqURI, accountName, queueName))
	userAgent := "app-" + logger.AppVersion
	pipelineOptions := azqueue.PipelineOptions{
		Telemetry: azqueue.TelemetryOptions{
			Value: userAgent,
		},
	}
	d.queueURL = azqueue.NewQueueURL(*u, azqueue.NewPipeline(credential, pipelineOptions))
	ctx := context.TODO()
	_, err = d.queueURL.Create(ctx, azqueue.Metadata{})
	if err != nil {
		return err
	}

	return nil
}

func (d *AzureQueueHelper) Write(data []byte, ttl *time.Duration) error {
	ctx := context.TODO()
	messagesURL := d.queueURL.NewMessagesURL()

	s, err := strconv.Unquote(string(data))
	if err != nil {
		s = string(data)
	}

	if ttl == nil {
		ttlToUse := defaultTTL
		ttl = &ttlToUse
	}
	_, err = messagesURL.Enqueue(ctx, s, time.Second*0, *ttl)

	return err
}

func (d *AzureQueueHelper) Read(ctx context.Context, consumer *consumer) error {
	messagesURL := d.queueURL.NewMessagesURL()
	res, err := messagesURL.Dequeue(ctx, 1, time.Second*30)
	if err != nil {
		return err
	}
	if res.NumMessages() == 0 {
		// Queue was empty so back off by 10 seconds before trying again
		time.Sleep(10 * time.Second)

		return nil
	}
	mt := res.Message(0).Text

	var data []byte

	if d.decodeBase64 {
		decoded, decodeError := base64.StdEncoding.DecodeString(mt)
		if decodeError != nil {
			return decodeError
		}
		data = decoded
	} else {
		data = []byte(mt)
	}

	_, err = consumer.callback(&bindings.ReadResponse{
		Data:     data,
		Metadata: map[string]string{},
	})
	if err != nil {
		return err
	}
	messageIDURL := messagesURL.NewMessageIDURL(res.Message(0).ID)
	pr := res.Message(0).PopReceipt
	_, err = messageIDURL.Delete(ctx, pr)
	if err != nil {
		return err
	}

	return nil
}

// NewAzureQueueHelper creates new helper.
func NewAzureQueueHelper(logger logger.Logger) QueueHelper {
	return &AzureQueueHelper{
		reqURI: "https://%s.queue.core.windows.net/%s",
		logger: logger,
	}
}

// AzureStorageQueues is an input/output binding reading from and sending events to Azure Storage queues.
type AzureStorageQueues struct {
	metadata *storageQueuesMetadata
	helper   QueueHelper

	logger logger.Logger
}

type storageQueuesMetadata struct {
	AccountKey   string `json:"storageAccessKey"`
	QueueName    string `json:"queue"`
	AccountName  string `json:"storageAccount"`
	DecodeBase64 string `json:"decodeBase64"`
	ttl          *time.Duration
}

// NewAzureStorageQueues returns a new AzureStorageQueues instance.
func NewAzureStorageQueues(logger logger.Logger) *AzureStorageQueues {
	return &AzureStorageQueues{helper: NewAzureQueueHelper(logger), logger: logger}
}

// Init parses connection properties and creates a new Storage Queue client.
func (a *AzureStorageQueues) Init(metadata bindings.Metadata) error {
	meta, err := a.parseMetadata(metadata)
	if err != nil {
		return err
	}
	a.metadata = meta

	decodeBase64 := false
	if a.metadata.DecodeBase64 == "true" {
		decodeBase64 = true
	}

	err = a.helper.Init(a.metadata.AccountName, a.metadata.AccountKey, a.metadata.QueueName, decodeBase64)
	if err != nil {
		return err
	}

	return nil
}

func (a *AzureStorageQueues) parseMetadata(metadata bindings.Metadata) (*storageQueuesMetadata, error) {
	b, err := json.Marshal(metadata.Properties)
	if err != nil {
		return nil, err
	}
	var m storageQueuesMetadata
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}

	ttl, ok, err := contrib_metadata.TryGetTTL(metadata.Properties)
	if err != nil {
		return nil, err
	}

	if ok {
		m.ttl = &ttl
	}

	return &m, nil
}

func (a *AzureStorageQueues) Operations() []bindings.OperationKind {
	return []bindings.OperationKind{bindings.CreateOperation}
}

func (a *AzureStorageQueues) Invoke(req *bindings.InvokeRequest) (*bindings.InvokeResponse, error) {
	ttlToUse := a.metadata.ttl
	ttl, ok, err := contrib_metadata.TryGetTTL(req.Metadata)
	if err != nil {
		return nil, err
	}

	if ok {
		ttlToUse = &ttl
	}

	err = a.helper.Write(req.Data, ttlToUse)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (a *AzureStorageQueues) Read(handler func(*bindings.ReadResponse) ([]byte, error)) error {
	c := consumer{
		callback: handler,
	}
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			err := a.helper.Read(ctx, &c)
			if err != nil {
				a.logger.Errorf("error from c: %s", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
	cancel()
	wg.Wait()

	return nil
}
