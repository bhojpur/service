package mysql

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
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"io/ioutil"

	"github.com/go-sql-driver/mysql"

	"github.com/bhojpur/service/pkg/utils/logger"
)

type mySQLFactory struct {
	logger logger.Logger
}

func newMySQLFactory(logger logger.Logger) *mySQLFactory {
	return &mySQLFactory{
		logger: logger,
	}
}

func (m *mySQLFactory) Open(connectionString string) (*sql.DB, error) {
	return sql.Open("mysql", connectionString)
}

func (m *mySQLFactory) RegisterTLSConfig(pemPath string) error {
	rootCertPool := x509.NewCertPool()
	pem, readErr := ioutil.ReadFile(pemPath)

	if readErr != nil {
		m.logger.Errorf("Error reading PEM file from $s", pemPath)

		return readErr
	}

	ok := rootCertPool.AppendCertsFromPEM(pem)

	if !ok {
		return fmt.Errorf("failed to append PEM")
	}

	mysql.RegisterTLSConfig("custom", &tls.Config{RootCAs: rootCertPool, MinVersion: tls.VersionTLS12})

	return nil
}
