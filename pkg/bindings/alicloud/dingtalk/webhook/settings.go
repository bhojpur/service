package webhook

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

// DingTalk webhook are a simple way to post messages from apps into DingTalk
//
// See https://developers.dingtalk.com/document/app/custom-robot-access for details

import (
	"errors"

	"github.com/bhojpur/service/pkg/utils/config"
)

type Settings struct {
	ID     string `mapstructure:"id"`
	URL    string `mapstructure:"url"`
	Secret string `mapstructure:"secret"`
}

func (s *Settings) Decode(in interface{}) error {
	return config.Decode(in, s)
}

func (s *Settings) Validate() error {
	if s.ID == "" {
		return errors.New("webhook error: missing webhook id")
	}
	if s.URL == "" {
		return errors.New("webhook error: missing webhook url")
	}

	return nil
}
