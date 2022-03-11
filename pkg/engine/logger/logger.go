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
	"os"
	"strings"

	"github.com/bhojpur/service/pkg/engine/core/log"
)

var logger = Default(isEnableDebug())

// EnableDebug enables the development model for logging.
func EnableDebug() {
	logger = Default(true)
}

// Printf prints a formated message without a specified level.
func Printf(format string, v ...interface{}) {
	logger.Printf(format, v...)
}

// Debugf logs a message at DebugLevel.
func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

// Infof logs a message at InfoLevel.
func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

// Warnf logs a message at WarnLevel.
func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

// Errorf logs a message at ErrorLevel.
func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

// isEnableDebug indicates whether the debug is enabled.
func isEnableDebug() bool {
	return os.Getenv("BHOJPUR_SERVICE_ENABLE_DEBUG") == "true"
}

// isJSONFormat indicates whether the log is in JSON format.
func isJSONFormat() bool {
	return os.Getenv("BHOJPUR_SERVICE_LOG_FORMAT") == "json"
}

func logFormat() string {
	return os.Getenv("BHOJPUR_SERVICE_LOG_FORMAT")
}

func logLevel() log.Level {
	envLevel := strings.ToLower(os.Getenv("BHOJPUR_SERVICE_LOG_LEVEL"))
	level := log.ErrorLevel
	switch envLevel {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warn":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	}
	return level
}

func output() string {
	return strings.ToLower(os.Getenv("BHOJPUR_SERVICE_LOG_OUTPUT"))
}

func errorOutput() string {
	return strings.ToLower(os.Getenv("BHOJPUR_SERVICE_LOG_ERROR_OUTPUT"))
}
