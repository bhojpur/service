package golang

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
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"strings"

	"os/exec"
	"path/filepath"
	"runtime"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/imports"

	"github.com/bhojpur/service/pkg/serverless"
	"github.com/bhojpur/service/pkg/utils"
)

// GolangServerless defines Go implementation of Serverless stream function interface.
type GolangServerless struct {
	opts    *serverless.Options
	source  string
	target  string
	tempDir string
}

// Init initializes the serverless stream function
func (s *GolangServerless) Init(opts *serverless.Options) error {
	// now := time.Now()
	// msg := "Init: serverless function..."
	// initSpinning := utils.Spinner(os.Stdout, msg)
	// defer initSpinning(utils.Failure)

	s.opts = opts
	if !utils.Exists(s.opts.Filename) {
		return fmt.Errorf("the file %s doesn't exist", s.opts.Filename)
	}

	// generate source code
	source := utils.GetBinContents(s.opts.Filename)
	if len(source) < 1 {
		return fmt.Errorf(`"%s" content is empty`, s.opts.Filename)
	}

	// append main function
	ctx := Context{
		Name: s.opts.Name,
		Host: s.opts.Host,
		Port: s.opts.Port,
	}

	// determine: Reactive Stream serverless or raw bytes serverless function.
	isRx := strings.Contains(string(source), "rx.Stream")
	mainFuncTmpl := ""
	if isRx {
		mainFuncTmpl = string(MainFuncRxTmpl)
	} else {
		mainFuncTmpl = string(MainFuncRawBytesTmpl)
	}

	mainFunc, err := RenderTmpl(mainFuncTmpl, &ctx)
	if err != nil {
		return fmt.Errorf("Init: %s", err)
	}
	source = append(source, mainFunc...)
	// utils.InfoStatusEvent(os.Stdout, "merge source elapse: %v", time.Since(now))
	// Create the AST by parsing src
	fset := token.NewFileSet()
	astf, err := parser.ParseFile(fset, "", source, 0)
	if err != nil {
		return fmt.Errorf("Init: parse source file err %s", err)
	}
	// Add import packages
	astutil.AddNamedImport(fset, astf, "", "github.com/bhojpur/service")
	// astutil.AddNamedImport(fset, astf, "stdlog", "log")
	// utils.InfoStatusEvent(os.Stdout, "import elapse: %v", time.Since(now))
	// Generate the code
	code, err := generateCode(fset, astf)
	if err != nil {
		return fmt.Errorf("Init: generate code err %s", err)
	}
	// Create a temp folder.
	tempDir, err := ioutil.TempDir("", "bhojpur_")
	if err != nil {
		return err
	}
	s.tempDir = tempDir
	tempFile := filepath.Join(tempDir, "app.go")
	// Fix imports
	fixedSource, err := imports.Process(tempFile, code, nil)
	if err != nil {
		return fmt.Errorf("Init: imports %s", err)
	}
	// utils.InfoStatusEvent(os.Stdout, "fix import elapse: %v", time.Since(now))
	if err := utils.PutContents(tempFile, fixedSource); err != nil {
		return fmt.Errorf("Init: write file err %s", err)
	}
	// utils.InfoStatusEvent(os.Stdout, "final write file elapse: %v", time.Since(now))
	// mod
	name := strings.ReplaceAll(opts.Name, " ", "_")
	cmd := exec.Command("go", "mod", "init", name)
	cmd.Dir = tempDir
	env := os.Environ()
	env = append(env, fmt.Sprintf("GO111MODULE=%s", "on"))
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("Init: go mod init err %s", out)
		return err
	}

	// TODO: check if is already built in temp dir by MD5
	s.source = tempFile
	return nil
}

// Build compiles the serverless to executable
func (s *GolangServerless) Build(clean bool) error {
	// check if the file exists
	appPath := s.source
	if _, err := os.Stat(appPath); os.IsNotExist(err) {
		return fmt.Errorf("the file %s doesn't exist", appPath)
	}
	// env
	env := os.Environ()
	env = append(env, fmt.Sprintf("GO111MODULE=%s", "on"))
	// use custom go.mod
	if s.opts.ModFile != "" {
		mfile, _ := filepath.Abs(s.opts.ModFile)
		if !utils.Exists(mfile) {
			return fmt.Errorf("the mod file %s doesn't exist", mfile)
		}
		// go.mod
		utils.WarningStatusEvent(os.Stdout, "Use custom go.mod: %s", mfile)
		tempMod := filepath.Join(s.tempDir, "go.mod")
		utils.Copy(mfile, tempMod)
		// source := file.GetContents(tempMod)
		// utils.InfoStatusEvent(os.Stdout, "go.mod: %s", source)
		// mod download
		cmd := exec.Command("go", "mod", "tidy")
		cmd.Env = env
		cmd.Dir = s.tempDir
		out, err := cmd.CombinedOutput()
		if err != nil {
			err = fmt.Errorf("Build: go mod tidy err %s", out)
			return err
		}
	} else {
		// Upgrade modules that provide packages imported by packages in the main module
		cmd := exec.Command("go", "get", "-d", "-u", "./...")
		cmd.Dir = s.tempDir
		cmd.Env = env
		out, err := cmd.CombinedOutput()
		if err != nil {
			err = fmt.Errorf("Build: go get err %s", out)
			return err
		}
	}
	// build
	goos := runtime.GOOS
	dir, _ := filepath.Split(s.opts.Filename)
	sl, _ := filepath.Abs(dir + "sl.basm")

	// clean build
	if clean {
		defer func() {
			utils.Remove(s.tempDir)
		}()
	}
	s.target = sl
	// fmt.Printf("goos=%s\n", goos)
	if goos == "windows" {
		sl, _ = filepath.Abs(dir + "sl.exe")
		s.target = sl
	}
	// go build
	cmd := exec.Command("go", "build", "-ldflags", "-s -w", "-o", sl, appPath)
	cmd.Env = env
	cmd.Dir = s.tempDir
	// utils.InfoStatusEvent(os.Stdout, "Build: cmd: %+v", cmd)
	// source := file.GetContents(s.source)
	// utils.InfoStatusEvent(os.Stdout, "source: %s", source)
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("Build: failure %s", out)
		return err
	}
	return nil
}

// Run compiles and runs the serverless stream function
func (s *GolangServerless) Run(verbose bool) error {
	utils.InfoStatusEvent(os.Stdout, "Run Go serverless: %s", s.target)
	cmd := exec.Command(s.target)
	if verbose {
		cmd.Env = []string{"BHOJPUR_SERVICE_LOG_LEVEL=debug"}
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (s *GolangServerless) Executable() bool {
	return false
}

func generateCode(fset *token.FileSet, file *ast.File) ([]byte, error) {
	var output []byte
	buffer := bytes.NewBuffer(output)
	if err := printer.Fprint(buffer, fset, file); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func init() {
	serverless.Register(&GolangServerless{}, ".go")
}
