package postgresql

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

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/state"
	"github.com/bhojpur/service/pkg/utils/logger"
)

const (
	fakeConnectionString = "not a real connection"
)

// Fake implementation of interface postgressql.dbaccess.
type fakeDBaccess struct {
	logger       logger.Logger
	initExecuted bool
	setExecuted  bool
	getExecuted  bool
}

func (m *fakeDBaccess) Init(metadata state.Metadata) error {
	m.initExecuted = true

	return nil
}

func (m *fakeDBaccess) Set(req *state.SetRequest) error {
	m.setExecuted = true

	return nil
}

func (m *fakeDBaccess) Get(req *state.GetRequest) (*state.GetResponse, error) {
	m.getExecuted = true

	return nil, nil
}

func (m *fakeDBaccess) Delete(req *state.DeleteRequest) error {
	return nil
}

func (m *fakeDBaccess) ExecuteMulti(sets []state.SetRequest, deletes []state.DeleteRequest) error {
	return nil
}

func (m *fakeDBaccess) Query(req *state.QueryRequest) (*state.QueryResponse, error) {
	return nil, nil
}

func (m *fakeDBaccess) Close() error {
	return nil
}

// Proves that the Init method runs the init method.
func TestInitRunsDBAccessInit(t *testing.T) {
	t.Parallel()
	_, fake := createPostgreSQLWithFake(t)
	assert.True(t, fake.initExecuted)
}

func TestMultiWithNoRequestsReturnsNil(t *testing.T) {
	t.Parallel()
	var operations []state.TransactionalStateOperation
	pgs := createPostgreSQL(t)
	err := pgs.Multi(&state.TransactionalStateRequest{
		Operations: operations,
	})
	assert.Nil(t, err)
}

func TestInvalidMultiAction(t *testing.T) {
	t.Parallel()
	var operations []state.TransactionalStateOperation

	operations = append(operations, state.TransactionalStateOperation{
		Operation: "Something invalid",
		Request:   createSetRequest(),
	})

	pgs := createPostgreSQL(t)
	err := pgs.Multi(&state.TransactionalStateRequest{
		Operations: operations,
	})
	assert.NotNil(t, err)
}

func TestValidSetRequest(t *testing.T) {
	t.Parallel()
	var operations []state.TransactionalStateOperation

	operations = append(operations, state.TransactionalStateOperation{
		Operation: state.Upsert,
		Request:   createSetRequest(),
	})

	pgs := createPostgreSQL(t)
	err := pgs.Multi(&state.TransactionalStateRequest{
		Operations: operations,
	})
	assert.Nil(t, err)
}

func TestInvalidMultiSetRequest(t *testing.T) {
	t.Parallel()
	var operations []state.TransactionalStateOperation

	operations = append(operations, state.TransactionalStateOperation{
		Operation: state.Upsert,
		Request:   createDeleteRequest(), // Delete request is not valid for Upsert operation
	})

	pgs := createPostgreSQL(t)
	err := pgs.Multi(&state.TransactionalStateRequest{
		Operations: operations,
	})
	assert.NotNil(t, err)
}

func TestValidMultiDeleteRequest(t *testing.T) {
	t.Parallel()
	var operations []state.TransactionalStateOperation

	operations = append(operations, state.TransactionalStateOperation{
		Operation: state.Delete,
		Request:   createDeleteRequest(),
	})

	pgs := createPostgreSQL(t)
	err := pgs.Multi(&state.TransactionalStateRequest{
		Operations: operations,
	})
	assert.Nil(t, err)
}

func TestInvalidMultiDeleteRequest(t *testing.T) {
	t.Parallel()
	var operations []state.TransactionalStateOperation

	operations = append(operations, state.TransactionalStateOperation{
		Operation: state.Delete,
		Request:   createSetRequest(), // Set request is not valid for Delete operation
	})

	pgs := createPostgreSQL(t)
	err := pgs.Multi(&state.TransactionalStateRequest{
		Operations: operations,
	})
	assert.NotNil(t, err)
}

func createSetRequest() state.SetRequest {
	return state.SetRequest{
		Key:   randomKey(),
		Value: randomJSON(),
	}
}

func createDeleteRequest() state.DeleteRequest {
	return state.DeleteRequest{
		Key: randomKey(),
	}
}

func createPostgreSQLWithFake(t *testing.T) (*PostgreSQL, *fakeDBaccess) {
	pgs := createPostgreSQL(t)
	fake := pgs.dbaccess.(*fakeDBaccess)

	return pgs, fake
}

func createPostgreSQL(t *testing.T) *PostgreSQL {
	logger := logger.NewLogger("test")

	dba := &fakeDBaccess{
		logger: logger,
	}

	pgs := newPostgreSQLStateStore(logger, dba)
	assert.NotNil(t, pgs)

	metadata := &state.Metadata{
		Properties: map[string]string{connectionStringKey: fakeConnectionString},
	}

	err := pgs.Init(*metadata)

	assert.Nil(t, err)
	assert.NotNil(t, pgs.dbaccess)

	return pgs
}
