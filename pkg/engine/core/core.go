package core

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
	"math/rand"
	"sync"
	"time"
)

var (
	once sync.Once
)

// ConnState represents the state of a connection.
const (
	ConnStateReady          ConnState = "Ready"
	ConnStateDisconnected   ConnState = "Disconnected"
	ConnStateConnecting     ConnState = "Connecting"
	ConnStateConnected      ConnState = "Connected"
	ConnStateAuthenticating ConnState = "Authenticating"
	ConnStateAccepted       ConnState = "Accepted"
	ConnStateRejected       ConnState = "Rejected"
	ConnStatePing           ConnState = "Ping"
	ConnStatePong           ConnState = "Pong"
	ConnStateTransportData  ConnState = "TransportData"
	ConnStateAborted        ConnState = "Aborted"
)

// Prefix is the prefix for logger.
const (
	ClientLogPrefix     = "\033[36m[bhojpur:client]\033[0m "
	ServerLogPrefix     = "\033[32m[bhojpur:server]\033[0m "
	ParseFrameLogPrefix = "\033[36m[bhojpur:stream_parser]\033[0m "
)

func init() {
	rand.Seed(time.Now().Unix())
}
