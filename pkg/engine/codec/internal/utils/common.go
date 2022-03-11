package utils

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

import "reflect"

// MSB `1000 0000`
const MSB byte = 0x80

// DropMSB `0111 1111`
const DropMSB = 0x3F

// DropMSBArrayFlag `0011 1111`
const DropMSBArrayFlag = 0x3F

// SliceFlag `0100 0000`, Value Slice
const SliceFlag = 0x40

// KeyOfSliceItem TLV sid
const KeyOfSliceItem = 0x00

// KeyStringOfSliceItem
const KeyStringOfSliceItem = "0x00"

// RootToken
const RootToken byte = 0x01

// TypeOfByteSlice Type of []byte
var TypeOfByteSlice = reflect.TypeOf([]byte{})

// TypeOfStringSlice Type of []string{}
var TypeOfStringSlice = reflect.TypeOf([]string{})

// TypeOfInt32Slice Type of []int32{}
var TypeOfInt32Slice = reflect.TypeOf([]int32{})

// TypeOfUint32Slice Type of []uint32{}
var TypeOfUint32Slice = reflect.TypeOf([]uint32{})

// TypeOfInt64Slice Type of []int64{}
var TypeOfInt64Slice = reflect.TypeOf([]int64{})

// TypeOfUint64Slice Type of []uint64{}
var TypeOfUint64Slice = reflect.TypeOf([]uint64{})

// TypeOfFloat32Slice Type of []float32{}
var TypeOfFloat32Slice = reflect.TypeOf([]float32{})

// TypeOfFloat64Slice Type of []float64{}
var TypeOfFloat64Slice = reflect.TypeOf([]float64{})

// TypeOfBoolSlice Type of []bool{}
var TypeOfBoolSlice = reflect.TypeOf([]bool{})
