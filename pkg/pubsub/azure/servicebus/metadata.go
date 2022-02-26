package servicebus

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

// Reference for settings:
// https://github.com/Azure/azure-service-bus-go/blob/54b2faa53e5216616e59725281be692acc120c34/subscription_manager.go#L101
type metadata struct {
	ConnectionString                string `json:"connectionString"`
	ConsumerID                      string `json:"consumerID"`
	TimeoutInSec                    int    `json:"timeoutInSec"`
	HandlerTimeoutInSec             int    `json:"handlerTimeoutInSec"`
	LockRenewalInSec                int    `json:"lockRenewalInSec"`
	MaxActiveMessages               int    `json:"maxActiveMessages"`
	MaxActiveMessagesRecoveryInSec  int    `json:"maxActiveMessagesRecoveryInSec"`
	MaxReconnectionAttempts         int    `json:"maxReconnectionAttempts"`
	ConnectionRecoveryInSec         int    `json:"connectionRecoveryInSec"`
	DisableEntityManagement         bool   `json:"disableEntityManagement"`
	MaxDeliveryCount                *int   `json:"maxDeliveryCount"`
	LockDurationInSec               *int   `json:"lockDurationInSec"`
	DefaultMessageTimeToLiveInSec   *int   `json:"defaultMessageTimeToLiveInSec"`
	AutoDeleteOnIdleInSec           *int   `json:"autoDeleteOnIdleInSec"`
	MaxConcurrentHandlers           *int   `json:"maxConcurrentHandlers"`
	PrefetchCount                   *int   `json:"prefetchCount"`
	PublishMaxRetries               int    `json:"publishMaxRetries"`
	PublishInitialRetryIntervalInMs int    `json:"publishInitialRetryInternalInMs"`
	NamespaceName                   string `json:"namespaceName,omitempty"`
}
