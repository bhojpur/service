package auth

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
	"bytes"

	"github.com/bhojpur/service/pkg/engine/core/auth"
)

var _ auth.Credential = (*AppKeyCredential)(nil)

type AppKeyCredential struct {
	appID    string
	authType auth.AuthType
	payload  []byte
}

func NewAppKeyCredential(appID string, appSecret string) *AppKeyCredential {
	var buf bytes.Buffer
	buf.WriteString(appID)
	buf.WriteString(appSecret)
	payload := buf.Bytes()

	return &AppKeyCredential{
		appID:    appID,
		authType: auth.AuthTypeAppKey,
		payload:  payload,
	}
}

func (a *AppKeyCredential) AppID() string {
	return a.appID
}

func (a *AppKeyCredential) Type() auth.AuthType {
	return auth.AuthType(a.authType)
}

func (a *AppKeyCredential) Payload() []byte {
	return a.payload
}
