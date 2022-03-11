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

	"github.com/bhojpur/service/pkg/engine/codec/internal/mark"
	"github.com/bhojpur/service/pkg/utils/encoding"
)

func parsePayload(b []byte) (endPos int, ifNodePacket bool, np *NodePacket, pp *PrimitivePacket, err error) {
	if len(b) == 0 {
		return 0, false, nil, nil, errors.New("parsePacket params can not be nil")
	}

	pos := 0
	// NodePacket
	if ok := isNodePacket(b[pos]); ok {
		np, endPos, err := DecodeNodePacket(b)
		return endPos, true, np, nil, err
	}

	pp, endPos, _, err = DecodePrimitivePacket(b)
	return endPos, false, nil, pp, err
}

// DecodeNodePacket parse out whole buffer to a NodePacket
func DecodeNodePacket(buf []byte) (pct *NodePacket, endPos int, err error) {
	// fmt.Println(hex.Dump(buf))
	pct = &NodePacket{}

	if len(buf) == 0 {
		return pct, 0, nil
	}

	nodeArr := make([]NodePacket, 0)
	primArr := make([]PrimitivePacket, 0)

	pos := 0

	// `Tag`
	tag := mark.NewTag(buf[pos])
	pct.basePacket.tag = tag
	pos++

	// `Length`: the type is `varint`
	tmpBuf := buf[pos:]
	var vallen int32
	codec := encoding.VarCodec{}
	err = codec.DecodePVarInt32(tmpBuf, &vallen)
	// _len, vallen, err := encoding.Upvarint(buf, pos)
	if err != nil {
		return nil, 0, err
	}

	pct.basePacket.length = uint32(vallen)
	pos += codec.Size

	// `raw` is pct.Length() length
	vl := int(vallen)
	endPos = pos + vl
	pct.basePacket.valBuf = make([]byte, vl)
	copy(pct.basePacket.valBuf, buf[pos:endPos])

	// Parse value to Packet
	for {
		if pos >= endPos || pos >= len(buf) {
			break
		}
		_p, isNode, np, pp, err := parsePayload(buf[pos:endPos])
		pos += _p
		if err != nil {
			return nil, 0, err
		}
		if isNode {
			nodeArr = append(nodeArr, *np)
		} else {
			primArr = append(primArr, *pp)
		}
	}

	pct.NodePackets = nodeArr
	pct.PrimitivePackets = primArr

	return pct, endPos, nil
}
