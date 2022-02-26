package jetstream

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
	"fmt"
	"strconv"
	"time"

	"github.com/bhojpur/service/pkg/pubsub"
)

type metadata struct {
	natsURL string
	jwt     string
	seedKey string

	name           string
	durableName    string
	queueGroupName string
	startSequence  uint64
	startTime      time.Time
	deliverAll     bool
	flowControl    bool
}

func parseMetadata(psm pubsub.Metadata) (metadata, error) {
	var m metadata

	if v, ok := psm.Properties["natsURL"]; ok && v != "" {
		m.natsURL = v
	} else {
		return metadata{}, fmt.Errorf("missing nats URL")
	}

	m.jwt = psm.Properties["jwt"]
	m.seedKey = psm.Properties["seedKey"]

	if m.jwt != "" && m.seedKey == "" {
		return metadata{}, fmt.Errorf("missing seed key")
	}

	if m.jwt == "" && m.seedKey != "" {
		return metadata{}, fmt.Errorf("missing jwt")
	}

	if m.name = psm.Properties["name"]; m.name == "" {
		m.name = "bhojpur.net - pubsub.jetstream"
	}

	m.durableName = psm.Properties["durableName"]
	m.queueGroupName = psm.Properties["queueGroupName"]

	if v, err := strconv.ParseUint(psm.Properties["startSequence"], 10, 64); err == nil {
		m.startSequence = v
	}

	if v, err := strconv.ParseInt(psm.Properties["startTime"], 10, 64); err == nil {
		m.startTime = time.Unix(v, 0)
	}

	if v, err := strconv.ParseBool(psm.Properties["deliverAll"]); err == nil {
		m.deliverAll = v
	}

	if v, err := strconv.ParseBool(psm.Properties["flowControl"]); err == nil {
		m.flowControl = v
	}

	return m, nil
}
