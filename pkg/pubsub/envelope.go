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
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"

	svc_contenttype "github.com/bhojpur/service/pkg/contenttype"
	svc_metadata "github.com/bhojpur/service/pkg/metadata"
)

const (
	// DefaultCloudEventType is the default event type for a Bhojpur Service published event.
	DefaultCloudEventType = "net.bhojpur.event.sent"
	// CloudEventsSpecVersion is the specversion used by Bhojpur Service for the cloud events implementation.
	CloudEventsSpecVersion = "1.0"
	// DefaultCloudEventSource is the default event source.
	DefaultCloudEventSource = "Bhojpur"
	// DefaultCloudEventDataContentType is the default content-type for the data attribute.
	DefaultCloudEventDataContentType = "text/plain"
	TraceIDField                     = "traceid"
	TraceStateField                  = "tracestate"
	TopicField                       = "topic"
	PubsubField                      = "pubsubname"
	ExpirationField                  = "expiration"
	DataContentTypeField             = "datacontenttype"
	DataField                        = "data"
	DataBase64Field                  = "data_base64"
	SpecVersionField                 = "specversion"
	TypeField                        = "type"
	SourceField                      = "source"
	IDField                          = "id"
	SubjectField                     = "subject"
)

// NewCloudEventsEnvelope returns a map representation of a cloudevents JSON.
func NewCloudEventsEnvelope(id, source, eventType, subject string, topic string, pubsubName string,
	dataContentType string, data []byte, traceID string, traceState string) map[string]interface{} {
	// defaults
	if id == "" {
		id = uuid.New().String()
	}
	if source == "" {
		source = DefaultCloudEventSource
	}
	if eventType == "" {
		eventType = DefaultCloudEventType
	}
	if dataContentType == "" {
		dataContentType = DefaultCloudEventDataContentType
	}

	var ceData interface{}
	ceDataField := DataField
	var err error
	if svc_contenttype.IsJSONContentType(dataContentType) {
		err = jsoniter.Unmarshal(data, &ceData)
	} else if svc_contenttype.IsBinaryContentType(dataContentType) {
		ceData = base64.StdEncoding.EncodeToString(data)
		ceDataField = DataBase64Field
	} else {
		ceData = string(data)
	}

	if err != nil {
		ceData = string(data)
	}

	ce := map[string]interface{}{
		IDField:              id,
		SpecVersionField:     CloudEventsSpecVersion,
		DataContentTypeField: dataContentType,
		SourceField:          source,
		TypeField:            eventType,
		TopicField:           topic,
		PubsubField:          pubsubName,
		TraceIDField:         traceID,
		TraceStateField:      traceState,
	}

	ce[ceDataField] = ceData

	if subject != "" {
		ce[SubjectField] = subject
	}

	return ce
}

// FromCloudEvent returns a map representation of an existing cloudevents JSON.
func FromCloudEvent(cloudEvent []byte, topic, pubsub, traceID string, traceState string) (map[string]interface{}, error) {
	var m map[string]interface{}
	err := jsoniter.Unmarshal(cloudEvent, &m)
	if err != nil {
		return m, err
	}

	m[TraceIDField] = traceID
	m[TraceStateField] = traceState
	m[TopicField] = topic
	m[PubsubField] = pubsub

	// specify default value if it's unspecified from the original CloudEvent
	if m[SourceField] == nil {
		m[SourceField] = DefaultCloudEventSource
	}

	if m[TypeField] == nil {
		m[TypeField] = DefaultCloudEventType
	}

	if m[SpecVersionField] == nil {
		m[SpecVersionField] = CloudEventsSpecVersion
	}

	return m, nil
}

// FromRawPayload returns a CloudEvent for a raw payload on subscriber's end.
func FromRawPayload(data []byte, topic, pubsub string) map[string]interface{} {
	// Limitations of generating the CloudEvent on the subscriber side based on raw payload:
	// - The CloudEvent ID will be random, so the same message can be redelivered as a different ID.
	// - TraceID is not useful since it is random and not from publisher side.
	// - Data is always returned as `data_base64` since we don't know the actual content type.
	return map[string]interface{}{
		IDField:              uuid.New().String(),
		SpecVersionField:     CloudEventsSpecVersion,
		DataContentTypeField: "application/octet-stream",
		SourceField:          DefaultCloudEventSource,
		TypeField:            DefaultCloudEventType,
		TopicField:           topic,
		PubsubField:          pubsub,
		DataBase64Field:      base64.StdEncoding.EncodeToString(data),
	}
}

// HasExpired determines if the current cloud event has expired.
func HasExpired(cloudEvent map[string]interface{}) bool {
	e, ok := cloudEvent[ExpirationField]
	if ok && e != "" {
		expiration, err := time.Parse(time.RFC3339, fmt.Sprintf("%s", e))
		if err != nil {
			return false
		}

		return expiration.UTC().Before(time.Now().UTC())
	}

	return false
}

// ApplyMetadata will process metadata to modify the cloud event based on the component's feature set.
func ApplyMetadata(cloudEvent map[string]interface{}, componentFeatures []Feature, metadata map[string]string) {
	ttl, hasTTL, _ := svc_metadata.TryGetTTL(metadata)
	if hasTTL && !FeatureMessageTTL.IsPresent(componentFeatures) {
		// Bhojpur Service only handles Message TTL if component does not.
		now := time.Now().UTC()
		// The maximum ttl is maxInt64, which is not enough to overflow time, for now.
		// As of the time this code was written (2020 Dec 28th),
		// the maximum time of now() adding maxInt64 is ~ "2313-04-09T23:30:26Z".
		// Max time in golang is currently 292277024627-12-06T15:30:07.999999999Z.
		// So, we have some time before the overflow below happens :)
		expiration := now.Add(ttl)
		cloudEvent[ExpirationField] = expiration.Format(time.RFC3339)
	}
}
