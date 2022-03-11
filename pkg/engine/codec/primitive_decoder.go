package codec

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
	"errors"
	"fmt"

	"github.com/bhojpur/service/pkg/utils/encoding"

	"github.com/bhojpur/service/pkg/engine/codec/internal/mark"
	"github.com/bhojpur/service/pkg/engine/codec/internal/utils"
)

// DecodePrimitivePacket parse out whole buffer to a PrimitivePacket
//
// Examples:
// [0x01, 0x01, 0x01] -> Key=0x01, Value=0x01
// [0x41, 0x06, 0x03, 0x01, 0x61, 0x04, 0x01, 0x62] -> key=0x03, value=0x61; key=0x04, value=0x62
func DecodePrimitivePacket(buf []byte) (packet *PrimitivePacket, endPos int, sizeL int, err error) {
	logger := utils.Logger.WithPrefix(utils.DefaultLogger, "BasePacket::Decode")
	logger.Debugf("buf=%#X", buf)

	if buf == nil || len(buf) < primitivePacketBufferMinimalLength {
		return nil, 0, 0, errors.New("invalid Bhojpur Service packet minimal size")
	}

	p := &PrimitivePacket{valBuf: buf}

	var pos = 0
	// first byte is `Tag`
	p.tag = mark.NewTag(buf[pos])
	pos++

	// read `Varint` from buf for `Length of value`
	tmpBuf := buf[pos:]
	var bufLen int32
	codec := encoding.VarCodec{}
	err = codec.DecodePVarInt32(tmpBuf, &bufLen)
	if err != nil {
		return nil, 0, 0, err
	}
	sizeL = codec.Size

	if sizeL < 1 {
		return nil, 0, sizeL, errors.New("malformed, size of Length can not smaller than 1")
	}

	//p.length value，bufLen
	//p.length = uint32(len)
	//pos += int(bufLen)
	p.length = uint32(bufLen)
	pos += sizeL

	endPos = pos + int(p.length)

	logger.Debugf(">>> sizeL=%v, length=%v, pos=%v, endPos=%v", sizeL, p.length, pos, endPos)

	if pos > endPos || endPos > len(buf) || pos > len(buf) {
		return nil, 0, sizeL, fmt.Errorf("beyond the boundary, pos=%v, endPos=%v", pos, endPos)
	}
	p.valBuf = buf[pos:endPos]
	logger.Debugf("valBuf = %#X", p.valBuf)

	return p, endPos, sizeL, nil
}
