package serverless

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
	"path/filepath"
	"sync"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]Serverless)
)

// Serverless defines the interface for serverless
type Serverless interface {
	// Init initializes the serverless
	Init(opts *Options) error

	// Build compiles the serverless to executable
	Build(clean bool) error

	// Run compiles and runs the serverless
	Run(verbose bool) error

	Executable() bool
}

// Register will register a serverless to drivers collections safely
func Register(s Serverless, exts ...string) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if s == nil {
		panic("serverless: Register serverless is nil")
	}
	for _, ext := range exts {
		if _, dup := drivers[ext]; dup {
			panic("serverless: Register called twice for source " + ext)
		}
		drivers[ext] = s
	}
}

// Create returns a new serverless instance with options.
func Create(opts *Options) (Serverless, error) {
	ext := filepath.Ext(opts.Filename)

	driversMu.RLock()
	s, ok := drivers[ext]
	driversMu.RUnlock()
	if ok {
		if err := s.Init(opts); err != nil {
			return nil, err
		}
		return s, nil
	}

	return nil, fmt.Errorf(`serverless: unsupport "%s" source (forgotten import?)`, ext)
}
