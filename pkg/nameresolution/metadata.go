package nameresolution

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

const (
	// MDNSInstanceName is the Bhojpur Service instance name which is broadcasted.
	MDNSInstanceName string = "name"
	// MDNSInstanceAddress is the address of the Bhojpur Service instance.
	MDNSInstanceAddress string = "address"
	// MDNSInstancePort is the port of Bhojpur Service instance.
	MDNSInstancePort string = "port"
	// MDNSInstanceID is an optional unique Bhojpur Service instance ID.
	MDNSInstanceID string = "instance"

	// HostAddress is the address of the Bhojpur Service instance.
	HostAddress string = "HOST_ADDRESS"
	// AppHTTPPort is the Bhojpur Application API http port.
	AppHTTPPort string = "APP_HTTP_PORT"
	// SvcPort is the Bhojpur Service internal gRPC port (sidecar to sidecar).
	SvcPort string = "SVC_PORT"
	// AppPort is the port of the application, http/grpc depending on mode.
	AppPort string = "APP_PORT"
	// AppID is the ID of the application.
	AppID string = "APP_ID"
)

// Metadata contains a name resolution specific set of metadata properties.
type Metadata struct {
	Properties    map[string]string `json:"properties"`
	Configuration interface{}
}
