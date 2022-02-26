// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/bhojpur/service/pkg/middleware/http/oauth2clientcredentials (interfaces: TokenProviderInterface)

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

// Package mock_oauth2clientcredentials is a generated GoMock package.
package mock_oauth2clientcredentials

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	oauth2 "golang.org/x/oauth2"
	clientcredentials "golang.org/x/oauth2/clientcredentials"
)

// MockTokenProviderInterface is a mock of TokenProviderInterface interface
type MockTokenProviderInterface struct {
	ctrl     *gomock.Controller
	recorder *MockTokenProviderInterfaceMockRecorder
}

// MockTokenProviderInterfaceMockRecorder is the mock recorder for MockTokenProviderInterface
type MockTokenProviderInterfaceMockRecorder struct {
	mock *MockTokenProviderInterface
}

// NewMockTokenProviderInterface creates a new mock instance
func NewMockTokenProviderInterface(ctrl *gomock.Controller) *MockTokenProviderInterface {
	mock := &MockTokenProviderInterface{ctrl: ctrl}
	mock.recorder = &MockTokenProviderInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTokenProviderInterface) EXPECT() *MockTokenProviderInterfaceMockRecorder {
	return m.recorder
}

// GetToken mocks base method
func (m *MockTokenProviderInterface) GetToken(arg0 *clientcredentials.Config) (*oauth2.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetToken", arg0)
	ret0, _ := ret[0].(*oauth2.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetToken indicates an expected call of GetToken
func (mr *MockTokenProviderInterfaceMockRecorder) GetToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetToken", reflect.TypeOf((*MockTokenProviderInterface)(nil).GetToken), arg0)
}