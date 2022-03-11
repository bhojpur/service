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
	"fmt"
	"io"

	//"github.com/bhojpur/service/pkg/engine/codec"
	"github.com/bhojpur/service/pkg/engine/core/frame"
)

// ParseFrame parses the frame from QUIC stream.
func ParseFrame(stream io.Reader) (frame.Frame, error) {
	//codec.FromStream(stream)
	buf := make([]byte, 3*1024)
	_, err := stream.Read(buf)
	if err != nil {
		return nil, err
	}
	// if len(buf) > 512 {
	// 	logger.Debugf("%s🔗 parsed out total %d bytes: \n\thead 64 bytes are: [%# x], \n\ttail 64 bytes are: [%#x]", ParseFrameLogPrefix, len(buf), buf[0:64], buf[len(buf)-64:])
	// } else {
	// 	logger.Debugf("%s🔗 parsed out: [%# x]", ParseFrameLogPrefix, buf)
	// }

	frameType := buf[0]
	// determine the frame type
	switch frameType {
	case 0x80 | byte(frame.TagOfHandshakeFrame):
		handshakeFrame, err := readHandshakeFrame(buf)
		// logger.Debugf("%sHandshakeFrame: name=%s, type=%s", ParseFrameLogPrefix, handshakeFrame.Name, handshakeFrame.Type())
		return handshakeFrame, err
	case 0x80 | byte(frame.TagOfDataFrame):
		data, err := readDataFrame(buf)
		// logger.Debugf("%sDataFrame: tid=%s, tag=%#x, len(carriage)=%d", ParseFrameLogPrefix, data.TransactionID(), data.GetDataTag(), len(data.GetCarriage()))
		return data, err
	case 0x80 | byte(frame.TagOfAcceptedFrame):
		return frame.DecodeToAcceptedFrame(buf)
	case 0x80 | byte(frame.TagOfRejectedFrame):
		return frame.DecodeToRejectedFrame(buf)
	default:
		return nil, fmt.Errorf("unknown frame type, buf[0]=%#x", buf[0])
	}
}

func readHandshakeFrame(buf []byte) (*frame.HandshakeFrame, error) {
	// parse to HandshakeFrame
	// handshake, err := frame.DecodeToHandshakeFrame(buf)
	// if err != nil {
	// 	logger.Errorf("%sreadHandshakeFrame: err=%v", ParseFrameLogPrefix, err)
	// 	return nil
	// }
	// return handshake
	return frame.DecodeToHandshakeFrame(buf)
}

func readDataFrame(buf []byte) (*frame.DataFrame, error) {
	// parse to DataFrame
	// data, err := frame.DecodeToDataFrame(buf)
	// if err != nil {
	// 	logger.Errorf("%sreadDataFrame: err=%v", ParseFrameLogPrefix, err)
	// 	return err
	// }
	// return data
	return frame.DecodeToDataFrame(buf)
}
