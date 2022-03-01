package couchbase

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
	"strconv"

	"github.com/couchbase/gocb/v2"
	jsoniter "github.com/json-iterator/go"

	"github.com/bhojpur/service/pkg/state"
	"github.com/bhojpur/service/pkg/state/utils"
	"github.com/bhojpur/service/pkg/utils/logger"
)

const (
	couchbaseURL = "couchbaseURL"
	username     = "username"
	password     = "password"
	bucketName   = "bucketName"

	// see https://docs.couchbase.com/go-sdk/1.6/durability.html#configuring-durability
	numReplicasDurableReplication = "numReplicasDurableReplication"
	numReplicasDurablePersistence = "numReplicasDurablePersistence"
)

// Couchbase is a couchbase state store.
type Couchbase struct {
	state.DefaultBulkStore
	bucket                        *gocb.Bucket
	bucketName                    string // TODO: having bucket name sent as part of request (get,set etc.) metadata would be more flexible
	numReplicasDurableReplication uint
	numReplicasDurablePersistence uint
	json                          jsoniter.API

	features []state.Feature
	logger   logger.Logger
}

// NewCouchbaseStateStore returns a new couchbase state store.
func NewCouchbaseStateStore(logger logger.Logger) *Couchbase {
	s := &Couchbase{
		json:     jsoniter.ConfigFastest,
		features: []state.Feature{state.FeatureETag},
		logger:   logger,
	}
	s.DefaultBulkStore = state.NewDefaultBulkStore(s)

	return s
}

func validateMetadata(metadata state.Metadata) error {
	if metadata.Properties[couchbaseURL] == "" {
		return errors.New("couchbase error: couchbase URL is missing")
	}

	if metadata.Properties[username] == "" {
		return errors.New("couchbase error: couchbase username is missing")
	}

	if metadata.Properties[password] == "" {
		return errors.New("couchbase error: couchbase password is missing")
	}

	if metadata.Properties[bucketName] == "" {
		return errors.New("couchbase error: couchbase bucket name is missing")
	}

	v := metadata.Properties[numReplicasDurableReplication]
	if v != "" {
		_, err := strconv.ParseUint(v, 10, 0)
		if err != nil {
			return fmt.Errorf("couchbase error: %v", err)
		}
	}

	v = metadata.Properties[numReplicasDurablePersistence]
	if v != "" {
		_, err := strconv.ParseUint(v, 10, 0)
		if err != nil {
			return fmt.Errorf("couchbase error: %v", err)
		}
	}

	return nil
}

// Init does metadata and connection parsing.
func (cbs *Couchbase) Init(metadata state.Metadata) error {
	err := validateMetadata(metadata)
	if err != nil {
		return err
	}
	cbs.bucketName = metadata.Properties[bucketName]
	connString := metadata.Properties[couchbaseURL]
	c, err := gocb.Connect(connString, gocb.ClusterOptions{
		Username: metadata.Properties[username],
		Password: metadata.Properties[password],
	})
	//c, err := gocb.Connect(metadata.Properties[couchbaseURL])
	if err != nil {
		return fmt.Errorf("couchbase error: unable to connect to couchbase at %s - %v ", connString, err)
	}

	// with RBAC, bucket-passwords are no longer used - https://docs.couchbase.com/go-sdk/1.6/sdk-authentication-overview.html#authenticating-with-legacy-sdk-versions
	bktopt := gocb.GetBucketOptions{}
	bktset, err := c.Buckets().GetBucket(cbs.bucketName, &bktopt)
	if err != nil {
		return fmt.Errorf("couchbase error: failed to open bucket %s - %v", cbs.bucketName, err)
	}
	if bktset == nil {
		return fmt.Errorf("couchbase error: failed to read settings of bucket %s - %v", cbs.bucketName, err)
	}
	//cbs.bucket = bucket

	r := metadata.Properties[numReplicasDurableReplication]
	if r != "" {
		_r, _ := strconv.ParseUint(r, 10, 0)
		cbs.numReplicasDurableReplication = uint(_r)
	}

	p := metadata.Properties[numReplicasDurablePersistence]
	if p != "" {
		_p, _ := strconv.ParseUint(p, 10, 0)
		cbs.numReplicasDurablePersistence = uint(_p)
	}

	return nil
}

