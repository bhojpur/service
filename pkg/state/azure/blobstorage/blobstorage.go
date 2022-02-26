package blobstorage

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

/*
Azure Blob Storage state store.

Sample configuration in yaml:

	apiVersion: bhojpur.net/v1alpha1
	kind: Component
	metadata:
	  name: statestore
	spec:
	  type: state.azure.blobstorage
	  metadata:
	  - name: accountName
		value: <storage account name>
	  - name: accountKey
		value: <key>
	  - name: containerName
		value: <container Name>

Concurrency is supported with ETags according to
https://docs.microsoft.com/en-us/azure/storage/common/storage-concurrency#managing-concurrency-in-blob-storage
*/

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/agrea/ptr"
	jsoniter "github.com/json-iterator/go"

	azauth "github.com/bhojpur/service/pkg/authentication/azure"
	"github.com/bhojpur/service/pkg/state"
	"github.com/bhojpur/service/pkg/utils/logger"
)

const (
	keyDelimiter       = "||"
	accountNameKey     = "accountName"
	containerNameKey   = "containerName"
	endpointKey        = "endpoint"
	contentType        = "ContentType"
	contentMD5         = "ContentMD5"
	contentEncoding    = "ContentEncoding"
	contentLanguage    = "ContentLanguage"
	contentDisposition = "ContentDisposition"
	cacheControl       = "CacheControl"
)

// StateStore Type.
type StateStore struct {
	state.DefaultBulkStore
	containerURL azblob.ContainerURL
	json         jsoniter.API

	features []state.Feature
	logger   logger.Logger
}

type blobStorageMetadata struct {
	accountName   string
	containerName string
}

// Init the connection to blob storage, optionally creates a blob container if it doesn't exist.
func (r *StateStore) Init(metadata state.Metadata) error {
	meta, err := getBlobStorageMetadata(metadata.Properties)
	if err != nil {
		return err
	}

	credential, env, err := azauth.GetAzureStorageCredentials(r.logger, meta.accountName, metadata.Properties)
	if err != nil {
		return fmt.Errorf("invalid credentials with error: %s", err.Error())
	}

	userAgent := "app-" + logger.AppVersion
	options := azblob.PipelineOptions{
		Telemetry: azblob.TelemetryOptions{Value: userAgent},
	}
	p := azblob.NewPipeline(credential, options)

	var URL *url.URL
	customEndpoint, ok := metadata.Properties[endpointKey]
	if ok && customEndpoint != "" {
		URL, err = url.Parse(fmt.Sprintf("%s/%s/%s", customEndpoint, meta.accountName, meta.containerName))
	} else {
		URL, err = url.Parse(fmt.Sprintf("https://%s.blob.%s/%s", meta.accountName, env.StorageEndpointSuffix, meta.containerName))
	}
	if err != nil {
		return err
	}
	containerURL := azblob.NewContainerURL(*URL, p)

	_, err = net.LookupHost(URL.Hostname())
	if err != nil {
		return err
	}

	ctx := context.Background()
	_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	r.logger.Debugf("error creating container: %s", err)

	r.containerURL = containerURL
	r.logger.Debugf("using container '%s'", meta.containerName)

	return nil
}

// Features returns the features available in this state store.
func (r *StateStore) Features() []state.Feature {
	return r.features
}

// Delete the state.
func (r *StateStore) Delete(req *state.DeleteRequest) error {
	r.logger.Debugf("delete %s", req.Key)

	return r.deleteFile(req)
}

// Get the state.
func (r *StateStore) Get(req *state.GetRequest) (*state.GetResponse, error) {
	r.logger.Debugf("fetching %s", req.Key)
	data, etag, contentType, err := r.readFile(req)
	if err != nil {
		r.logger.Debugf("error %s", err)

		if isNotFoundError(err) {
			return &state.GetResponse{}, nil
		}

		return &state.GetResponse{}, err
	}

	return &state.GetResponse{
		Data:        data,
		ETag:        ptr.String(etag),
		ContentType: contentType,
	}, err
}

// Set the state.
func (r *StateStore) Set(req *state.SetRequest) error {
	r.logger.Debugf("saving %s", req.Key)

	return r.writeFile(req)
}

func (r *StateStore) Ping() error {
	accessConditions := azblob.BlobAccessConditions{}

	if _, err := r.containerURL.GetProperties(context.Background(), accessConditions.LeaseAccessConditions); err != nil {
		return fmt.Errorf("blob storage: error connecting to Blob storage at %s: %s", r.containerURL.URL().Host, err)
	}

	return nil
}

// NewAzureBlobStorageStore instance.
func NewAzureBlobStorageStore(logger logger.Logger) *StateStore {
	s := &StateStore{
		json:     jsoniter.ConfigFastest,
		features: []state.Feature{state.FeatureETag},
		logger:   logger,
	}
	s.DefaultBulkStore = state.NewDefaultBulkStore(s)

	return s
}

func getBlobStorageMetadata(metadata map[string]string) (*blobStorageMetadata, error) {
	meta := blobStorageMetadata{}

	if val, ok := metadata[accountNameKey]; ok && val != "" {
		meta.accountName = val
	} else {
		return nil, fmt.Errorf("missing or empty %s field from metadata", accountNameKey)
	}

	if val, ok := metadata[containerNameKey]; ok && val != "" {
		meta.containerName = val
	} else {
		return nil, fmt.Errorf("missing or empty %s field from metadata", containerNameKey)
	}

	return &meta, nil
}

