package internal

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
	"strings"
)

const (
	channelPrefix = "__keyspace@0__:"
	separator     = "||"
)

func GetRedisValueAndVersion(redisValue string) (string, string) {
	valueAndRevision := strings.Split(redisValue, separator)
	if len(valueAndRevision) == 0 {
		return "", ""
	}
	if len(valueAndRevision) == 1 {
		return valueAndRevision[0], ""
	}
	return valueAndRevision[0], valueAndRevision[1]
}

func ParseRedisKeyFromEvent(eventChannel string) (string, error) {
	index := strings.Index(eventChannel, channelPrefix)
	if index == -1 {
		return "", fmt.Errorf("wrong format of event channel, it should start with '%s': eventChannel=%s", channelPrefix, eventChannel)
	}

	return eventChannel[len(channelPrefix):], nil
}
