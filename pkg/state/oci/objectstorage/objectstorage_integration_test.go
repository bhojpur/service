package objectstorage

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

// run the test
// go test -v github.com/bhojpur/service/pkg/state/oci/objectstorage.

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/service/pkg/state"
	"github.com/bhojpur/service/pkg/utils/logger"
)

const (
	configFilePathEnvKey    = "APP_TEST_OCI_CONFIG_FILE_PATH"
	bucketNameEnvKey        = "APP_TEST_OCI_BUCKET_NAME"
	compartmentOCIDEnvKey   = "APP_TEST_OCI_COMPARTMENT_OCID"
	configFileProfileEnvKey = "APP_TEST_OCI_CONFIG_PROFILE"
)

func getConfigFilePathString() string {
	return os.Getenv(configFilePathEnvKey)
}

func getConfigProfile() string {
	return os.Getenv(configFileProfileEnvKey)
}

func getCompartmentOCID() string {
	return os.Getenv(compartmentOCIDEnvKey)
}

func getBucketName() string {
	return os.Getenv(bucketNameEnvKey)
}

func TestOCIObjectStorageIntegration(t *testing.T) {
	ociObjectStorageConfiguration := map[string]string{
		"configFileAuthentication": "true",
	}
	configFilePath := getConfigFilePathString()
	if configFilePath == "" {
		// first run export APP_TEST_OCI_CONFIG_FILE_PATH="/home/app/.oci/config".
		t.Skipf("OCI ObjectStorage state integration tests skipped. To enable define the configuration file path string using environment variable '%s' (example 'export %s=\"/home/app/.oci/config\")", configFilePathEnvKey, configFilePathEnvKey)
	}
	ociObjectStorageConfiguration["configFilePath"] = configFilePath
	ociObjectStorageConfiguration["configFileProfile"] = getConfigProfile()

	compartmentOCID := getCompartmentOCID()
	if compartmentOCID == "" {
		t.Skipf("OCI ObjectStorage state integration tests skipped. To enable define the compartment OCID string using environment variable '%s' (example 'export %s=\"ocid1.compartment.oc1..aaaaaaaacsssekayq4d34nl5h3eqs5e6ak3j5s4jhlws7rr5pxmt3zrq\")", compartmentOCIDEnvKey, compartmentOCIDEnvKey)
	}
	ociObjectStorageConfiguration["compartmentOCID"] = compartmentOCID

	bucketName := getBucketName()
	if bucketName == "" {
		bucketName = "APP_TEST_BUCKET"
	}
	ociObjectStorageConfiguration["bucketName"] = bucketName

	t.Run("Test Get", func(t *testing.T) {
		t.Parallel()
		testGet(t, ociObjectStorageConfiguration)
	})
	t.Run("Test Set", func(t *testing.T) {
		t.Parallel()
		testSet(t, ociObjectStorageConfiguration)
	})
	t.Run("Test Delete", func(t *testing.T) {
		t.Parallel()
		testDelete(t, ociObjectStorageConfiguration)
	})
	t.Run("Test Ping", func(t *testing.T) {
		t.Parallel()
		testPing(t, ociObjectStorageConfiguration)
	})
}

