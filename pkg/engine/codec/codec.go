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
	"reflect"

	"github.com/bhojpur/service/pkg/engine/codec/internal/utils"
)

// Codec encode the user's data according to the Bhojpur Service encoding rules
type Codec interface {
	// Marshal encode interface to []byte
	Marshal(input interface{}) ([]byte, error)
}

// NewCodec create a Codec interface
func NewCodec(observe byte) Codec {
	return &bhojpurCodec{
		observe: observe,
	}
}

// bhojpurCodec is implementation of the Codec interface
type bhojpurCodec struct {
	observe byte
}

// Marshal encode interface to []byte
func (c bhojpurCodec) Marshal(input interface{}) ([]byte, error) {
	if c.isStruct(input) {
		return newStructEncoder(c.observe,
			structEncoderOptionRoot(utils.RootToken),
			structEncoderOptionForbidUserKey(utils.ForbidUserKey),
			structEncoderOptionAllowSignalKey(utils.AllowSignalKey)).
			Encode(input)
	}
	return newBasicEncoder(c.observe,
		basicEncoderOptionRoot(utils.RootToken),
		basicEncoderOptionForbidUserKey(utils.ForbidUserKey),
		basicEncoderOptionAllowSignalKey(utils.AllowSignalKey)).
		Encode(input)
}

// isStruct determine whether an interface is a structure
func (c bhojpurCodec) isStruct(mold interface{}) bool {
	isStruct := false

	moldValue := reflect.Indirect(reflect.ValueOf(mold))
	moldType := moldValue.Type()
	switch moldType.Kind() {
	case reflect.Struct:
		isStruct = true
	case reflect.Slice:
		if moldType.Elem().Kind() == reflect.Struct {
			isStruct = true
		}
	}

	return isStruct
}
