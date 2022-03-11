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
	"github.com/bhojpur/service/pkg/utils"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the Bhojpur Service Stream Function",
	Long:  "Build the Bhojpur Service Stream Function as binary file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			opts.Filename = args[0]
		}
		utils.InfoStatusEvent(os.Stdout, "Bhojpur Service Stream Function file: %v", opts.Filename)
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
		utils.InfoStatusEvent(os.Stdout,
			"Starting the Bhojpur Service Stream Function instance with Name: %s. Host: %s. Port: %d.",
			opts.Name,
			opts.Host,
			opts.Port,
		)
		// build
		utils.PendingStatusEvent(os.Stdout, "Bhojpur Service Stream Function function building...")
		if err := s.Build(true); err != nil {
			utils.FailureStatusEvent(os.Stdout, err.Error())
			return
		}
		utils.SuccessStatusEvent(os.Stdout, "Success! Bhojpur Service Stream Function build.")
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringVarP(&opts.Filename, "file-name", "f", "app.go", "Stream function file (default is app.go)")
	buildCmd.Flags().StringVarP(&url, "url", "u", "localhost:9000", "Bhojpur Serice-Processor endpoint addr")
	buildCmd.Flags().StringVarP(&opts.Name, "name", "n", "", "Bhojpur Service stream function app name (required). It should match the specific service name in Bhojpur Service-Processor config (workflow.yaml)")
	buildCmd.MarkFlagRequired("name")
	buildCmd.Flags().StringVarP(&opts.ModFile, "modfile", "m", "", "custom go.mod")
}
