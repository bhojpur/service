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

// Nacos is an easy-to-use dynamic service discovery, configuration and service management platform
//
// See https://github.com/nacos-group/nacos-sdk-go/

import (
	"errors"
	"fmt"
	"time"

	"github.com/bhojpur/service/pkg/utils/config"
)

type Settings struct {
	NameServer           string        `mapstructure:"nameServer"`
	Endpoint             string        `mapstructure:"endpoint"`
	RegionID             string        `mapstructure:"region"`
	NamespaceID          string        `mapstructure:"namespace"`
	AccessKey            string        `mapstructure:"accessKey"`
	SecretKey            string        `mapstructure:"secretKey"`
	Timeout              time.Duration `mapstructure:"timeout"`
	CacheDir             string        `mapstructure:"cacheDir"`
	UpdateThreadNum      int           `mapstructure:"updateThreadNum"`
	NotLoadCacheAtStart  bool          `mapstructure:"notLoadCacheAtStart"`
	UpdateCacheWhenEmpty bool          `mapstructure:"updateCacheWhenEmpty"`
	Username             string        `mapstructure:"username"`
	Password             string        `mapstructure:"password"`
	LogDir               string        `mapstructure:"logDir"`
	// TODO: implement LogRollingConfig
	//RotateTime         string        `mapstructure:"rotateTime"`
	//MaxAge             int           `mapstructure:"maxAge"`
	LogLevel string `mapstructure:"logLevel"`
	Config   string `mapstructure:"config"`
	Watches  string `mapstructure:"watches"`
}

func (s *Settings) Decode(in interface{}) error {
	return config.Decode(in, s)
}

func (s *Settings) Validate() error {
	if s.Timeout <= 0 {
		return fmt.Errorf("invalid timeout %s", s.Timeout)
	}

	if s.Endpoint == "" && s.NameServer == "" {
		return errors.New("either endpoint or nameserver must be configured")
	}

	return nil
}
