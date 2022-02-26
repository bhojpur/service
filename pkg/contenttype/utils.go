package contenttype

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

import "strings"

const (
	// CloudEventContentType is the content type for cloud event.
	CloudEventContentType = "application/cloudevents+json"
	// JSONContentType is the content type for JSON.
	JSONContentType = "application/json"
)

// IsCloudEventContentType checks for content type.
func IsCloudEventContentType(contentType string) bool {
	return isContentType(contentType, CloudEventContentType)
}

// IsJSONContentType checks for content type.
func IsJSONContentType(contentType string) bool {
	return isContentType(contentType, JSONContentType)
}

// IsStringContentType determines if content type is string.
func IsStringContentType(contentType string) bool {
	if strings.HasPrefix(strings.ToLower(contentType), "text/") {
		return true
	}

	return isContentType(contentType, "application/xml")
}

// IsBinaryContentType determines if content type is byte[].
func IsBinaryContentType(contentType string) bool {
	return isContentType(contentType, "application/octet-stream")
}

func isContentType(contentType string, expected string) bool {
	lowerContentType := strings.ToLower(contentType)
	if lowerContentType == expected {
		return true
	}

	semiColonPos := strings.Index(lowerContentType, ";")
	if semiColonPos >= 0 {
		return lowerContentType[0:semiColonPos] == expected
	}

	return false
}
