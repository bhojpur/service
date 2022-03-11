package frame

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
	"strconv"
)

// debugFrameSize print frame data size on debug mode
var debugFrameSize = 16

// Kinds of frames transferable within Bhojpur Service
const (
	// DataFrame
	TagOfDataFrame Type = 0x3F
	// MetaFrame of DataFrame
	TagOfMetaFrame     Type = 0x2F
	TagOfMetadata      Type = 0x03
	TagOfTransactionID Type = 0x01
	TagOfIssuer        Type = 0x02
	// PayloadFrame of DataFrame
	TagOfPayloadFrame Type = 0x2E

	TagOfTokenFrame Type = 0x3E
	// HandshakeFrame
	TagOfHandshakeFrame           Type = 0x3D
	TagOfHandshakeName            Type = 0x01
	TagOfHandshakeType            Type = 0x02
	TagOfHandshakeAppID           Type = 0x03
	TagOfHandshakeAuthType        Type = 0x04
	TagOfHandshakeAuthPayload     Type = 0x05
	TagOfHandshakeObserveDataTags Type = 0x06

	TagOfPingFrame     Type = 0x3C
	TagOfPongFrame     Type = 0x3B
	TagOfAcceptedFrame Type = 0x3A
	TagOfRejectedFrame Type = 0x39
)

// Type represents the type of frame.
type Type uint8

// Frame is the inferface for frame.
type Frame interface {
	// Type gets the type of Frame.
	Type() Type

	// Encode the frame into []byte.
	Encode() []byte
}

func (f Type) String() string {
	switch f {
	case TagOfDataFrame:
		return "DataFrame"
	case TagOfTokenFrame:
		return "TokenFrame"
	case TagOfHandshakeFrame:
		return "HandshakeFrame"
	case TagOfPingFrame:
		return "PingFrame"
	case TagOfPongFrame:
		return "PongFrame"
	case TagOfAcceptedFrame:
		return "AcceptedFrame"
	case TagOfRejectedFrame:
		return "RejectedFrame"
	case TagOfMetaFrame:
		return "MetaFrame"
	case TagOfPayloadFrame:
		return "PayloadFrame"
	// case TagOfTransactionID:
	// 	return "TransactionID"
	case TagOfHandshakeName:
		return "HandshakeName"
	case TagOfHandshakeType:
		return "HandshakeType"
	default:
		return "UnknownFrame"
	}
}

// Shortly reduce data size for easy viewing
func Shortly(data []byte) []byte {
	if len(data) > debugFrameSize {
		return data[:debugFrameSize]
	}
	return data
}

func init() {
	if envFrameSize := os.Getenv("BHOJPUR_SERVICE_DEBUG_FRAME_SIZE"); envFrameSize != "" {
		if val, err := strconv.Atoi(envFrameSize); err == nil {
			debugFrameSize = val
		}
	}
}
