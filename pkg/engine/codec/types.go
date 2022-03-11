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
	"fmt"

	"github.com/bhojpur/service/pkg/engine/codec/internal/utils"
)

// ToObject decode bytes to interface
func ToObject(v []byte, output interface{}) error {
	_, err := newStructDecoder(output).Decode(v) // nolint
	return err
}

// ToInt32 decode bytes to int32
func ToInt32(v []byte) (int32, error) {
	primitivePacket, _, _, err := DecodePrimitivePacket(v)
	if err != nil {
		return 0, err
	}
	value, err := primitivePacket.ToInt32()
	if err != nil {
		return 0, nil
	}
	return value, nil
}

// ToInt32Slice decode bytes to []int32
func ToInt32Slice(v []byte) ([]int32, error) {
	packet, _, err := DecodeNodePacket(v)
	if err != nil {
		return nil, err
	}
	if !packet.IsSlice() {
		return nil, fmt.Errorf("v not a slice: %v", utils.FormatBytes(v))
	}

	result := make([]int32, 0)
	for _, p := range packet.PrimitivePackets {
		v, _ := p.ToInt32()
		result = append(result, v)
	}
	return result, nil
}

// ToUInt32 decode bytes to uint32
func ToUInt32(v []byte) (uint32, error) {
	primitivePacket, _, _, err := DecodePrimitivePacket(v)
	if err != nil {
		return 0, err
	}
	value, err := primitivePacket.ToUInt32()
	if err != nil {
		return 0, nil
	}
	return value, nil
}

// ToUInt32Slice decode bytes to []uint32
func ToUInt32Slice(v []byte) ([]uint32, error) {
	packet, _, err := DecodeNodePacket(v)
	if err != nil {
		return nil, err
	}
	if !packet.IsSlice() {
		return nil, fmt.Errorf("v not a slice: %v", utils.FormatBytes(v))
	}
	result := make([]uint32, 0)
	for _, p := range packet.PrimitivePackets {
		v, _ := p.ToUInt32()
		result = append(result, v)
	}
	return result, nil
}

// ToInt64 decode bytes to int64
func ToInt64(v []byte) (int64, error) {
	primitivePacket, _, _, err := DecodePrimitivePacket(v)
	if err != nil {
		return 0, err
	}
	value, err := primitivePacket.ToInt64()
	if err != nil {
		return 0, nil
	}
	return value, nil
}

// ToInt64Slice decode bytes to []int64
func ToInt64Slice(v []byte) ([]int64, error) {
	packet, _, err := DecodeNodePacket(v)
	if err != nil {
		return nil, err
	}
	if !packet.IsSlice() {
		return nil, fmt.Errorf("v not a slice: %v", utils.FormatBytes(v))
	}
	result := make([]int64, 0)
	for _, p := range packet.PrimitivePackets {
		v, _ := p.ToInt64()
		result = append(result, v)
	}
	return result, nil
}

// ToUInt64 decode bytes to uint64
func ToUInt64(v []byte) (uint64, error) {
	primitivePacket, _, _, err := DecodePrimitivePacket(v)
	if err != nil {
		return 0, err
	}
	value, err := primitivePacket.ToUInt64()
	if err != nil {
		return 0, nil
	}
	return value, nil
}

// ToUInt64Slice decode bytes to []uint64
func ToUInt64Slice(v []byte) ([]uint64, error) {
	packet, _, err := DecodeNodePacket(v)
	if err != nil {
		return nil, err
	}
	if !packet.IsSlice() {
		return nil, fmt.Errorf("v not a slice: %v", utils.FormatBytes(v))
	}
	result := make([]uint64, 0)
	for _, p := range packet.PrimitivePackets {
		v, _ := p.ToUInt64()
		result = append(result, v)
	}
	return result, nil
}

// ToFloat32 decode bytes to float32
func ToFloat32(v []byte) (float32, error) {
	primitivePacket, _, _, err := DecodePrimitivePacket(v)
	if err != nil {
		return 0, err
	}
	value, err := primitivePacket.ToFloat32()
	if err != nil {
		return 0, nil
	}
	return value, nil
}

// ToFloat32Slice decode bytes to []float32
func ToFloat32Slice(v []byte) ([]float32, error) {
	packet, _, err := DecodeNodePacket(v)
	if err != nil {
		return nil, err
	}
	if !packet.IsSlice() {
		return nil, fmt.Errorf("v not a slice: %v", utils.FormatBytes(v))
	}
	result := make([]float32, 0)
	for _, p := range packet.PrimitivePackets {
		v, _ := p.ToFloat32()
		result = append(result, v)
	}
	return result, nil
}

// ToFloat64 decode bytes to float64
func ToFloat64(v []byte) (float64, error) {
	primitivePacket, _, _, err := DecodePrimitivePacket(v)
	if err != nil {
		return 0, err
	}
	value, err := primitivePacket.ToFloat64()
	if err != nil {
		return 0, nil
	}
	return value, nil
}

// ToFloat64Slice decode bytes to []float64
func ToFloat64Slice(v []byte) ([]float64, error) {
	packet, _, err := DecodeNodePacket(v)
	if err != nil {
		return nil, err
	}
	if !packet.IsSlice() {
		return nil, fmt.Errorf("v not a slice: %v", utils.FormatBytes(v))
	}
	result := make([]float64, 0)
	for _, p := range packet.PrimitivePackets {
		v, _ := p.ToFloat64()
		result = append(result, v)
	}
	return result, nil
}

// ToBool decode bytes to bool
func ToBool(v []byte) (bool, error) {
	primitivePacket, _, _, err := DecodePrimitivePacket(v)
	if err != nil {
		return false, err
	}
	value, err := primitivePacket.ToBool()
	if err != nil {
		return false, nil
	}
	return value, nil
}

// ToBoolSlice decode bytes to []bool
func ToBoolSlice(v []byte) ([]bool, error) {
	packet, _, err := DecodeNodePacket(v)
	if err != nil {
		return nil, err
	}
	if !packet.IsSlice() {
		return nil, fmt.Errorf("v not a slice: %v", utils.FormatBytes(v))
	}
	result := make([]bool, 0)
	for _, p := range packet.PrimitivePackets {
		v, _ := p.ToBool()
		result = append(result, v)
	}
	return result, nil
}

// ToUTF8String decode bytes to string
func ToUTF8String(v []byte) (string, error) {
	primitivePacket, _, _, err := DecodePrimitivePacket(v)
	if err != nil {
		return "", err
	}
	value, err := primitivePacket.ToUTF8String()
	if err != nil {
		return "", nil
	}
	return value, nil
}

// ToUTF8StringSlice decode bytes to []string
func ToUTF8StringSlice(v []byte) ([]string, error) {
	packet, _, err := DecodeNodePacket(v)
	if err != nil {
		return nil, err
	}
	if !packet.IsSlice() {
		return nil, fmt.Errorf("v not a slice: %v", utils.FormatBytes(v))
	}
	result := make([]string, 0)
	for _, p := range packet.PrimitivePackets {
		v, _ := p.ToUTF8String()
		result = append(result, v)
	}
	return result, nil
}

func ToBytes(v []byte) ([]byte, error) {
	primitivePacket, _, _, err := DecodePrimitivePacket(v)
	if err != nil {
		return nil, err
	}
	return primitivePacket.valBuf, nil
}
