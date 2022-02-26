package logger

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
)

func clearLoggers() {
	globalLoggers = map[string]Logger{}
}

func TestNewLogger(t *testing.T) {
	testLoggerName := "app.test"

	t.Run("create new logger instance", func(t *testing.T) {
		clearLoggers()

		// act
		NewLogger(testLoggerName)
		_, ok := globalLoggers[testLoggerName]

		// assert
		assert.True(t, ok)
	})

	t.Run("return the existing logger instance", func(t *testing.T) {
		clearLoggers()

		// act
		oldLogger := NewLogger(testLoggerName)
		newLogger := NewLogger(testLoggerName)

		// assert
		assert.Equal(t, oldLogger, newLogger)
	})
}

func TestToLogLevel(t *testing.T) {
	t.Run("convert debug to DebugLevel", func(t *testing.T) {
		assert.Equal(t, DebugLevel, toLogLevel("debug"))
	})

	t.Run("convert info to InfoLevel", func(t *testing.T) {
		assert.Equal(t, InfoLevel, toLogLevel("info"))
	})

	t.Run("convert warn to WarnLevel", func(t *testing.T) {
		assert.Equal(t, WarnLevel, toLogLevel("warn"))
	})

	t.Run("convert error to ErrorLevel", func(t *testing.T) {
		assert.Equal(t, ErrorLevel, toLogLevel("error"))
	})

	t.Run("convert fatal to FatalLevel", func(t *testing.T) {
		assert.Equal(t, FatalLevel, toLogLevel("fatal"))
	})

	t.Run("undefined loglevel", func(t *testing.T) {
		assert.Equal(t, UndefinedLevel, toLogLevel("undefined"))
	})
}
