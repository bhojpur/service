package state

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

// Store is an interface to perform operations on store.
type Store interface {
	BulkStore
	Init(metadata Metadata) error
	Features() []Feature
	Delete(req *DeleteRequest) error
	Get(req *GetRequest) (*GetResponse, error)
	Set(req *SetRequest) error
	Ping() error
}

// BulkStore is an interface to perform bulk operations on store.
type BulkStore interface {
	BulkDelete(req []DeleteRequest) error
	BulkGet(req []GetRequest) (bool, []BulkGetResponse, error)
	BulkSet(req []SetRequest) error
}

// DefaultBulkStore is a default implementation of BulkStore.
type DefaultBulkStore struct {
	s Store
}

// NewDefaultBulkStore build a default bulk store.
func NewDefaultBulkStore(store Store) DefaultBulkStore {
	defaultBulkStore := DefaultBulkStore{}
	defaultBulkStore.s = store

	return defaultBulkStore
}

// Features returns the features of the encapsulated store.
func (b *DefaultBulkStore) Features() []Feature {
	return b.s.Features()
}

// BulkGet performs a bulks get operations.
func (b *DefaultBulkStore) BulkGet(req []GetRequest) (bool, []BulkGetResponse, error) {
	// by default, the store doesn't support bulk get
	// return false so Bhojpur Application will fallback to call get() method one by one
	return false, nil, nil
}

// BulkSet performs a bulks save operation.
func (b *DefaultBulkStore) BulkSet(req []SetRequest) error {
	for i := range req {
		err := b.s.Set(&req[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// BulkDelete performs a bulk delete operation.
func (b *DefaultBulkStore) BulkDelete(req []DeleteRequest) error {
	for i := range req {
		err := b.s.Delete(&req[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// Querier is an interface to execute queries.
type Querier interface {
	Query(req *QueryRequest) (*QueryResponse, error)
}
