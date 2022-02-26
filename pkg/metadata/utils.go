package metadata

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
	"math"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	// TTLMetadataKey defines the metadata key for setting a time to live (in seconds).
	TTLMetadataKey = "ttlInSeconds"

	// RawPayloadKey defines the metadata key for forcing raw payload in pubsub.
	RawPayloadKey = "rawPayload"

	// PriorityMetadataKey defines the metadata key for setting a priority.
	PriorityMetadataKey = "priority"

	// ContentType defines the metadata key for the content type.
	ContentType = "contentType"
)

// TryGetTTL tries to get the ttl as a time.Duration value for pubsub, binding and any other building block.
func TryGetTTL(props map[string]string) (time.Duration, bool, error) {
	if val, ok := props[TTLMetadataKey]; ok && val != "" {
		valInt64, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return 0, false, errors.Wrapf(err, "%s value must be a valid integer: actual is '%s'", TTLMetadataKey, val)
		}

		if valInt64 <= 0 {
			return 0, false, fmt.Errorf("%s value must be higher than zero: actual is %d", TTLMetadataKey, valInt64)
		}

		duration := time.Duration(valInt64) * time.Second
		if duration < 0 {
			// Overflow
			duration = math.MaxInt64
		}

		return duration, true, nil
	}

	return 0, false, nil
}

// TryGetPriority tries to get the priority for binding and any other building block.
func TryGetPriority(props map[string]string) (uint8, bool, error) {
	if val, ok := props[PriorityMetadataKey]; ok && val != "" {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			return 0, false, errors.Wrapf(err, "%s value must be a valid integer: actual is '%s'", PriorityMetadataKey, val)
		}

		priority := uint8(intVal)
		if intVal < 0 {
			priority = 0
		} else if intVal > 255 {
			priority = math.MaxUint8
		}

		return priority, true, nil
	}

	return 0, false, nil
}

// IsRawPayload determines if payload should be used as-is.
func IsRawPayload(props map[string]string) (bool, error) {
	if val, ok := props[RawPayloadKey]; ok && val != "" {
		boolVal, err := strconv.ParseBool(val)
		if err != nil {
			return false, errors.Wrapf(err, "%s value must be a valid boolean: actual is '%s'", RawPayloadKey, val)
		}

		return boolVal, nil
	}

	return false, nil
}

func TryGetContentType(props map[string]string) (string, bool) {
	if val, ok := props[ContentType]; ok && val != "" {
		return val, true
	}

	return "", false
}
