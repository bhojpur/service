package kubernetes

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

	"github.com/bhojpur/service/pkg/nameresolution"
	"github.com/bhojpur/service/pkg/utils/config"
	"github.com/bhojpur/service/pkg/utils/logger"
)

const (
	DefaultClusterDomain = "cluster.local"
	ClusterDomainKey     = "clusterDomain"
)

type resolver struct {
	logger        logger.Logger
	clusterDomain string
}

// NewResolver creates Kubernetes name resolver.
func NewResolver(logger logger.Logger) nameresolution.Resolver {
	return &resolver{
		logger:        logger,
		clusterDomain: DefaultClusterDomain,
	}
}

// Init initializes Kubernetes name resolver.
func (k *resolver) Init(metadata nameresolution.Metadata) error {
	configInterface, err := config.Normalize(metadata.Configuration)
	if err != nil {
		return err
	}
	if config, ok := configInterface.(map[string]string); ok {
		clusterDomain := config[ClusterDomainKey]
		if clusterDomain != "" {
			k.clusterDomain = clusterDomain
		}
	}

	return nil
}

// ResolveID resolves name to address in Kubernetes.
func (k *resolver) ResolveID(req nameresolution.ResolveRequest) (string, error) {
	// Bhojpur Service requires this formatting for Kubernetes services
	return fmt.Sprintf("%s-app.%s.svc.%s:%d", req.ID, req.Namespace, k.clusterDomain, req.Port), nil
}
