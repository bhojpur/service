package rocketmq

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
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/utils/logger"
)

func TestInputBindingRead(t *testing.T) { //nolint:paralleltest
	if !isLiveTest() {
		return
	}
	m := bindings.Metadata{} //nolint:exhaustivestruct
	m.Properties = getTestMetadata()
	r := NewAliCloudRocketMQ(logger.NewLogger("test"))
	err := r.Init(m)
	require.NoError(t, err)

	var count int32
	handler := func(in *bindings.ReadResponse) ([]byte, error) {
		require.Equal(t, "hello", string(in.Data))
		atomic.AddInt32(&count, 1)

		return nil, nil
	}
	go func() {
		err = r.Read(handler)
		require.NoError(t, err)
	}()

	time.Sleep(5 * time.Second)
	atomic.StoreInt32(&count, 0)
	req := &bindings.InvokeRequest{Data: []byte("hello"), Operation: bindings.CreateOperation, Metadata: map[string]string{}}
	_, err = r.Invoke(req)
	require.NoError(t, err)

	time.Sleep(10 * time.Second)
	for i := 0; i < 30; i++ {
		if atomic.LoadInt32(&count) > 0 {
			break
		}
		time.Sleep(time.Second)
	}
	assert.Equal(t, int32(1), atomic.LoadInt32(&count))
}

func isLiveTest() bool {
	return os.Getenv("RUN_LIVE_ROCKETMQ_TEST") == "true"
}

func getTestMetadata() map[string]string {
	return map[string]string{
		"accessProto":        "tcp",
		"nameServer":         "http://xx.mq-internet-access.mq-internet.aliyuncs.com:80",
		"consumerGroup":      "GID_APP-MQ-TCP",
		"topics":             "TOPIC_TEST",
		"accessKey":          "xx",
		"secretKey":          "xx",
		"instanceId":         "MQ_INST_xx",
		"consumerBatchSize":  "1",
		"consumerThreadNums": "5",
	}
}
