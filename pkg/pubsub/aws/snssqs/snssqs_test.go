package snssqs

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
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bhojpur/service/pkg/pubsub"
	"github.com/bhojpur/service/pkg/utils/logger"
)

type testUnitFixture struct {
	metadata pubsub.Metadata
	name     string
}

func Test_parseTopicArn(t *testing.T) {
	t.Parallel()
	// no further guarantees are made about this function.
	r := require.New(t)
	tSnsMessage := &snsMessage{TopicArn: "arn:aws:sqs:us-east-1:000000000000:qqnoob"}
	r.Equal("qqnoob", tSnsMessage.parseTopicArn())
}

// Verify that all metadata ends up in the correct spot.
func Test_getSnsSqsMetatdata_AllConfiguration(t *testing.T) {
	t.Parallel()
	r := require.New(t)
	l := logger.NewLogger("SnsSqs unit test")
	l.SetOutputLevel(logger.DebugLevel)
	ps := snsSqs{
		logger: l,
	}

	md, err := ps.getSnsSqsMetatdata(pubsub.Metadata{Properties: map[string]string{
		"consumerID":               "consumer",
		"Endpoint":                 "endpoint",
		"accessKey":                "a",
		"secretKey":                "s",
		"sessionToken":             "t",
		"region":                   "r",
		"sqsDeadLettersQueueName":  "q",
		"messageVisibilityTimeout": "2",
		"messageRetryLimit":        "3",
		"messageWaitTimeSeconds":   "4",
		"messageMaxNumber":         "5",
		"messageReceiveLimit":      "6",
	}})

	r.NoError(err)

	r.Equal("consumer", md.sqsQueueName)
	r.Equal("endpoint", md.Endpoint)
	r.Equal("a", md.AccessKey)
	r.Equal("s", md.SecretKey)
	r.Equal("t", md.SessionToken)
	r.Equal("r", md.Region)
	r.Equal("q", md.sqsDeadLettersQueueName)
	r.Equal(int64(2), md.messageVisibilityTimeout)
	r.Equal(int64(3), md.messageRetryLimit)
	r.Equal(int64(4), md.messageWaitTimeSeconds)
	r.Equal(int64(5), md.messageMaxNumber)
	r.Equal(int64(6), md.messageReceiveLimit)
}

func Test_getSnsSqsMetatdata_defaults(t *testing.T) {
	t.Parallel()
	r := require.New(t)
	l := logger.NewLogger("SnsSqs unit test")
	l.SetOutputLevel(logger.DebugLevel)
	ps := snsSqs{
		logger: l,
	}

	md, err := ps.getSnsSqsMetatdata(pubsub.Metadata{Properties: map[string]string{
		"consumerID": "c",
		"accessKey":  "a",
		"secretKey":  "s",
		"region":     "r",
	}})

	r.NoError(err)

	r.Equal("c", md.sqsQueueName)
	r.Equal("", md.Endpoint)
	r.Equal("a", md.AccessKey)
	r.Equal("s", md.SecretKey)
	r.Equal("", md.SessionToken)
	r.Equal("r", md.Region)
	r.Equal(int64(10), md.messageVisibilityTimeout)
	r.Equal(int64(10), md.messageRetryLimit)
	r.Equal(int64(1), md.messageWaitTimeSeconds)
	r.Equal(int64(10), md.messageMaxNumber)
	r.Equal(false, md.disableEntityManagement)
	r.Equal(float64(5), md.assetsManagementTimeoutSeconds)
	r.Equal(false, md.disableDeleteOnRetryLimit)
}

func Test_getSnsSqsMetatdata_legacyaliases(t *testing.T) {
	t.Parallel()
	r := require.New(t)
	l := logger.NewLogger("SnsSqs unit test")
	l.SetOutputLevel(logger.DebugLevel)
	ps := snsSqs{
		logger: l,
	}

	md, err := ps.getSnsSqsMetatdata(pubsub.Metadata{Properties: map[string]string{
		"consumerID":   "consumer",
		"awsAccountID": "acctId",
		"awsSecret":    "secret",
		"awsRegion":    "region",
	}})

	r.NoError(err)

	r.Equal("consumer", md.sqsQueueName)
	r.Equal("", md.Endpoint)
	r.Equal("acctId", md.AccessKey)
	r.Equal("secret", md.SecretKey)
	r.Equal("region", md.Region)
	r.Equal(int64(10), md.messageVisibilityTimeout)
	r.Equal(int64(10), md.messageRetryLimit)
	r.Equal(int64(1), md.messageWaitTimeSeconds)
	r.Equal(int64(10), md.messageMaxNumber)
}

