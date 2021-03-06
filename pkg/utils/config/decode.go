package config

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
	"reflect"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

var (
	typeDuration      = reflect.TypeOf(time.Duration(5))             // nolint: gochecknoglobals
	typeTime          = reflect.TypeOf(time.Time{})                  // nolint: gochecknoglobals
	typeStringDecoder = reflect.TypeOf((*StringDecoder)(nil)).Elem() // nolint: gochecknoglobals
)

// StringDecoder is used as a way for custom types (or alias types) to
// override the basic decoding function in the `decodeString`
// DecodeHook. `encoding.TextMashaller` was not used because it
// matches many Go types and would have potentially unexpected results.
// Specifying a custom decoding func should be very intentional.
type StringDecoder interface {
	DecodeString(value string) error
}

// Decode decodes generic map values from `input` to `output`, while providing helpful error information.
// `output` must be a pointer to a Go struct that contains `mapstructure` struct tags on fields that should
// be decoded. This function is useful when decoding values from configuration files parsed as
// `map[string]interface{}` or component metadata as `map[string]string`.
//
// Most of the heavy lifting is handled by the mapstructure library. A custom decoder is used to handle
// decoding string values to the supported primitives.
func Decode(input interface{}, output interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{ // nolint:exhaustivestruct
		Result:     output,
		DecodeHook: decodeString,
	})
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}

// nolint:cyclop
func decodeString(
	f reflect.Type,
	t reflect.Type,
	data interface{}) (interface{}, error) {
	if t.Kind() == reflect.String && f.Kind() != reflect.String {
		return fmt.Sprintf("%v", data), nil
	}
	if f.Kind() == reflect.Ptr {
		f = f.Elem()
		data = reflect.ValueOf(data).Elem().Interface()
	}
	if f.Kind() != reflect.String {
		return data, nil
	}

	dataString, ok := data.(string)
	if !ok {
		return nil, errors.Errorf("expected string: got %s", reflect.TypeOf(data))
	}

	var result interface{}
	var decoder StringDecoder

	if t.Implements(typeStringDecoder) {
		result = reflect.New(t.Elem()).Interface()
		decoder = result.(StringDecoder)
	} else if reflect.PtrTo(t).Implements(typeStringDecoder) {
		result = reflect.New(t).Interface()
		decoder = result.(StringDecoder)
	}

	if decoder != nil {
		if err := decoder.DecodeString(dataString); err != nil {
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}

			return nil, errors.Errorf("invalid %s %q: %v", t.Name(), dataString, err)
		}

		return result, nil
	}

	switch t {
	case typeDuration:
		// Check for simple integer values and treat them
		// as milliseconds
		if val, err := strconv.Atoi(dataString); err == nil {
			return time.Duration(val) * time.Millisecond, nil
		}

		// Convert it by parsing
		d, err := time.ParseDuration(dataString)

		return d, invalidError(err, "duration", dataString)
	case typeTime:
		// Convert it by parsing
		t, err := time.Parse(time.RFC3339Nano, dataString)
		if err == nil {
			return t, nil
		}
		t, err = time.Parse(time.RFC3339, dataString)

		return t, invalidError(err, "time", dataString)
	}

	switch t.Kind() { // nolint: exhaustive
	case reflect.Uint:
		val, err := strconv.ParseUint(dataString, 10, 64)

		return uint(val), invalidError(err, "uint", dataString)
	case reflect.Uint64:
		val, err := strconv.ParseUint(dataString, 10, 64)

		return val, invalidError(err, "uint64", dataString)
	case reflect.Uint32:
		val, err := strconv.ParseUint(dataString, 10, 32)

		return uint32(val), invalidError(err, "uint32", dataString)
	case reflect.Uint16:
		val, err := strconv.ParseUint(dataString, 10, 16)

		return uint16(val), invalidError(err, "uint16", dataString)
	case reflect.Uint8:
		val, err := strconv.ParseUint(dataString, 10, 8)

		return uint8(val), invalidError(err, "uint8", dataString)

	case reflect.Int:
		val, err := strconv.ParseInt(dataString, 10, 64)

		return int(val), invalidError(err, "int", dataString)
	case reflect.Int64:
		val, err := strconv.ParseInt(dataString, 10, 64)

		return val, invalidError(err, "int64", dataString)
	case reflect.Int32:
		val, err := strconv.ParseInt(dataString, 10, 32)

		return int32(val), invalidError(err, "int32", dataString)
	case reflect.Int16:
		val, err := strconv.ParseInt(dataString, 10, 16)

		return int16(val), invalidError(err, "int16", dataString)
	case reflect.Int8:
		val, err := strconv.ParseInt(dataString, 10, 8)

		return int8(val), invalidError(err, "int8", dataString)

	case reflect.Float32:
		val, err := strconv.ParseFloat(dataString, 32)

		return float32(val), invalidError(err, "float32", dataString)
	case reflect.Float64:
		val, err := strconv.ParseFloat(dataString, 64)

		return val, invalidError(err, "float64", dataString)

	case reflect.Bool:
		val, err := strconv.ParseBool(dataString)

		return val, invalidError(err, "bool", dataString)

	default:
		return data, nil
	}
}

func invalidError(err error, msg, value string) error {
	if err == nil {
		return nil
	}

	return errors.Errorf("invalid %s %q", msg, value)
}