func (r *StateStore) readFile(req *state.GetRequest) ([]byte, string, *string, error) {
	blobURL := r.containerURL.NewBlockBlobURL(getFileName(req.Key))

	cpk := azblob.NewClientProvidedKeyOptions(&req.Key, nil, nil)
	resp, err := blobURL.Download(context.Background(), 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, cpk)
	if err != nil {
		r.logger.Debugf("download file %s, err %s", req.Key, err)

		return nil, "", nil, err
	}

	bodyStream := resp.Body(azblob.RetryReaderOptions{})
	data := bytes.Buffer{}
	_, err = data.ReadFrom(bodyStream)

	if err != nil {
		r.logger.Debugf("read file %s, err %s", req.Key, err)

		return nil, "", nil, err
	}

	contentType := resp.ContentType()
	return data.Bytes(), string(resp.ETag()), &contentType, nil
}

func (r *StateStore) writeFile(req *state.SetRequest) error {
	accessConditions := azblob.BlobAccessConditions{}

	if req.Options.Concurrency == state.FirstWrite && req.ETag != nil {
		var etag string
		if req.ETag != nil {
			etag = *req.ETag
		}
		accessConditions.IfMatch = azblob.ETag(etag)
	}

	blobURL := r.containerURL.NewBlockBlobURL(getFileName(req.Key))

	blobHTTPHeaders, err := r.createBlobHTTPHeadersFromRequest(req)
	if err != nil {
		return err
	}
	_, err = azblob.UploadBufferToBlockBlob(context.Background(), r.marshal(req), blobURL, azblob.UploadToBlockBlobOptions{
		Parallelism:      16,
		Metadata:         req.Metadata,
		AccessConditions: accessConditions,
		BlobHTTPHeaders:  blobHTTPHeaders,
	})
	if err != nil {
		r.logger.Debugf("write file %s, err %s", req.Key, err)

		if req.ETag != nil {
			return state.NewETagError(state.ETagMismatch, err)
		}

		return err
	}

	return nil
}

func (r *StateStore) createBlobHTTPHeadersFromRequest(req *state.SetRequest) (azblob.BlobHTTPHeaders, error) {
	var blobHTTPHeaders azblob.BlobHTTPHeaders
	if val, ok := req.Metadata[contentType]; ok && val != "" {
		blobHTTPHeaders.ContentType = val
		delete(req.Metadata, contentType)
	}

	if req.ContentType != nil {
		if blobHTTPHeaders.ContentType != "" {
			r.logger.Warnf("ContentType received from request Metadata %s, as well as ContentType property %s, choosing value from contentType property", blobHTTPHeaders.ContentType, *req.ContentType)
		}
		blobHTTPHeaders.ContentType = *req.ContentType
	}

	if val, ok := req.Metadata[contentMD5]; ok && val != "" {
		sDec, err := b64.StdEncoding.DecodeString(val)
		if err != nil || len(sDec) != 16 {
			return azblob.BlobHTTPHeaders{}, fmt.Errorf("the MD5 value specified in Content MD5 is invalid, MD5 value must be 128 bits and base64 encoded")
		}
		blobHTTPHeaders.ContentMD5 = sDec
		delete(req.Metadata, contentMD5)
	}
	if val, ok := req.Metadata[contentEncoding]; ok && val != "" {
		blobHTTPHeaders.ContentEncoding = val
		delete(req.Metadata, contentEncoding)
	}
	if val, ok := req.Metadata[contentLanguage]; ok && val != "" {
		blobHTTPHeaders.ContentLanguage = val
		delete(req.Metadata, contentLanguage)
	}
	if val, ok := req.Metadata[contentDisposition]; ok && val != "" {
		blobHTTPHeaders.ContentDisposition = val
		delete(req.Metadata, contentDisposition)
	}
	if val, ok := req.Metadata[cacheControl]; ok && val != "" {
		blobHTTPHeaders.CacheControl = val
		delete(req.Metadata, cacheControl)
	}
	return blobHTTPHeaders, nil
}

func (r *StateStore) deleteFile(req *state.DeleteRequest) error {
	blobURL := r.containerURL.NewBlockBlobURL(getFileName(req.Key))
	accessConditions := azblob.BlobAccessConditions{}

	if req.Options.Concurrency == state.FirstWrite && req.ETag != nil {
		var etag string
		if req.ETag != nil {
			etag = *req.ETag
		}
		accessConditions.IfMatch = azblob.ETag(etag)
	}

	_, err := blobURL.Delete(context.Background(), azblob.DeleteSnapshotsOptionNone, accessConditions)
	if err != nil {
		r.logger.Debugf("delete file %s, err %s", req.Key, err)

		if req.ETag != nil {
			return state.NewETagError(state.ETagMismatch, err)
		}

		return err
	}

	return nil
}

func getFileName(key string) string {
	pr := strings.Split(key, keyDelimiter)
	if len(pr) != 2 {
		return pr[0]
	}

	return pr[1]
}

func (r *StateStore) marshal(req *state.SetRequest) []byte {
	var v string
	b, ok := req.Value.([]byte)
	if ok {
		v = string(b)
	} else {
		v, _ = jsoniter.MarshalToString(req.Value)
	}

	return []byte(v)
}

func isNotFoundError(err error) bool {
	azureError, ok := err.(azblob.StorageError)

	return ok && azureError.ServiceCode() == azblob.ServiceCodeBlobNotFound
}