func testMetadataParsingShouldFail(t *testing.T, metadata pubsub.Metadata, l logger.Logger) {
	t.Parallel()
	r := require.New(t)

	ps := snsSqs{
		logger: l,
	}

	md, err := ps.getSnsSqsMetatdata(metadata)

	r.Error(err)
	r.Nil(md)
}

func Test_getSnsSqsMetatdata_invalidMetadataSetup(t *testing.T) {
	t.Parallel()

	fixtures := []testUnitFixture{
		{
			metadata: pubsub.Metadata{Properties: map[string]string{
				"consumerID": "consumer",
				"Endpoint":   "endpoint",
				"AccessKey":  "acctId",
				"SecretKey":  "secret",
				"awsToken":   "token",
				"Region":     "region",
				"fifo":       "none bool",
			}},
			name: "fifo not set to boolean",
		},
		{
			metadata: pubsub.Metadata{Properties: map[string]string{
				"consumerID":          "consumer",
				"Endpoint":            "endpoint",
				"AccessKey":           "acctId",
				"SecretKey":           "secret",
				"awsToken":            "token",
				"Region":              "region",
				"messageReceiveLimit": "100",
			}},
			name: "deadletters receive limit without deadletters queue name",
		},
		{
			metadata: pubsub.Metadata{Properties: map[string]string{
				"consumerID":              "consumer",
				"Endpoint":                "endpoint",
				"AccessKey":               "acctId",
				"SecretKey":               "secret",
				"awsToken":                "token",
				"Region":                  "region",
				"sqsDeadLettersQueueName": "my-queue",
			}},
			name: "deadletters message queue without deadletters receive limit",
		},
		{
			metadata: pubsub.Metadata{Properties: map[string]string{
				"consumerID":                "consumer",
				"Endpoint":                  "endpoint",
				"AccessKey":                 "acctId",
				"SecretKey":                 "secret",
				"awsToken":                  "token",
				"Region":                    "region",
				"sqsDeadLettersQueueName":   "my-queue",
				"messageReceiveLimit":       "9",
				"disableDeleteOnRetryLimit": "true",
			}},
			name: "deadletters message queue with disableDeleteOnRetryLimit",
		},
		{
			metadata: pubsub.Metadata{Properties: map[string]string{
				"consumerID":       "consumer",
				"Endpoint":         "endpoint",
				"AccessKey":        "acctId",
				"SecretKey":        "secret",
				"awsToken":         "token",
				"Region":           "region",
				"messageMaxNumber": "-100",
			}},
			name: "illigal message max number (negative, too low)",
		},
		{
			metadata: pubsub.Metadata{Properties: map[string]string{
				"consumerID":       "consumer",
				"Endpoint":         "endpoint",
				"AccessKey":        "acctId",
				"SecretKey":        "secret",
				"awsToken":         "token",
				"Region":           "region",
				"messageMaxNumber": "100",
			}},
			name: "illigal message max number (too high)",
		},
		{
			metadata: pubsub.Metadata{Properties: map[string]string{
				"consumerID":             "consumer",
				"Endpoint":               "endpoint",
				"AccessKey":              "acctId",
				"SecretKey":              "secret",
				"awsToken":               "token",
				"Region":                 "region",
				"messageWaitTimeSeconds": "0",
			}},
			name: "invalid wait time seconds (too low)",
		},
		{
			metadata: pubsub.Metadata{Properties: map[string]string{
				"consumerID":               "consumer",
				"Endpoint":                 "endpoint",
				"AccessKey":                "acctId",
				"SecretKey":                "secret",
				"awsToken":                 "token",
				"Region":                   "region",
				"messageVisibilityTimeout": "-100",
			}},
			name: "invalid message visibility",
		},
		{
			metadata: pubsub.Metadata{Properties: map[string]string{
				"consumerID":        "consumer",
				"Endpoint":          "endpoint",
				"AccessKey":         "acctId",
				"SecretKey":         "secret",
				"awsToken":          "token",
				"Region":            "region",
				"messageRetryLimit": "-100",
			}},
			name: "invalid message retry limit",
		},
		// disableEntityManagement
		{
			metadata: pubsub.Metadata{Properties: map[string]string{
				"consumerID":              "consumer",
				"Endpoint":                "endpoint",
				"AccessKey":               "acctId",
				"SecretKey":               "secret",
				"awsToken":                "token",
				"Region":                  "region",
				"messageRetryLimit":       "10",
				"disableEntityManagement": "y",
			}},
			name: "invalid message disableEntityManagement",
		},
	}

	l := logger.NewLogger("SnsSqs unit test")
	l.SetOutputLevel(logger.DebugLevel)

	for _, tc := range fixtures {
		t.Run(tc.name, func(t *testing.T) {
			testMetadataParsingShouldFail(t, tc.metadata, l)
		})
	}
}

