package cmd

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

	"github.com/bhojpur/service/pkg/serverless"
	_ "github.com/bhojpur/service/pkg/serverless/exec"
	_ "github.com/bhojpur/service/pkg/serverless/golang"
	_ "github.com/bhojpur/service/pkg/serverless/js"
	"github.com/bhojpur/service/pkg/utils"
	"github.com/spf13/cobra"
)

const (
	runtimeWaitTimeoutInSeconds = 60
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a Bhojpur Service Stream Function",
	Long:  "Run a Bhojpur Service Stream Function",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			opts.Filename = args[0]
		}
		// os signal
		// sigCh := make(chan os.Signal, 1)
		// signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
		// Serverless
		utils.InfoStatusEvent(os.Stdout, "Bhojpur Service Stream Function file: %v", opts.Filename)
		if !utils.IsExec(opts.Filename) && opts.Name == "" {
			utils.FailureStatusEvent(os.Stdout, "Bhojpur Service Stream Function's Name is empty, please set name used by `-n` flag")
			return
		}
		// resolve serverless
		utils.PendingStatusEvent(os.Stdout, "Create Bhojpur Service Stream Function instance...")
		if err := parseURL(url, &opts); err != nil {
			utils.FailureStatusEvent(os.Stdout, err.Error())
			return
		}
		s, err := serverless.Create(&opts)
		if err != nil {
			utils.FailureStatusEvent(os.Stdout, err.Error())
			return
		}
		if !s.Executable() {
			utils.InfoStatusEvent(os.Stdout,
				"Starting the Bhojpur Service Stream Function instance with Name: %s. Host: %s. Port: %d.",
				opts.Name,
				opts.Host,
				opts.Port,
			)
			// build
			utils.PendingStatusEvent(os.Stdout, "Bhojpur Service Stream Function building...")
			if err := s.Build(true); err != nil {
				utils.FailureStatusEvent(os.Stdout, err.Error())
				return
			}
			utils.SuccessStatusEvent(os.Stdout, "Success! Bhojpur Service Stream Function build.")
		} else { // executable
			utils.InfoStatusEvent(os.Stdout,
				"Starting the Bhojpur Service Stream Function instance with executable file: %s. Host: %s. Port: %d.",
				opts.Filename,
				opts.Host,
				opts.Port,
			)
		}
		// run
		utils.InfoStatusEvent(os.Stdout, "Bhojpur Service Stream Function is running...")
		if err := s.Run(verbose); err != nil {
			utils.FailureStatusEvent(os.Stdout, err.Error())
			return
		}
		// Exit
		// <-sigCh
		// utils.WarningStatusEvent(os.Stdout, "Terminated signal received: shutting down")
		// utils.InfoStatusEvent(os.Stdout, "Exited Bhojpur Service Stream Function instance.")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&opts.Filename, "file-name", "f", "app.go", "Stream function file")
	// runCmd.Flags().StringVarP(&opts.Lang, "lang", "l", "go", "source language")
	runCmd.Flags().StringVarP(&url, "url", "u", "localhost:9000", "Bhojpur Service-Processor endpoint addr")
	runCmd.Flags().StringVarP(&opts.Name, "name", "n", "", "Bhojpur Service stream function name. It should match the specific service name in Bhojpur Service-Processor config (workflow.yaml)")
	runCmd.Flags().StringVarP(&opts.ModFile, "modfile", "m", "", "custom go.mod")
	// runCmd.MarkFlagRequired("name")

}
