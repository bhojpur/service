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
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const fakeLoggerName = "fakeLogger"

func getTestLogger(buf io.Writer) *appLogger {
	l := newAppLogger(fakeLoggerName)
	l.format.Logger.SetOutput(buf)

	return l
}

func TestEnableJSON(t *testing.T) {
	var buf bytes.Buffer
	testLogger := getTestLogger(&buf)

	expectedHost, _ := os.Hostname()
	testLogger.EnableJSONOutput(true)
	_, okJSON := testLogger.format.Logger.Formatter.(*logrus.JSONFormatter)
	assert.True(t, okJSON)
	assert.Equal(t, "fakeLogger", testLogger.format.Data[logFieldScope])
	assert.Equal(t, LogTypeLog, testLogger.format.Data[logFieldType])
	assert.Equal(t, expectedHost, testLogger.format.Data[logFieldInstance])

	testLogger.EnableJSONOutput(false)
	_, okText := testLogger.format.Logger.Formatter.(*logrus.TextFormatter)
	assert.True(t, okText)
	assert.Equal(t, "fakeLogger", testLogger.format.Data[logFieldScope])
	assert.Equal(t, LogTypeLog, testLogger.format.Data[logFieldType])
	assert.Equal(t, expectedHost, testLogger.format.Data[logFieldInstance])
}

func TestJSONLoggerFields(t *testing.T) {
	tests := []struct {
		name        string
		outputLevel LogLevel
		level       string
		appID       string
		message     string
		instance    string
		fn          func(*appLogger, string)
	}{
		{
			"info()",
			InfoLevel,
			"info",
			"app_app",
			"King App",
			"app-pod",
			func(l *appLogger, msg string) {
				l.Info(msg)
			},
		},
		{
			"infof()",
			InfoLevel,
			"info",
			"App_app",
			"King App",
			"app-pod",
			func(l *appLogger, msg string) {
				l.Infof("%s", msg)
			},
		},
		{
			"debug()",
			DebugLevel,
			"debug",
			"app_app",
			"King App",
			"app-pod",
			func(l *appLogger, msg string) {
				l.Debug(msg)
			},
		},
		{
			"debugf()",
			DebugLevel,
			"debug",
			"app_app",
			"King App",
			"app-pod",
			func(l *appLogger, msg string) {
				l.Debugf("%s", msg)
			},
		},
		{
			"error()",
			InfoLevel,
			"error",
			"app_app",
			"King App",
			"app-pod",
			func(l *appLogger, msg string) {
				l.Error(msg)
			},
		},
		{
			"errorf()",
			InfoLevel,
			"error",
			"app_app",
			"King App",
			"app-pod",
			func(l *appLogger, msg string) {
				l.Errorf("%s", msg)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			testLogger := getTestLogger(&buf)
			testLogger.EnableJSONOutput(true)
			testLogger.SetAppID(tt.appID)
			AppVersion = tt.appID
			testLogger.SetOutputLevel(tt.outputLevel)
			testLogger.format.Data[logFieldInstance] = tt.instance

			tt.fn(testLogger, tt.message)

			b, _ := buf.ReadBytes('\n')
			var o map[string]interface{}
			assert.NoError(t, json.Unmarshal(b, &o))

			// assert
			assert.Equal(t, tt.appID, o[logFieldAppID])
			assert.Equal(t, tt.instance, o[logFieldInstance])
			assert.Equal(t, tt.level, o[logFieldLevel])
			assert.Equal(t, LogTypeLog, o[logFieldType])
			assert.Equal(t, fakeLoggerName, o[logFieldScope])
			assert.Equal(t, tt.message, o[logFieldMessage])
			_, err := time.Parse(time.RFC3339, o[logFieldTimeStamp].(string))
			assert.NoError(t, err)
		})
	}
}

func TestWithTypeFields(t *testing.T) {
	var buf bytes.Buffer
	testLogger := getTestLogger(&buf)
	testLogger.EnableJSONOutput(true)
	testLogger.SetAppID("app_app")
	testLogger.SetOutputLevel(InfoLevel)

	// WithLogType will return new Logger with request log type
	// Meanwhile, testLogger uses the default logtype
	loggerWithRequestType := testLogger.WithLogType(LogTypeRequest)
	loggerWithRequestType.Info("call user app")

	b, _ := buf.ReadBytes('\n')
	var o map[string]interface{}
	assert.NoError(t, json.Unmarshal(b, &o))

	assert.Equalf(t, LogTypeRequest, o[logFieldType], "new logger must be %s type", LogTypeRequest)

	// Log our via testLogger to ensure that testLogger still uses the default logtype
	testLogger.Info("testLogger with log LogType")

	b, _ = buf.ReadBytes('\n')
	assert.NoError(t, json.Unmarshal(b, &o))

	assert.Equalf(t, LogTypeLog, o[logFieldType], "testLogger must be %s type", LogTypeLog)
}

func TestToLogrusLevel(t *testing.T) {
	t.Run("app DebugLevel to Logrus.DebugLevel", func(t *testing.T) {
		assert.Equal(t, logrus.DebugLevel, toLogrusLevel(DebugLevel))
	})

	t.Run("app InfoLevel to Logrus.InfoLevel", func(t *testing.T) {
		assert.Equal(t, logrus.InfoLevel, toLogrusLevel(InfoLevel))
	})

	t.Run("app WarnLevel to Logrus.WarnLevel", func(t *testing.T) {
		assert.Equal(t, logrus.WarnLevel, toLogrusLevel(WarnLevel))
	})

	t.Run("app ErrorLevel to Logrus.ErrorLevel", func(t *testing.T) {
		assert.Equal(t, logrus.ErrorLevel, toLogrusLevel(ErrorLevel))
	})

	t.Run("app FatalLevel to Logrus.FatalLevel", func(t *testing.T) {
		assert.Equal(t, logrus.FatalLevel, toLogrusLevel(FatalLevel))
	})
}
