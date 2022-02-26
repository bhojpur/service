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
	"encoding/json"
	"strconv"
	"strings"
	"time"

	r "github.com/dancannon/gorethink"
	"github.com/pkg/errors"

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/utils/logger"
)

// Binding represents RethinkDB change change state input binding which fires handler with
// both the previous and current state store content each time there is a change.
type Binding struct {
	logger  logger.Logger
	session *r.Session
	config  StateConfig
	stopCh  chan bool
}

// StateConfig is the binding config.
type StateConfig struct {
	r.ConnectOpts
	Table string `json:"table"`
}

var _ = bindings.InputBinding(&Binding{})

// NewRethinkDBStateChangeBinding returns a new RethinkDB actor event input binding.
func NewRethinkDBStateChangeBinding(logger logger.Logger) *Binding {
	return &Binding{
		logger: logger,
		stopCh: make(chan bool),
	}
}

// Init initializes the RethinkDB binding.
func (b *Binding) Init(metadata bindings.Metadata) error {
	cfg, err := metadataToConfig(metadata.Properties, b.logger)
	if err != nil {
		return errors.Wrap(err, "unable to parse metadata properties")
	}
	b.config = cfg

	ses, err := r.Connect(b.config.ConnectOpts)
	if err != nil {
		return errors.Wrap(err, "error connecting to the database")
	}
	b.session = ses

	return nil
}

// Read triggers the RethinkDB scheduler.
func (b *Binding) Read(handler func(*bindings.ReadResponse) ([]byte, error)) error {
	b.logger.Infof("subscribing to state changes in %s.%s...", b.config.Database, b.config.Table)
	cursor, err := r.DB(b.config.Database).Table(b.config.Table).Changes(r.ChangesOpts{
		IncludeTypes: true,
	}).Run(b.session)
	if err != nil {
		errors.Wrapf(err, "error connecting to table %s", b.config.Table)
	}

	go func() {
		for {
			var change interface{}
			ok := cursor.Next(&change)
			if !ok {
				b.logger.Errorf("error detecting change: %v", cursor.Err())

				break
			}

			data, err := json.Marshal(change)
			if err != nil {
				b.logger.Errorf("error marshalling change handler: %v", err)
			}
			b.logger.Debugf("event: %s", string(data))

			resp := &bindings.ReadResponse{
				Data: data,
				Metadata: map[string]string{
					"store-address":  b.config.Address,
					"store-database": b.config.Database,
					"store-table":    b.config.Table,
				},
			}

			if _, err := handler(resp); err != nil {
				b.logger.Errorf("error invoking change handler: %v", err)

				continue
			}
		}
	}()

	done := <-b.stopCh
	b.logger.Errorf("done: %b", done)
	defer cursor.Close()

	return nil
}

func metadataToConfig(cfg map[string]string, logger logger.Logger) (StateConfig, error) {
	c := StateConfig{}
	for k, v := range cfg {
		switch k {
		case "address": // string
			c.Address = v
		case "addresses": // []string
			c.Addresses = strings.Split(v, ",")
		case "database": // string
			c.Database = v
		case "username": // string
			c.Username = v
		case "password": // string
			c.Password = v
		case "authkey": // string
			c.AuthKey = v
		case "table": // string
			c.Table = v
		case "timeout": // time.Duration
			d, err := time.ParseDuration(v)
			if err != nil {
				return c, errors.Wrapf(err, "invalid timeout format: %v", v)
			}
			c.Timeout = d
		case "write_timeout": // time.Duration
			d, err := time.ParseDuration(v)
			if err != nil {
				return c, errors.Wrapf(err, "invalid write timeout format: %v", v)
			}
			c.WriteTimeout = d
		case "read_timeout": // time.Duration
			d, err := time.ParseDuration(v)
			if err != nil {
				return c, errors.Wrapf(err, "invalid read timeout format: %v", v)
			}
			c.ReadTimeout = d
		case "keep_alive_timeout": // time.Duration
			d, err := time.ParseDuration(v)
			if err != nil {
				return c, errors.Wrapf(err, "invalid keep alive timeout format: %v", v)
			}
			c.KeepAlivePeriod = d
		case "initial_cap": // int
			i, err := strconv.Atoi(v)
			if err != nil {
				return c, errors.Wrapf(err, "invalid keep initial cap format: %v", v)
			}
			c.InitialCap = i
		case "max_open": // int
			i, err := strconv.Atoi(v)
			if err != nil {
				return c, errors.Wrapf(err, "invalid keep max open format: %v", v)
			}
			c.MaxOpen = i
		case "discover_hosts": // bool
			b, err := strconv.ParseBool(v)
			if err != nil {
				return c, errors.Wrapf(err, "invalid discover hosts format: %v", v)
			}
			c.DiscoverHosts = b
		case "use-open-tracing": // bool
			b, err := strconv.ParseBool(v)
			if err != nil {
				return c, errors.Wrapf(err, "invalid use open tracing format: %v", v)
			}
			c.UseOpentracing = b
		case "max_idle": // int
			i, err := strconv.Atoi(v)
			if err != nil {
				return c, errors.Wrapf(err, "invalid keep max idle format: %v", v)
			}
			c.InitialCap = i
		default:
			logger.Infof("unrecognized metadata: %s", k)
		}
	}

	return c, nil
}
