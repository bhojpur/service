package tester

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

// BasicTestData is data of basic test
type BasicTestData struct {
	Vstring  string  `bhojpur:"0x10"`
	Vint32   int32   `bhojpur:"0x11"`
	Vint64   int64   `bhojpur:"0x12"`
	Vuint32  uint32  `bhojpur:"0x13"`
	Vuint64  uint64  `bhojpur:"0x14"`
	Vfloat32 float32 `bhojpur:"0x15"`
	Vfloat64 float64 `bhojpur:"0x16"`
	Vbool    bool    `bhojpur:"0x17"`
}

// EmbeddedTestData is data of embedded test
type EmbeddedTestData struct {
	BasicTestData `bhojpur:"0x1a"`
	Vaction       string `bhojpur:"0x1b"`
}

// EmbeddedMoreTestData is data of embedded more test
type EmbeddedMoreTestData struct {
	EmbeddedTestData `bhojpur:"0x1c"`
	Vanimal          string `bhojpur:"0x1d"`
}

// NamedTestData is data of named test
type NamedTestData struct {
	Base    BasicTestData `bhojpur:"0x1e"`
	Vaction string        `bhojpur:"0x1f"`
}

// NamedMoreTestData is data of named more test
type NamedMoreTestData struct {
	MyNest  NamedTestData `bhojpur:"0x2a"`
	Vanimal string        `bhojpur:"0x2b"`
}

// ArrayTestData is data of array test
type ArrayTestData struct {
	Vfoo          string     `bhojpur:"0x20"`
	Vbar          [2]string  `bhojpur:"0x21"`
	Vint32Array   [2]int32   `bhojpur:"0x22"`
	Vint64Array   [2]int64   `bhojpur:"0x23"`
	Vuint32Array  [2]uint32  `bhojpur:"0x24"`
	Vuint64Array  [2]uint64  `bhojpur:"0x25"`
	Vfloat32Array [2]float32 `bhojpur:"0x26"`
	Vfloat64Array [2]float64 `bhojpur:"0x27"`
}

// SliceTestData is data of slice test
type SliceTestData struct {
	Vfoo          string    `bhojpur:"0x30"`
	Vbar          []string  `bhojpur:"0x31"`
	Vint32Slice   []int32   `bhojpur:"0x32"`
	Vint64Slice   []int64   `bhojpur:"0x33"`
	Vuint32Slice  []uint32  `bhojpur:"0x34"`
	Vuint64Slice  []uint64  `bhojpur:"0x35"`
	Vfloat32Slice []float32 `bhojpur:"0x36"`
	Vfloat64Slice []float64 `bhojpur:"0x37"`
}

// SliceStructTestData is data of slice struct test
type SliceStructTestData struct {
	Vstring          string                 `bhojpur:"0x2e"`
	BaseList         []BasicTestData        `bhojpur:"0x2f"`
	NamedMoreList    []NamedMoreTestData    `bhojpur:"0x3a"`
	EmbeddedMoreList []EmbeddedMoreTestData `bhojpur:"0x3b"`
}

// ArrayStructTestData is data of array struct test
type ArrayStructTestData struct {
	Vstring          string                  `bhojpur:"0x2e"`
	BaseList         [2]BasicTestData        `bhojpur:"0x2f"`
	NamedMoreList    [2]NamedMoreTestData    `bhojpur:"0x3a"`
	EmbeddedMoreList [2]EmbeddedMoreTestData `bhojpur:"0x3b"`
}

// NestedTestData is data of nested test
type NestedTestData struct {
	SubNested Sub1NestedTestData `bhojpur:"0x3a"`
}

// Sub1NestedTestData is data of sub1 nested test
type Sub1NestedTestData struct {
	SubNested Sub2NestedTestData `bhojpur:"0x3b"`
}

// Sub2NestedTestData is data of sub2 nested test
type Sub2NestedTestData struct {
	SubNested Sub3NestedTestData `bhojpur:"0x3c"`
}

// Sub3NestedTestData is data of sub3 nested test
type Sub3NestedTestData struct {
	BasicList []BasicTestData `bhojpur:"0x3d"`
}
