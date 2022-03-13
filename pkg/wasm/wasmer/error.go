package wasmer

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

// #include <wasmer.h>
import "C"
import (
	"runtime"
	"unsafe"
)

// Error represents a Wasmer runtime error.
type Error struct {
	message string
}

func newErrorWith(message string) *Error {
	return &Error{
		message: message,
	}
}

func _newErrorFromWasmer() *Error {
	var errorLength = C.wasmer_last_error_length()

	if errorLength == 0 {
		return newErrorWith("(no error from Wasmer)")
	}

	var errorMessage = make([]C.char, errorLength)
	var errorMessagePointer = (*C.char)(unsafe.Pointer(&errorMessage[0]))

	var errorResult = C.wasmer_last_error_message(errorMessagePointer, errorLength)

	if errorResult == -1 {
		return newErrorWith("(failed to read last error from Wasmer)")
	}

	return newErrorWith(C.GoStringN(errorMessagePointer, errorLength-1))
}

func maybeNewErrorFromWasmer(block func() bool) *Error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if block() /* has failed */ {
		return _newErrorFromWasmer()
	}

	return nil
}

// Error returns the Error's message.
func (error *Error) Error() string {
	return error.message
}

// TrapError represents a trap produced during Wasm execution.
//
// See also
//
// Specification: https://webassembly.github.io/spec/core/intro/overview.html#trap
type TrapError struct {
	message string
	origin  *Frame
	trace   []*Frame
}

func newErrorFromTrap(pointer *C.wasm_trap_t) *TrapError {
	trap := newTrap(pointer, nil)

	return &TrapError{
		message: trap.Message(),
		origin:  trap.Origin(),
		trace:   trap.Trace().frames,
	}
}

// Error returns the TrapError's message.
func (self *TrapError) Error() string {
	return self.message
}

// Origin returns the TrapError's origin as a Frame.
func (self *TrapError) Origin() *Frame {
	return self.origin
}

// Trace returns the TrapError's trace as a Frame array.
func (self *TrapError) Trace() []*Frame {
	return self.trace
}