// Features returns the features available in this state store.
func (cbs *Couchbase) Features() []state.Feature {
	return cbs.features
}

// Set stores value for a key to couchbase. It honors ETag (for concurrency) and consistency settings.
func (cbs *Couchbase) Set(req *state.SetRequest) error {
	err := state.CheckRequestOptions(req.Options)
	if err != nil {
		return err
	}
	value, err := utils.Marshal(req.Value, cbs.json.Marshal)
	if err != nil {
		return fmt.Errorf("couchbase error: failed to convert value %v", err)
	}
	if value == nil {
		return err
	}
	// nolint:nestif
	// key already exists (use Replace)
	if req.ETag != nil {
		// compare-and-swap (CAS) for managing concurrent modifications - https://docs.couchbase.com/go-sdk/current/concurrent-mutations-cluster.html
		cas, cerr := eTagToCas(*req.ETag)
		if cerr != nil {
			return err
		}
		if cas == 0 {
			return cerr
		}
		if req.Options.Consistency == state.Strong {
			//_, err = cbs.bucket.ReplaceDura(req.Key, value, cas, 0, cbs.numReplicasDurableReplication, cbs.numReplicasDurablePersistence)
		} else {
			//_, err = cbs.bucket.Replace(req.Key, value, cas, 0)
		}
	} else {
		// key does not exist: replace or insert (with Upsert)
		if req.Options.Consistency == state.Strong {
			//_, err = cbs.bucket.UpsertDura(req.Key, value, 0, cbs.numReplicasDurableReplication, cbs.numReplicasDurablePersistence)
		} else {
			//_, err = cbs.bucket.Upsert(req.Key, value, 0)
		}
	}

	if err != nil {
		if req.ETag != nil {
			return state.NewETagError(state.ETagMismatch, err)
		}

		return fmt.Errorf("couchbase error: failed to set value for key %s - %v", req.Key, err)
	}

	return nil
}

// Get retrieves state from couchbase with a key.
func (cbs *Couchbase) Get(req *state.GetRequest) (*state.GetResponse, error) {
	var getopt gocb.GetOptions
	cas, err := cbs.bucket.DefaultCollection().Get(req.Key, &getopt)
	if err != nil {
		if err == gocb.ErrCollectionNotFound {
			return &state.GetResponse{}, nil
		}

		return nil, fmt.Errorf("couchbase error: failed to get value for key %s - %v", req.Key, err)
	}

	var data interface{}
	cas.Content(&data)
	return &state.GetResponse{
		//Data: data([]byte),
		//ETag: ptr.String(strconv.FormatUint(uint64(cas), 10)),
	}, nil
}

// Delete performs a delete operation.
func (cbs *Couchbase) Delete(req *state.DeleteRequest) error {
	err := state.CheckRequestOptions(req.Options)
	if err != nil {
		return err
	}

	var cas gocb.Cas = 0

	if req.ETag != nil {
		cas, err = eTagToCas(*req.ETag)
		if err != nil {
			return err
		}
		if cas == 0 {
			return err
		}
	}
	if req.Options.Consistency == state.Strong {
		//_, err = cbs.bucket.RemoveDura(req.Key, cas, cbs.numReplicasDurableReplication, cbs.numReplicasDurablePersistence)
	} else {
		//_, err = cbs.bucket.Remove(req.Key, cas)
	}
	if err != nil {
		if req.ETag != nil {
			return state.NewETagError(state.ETagMismatch, err)
		}

		return fmt.Errorf("couchbase error: failed to delete key %s - %v", req.Key, err)
	}

	return nil
}

func (cbs *Couchbase) Ping() error {
	return nil
}

// converts string etag sent by the application into a gocb.Cas object, which can then be used for optimistic locking for set and delete operations.
func eTagToCas(eTag string) (gocb.Cas, error) {
	var cas gocb.Cas = 0
	// CAS is a 64-bit integer - https://docs.couchbase.com/go-sdk/current/concurrent-mutations-cluster.html#cas-value-format
	temp, err := strconv.ParseUint(eTag, 10, 64)
	if err != nil {
		return cas, state.NewETagError(state.ETagInvalid, err)
	}
	cas = gocb.Cas(temp)

	return cas, nil
}
