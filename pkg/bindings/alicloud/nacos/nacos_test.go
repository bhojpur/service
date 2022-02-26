package nacos

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
	"os"
	"path"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bhojpur/service/pkg/bindings"
	"github.com/bhojpur/service/pkg/utils/logger"
)

func TestInputBindingRead(t *testing.T) { //nolint:paralleltest
	m := bindings.Metadata{Name: "test", Properties: nil}
	var err error
	m.Properties, err = getNacosLocalCacheMetadata()
	require.NoError(t, err)
	n := NewNacos(logger.NewLogger("test"))
	err = n.Init(m)
	require.NoError(t, err)
	var count int32
	ch := make(chan bool, 1)

	handler := func(in *bindings.ReadResponse) ([]byte, error) {
		require.Equal(t, "hello", string(in.Data))
		atomic.AddInt32(&count, 1)
		ch <- true

		return nil, nil
	}

	go func() {
		err = n.Read(handler)
		require.NoError(t, err)
	}()

	select {
	case <-ch:
		require.Equal(t, int32(1), atomic.LoadInt32(&count))
	case <-time.After(time.Second):
		require.FailNow(t, "read timeout")
	}
}

func getNacosLocalCacheMetadata() (map[string]string, error) {
	tmpDir := "/tmp/config"
	dataID := "test"
	group := "DEFAULT_GROUP"

	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("create dir failed. %w", err)
	}

	cfgFile := path.Join(tmpDir, fmt.Sprintf("%s@@%s@@", dataID, group))
	file, err := os.OpenFile(cfgFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil || file == nil {
		return nil, fmt.Errorf("open %s failed. %w", cfgFile, err)
	}

	defer func() {
		_ = file.Close()
	}()

	if _, err = file.WriteString("hello"); err != nil {
		return nil, fmt.Errorf("write file failed. %w", err)
	}

	return map[string]string{
		"cacheDir":   "/tmp", // default
		"nameServer": "localhost:8080/fake",
		"watches":    fmt.Sprintf("%s:%s", dataID, group),
	}, nil
}
