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
	"testing"
)

// JSON：{
// '0x04': {
//   '0x01': -1,
// },
// Bhojpur Service Codec should ->
// 0x84 (is a node, sequence id=4)
//   0x03 (node value length is 4 bytes)
//     0x01, 0x01, 0x7F (pvarint: -1)
func TestSimple1Node(t *testing.T) {
	buf := []byte{0x84, 0x03, 0x01, 0x01, 0x7F}
	res, packetLength, err := DecodeNodePacket(buf)
	if err != nil {
		t.Errorf("err should be nil, actual = %v", err)
	}

	if len(res.PrimitivePackets) != 1 {
		t.Errorf("len(res.nodes) actual = %v, and expected = %v", len(res.NodePackets), 1)
	}

	if res.SeqID() != 0x04 {
		t.Errorf("res.SeqID actual = %v, and expected = %v", res.SeqID(), 0x04)
	}

	v1, err := res.PrimitivePackets[0].ToInt32()
	if err != nil {
		t.Error(err)
	}

	if v1 != int32(-1) {
		t.Errorf("n1 value actual = %v, and expected = %v", v1, -1)
	}

	if packetLength != 5 {
		t.Errorf("packetLength actual = %v, and Expected = %v", packetLength, 5)
	}
}

// node
// JSON
// '0x03': {
//   '0x01': -1,
//   '0x02':  1,
// },
// Bhojpur Service codec should ->
// 0x83 (is a node, sequence id=3)
//   0x06 (node value length is 8 bytes)
//     0x01, 0x01, 0x7F (pvarint: -1)
//     0x02, 0x01, 0x01 (pvarint: 1)
func TestSimple2Nodes(t *testing.T) {
	buf := []byte{0x83, 0x06, 0x01, 0x01, 0x7F, 0x02, 0x01, 0x01}
	res, packetLength, err := DecodeNodePacket(buf)
	if err != nil {
		t.Errorf("err should be nil, actual = %v", err)
	}

	if len(res.PrimitivePackets) != 2 {
		t.Errorf("len(res.nodes) actual = %v, and expected = %v", len(res.NodePackets), 2)
	}

	v1, err := res.PrimitivePackets[0].ToInt32()
	if err != nil {
		t.Error(err)
	}

	if v1 != int32(-1) {
		t.Errorf("n1 value actual = %v, and expected = %v", v1, -1)
	}

	v2, err := res.PrimitivePackets[1].ToInt32()
	if err != nil {
		t.Error(err)
	}

	if v2 != int32(1) {
		t.Errorf("n1 value actual = %v, and expected = %v", v2, 1)
	}

	if packetLength != 8 {
		t.Errorf("packetLength actual = %v, and Expected = %v", packetLength, 8)
	}
}

// node
// JSON
// '0x05': {
//	'0x04': {
//     '0x01': -1,
//     '0x02':  1,
//  },
//	'0x03': {
//     '0x01': -2,
//  },
// }
// Bhojpur Service codec should ->
// 0x85
//   0x0D(node value length is 16 bytes)
//     0x84 (is a node, sequence id=3)
//       0x06 (node value length is 8 bytes)
//         0x01, 0x01, 0x7F (varint: -1)
//         0x02, 0x01, 0x43 (string: "C")
//     0x83 (is a node, sequence id=4)
//       0x03 (node value length is 4 bytes)
//         0x01, 0x01, 0x7E (varint: -2)
func TestComplexNodes(t *testing.T) {
	buf := []byte{0x85, 0x0D, 0x84, 0x06, 0x01, 0x01, 0x7F, 0x02, 0x01, 0x43, 0x83, 0x03, 0x01, 0x01, 0x7E}
	res, packetLength, err := DecodeNodePacket(buf)
	if err != nil {
		t.Errorf("err should be nil, actual = %v", err)
	}

	if packetLength != len(buf) {
		t.Errorf("packetLength actual = %v, and expected = %v", packetLength, len(buf))
	}

	if len(res.NodePackets) != 2 {
		t.Errorf("res.NodePackets actual = %v, and expected = %v", len(res.NodePackets), 2)
	}

	if len(res.PrimitivePackets) != 0 {
		t.Errorf("res.PrimitivePackets actual = %v, and expected = %v", len(res.PrimitivePackets), 0)
	}

	n1 := res.NodePackets[0]
	if len(n1.PrimitivePackets) != 2 {
		t.Errorf("n1.PrimitivePackets actual = %v, and expected = %v", len(n1.PrimitivePackets), 2)
	}

	n1p1, _ := n1.PrimitivePackets[0].ToInt32()
	n1p2, _ := n1.PrimitivePackets[1].ToUTF8String()

	n2 := res.NodePackets[1]
	if len(n2.PrimitivePackets) != 1 {
		t.Errorf("n2.PrimitivePackets actual = %v, and expected = %v", len(n2.PrimitivePackets), 1)
	}

	n2p1, _ := n2.PrimitivePackets[0].ToInt32()

	if n1p1 != -1 || n1p2 != "C" || n2p1 != -2 {
		t.Errorf("n1p1=%v, n1p2=%v, n2p1=%v", n1p1, n1p2, n2p1)
	}
}
