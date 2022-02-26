package statechange

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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/utils/logger"
)

func getTestMetadata() map[string]string {
	return map[string]string{
		"address":  "127.0.0.1:28015",
		"database": "app",
		"username": "admin",
		"password": "rethinkdb",
		"table":    "appstate",
	}
}

func getNewRethinkActorBinding() *Binding {
	l := logger.NewLogger("test")
	if os.Getenv("DEBUG") != "" {
		l.SetOutputLevel(logger.DebugLevel)
	}

	return NewRethinkDBStateChangeBinding(l)
}

/*
go test github.com/bhojpur/service/pkg/bindings/rethinkdb/statechange \
	-run ^TestBinding$ -count 1
*/

func TestBinding(t *testing.T) {
	if os.Getenv("RUN_LIVE_RETHINKDB_TEST") != "true" {
		t.SkipNow()
	}
	testDuration := 10 * time.Second
	testDurationStr := os.Getenv("RETHINKDB_TEST_DURATION")
	if testDurationStr != "" {
		d, err := time.ParseDuration(testDurationStr)
		if err != nil {
			t.Fatalf("invalid test duration: %s, expected time.Duration", testDurationStr)
		}
		testDuration = d
	}

	m := bindings.Metadata{
		Name:       "test",
		Properties: getTestMetadata(),
	}
	assert.NotNil(t, m.Properties)

	b := getNewRethinkActorBinding()
	err := b.Init(m)
	assert.NoErrorf(t, err, "error initializing")

	go func() {
		err = b.Read(func(res *bindings.ReadResponse) ([]byte, error) {
			assert.NotNil(t, res)
			t.Logf("state change event:\n%s", string(res.Data))

			return nil, nil
		})
		assert.NoErrorf(t, err, "error on read")
	}()

	testTimer := time.AfterFunc(testDuration, func() {
		t.Log("done")
		b.stopCh <- true
	})
	defer testTimer.Stop()
	<-b.stopCh
}
