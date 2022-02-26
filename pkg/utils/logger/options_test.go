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

func TestOptions(t *testing.T) {
	t.Run("default options", func(t *testing.T) {
		o := DefaultOptions()
		assert.Equal(t, defaultJSONOutput, o.JSONFormatEnabled)
		assert.Equal(t, undefinedAppID, o.appID)
		assert.Equal(t, defaultOutputLevel, o.OutputLevel)
	})

	t.Run("set app ID", func(t *testing.T) {
		o := DefaultOptions()
		assert.Equal(t, undefinedAppID, o.appID)

		o.SetAppID("app-app")
		assert.Equal(t, "app-app", o.appID)
	})

	t.Run("attaching log related cmd flags", func(t *testing.T) {
		o := DefaultOptions()

		logLevelAsserted := false
		testStringVarFn := func(p *string, name string, value string, usage string) {
			if name == "log-level" && value == defaultOutputLevel {
				logLevelAsserted = true
			}
		}

		logAsJSONAsserted := false
		testBoolVarFn := func(p *bool, name string, value bool, usage string) {
			if name == "log-as-json" && value == defaultJSONOutput {
				logAsJSONAsserted = true
			}
		}

		o.AttachCmdFlags(testStringVarFn, testBoolVarFn)

		// assert
		assert.True(t, logLevelAsserted)
		assert.True(t, logAsJSONAsserted)
	})
}

func TestApplyOptionsToLoggers(t *testing.T) {
	testOptions := Options{
		JSONFormatEnabled: true,
		appID:             "app-app",
		OutputLevel:       "debug",
	}

	// Create two loggers
	testLoggers := []Logger{
		NewLogger("testLogger0"),
		NewLogger("testLogger1"),
	}

	for _, l := range testLoggers {
		l.EnableJSONOutput(false)
		l.SetOutputLevel(InfoLevel)
	}

	assert.NoError(t, ApplyOptionsToLoggers(&testOptions))

	for _, l := range testLoggers {
		assert.Equal(
			t,
			"app-app",
			(l.(*appLogger)).format.Data[logFieldAppID])
		assert.Equal(
			t,
			toLogrusLevel(DebugLevel),
			(l.(*appLogger)).format.Logger.GetLevel())
	}
}
