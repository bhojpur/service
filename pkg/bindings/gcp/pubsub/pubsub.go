package pubsub

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

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/utils/logger"
)

const (
	id          = "id"
	publishTime = "publishTime"
	topic       = "topic"
)

// GCPPubSub is an input/output binding for GCP Pub Sub.
type GCPPubSub struct {
	client   *pubsub.Client
	metadata *pubSubMetadata
	logger   logger.Logger
}

type pubSubMetadata struct {
	Topic               string `json:"topic"`
	Subscription        string `json:"subscription"`
	Type                string `json:"type"`
	ProjectID           string `json:"project_id"`
	PrivateKeyID        string `json:"private_key_id"`
	PrivateKey          string `json:"private_key"`
	ClientEmail         string `json:"client_email"`
	ClientID            string `json:"client_id"`
	AuthURI             string `json:"auth_uri"`
	TokenURI            string `json:"token_uri"`
	AuthProviderCertURL string `json:"auth_provider_x509_cert_url"`
	ClientCertURL       string `json:"client_x509_cert_url"`
}

// NewGCPPubSub returns a new GCPPubSub instance.
func NewGCPPubSub(logger logger.Logger) *GCPPubSub {
	return &GCPPubSub{logger: logger}
}

// Init parses metadata and creates a new Pub Sub client.
func (g *GCPPubSub) Init(metadata bindings.Metadata) error {
	b, err := g.parseMetadata(metadata)
	if err != nil {
		return err
	}

	var pubsubMeta pubSubMetadata
	err = json.Unmarshal(b, &pubsubMeta)
	if err != nil {
		return err
	}
	clientOptions := option.WithCredentialsJSON(b)
	ctx := context.Background()
	pubsubClient, err := pubsub.NewClient(ctx, pubsubMeta.ProjectID, clientOptions)
	if err != nil {
		return fmt.Errorf("error creating pubsub client: %s", err)
	}

	g.client = pubsubClient
	g.metadata = &pubsubMeta

	return nil
}

func (g *GCPPubSub) parseMetadata(metadata bindings.Metadata) ([]byte, error) {
	b, err := json.Marshal(metadata.Properties)

	return b, err
}

func (g *GCPPubSub) Read(handler func(*bindings.ReadResponse) ([]byte, error)) error {
	sub := g.client.Subscription(g.metadata.Subscription)
	err := sub.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
		_, err := handler(&bindings.ReadResponse{
			Data:     m.Data,
			Metadata: map[string]string{id: m.ID, publishTime: m.PublishTime.String()},
		})
		if err == nil {
			m.Ack()
		}
	})

	return err
}

func (g *GCPPubSub) Operations() []bindings.OperationKind {
	return []bindings.OperationKind{bindings.CreateOperation}
}

func (g *GCPPubSub) Invoke(req *bindings.InvokeRequest) (*bindings.InvokeResponse, error) {
	topicName := g.metadata.Topic
	if val, ok := req.Metadata[topic]; ok && val != "" {
		topicName = val
	}

	t := g.client.Topic(topicName)
	ctx := context.Background()
	_, err := t.Publish(ctx, &pubsub.Message{
		Data: req.Data,
	}).Get(ctx)

	return nil, err
}

func (g *GCPPubSub) Close() error {
	return g.client.Close()
}
