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
	"time"

	"github.com/sirupsen/logrus"
)

// appLogger is the implemention for logrus.
type appLogger struct {
	// name is the name of logger that is published to log as a scope
	name string
	// logger is the instance of logrus logger
	logger logrus.Logger
	// format is the instance of custom logging fomatter
	format *logrus.Entry
}

var AppVersion string = "unknown"

func newAppLogger(name string) *appLogger {
	newLogger := logrus.New()
	newLogger.SetOutput(os.Stdout)

	dl := &appLogger{
		name:   name,
		logger: *newLogger,
		format: newLogger.WithFields(logrus.Fields{
			logFieldScope: name,
			logFieldType:  LogTypeLog,
		}),
	}

	dl.EnableJSONOutput(defaultJSONOutput)

	return dl
}

// EnableJSONOutput enables JSON formatted output log.
func (l *appLogger) EnableJSONOutput(enabled bool) {
	var formatter logrus.Formatter

	fieldMap := logrus.FieldMap{
		// If time field name is conflicted, logrus adds "fields." prefix.
		// So rename to unused field @time to avoid the confliction.
		logrus.FieldKeyTime:  logFieldTimeStamp,
		logrus.FieldKeyLevel: logFieldLevel,
		logrus.FieldKeyMsg:   logFieldMessage,
	}

	hostname, _ := os.Hostname()
	l.format.Data = logrus.Fields{
		logFieldScope:    l.format.Data[logFieldScope],
		logFieldType:     LogTypeLog,
		logFieldInstance: hostname,
		logFieldAppVer:   AppVersion,
	}

	if enabled {
		formatter = &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
			FieldMap:        fieldMap,
		}
	} else {
		formatter = &logrus.TextFormatter{
			TimestampFormat: time.RFC3339Nano,
			FieldMap:        fieldMap,
		}
	}

	l.format.Logger.SetFormatter(formatter)
}

// SetAppID sets app_id field in the log. Default value is empty string.
func (l *appLogger) SetAppID(id string) {
	l.format = l.logger.WithField(logFieldAppID, id)
}

func toLogrusLevel(lvl LogLevel) logrus.Level {
	// ignore error because it will never happens
	l, _ := logrus.ParseLevel(string(lvl))
	return l
}

// SetOutputLevel sets log output level.
func (l *appLogger) SetOutputLevel(outputLevel LogLevel) {
	l.format.Logger.SetLevel(toLogrusLevel(outputLevel))
}

// WithLogType specify the log_type field in log. Default value is LogTypeLog.
func (l *appLogger) WithLogType(logType string) Logger {
	return &appLogger{
		name:   l.name,
		format: l.logger.WithField(logFieldType, logType),
	}
}

// Info logs a message at level Info.
func (l *appLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

// Infof logs a message at level Info.
func (l *appLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Debug logs a message at level Debug.
func (l *appLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

// Debugf logs a message at level Debug.
func (l *appLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// Warn logs a message at level Warn.
func (l *appLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

// Warnf logs a message at level Warn.
func (l *appLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

// Error logs a message at level Error.
func (l *appLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

// Errorf logs a message at level Error.
func (l *appLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// Fatal logs a message at level Fatal then the process will exit with status set to 1.
func (l *appLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

// Fatalf logs a message at level Fatal then the process will exit with status set to 1.
func (l *appLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}