func Test_parseInt64(t *testing.T) {
	t.Parallel()
	r := require.New(t)
	number, err := parseInt64("applesauce", "propertyName")
	r.EqualError(err, "parsing propertyName failed with: strconv.Atoi: parsing \"applesauce\": invalid syntax")
	r.Equal(int64(-1), number)

	number, _ = parseInt64("1000", "")
	r.Equal(int64(1000), number)

	number, _ = parseInt64("-1000", "")
	r.Equal(int64(-1000), number)

	// Expecting that this function doesn't panic.
	_, err = parseInt64("999999999999999999999999999999999999999999999999999999999999999999999999999", "")
	r.Error(err)
}

func Test_replaceNameToAWSSanitizedName(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	s := `Some_invalid-name // for an AWS resource &*()*&&^Some invalid name // for an AWS resource &*()*&&^Some invalid 
		name // for an AWS resource &*()*&&^Some invalid name // for an AWS resource &*()*&&^Some invalid name // for an
		AWS resource &*()*&&^Some invalid name // for an AWS resource &*()*&&^`
	v := nameToAWSSanitizedName(s, false)
	r.Equal(80, len(v))
	r.Equal("Some_invalid-nameforanAWSresourceSomeinvalidnameforanAWSresourceSomeinvalidnamef", v)
}

func Test_replaceNameToAWSSanitizedFifoName_Trimmed(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	s := `Some_invalid-name // for an AWS resource &*()*&&^Some invalid name // for an AWS resource &*()*&&^Some invalid 
		name // for an AWS resource &*()*&&^Some invalid name // for an AWS resource &*()*&&^Some invalid name // for an
		AWS resource &*()*&&^Some invalid name // for an AWS resource &*()*&&^`
	v := nameToAWSSanitizedName(s, true)
	r.Equal(80, len(v))
	r.Equal("Some_invalid-nameforanAWSresourceSomeinvalidnameforanAWSresourceSomeinvalid.fifo", v)
}

func Test_replaceNameToAWSSanitizedFifoName_NonTrimmed(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	s := `012345678901234567890123456789012345678901234567890123456789012345678901234`
	v := nameToAWSSanitizedName(s, true)
	r.Equal(80, len(v))
	r.Equal("012345678901234567890123456789012345678901234567890123456789012345678901234.fifo", v)
}

func Test_replaceNameToAWSSanitizedExistingFifoName_NonTrimmed(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	s := `012345678901234567890123456789012345678901234567890123456789012345678901234.fifo`
	v := nameToAWSSanitizedName(s, true)
	r.Equal(80, len(v))
	r.Equal("012345678901234567890123456789012345678901234567890123456789012345678901234.fifo", v)
}

func Test_replaceNameToAWSSanitizedExistingFifoName_NonMax(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	s := `0123456789`
	v := nameToAWSSanitizedName(s, true)
	r.Equal(len(s)+len(".fifo"), len(v))
	r.Equal("0123456789.fifo", v)
}

func Test_replaceNameToAWSSanitizedExistingFifoName_NoFifoSetting(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	s := `012345678901234567890123456789012345678901234567890123456789012345678901234.fifo`
	v := nameToAWSSanitizedName(s, false)
	r.Equal(79, len(v))
	r.Equal("012345678901234567890123456789012345678901234567890123456789012345678901234fifo", v)
}

func Test_replaceNameToAWSSanitizedExistingFifoName_Trimmed(t *testing.T) {
	t.Parallel()
	r := require.New(t)

	s := `01234567890123456789012345678901234567890123456789012345678901234567890123456789.fifo`
	v := nameToAWSSanitizedName(s, true)
	r.Equal(80, len(v))
	r.Equal("012345678901234567890123456789012345678901234567890123456789012345678901234.fifo", v)
}
