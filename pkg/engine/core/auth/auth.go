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
	"github.com/bhojpur/service/pkg/engine/core/frame"
)

type AuthType byte

const (
	AuthTypeNone       AuthType = 0x0
	AuthTypeAppKey     AuthType = 0x1
	AuthTypePublicKey  AuthType = 0x2
	AuthTypePrivateKey AuthType = 0x3
)

func (a AuthType) String() string {
	switch a {
	case AuthTypeAppKey:
		return "AppKey"
	case AuthTypePublicKey:
		return "PublicKey"
	case AuthTypePrivateKey:
		return "PrivateKey"
	default:
		return "None"
	}
}

// Authentication for server
type Authentication interface {
	Type() AuthType
	Authenticate(f *frame.HandshakeFrame) bool
}

// Credential for client
type Credential interface {
	AppID() string
	Type() AuthType
	Payload() []byte
}

// None auth

var _ Authentication = (*AuthNone)(nil)

type AuthNone struct{}

func NewAuthNone() *AuthNone {
	return &AuthNone{}
}

func (a *AuthNone) Type() AuthType {
	return AuthTypeNone
}

func (a *AuthNone) Authenticate(f *frame.HandshakeFrame) bool {
	return true
}

var _ = Credential(&CredentialNone{})

type CredentialNone struct{}

func NewCredendialNone() *CredentialNone {
	return &CredentialNone{}
}

func (c *CredentialNone) AppID() string {
	return ""
}

func (c *CredentialNone) Type() AuthType {
	return AuthTypeNone
}

func (c *CredentialNone) Payload() []byte {
	return nil
}
