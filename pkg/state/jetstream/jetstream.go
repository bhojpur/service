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
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/bhojpur/service/pkg/state"
	"github.com/bhojpur/service/pkg/state/utils"
	"github.com/bhojpur/service/pkg/utils/logger"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
)

// StateStore is a NATS jetstream KV state store.
type StateStore struct {
	state.DefaultBulkStore
	nc     *nats.Conn
	json   jsoniter.API
	bucket nats.KeyValue
	logger logger.Logger
}

type jetstreamMetadata struct {
	name    string
	natsURL string
	jwt     string
	seedKey string
	bucket  string
}

// NewJetstreamStateStore returns a new NATS jetstream KV state store.
func NewJetstreamStateStore(logger logger.Logger) state.Store {
	s := &StateStore{
		json:   jsoniter.ConfigFastest,
		logger: logger,
	}
	s.DefaultBulkStore = state.NewDefaultBulkStore(s)

	return s
}

// Init does parse metadata and establishes connection to nats broker.
func (js *StateStore) Init(metadata state.Metadata) error {
	meta, err := js.getMetadata(metadata)
	if err != nil {
		return err
	}

	var opts []nats.Option
	opts = append(opts, nats.Name(meta.name))

	// Set nats.UserJWT options when jwt and seed key is provided.
	if meta.jwt != "" && meta.seedKey != "" {
		opts = append(opts, nats.UserJWT(func() (string, error) {
			return meta.jwt, nil
		}, func(nonce []byte) ([]byte, error) {
			return sigHandler(meta.seedKey, nonce)
		}))
	}

	js.nc, err = nats.Connect(meta.natsURL, opts...)
	if err != nil {
		return err
	}

	jsc, err := js.nc.JetStream()
	if err != nil {
		return err
	}

	js.bucket, err = jsc.KeyValue(meta.bucket)
	if err != nil {
		return err
	}

	return nil
}

func (js *StateStore) Ping() error {
	return nil
}

func (js *StateStore) Features() []state.Feature {
	return nil
}

// Get retrieves state with a key.
func (js *StateStore) Get(req *state.GetRequest) (*state.GetResponse, error) {
	entry, err := js.bucket.Get(escape(req.Key))
	if err != nil {
		return nil, err
	}

	return &state.GetResponse{
		Data: entry.Value(),
	}, nil
}

// Set stores value for a key.
func (js *StateStore) Set(req *state.SetRequest) error {
	bt, _ := utils.Marshal(req.Value, js.json.Marshal)
	_, err := js.bucket.Put(escape(req.Key), bt)
	return err
}

// Delete performs a delete operation.
func (js *StateStore) Delete(req *state.DeleteRequest) error {
	return js.bucket.Delete(escape(req.Key))
}

func (js *StateStore) getMetadata(metadata state.Metadata) (jetstreamMetadata, error) {
	var m jetstreamMetadata

	if v, ok := metadata.Properties["natsURL"]; ok && v != "" {
		m.natsURL = v
	} else {
		return jetstreamMetadata{}, fmt.Errorf("missing nats URL")
	}

	m.jwt = metadata.Properties["jwt"]
	m.seedKey = metadata.Properties["seedKey"]

	if m.jwt != "" && m.seedKey == "" {
		return jetstreamMetadata{}, fmt.Errorf("missing seed key")
	}

	if m.jwt == "" && m.seedKey != "" {
		return jetstreamMetadata{}, fmt.Errorf("missing jwt")
	}

	if m.name = metadata.Properties["name"]; m.name == "" {
		m.name = "bhojpur.net - statestore.jetstream"
	}

	if m.bucket = metadata.Properties["bucket"]; m.bucket == "" {
		return jetstreamMetadata{}, fmt.Errorf("missing bucket")
	}

	return m, nil
}

// Handle NATS signature request for challenge response authentication.
func sigHandler(seedKey string, nonce []byte) ([]byte, error) {
	kp, err := nkeys.FromSeed([]byte(seedKey))
	if err != nil {
		return nil, err
	}
	// Wipe our key on exit.
	defer kp.Wipe()

	sig, _ := kp.Sign(nonce)
	return sig, nil
}

// Escape Bhojpur Service keys, because || is forbidden.
func escape(key string) string {
	return strings.ReplaceAll(key, "||", ".")
}