func testGet(t *testing.T, ociProperties map[string]string) {
	statestore := NewOCIObjectStorageStore(logger.NewLogger("logger"))
	meta := state.Metadata{}
	meta.Properties = ociProperties

	t.Run("Get an non-existing key", func(t *testing.T) {
		err := statestore.Init(meta)
		assert.Nil(t, err)
		getResponse, err := statestore.Get(&state.GetRequest{Key: "xyzq"})
		assert.Equal(t, &state.GetResponse{}, getResponse, "Response must be empty")
		assert.NoError(t, err, "Non-existing key must not be treated as error")
	})
	t.Run("Get an existing key", func(t *testing.T) {
		err := statestore.Init(meta)
		assert.Nil(t, err)
		err = statestore.Set(&state.SetRequest{Key: "test-key", Value: []byte("test-value")})
		assert.Nil(t, err)
		getResponse, err := statestore.Get(&state.GetRequest{Key: "test-key"})
		assert.Nil(t, err)
		assert.Equal(t, "test-value", string(getResponse.Data), "Value retrieved should be equal to value set")
		assert.NotNil(t, *getResponse.ETag, "ETag should be set")
	})
	t.Run("Get an existing composed key", func(t *testing.T) {
		err := statestore.Init(meta)
		assert.Nil(t, err)
		err = statestore.Set(&state.SetRequest{Key: "test-app||test-key", Value: []byte("test-value")})
		assert.Nil(t, err)
		getResponse, err := statestore.Get(&state.GetRequest{Key: "test-app||test-key"})
		assert.Nil(t, err)
		assert.Equal(t, "test-value", string(getResponse.Data), "Value retrieved should be equal to value set")
	})
	t.Run("Get an unexpired state element with TTL set", func(t *testing.T) {
		testKey := "unexpired-ttl-test-key"
		err := statestore.Init(meta)
		assert.Nil(t, err)
		err = statestore.Set(&state.SetRequest{Key: testKey, Value: []byte("test-value"), Metadata: (map[string]string{
			"ttlInSeconds": "100",
		})})
		assert.Nil(t, err)
		getResponse, err := statestore.Get(&state.GetRequest{Key: testKey})
		assert.Nil(t, err)
		assert.Equal(t, "test-value", string(getResponse.Data), "Value retrieved should be equal to value set despite TTL setting")
	})
	t.Run("Get a state element with TTL set to -1 (not expire)", func(t *testing.T) {
		testKey := "never-expiring-ttl-test-key"
		err := statestore.Init(meta)
		assert.Nil(t, err)
		err = statestore.Set(&state.SetRequest{Key: testKey, Value: []byte("test-value"), Metadata: (map[string]string{
			"ttlInSeconds": "-1",
		})})
		assert.Nil(t, err)
		getResponse, err := statestore.Get(&state.GetRequest{Key: testKey})
		assert.Nil(t, err)
		assert.Equal(t, "test-value", string(getResponse.Data), "Value retrieved should be equal (TTL setting of -1 means never expire)")
	})
	t.Run("Get an expired (TTL in the past) state element", func(t *testing.T) {
		err := statestore.Init(meta)
		assert.Nil(t, err)
		err = statestore.Set(&state.SetRequest{Key: "ttl-test-key", Value: []byte("test-value"), Metadata: (map[string]string{
			"ttlInSeconds": "1",
		})})
		assert.Nil(t, err)
		time.Sleep(time.Second * 2)
		getResponse, err := statestore.Get(&state.GetRequest{Key: "ttl-test-key"})
		assert.Equal(t, &state.GetResponse{}, getResponse, "Response must be empty")
		assert.NoError(t, err, "Expired element must not be treated as error")
	})
}

func testSet(t *testing.T, ociProperties map[string]string) {
	meta := state.Metadata{}
	meta.Properties = ociProperties
	statestore := NewOCIObjectStorageStore(logger.NewLogger("logger"))
	t.Run("Set without a key", func(t *testing.T) {
		err := statestore.Init(meta)
		assert.Nil(t, err)
		err = statestore.Set(&state.SetRequest{Value: []byte("test-value")})
		assert.Equal(t, err, fmt.Errorf("key for value to set was missing from request"), "Lacking Key results in error")
	})
	t.Run("Regular Set Operation", func(t *testing.T) {
		testKey := "local-test-key"
		err := statestore.Init(meta)
		assert.Nil(t, err)

		err = statestore.Set(&state.SetRequest{Key: testKey, Value: []byte("test-value")})
		assert.Nil(t, err, "Setting a value with a proper key should be errorfree")
		getResponse, err := statestore.Get(&state.GetRequest{Key: testKey})
		assert.Nil(t, err)
		assert.Equal(t, "test-value", string(getResponse.Data), "Value retrieved should be equal to value set")
		assert.NotNil(t, *getResponse.ETag, "ETag should be set")
	})
	t.Run("Regular Set Operation with composite key", func(t *testing.T) {
		testKey := "test-app||other-test-key"
		err := statestore.Init(meta)
		assert.Nil(t, err)

		err = statestore.Set(&state.SetRequest{Key: testKey, Value: []byte("test-value")})
		assert.Nil(t, err, "Setting a value with a proper composite key should be errorfree")
		getResponse, err := statestore.Get(&state.GetRequest{Key: testKey})
		assert.Nil(t, err)
		assert.Equal(t, "test-value", string(getResponse.Data), "Value retrieved should be equal to value set")
		assert.NotNil(t, *getResponse.ETag, "ETag should be set")
	})
	t.Run("Regular Set Operation with TTL", func(t *testing.T) {
		testKey := "test-key-with-ttl"
		err := statestore.Set(&state.SetRequest{Key: testKey, Value: []byte("test-value"), Metadata: (map[string]string{
			"ttlInSeconds": "500",
		})})
		assert.Nil(t, err, "Setting a value with a proper key and a correct TTL value should be errorfree")
		err = statestore.Set(&state.SetRequest{Key: testKey, Value: []byte("test-value"), Metadata: (map[string]string{
			"ttlInSeconds": "XXX",
		})})
		assert.NotNil(t, err, "Setting a value with a proper key and a incorrect TTL value should be produce an error")
	})

	t.Run("Testing Set & Concurrency (ETags)", func(t *testing.T) {
		testKey := "etag-test-key"
		err := statestore.Init(meta)
		assert.Nil(t, err)

		err = statestore.Set(&state.SetRequest{Key: testKey, Value: []byte("test-value")})
		assert.Nil(t, err, "Setting a value with a proper key should be errorfree")
		getResponse, _ := statestore.Get(&state.GetRequest{Key: testKey})
		etag := getResponse.ETag
		err = statestore.Set(&state.SetRequest{Key: testKey, Value: []byte("overwritten-value"), ETag: etag, Options: state.SetStateOption{
			Concurrency: state.FirstWrite,
		}})
		assert.Nil(t, err, "Updating value with proper etag should go fine")

		err = statestore.Set(&state.SetRequest{Key: testKey, Value: []byte("more-overwritten-value"), ETag: etag, Options: state.SetStateOption{
			Concurrency: state.FirstWrite,
		}})
		assert.NotNil(t, err, "Updating value with the old etag should be refused")

		// retrieve the latest etag - assigned by the previous set operation.
		getResponse, _ = statestore.Get(&state.GetRequest{Key: testKey})
		assert.NotNil(t, *getResponse.ETag, "ETag should be set")
		etag = getResponse.ETag
		err = statestore.Set(&state.SetRequest{Key: testKey, Value: []byte("more-overwritten-value"), ETag: etag, Options: state.SetStateOption{
			Concurrency: state.FirstWrite,
		}})
		assert.Nil(t, err, "Updating value with the latest etag should be accepted")
	})
}

