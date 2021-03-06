package rabbitmq

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
	"math"
	"strconv"
	"time"

	"github.com/streadway/amqp"

	"github.com/bhojpur/service/pkg/bindings"
	contrib_metadata "github.com/bhojpur/service/pkg/metadata"
	"github.com/bhojpur/service/pkg/utils/logger"
)

const (
	host                       = "host"
	queueName                  = "queueName"
	exclusive                  = "exclusive"
	durable                    = "durable"
	deleteWhenUnused           = "deleteWhenUnused"
	prefetchCount              = "prefetchCount"
	maxPriority                = "maxPriority"
	rabbitMQQueueMessageTTLKey = "x-message-ttl"
	rabbitMQMaxPriorityKey     = "x-max-priority"
	defaultBase                = 10
	defaultBitSize             = 0
)

// RabbitMQ allows sending/receiving data to/from RabbitMQ.
type RabbitMQ struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	metadata   rabbitMQMetadata
	logger     logger.Logger
	queue      amqp.Queue
}

// Metadata is the rabbitmq config.
type rabbitMQMetadata struct {
	Host             string `json:"host"`
	QueueName        string `json:"queueName"`
	Exclusive        bool   `json:"exclusive,string"`
	Durable          bool   `json:"durable,string"`
	DeleteWhenUnused bool   `json:"deleteWhenUnused,string"`
	PrefetchCount    int    `json:"prefetchCount"`
	MaxPriority      *uint8 `json:"maxPriority"` // Priority Queue deactivated if nil
	defaultQueueTTL  *time.Duration
}

// NewRabbitMQ returns a new rabbitmq instance.
func NewRabbitMQ(logger logger.Logger) *RabbitMQ {
	return &RabbitMQ{logger: logger}
}

// Init does metadata parsing and connection creation.
func (r *RabbitMQ) Init(metadata bindings.Metadata) error {
	err := r.parseMetadata(metadata)
	if err != nil {
		return err
	}

	conn, err := amqp.Dial(r.metadata.Host)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	ch.Qos(r.metadata.PrefetchCount, 0, true)
	r.connection = conn
	r.channel = ch

	q, err := r.declareQueue()
	if err != nil {
		return err
	}

	r.queue = q

	return nil
}

func (r *RabbitMQ) Operations() []bindings.OperationKind {
	return []bindings.OperationKind{bindings.CreateOperation}
}

func (r *RabbitMQ) Invoke(req *bindings.InvokeRequest) (*bindings.InvokeResponse, error) {
	pub := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         req.Data,
	}

	contentType, ok := contrib_metadata.TryGetContentType(req.Metadata)

	if ok {
		pub.ContentType = contentType
	}

	ttl, ok, err := contrib_metadata.TryGetTTL(req.Metadata)
	if err != nil {
		return nil, err
	}

	// The default time to live has been set in the queue
	// We allow overriding on each call, by setting a value in request metadata
	if ok {
		// RabbitMQ expects the duration in ms
		pub.Expiration = strconv.FormatInt(ttl.Milliseconds(), 10)
	}

	priority, ok, err := contrib_metadata.TryGetPriority(req.Metadata)
	if err != nil {
		return nil, err
	}

	if ok {
		pub.Priority = priority
	}

	err = r.channel.Publish("", r.metadata.QueueName, false, false, pub)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *RabbitMQ) parseMetadata(metadata bindings.Metadata) error {
	m := rabbitMQMetadata{}

	if val, ok := metadata.Properties[host]; ok && val != "" {
		m.Host = val
	} else {
		return errors.New("rabbitMQ binding error: missing host address")
	}

	if val, ok := metadata.Properties[queueName]; ok && val != "" {
		m.QueueName = val
	} else {
		return errors.New("rabbitMQ binding error: missing queue Name")
	}

	if val, ok := metadata.Properties[durable]; ok && val != "" {
		d, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("rabbitMQ binding error: can't parse durable field: %s", err)
		}
		m.Durable = d
	}

	if val, ok := metadata.Properties[deleteWhenUnused]; ok && val != "" {
		d, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("rabbitMQ binding error: can't parse deleteWhenUnused field: %s", err)
		}
		m.DeleteWhenUnused = d
	}

	if val, ok := metadata.Properties[prefetchCount]; ok && val != "" {
		parsedVal, err := strconv.ParseInt(val, defaultBase, defaultBitSize)
		if err != nil {
			return fmt.Errorf("rabbitMQ binding error: can't parse prefetchCount field: %s", err)
		}
		m.PrefetchCount = int(parsedVal)
	}

	if val, ok := metadata.Properties[exclusive]; ok && val != "" {
		d, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("rabbitMQ binding error: can't parse exclusive field: %s", err)
		}
		m.Exclusive = d
	}

	if val, ok := metadata.Properties[maxPriority]; ok && val != "" {
		parsedVal, err := strconv.ParseUint(val, defaultBase, defaultBitSize)
		if err != nil {
			return fmt.Errorf("rabbitMQ binding error: can't parse maxPriority field: %s", err)
		}

		maxPriority := uint8(parsedVal)
		if parsedVal > 255 {
			// Overflow
			maxPriority = math.MaxUint8
		}

		m.MaxPriority = &maxPriority
	}

	ttl, ok, err := contrib_metadata.TryGetTTL(metadata.Properties)
	if err != nil {
		return err
	}

	if ok {
		m.defaultQueueTTL = &ttl
	}

	r.metadata = m

	return nil
}

func (r *RabbitMQ) declareQueue() (amqp.Queue, error) {
	args := amqp.Table{}
	if r.metadata.defaultQueueTTL != nil {
		// Value in ms
		ttl := *r.metadata.defaultQueueTTL / time.Millisecond
		args[rabbitMQQueueMessageTTLKey] = int(ttl)
	}

	if r.metadata.MaxPriority != nil {
		args[rabbitMQMaxPriorityKey] = *r.metadata.MaxPriority
	}

	return r.channel.QueueDeclare(r.metadata.QueueName, r.metadata.Durable, r.metadata.DeleteWhenUnused, r.metadata.Exclusive, false, args)
}

func (r *RabbitMQ) Read(handler func(*bindings.ReadResponse) ([]byte, error)) error {
	msgs, err := r.channel.Consume(
		r.queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			_, err := handler(&bindings.ReadResponse{
				Data: d.Body,
			})
			if err == nil {
				r.channel.Ack(d.DeliveryTag, false)
			}
		}
	}()

	<-forever

	return nil
}
