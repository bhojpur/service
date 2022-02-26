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
	"reflect"
	"testing"
	"time"

	"github.com/bhojpur/service/pkg/pubsub"
)

func TestParseMetadata(t *testing.T) {
	testCases := []struct {
		desc      string
		input     pubsub.Metadata
		want      metadata
		expectErr bool
	}{
		{
			desc: "Valid Metadata",
			input: pubsub.Metadata{
				Properties: map[string]string{
					"natsURL":        "nats://localhost:4222",
					"name":           "myName",
					"durableName":    "myDurable",
					"queueGroupName": "myQueue",
					"startSequence":  "1",
					"startTime":      "1629328511",
					"deliverAll":     "true",
					"flowControl":    "true",
				},
			},
			want: metadata{
				natsURL:        "nats://localhost:4222",
				name:           "myName",
				durableName:    "myDurable",
				queueGroupName: "myQueue",
				startSequence:  1,
				startTime:      time.Unix(1629328511, 0),
				deliverAll:     true,
				flowControl:    true,
			},
			expectErr: false,
		},
		{
			desc: "Invalid metadata with missing seed key",
			input: pubsub.Metadata{
				Properties: map[string]string{
					"natsURL":        "nats://localhost:4222",
					"name":           "myName",
					"durableName":    "myDurable",
					"queueGroupName": "myQueue",
					"startSequence":  "1",
					"startTime":      "1629328511",
					"deliverAll":     "true",
					"flowControl":    "true",
					"jwt":            "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
				},
			},
			want:      metadata{},
			expectErr: true,
		},
		{
			desc: "Invalid metadata with missing jwt",
			input: pubsub.Metadata{
				Properties: map[string]string{
					"natsURL":        "nats://localhost:4222",
					"name":           "myName",
					"durableName":    "myDurable",
					"queueGroupName": "myQueue",
					"startSequence":  "1",
					"startTime":      "1629328511",
					"deliverAll":     "true",
					"flowControl":    "true",
					"seedKey":        "SUACS34K232OKPRDOMKC6QEWXWUDJTT6R6RZM2WPMURUS5Z3POU7BNIL4Y",
				},
			},
			want:      metadata{},
			expectErr: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got, err := parseMetadata(tC.input)
			if !tC.expectErr && err != nil {
				t.Fatal(err)
			}
			if tC.expectErr && err == nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tC.want) {
				t.Fatalf("unexpected metadata: got=%v, want=%v", got, tC.want)
			}
		})
	}
}
