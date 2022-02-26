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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/utils/logger"
)

func TestGetNamespace(t *testing.T) {
	t.Run("has namespace metadata", func(t *testing.T) {
		store := kubernetesSecretStore{logger: logger.NewLogger("test")}
		namespace := "a"

		ns, err := store.getNamespaceFromMetadata(map[string]string{"namespace": namespace})
		assert.Nil(t, err)
		assert.Equal(t, namespace, ns)
	})

	t.Run("has namespace env", func(t *testing.T) {
		store := kubernetesSecretStore{logger: logger.NewLogger("test")}
		os.Setenv("NAMESPACE", "b")

		ns, err := store.getNamespaceFromMetadata(map[string]string{})
		assert.Nil(t, err)
		assert.Equal(t, "b", ns)
	})

	t.Run("no namespace", func(t *testing.T) {
		store := kubernetesSecretStore{logger: logger.NewLogger("test")}
		os.Setenv("NAMESPACE", "")
		_, err := store.getNamespaceFromMetadata(map[string]string{})

		assert.NotNil(t, err)
		assert.Equal(t, "namespace is missing on metadata and NAMESPACE env variable", err.Error())
	})
}
