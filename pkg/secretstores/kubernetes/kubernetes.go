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
	"context"
	"errors"
	"os"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	kubeclient "github.com/bhojpur/service/pkg/authentication/kubernetes"
	"github.com/bhojpur/service/pkg/secretstores"
	"github.com/bhojpur/service/pkg/utils/logger"
)

type kubernetesSecretStore struct {
	kubeClient kubernetes.Interface
	logger     logger.Logger
}

// NewKubernetesSecretStore returns a new Kubernetes secret store.
func NewKubernetesSecretStore(logger logger.Logger) secretstores.SecretStore {
	return &kubernetesSecretStore{logger: logger}
}

// Init creates a Kubernetes client.
func (k *kubernetesSecretStore) Init(metadata secretstores.Metadata) error {
	client, err := kubeclient.GetKubeClient()
	if err != nil {
		return err
	}
	k.kubeClient = client

	return nil
}

// GetSecret retrieves a secret using a key and returns a map of decrypted string/string values.
func (k *kubernetesSecretStore) GetSecret(req secretstores.GetSecretRequest) (secretstores.GetSecretResponse, error) {
	resp := secretstores.GetSecretResponse{
		Data: map[string]string{},
	}
	namespace, err := k.getNamespaceFromMetadata(req.Metadata)
	if err != nil {
		return resp, err
	}

	secret, err := k.kubeClient.CoreV1().Secrets(namespace).Get(context.TODO(), req.Name, meta_v1.GetOptions{})
	if err != nil {
		return resp, err
	}

	for k, v := range secret.Data {
		resp.Data[k] = string(v)
	}

	return resp, nil
}

// BulkGetSecret retrieves all secrets in the store and returns a map of decrypted string/string values.
func (k *kubernetesSecretStore) BulkGetSecret(req secretstores.BulkGetSecretRequest) (secretstores.BulkGetSecretResponse, error) {
	resp := secretstores.BulkGetSecretResponse{
		Data: map[string]map[string]string{},
	}
	namespace, err := k.getNamespaceFromMetadata(req.Metadata)
	if err != nil {
		return resp, err
	}

	secrets, err := k.kubeClient.CoreV1().Secrets(namespace).List(context.TODO(), meta_v1.ListOptions{})
	if err != nil {
		return resp, err
	}

	for _, s := range secrets.Items {
		resp.Data[s.Name] = map[string]string{}
		for k, v := range s.Data {
			resp.Data[s.Name][k] = string(v)
		}
	}

	return resp, nil
}

func (k *kubernetesSecretStore) getNamespaceFromMetadata(metadata map[string]string) (string, error) {
	if val, ok := metadata["namespace"]; ok && val != "" {
		return val, nil
	}

	val := os.Getenv("NAMESPACE")
	if val != "" {
		return val, nil
	}

	return "", errors.New("namespace is missing on metadata and NAMESPACE env variable")
}
