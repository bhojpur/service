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

// JSON-like node:
// { '0x01': -1 }
// Bhojpur Service codec should ->
// 0x01 (is a node, sequence id=1)
//   0x01 (node value length is 1 byte)
//     0x01 (pvarint: -1)
func TestEncoderPrimitiveInt32(t *testing.T) {
	expected := []byte{0x01, 0x01, 0x7F}
	// 0x01 - SeqID=1
	var prim = NewPrimitivePacketEncoder(0x01)
	// Value = -1
	prim.SetInt32Value(-1)

	res := prim.Encode()

	for i, p := range res {
		if p != expected[i] {
			t.Errorf("i=%v, expected=%#x, actual=%#x", i, expected[i], res[i])
		}
	}
}

// JSON-like node:
// { '0x01': "bhojpur" }
// Bhojpur Service codec should ->
// 0x01 (is a node, sequence id=1)
//   0x04 (pvarint, node value length is 4 bytes)
//     0x59, 0x6F, 0x4D, 0x6F (utf-8 string: "bhojpur")
func TestEncoderPrimitiveString(t *testing.T) {
	expected := []byte{0x01, 0x04, 0x59, 0x6F, 0x4D, 0x6F}
	// 0x01 - SeqID=1
	var prim = NewPrimitivePacketEncoder(0x01)
	// Value = "bhojpur"
	prim.SetStringValue("bhojpur")

	res := prim.Encode()

	for i, p := range res {
		if p != expected[i] {
			t.Errorf("i=%v, expected=%v, actual=%v", i, expected[i], res[i])
		}
	}
}

// 0x81 : {
//   0x02: "bhojpur",
// },
func TestEncoderNode1(t *testing.T) {
	expected := []byte{0x81, 0x06, 0x02, 0x04, 0x59, 0x6F, 0x4D, 0x6F}
	var prim = NewPrimitivePacketEncoder(0x02)
	prim.SetStringValue("bhojpur")
	var node = NewNodePacketEncoder(0x01)
	node.AddPrimitivePacket(prim)
	res := node.Encode()

	for i, p := range res {
		if p != expected[i] {
			t.Errorf("i=%v, expected=%#x, actual=%#x", i, expected[i], res[i])
		}
	}
}

// type bar struct {
// 	Name string
// }

// type foo struct {
// 	ID int
// 	*bar
// }
//
// var obj = &foo{ID: 1, bar: &bar{Name: "C"}}
//
// encode obj as:
//
// 0x81: {
//   0x02: 1,
//   0x83 : {
//     0x04: "C",
//   },
// }
//
// to
//
// [0x81, 0x08, 0x02, 0x01, 0x01, 0x83, 0x03, 0x04, 0x01, 0x43]
func TestEncoderNode2(t *testing.T) {
	expected := []byte{0x81, 0x08, 0x02, 0x01, 0x01, 0x83, 0x03, 0x04, 0x01, 0x43}
	// 0x81 - node
	var node1 = NewNodePacketEncoder(0x01)
	// 0x02 - ID=1
	var prim1 = NewPrimitivePacketEncoder(0x02)
	prim1.SetInt32Value(1)
	node1.AddPrimitivePacket(prim1)

	// 0x83 - &bar{}
	var node2 = NewNodePacketEncoder(0x03)

	// 0x04 - Name: "C"
	var prim2 = NewPrimitivePacketEncoder(0x04)
	prim2.SetStringValue("C")
	node2.AddPrimitivePacket(prim2)

	node1.AddNodePacket(node2)

	res := node1.Encode()

	if len(expected) != len(res) {
		t.Errorf("len(expected)=%v, len(res)=%v", len(expected), len(res))
	}

	for i, p := range res {
		if p != expected[i] {
			t.Errorf("i=%v, expected=%#x, actual=%#x", i, expected[i], res[i])
		}
	}
}

// 0x01 (is a node, sequence id=1)
//   0x04 (pvarint, node value length is 4 bytes)
//     0x59, 0x6F, 0x4D, 0x6F (utf-8 string: "bhojpur")
func TestEncoderPrimitiveBinary(t *testing.T) {
	expected := []byte{0x01, 0x03, 0x01, 0x23, 0xFF}
	// 0x01 - SeqID=1
	var prim = NewPrimitivePacketEncoder(0x01)
	// Value = 0x0123FF
	prim.SetBytes([]byte{0x01, 0x23, 0xFF})

	res := prim.Encode()

	for i, p := range res {
		if p != expected[i] {
			t.Errorf("i=%v, expected=%v, actual=%v", i, expected[i], res[i])
		}
	}
}

// 0x01 (is a node, sequence id=1)
//   0x04 (pvarint, node value length is 4 bytes)
//     0x59, 0x6F, 0x4D, 0x6F (utf-8 string: "bhojpur")
func TestEncoderPrimitiveBinaryWithEmptyValue(t *testing.T) {
	expected := []byte{0x01, 0x00}
	// 0x01 - SeqID=1
	var prim = NewPrimitivePacketEncoder(0x01)
	// Value = 0x0123FF
	prim.SetBytes([]byte{})

	res := prim.Encode()

	for i, p := range res {
		if p != expected[i] {
			t.Errorf("i=%v, expected=%v, actual=%v", i, expected[i], res[i])
		}
	}
}
