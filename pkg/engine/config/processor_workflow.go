package config

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
	"errors"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// App represents a Bhojpur Service workflow application.
type App struct {
	Name string `yaml:"name"`
}

// Workflow represents a Bhojpur Service workflow.
type Workflow struct {
	Functions []App `yaml:"functions"`
}

// WorkflowConfig represents a Bhojpur Service workflow config.
type WorkflowConfig struct {
	// Name represents the name of the processor.
	Name string `yaml:"name"`
	// Host represents the listening host of the processor.
	Host string `yaml:"host"`
	// Port represents the listening port of the processor.
	Port int `yaml:"port"`
	// Workflow represents the sfn workflow.
	Workflow `yaml:",inline"`
}

// LoadWorkflowConfig the WorkflowConfig by path.
func LoadWorkflowConfig(path string) (*WorkflowConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}

	return load(buffer)
}

func load(data []byte) (*WorkflowConfig, error) {
	var config = &WorkflowConfig{}
	err := yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// ParseWorkflowConfig parses the config.
func ParseWorkflowConfig(config string) (*WorkflowConfig, error) {
	if !(strings.HasSuffix(config, ".yaml") || strings.HasSuffix(config, ".yml")) {
		return nil, errors.New(`Bhojpur Service workflow: the extension of workflow config is incorrect, it should ".yaml|.yml"`)
	}

	// parse workflow.yaml
	wfConf, err := LoadWorkflowConfig(config)
	if err != nil {
		return nil, err
	}

	// validate
	err = validateWorkflowConfig(wfConf)
	if err != nil {
		return nil, err
	}

	return wfConf, nil
}

func validateWorkflowConfig(wfConf *WorkflowConfig) error {
	if wfConf == nil {
		return errors.New("conf is nil")
	}

	m := map[string][]App{
		"Functions": wfConf.Functions,
	}

	missingParams := []string{}
	for k, apps := range m {
		for _, app := range apps {
			if app.Name == "" {
				missingParams = append(missingParams, k)
			}
		}
	}

	errMsg := ""
	if wfConf.Name == "" || wfConf.Host == "" || wfConf.Port <= 0 {
		errMsg = "Missing name, host or port in workflow config. "
	}

	if len(missingParams) > 0 {
		errMsg += "Missing name, host or port in " + strings.Join(missingParams, ", "+". ")
	}

	if errMsg != "" {
		return errors.New(errMsg)
	}

	return nil
}
