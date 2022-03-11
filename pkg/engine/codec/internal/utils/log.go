package utils

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
	"log"
	"os"
	"strings"
	"time"
)

// LogLevel describes the log level
type LogLevel uint8

const (
	// LogLevelNothing disable log
	LogLevelNothing LogLevel = iota
	// LogLevelError enables err logs
	LogLevelError
	// LogLevelInfo enables info logs (e.g. packets)
	LogLevelInfo
	// LogLevelDebug enables debug logs (e.g. packet contents)
	LogLevelDebug
)

const logEnv = "BHOJPUR_SERVICE_LOG_LEVEL"

// Logger logs log ...
type Logger interface {
	SetLogLevel(LogLevel)
	SetLogTimeFormat(format string)
	WithPrefix(prefix string) Logger
	Debug() bool

	Errorf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

// DefaultLogger is used by Bhojpur Service for logging.
var DefaultLogger Logger

type defaultLogger struct {
	prefix string

	logLevel   LogLevel
	timeFormat string
}

var _ Logger = &defaultLogger{}

// SetLogLevel sets the log level
func (l *defaultLogger) SetLogLevel(level LogLevel) {
	l.logLevel = level
}

// SetLogTimeFormat sets the format of the timestamp
// an empty string disables the logging of timestamps
func (l *defaultLogger) SetLogTimeFormat(format string) {
	log.SetFlags(0) // disable timestamp logging done by the log package
	l.timeFormat = format
}

// Debugf logs something
func (l *defaultLogger) Debugf(format string, args ...interface{}) {
	if l.logLevel == LogLevelDebug {
		l.logMessage(format, args...)
	}
}

// Infof logs something
func (l *defaultLogger) Infof(format string, args ...interface{}) {
	if l.logLevel >= LogLevelInfo {
		l.logMessage(format, args...)
	}
}

// Errorf logs something
func (l *defaultLogger) Errorf(format string, args ...interface{}) {
	if l.logLevel >= LogLevelError {
		l.logMessage(format, args...)
	}
}

func (l *defaultLogger) logMessage(format string, args ...interface{}) {
	var pre string

	if len(l.timeFormat) > 0 {
		pre = time.Now().Format(l.timeFormat) + " "
	}
	if len(l.prefix) > 0 {
		pre += l.prefix + " "
	}
	log.Printf(pre+format, args...)
}

func (l *defaultLogger) WithPrefix(prefix string) Logger {
	if len(l.prefix) > 0 {
		prefix = l.prefix + " " + prefix
	}
	return &defaultLogger{
		logLevel:   l.logLevel,
		timeFormat: l.timeFormat,
		prefix:     prefix,
	}
}

// Debug returns true if the log level is LogLevelDebug
func (l *defaultLogger) Debug() bool {
	return l.logLevel == LogLevelDebug
}

func init() {
	DefaultLogger = &defaultLogger{}
	DefaultLogger.SetLogLevel(readLoggingEnv())
}

func readLoggingEnv() LogLevel {
	lvl := strings.ToLower(os.Getenv(logEnv))
	switch lvl {
	case "":
		return LogLevelNothing
	case "debug":
		return LogLevelDebug
	case "info":
		return LogLevelInfo
	case "error":
		return LogLevelError
	default:
		fmt.Fprintln(os.Stderr, "invalid log level")
		return LogLevelNothing
	}
}
