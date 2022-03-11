//go:build !windows
// +build !windows

package engine

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
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/bhojpur/service/pkg/engine/logger"
)

// initialize when processor running as server. support inspection:
// - `kill -SIGUSR1 <pid>` inspect state()
// - `kill -SIGTERM <pid>` graceful shutdown
// - `kill -SIGUSR2 <pid>` inspect golang GC
func (z *processor) init() {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGUSR2, syscall.SIGUSR1, syscall.SIGINT)
		logger.Printf("%sListening SIGUSR1, SIGUSR2, SIGTERM/SIGINT...", processorLogPrefix)
		for p1 := range c {
			logger.Printf("Received signal: %s", p1)
			if p1 == syscall.SIGTERM || p1 == syscall.SIGINT {
				logger.Printf("graceful shutting down ... %s", p1)
				os.Exit(0)
				// close(sgnl)
			} else if p1 == syscall.SIGUSR2 {
				var m runtime.MemStats
				runtime.ReadMemStats(&m)
				fmt.Printf("\tNumGC = %v\n", m.NumGC)
			} else if p1 == syscall.SIGUSR1 {
				logger.Printf("print processor stats(): %d", z.Stats())
			}
		}
	}()
}
