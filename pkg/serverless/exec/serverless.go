package exec

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

	"os/exec"

	"github.com/bhojpur/service/pkg/serverless"
	"github.com/bhojpur/service/pkg/utils"
)

// ExecServerless defines executable file implementation of Serverless stream function interface.
type ExecServerless struct {
	target string
}

// Init initializes the serverless stream function
func (s *ExecServerless) Init(opts *serverless.Options) error {
	if !utils.Exists(opts.Filename) {
		return fmt.Errorf("the file %s doesn't exist", opts.Filename)
	}
	s.target = opts.Filename

	return nil
}

// Build compiles the serverless stream function to executable
func (s *ExecServerless) Build(clean bool) error {
	return nil
}

// Run compiles and runs the serverless stream function
func (s *ExecServerless) Run(verbose bool) error {
	utils.InfoStatusEvent(os.Stdout, "Run Execute serverless: %s", s.target)
	cmd := exec.Command(s.target)
	if verbose {
		cmd.Env = []string{"BHOJPUR_SERVICE_LOG_LEVEL=debug"}
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (s *ExecServerless) Executable() bool {
	return true
}

func init() {
	serverless.Register(&ExecServerless{}, ".basm", ".exe")
}
