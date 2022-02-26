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
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"

	"github.com/bhojpur/service/pkg/state"
)

type tLogger interface {
	Fatalf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

func runServerWithOptions(opts server.Options) *server.Server {
	return natsserver.RunServer(&opts)
}

func runServerOnPort(port int) *server.Server {
	opts := natsserver.DefaultTestOptions
	opts.Port = port
	opts.JetStream = true
	opts.Cluster.Name = "testing"
	return runServerWithOptions(opts)
}

func runDefaultServer() *server.Server {
	return runServerOnPort(nats.DefaultPort)
}

func newDefaultConnection(t tLogger) *nats.Conn {
	return newConnection(t, nats.DefaultPort)
}

func newConnection(t tLogger, port int) *nats.Conn {
	url := fmt.Sprintf("nats://127.0.0.1:%d", port)
	nc, err := nats.Connect(url)
	if err != nil {
		t.Fatalf("Failed to create default connection: %v\n", err)
		return nil
	}
	return nc
}

func connectAndCreateBucket(t *testing.T) (nats.KeyValue, *nats.Conn) {
	nc := newDefaultConnection(t)
	jsc, err := nc.JetStream()
	if err != nil {
		t.Fatalf("Could not open jetstream: %v\n", err)
		return nil, nil
	}
	kv, err := jsc.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: "test",
	})
	if err != nil {
		t.Fatalf("Could not open jetstream: %v\n", err)
		return nil, nil
	}
	return kv, nc
}

func TestDefaultConnection(t *testing.T) {
	s := runDefaultServer()
	defer s.Shutdown()

	_, nc := connectAndCreateBucket(t)
	defer nc.Close()
}

func TestSetGetAndDelete(t *testing.T) {
	s := runDefaultServer()
	defer s.Shutdown()

	_, nc := connectAndCreateBucket(t)
	nc.Close()

	store := NewJetstreamStateStore(nil)

	err := store.Init(state.Metadata{
		Properties: map[string]string{
			"natsURL": nats.DefaultURL,
			"bucket":  "test",
		},
	})
	if err != nil {
		t.Fatalf("Could not init: %v\n", err)
		return
	}

	tkey := "key"
	tData := map[string]string{
		"dkey": "dvalue",
	}

	err = store.Set(&state.SetRequest{
		Key:   tkey,
		Value: tData,
	})
	if err != nil {
		t.Fatalf("Could not set: %v\n", err)
		return
	}

	resp, err := store.Get(&state.GetRequest{
		Key: tkey,
	})
	if err != nil {
		t.Fatalf("Could not get: %v\n", err)
		return
	}
	rData := make(map[string]string)
	json.Unmarshal(resp.Data, &rData)
	if !reflect.DeepEqual(rData, tData) {
		t.Fatal("Response data does not match written data\n")
	}

	err = store.Delete(&state.DeleteRequest{
		Key: tkey,
	})
	if err != nil {
		t.Fatalf("Could not delete: %v\n", err)
		return
	}

	_, err = store.Get(&state.GetRequest{
		Key: tkey,
	})
	if err == nil {
		t.Fatal("Could get after delete\n")
		return
	}
}
