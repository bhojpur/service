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
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/bhojpur/service/pkg/utils"
	"github.com/bhojpur/service/pkg/serverless/golang"
)

var (
	name string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a Bhojpur Service Stream function",
	Long:  "Initialize a Bhojpur Service Stream function",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) >= 1 && args[0] != "" {
			name = args[0]
		}

		if name == "" {
			utils.FailureStatusEvent(os.Stdout, "Please input your app name")
			return
		}

		utils.PendingStatusEvent(os.Stdout, "Initializing the Stream Function...")
		// create app.go
		fname := filepath.Join(name, "app.go")
		if err := utils.PutContents(fname, golang.InitFuncTmpl); err != nil {
			utils.FailureStatusEvent(os.Stdout, "Write stream function into app.go file failure with the error: %v", err)
			return
		}

		utils.SuccessStatusEvent(os.Stdout, "Congratulations! You have initialized the stream function successfully.")
		utils.InfoStatusEvent(os.Stdout, "You can enjoy the Bhojpur Service Stream Function via the command: ")
		utils.InfoStatusEvent(os.Stdout, "\tDEV: \tsvcutl dev -n %s %s/app.go", "Noise", name)
		utils.InfoStatusEvent(os.Stdout, "\tPROD: \tFirst run source application, e.g: go run example/source/main.go\r\n\t\tSecond: svcutl run -n %s %s/app.go", name, name)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&name, "name", "n", "", "The name of Stream Function")
}
