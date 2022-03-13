package wasmedge

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

// #include <wasmedge/wasmedge.h>
import "C"

type Validator struct {
	_inner *C.WasmEdge_ValidatorContext
	_own   bool
}

func NewValidator() *Validator {
	validator := C.WasmEdge_ValidatorCreate(nil)
	if validator == nil {
		return nil
	}
	return &Validator{_inner: validator, _own: true}
}

func NewValidatorWithConfig(conf *Configure) *Validator {
	validator := C.WasmEdge_ValidatorCreate(conf._inner)
	if validator == nil {
		return nil
	}
	return &Validator{_inner: validator, _own: true}
}

func (self *Validator) Validate(ast *AST) error {
	res := C.WasmEdge_ValidatorValidate(self._inner, ast._inner)
	if !C.WasmEdge_ResultOK(res) {
		return newError(res)
	}
	return nil
}

func (self *Validator) Release() {
	if self._own {
		C.WasmEdge_ValidatorDelete(self._inner)
	}
	self._inner = nil
	self._own = false
}
