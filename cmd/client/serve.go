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

	"github.com/spf13/cobra"

	svcsvr "github.com/bhojpur/service/pkg/engine"
	"github.com/bhojpur/service/pkg/utils"
)

var meshConfURL string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run a Bhojpur Service-Processor",
	Long:  "Run a Bhojpur Service-Processor",
	Run: func(cmd *cobra.Command, args []string) {
		if config == "" {
			utils.FailureStatusEvent(os.Stdout, "Please input the file name of workflow config")
			return
		}
		// printBhojpurServerConf(conf)

		// endpoint := fmt.Sprintf("%s:%d", conf.Host, conf.Port)

		processor, err := svcsvr.NewProcessor(config)
		if err != nil {
			utils.FailureStatusEvent(os.Stdout, err.Error())
		}
		err = processor.ConfigMesh(meshConfURL)
		if err != nil {
			utils.FailureStatusEvent(os.Stdout, err.Error())
		}

		utils.InfoStatusEvent(os.Stdout, "Running Bhojpur Service-Processor...")
		err = processor.ListenAndServe()
		if err != nil {
			utils.FailureStatusEvent(os.Stdout, err.Error())
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVarP(&config, "config", "c", "workflow.yaml", "Workflow config file")
	serveCmd.Flags().StringVarP(&meshConfURL, "mesh-config", "m", "", "The URL of service mesh config")
	// serveCmd.MarkFlagRequired("config")
}

// func printBhojpurServerConf(wfConf *util.WorkflowConfig) {
// 	utils.InfoStatusEvent(os.Stdout, "Found %d stream functions in Bhojpur Service-Processor config", len(wfConf.Functions))
// 	for i, sfn := range wfConf.Functions {
// 		utils.InfoStatusEvent(os.Stdout, "Stream Function %d: %s", i+1, sfn.Name)
// 	}
// }