func testDelete(t *testing.T, ociProperties map[string]string) {
	m := state.Metadata{}
	m.Properties = ociProperties
	s := NewOCIObjectStorageStore(logger.NewLogger("logger"))
	t.Run("Delete without a key", func(t *testing.T) {
		err := s.Init(m)
		assert.Nil(t, err)
		err = s.Delete(&state.DeleteRequest{})
		assert.Equal(t, err, fmt.Errorf("key for value to delete was missing from request"), "Lacking Key results in error")
	})
	t.Run("Regular Delete Operation", func(t *testing.T) {
		testKey := "test-key"
		err := s.Init(m)
		assert.Nil(t, err)

		err = s.Set(&state.SetRequest{Key: testKey, Value: []byte("test-value")})
		assert.Nil(t, err, "Setting a value with a proper key should be errorfree")
		err = s.Delete(&state.DeleteRequest{Key: testKey})
		assert.Nil(t, err, "Deleting an existing value with a proper key should be errorfree")
	})
	t.Run("Regular Delete Operation for composite key", func(t *testing.T) {
		testKey := "test-app||some-test-key"
		err := s.Init(m)
		assert.Nil(t, err)

		err = s.Set(&state.SetRequest{Key: testKey, Value: []byte("test-value")})
		assert.Nil(t, err, "Setting a value with a proper composite key should be errorfree")
		err = s.Delete(&state.DeleteRequest{Key: testKey})
		assert.Nil(t, err, "Deleting an existing value with a proper composite key should be errorfree")
	})
	t.Run("Delete with an unknown key", func(t *testing.T) {
		err := s.Delete(&state.DeleteRequest{Key: "unknownKey"})
		assert.Contains(t, err.Error(), "404", "Unknown Key results in error: http status code 404, object not found")
	})

	t.Run("Testing Delete & Concurrency (ETags)", func(t *testing.T) {
		testKey := "etag-test-delete-key"
		err := s.Init(m)
		assert.Nil(t, err)
		// create document.
		err = s.Set(&state.SetRequest{Key: testKey, Value: []byte("test-value")})
		assert.Nil(t, err, "Setting a value with a proper key should be errorfree")
		getResponse, _ := s.Get(&state.GetRequest{Key: testKey})
		etag := getResponse.ETag

		incorrectETag := "someRandomETag"
		err = s.Delete(&state.DeleteRequest{Key: testKey, ETag: &incorrectETag, Options: state.DeleteStateOption{
			Concurrency: state.FirstWrite,
		}})
		assert.NotNil(t, err, "Deleting value with an incorrect etag should be prevented")

		err = s.Delete(&state.DeleteRequest{Key: testKey, ETag: etag, Options: state.DeleteStateOption{
			Concurrency: state.FirstWrite,
		}})
		assert.Nil(t, err, "Deleting value with proper etag should go fine")
	})
}

func testPing(t *testing.T, ociProperties map[string]string) {
	m := state.Metadata{}
	m.Properties = ociProperties
	s := NewOCIObjectStorageStore(logger.NewLogger("logger"))
	t.Run("Ping", func(t *testing.T) {
		err := s.Init(m)
		assert.Nil(t, err)
		err = s.Ping()
		assert.Nil(t, err, "Ping should be successful")
	})
}
