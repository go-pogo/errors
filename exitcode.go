// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"
)

// ExitCoder interfaces provide access to an exit code.
type ExitCoder interface {
	error
	ExitCode() int
}

// ExitCoderSetter interfaces provide access to an exit code which can be
// changed.
type ExitCoderSetter interface {
	ExitCoder
	SetExitCode(int)
}

// WithExitCode adds an exit status code to the error which may indicate a
// fatal error. The exit code can be supplied to [os.Exit] to terminate the
// program immediately.
func WithExitCode(err error, exitCode int) ExitCoder {
	if err == nil {
		return nil
	}

	//goland:noinspection GoTypeAssertionOnErrors
	if e, ok := err.(ExitCoderSetter); ok {
		e.SetExitCode(exitCode)
		return e
	}

	return &exitCodeError{
		embedError: &embedError{error: err},
		exitCode:   exitCode,
	}
}

// GetExitCode returns an exit status code if the error implements the
// [ExitCoder] interface, otherwise it returns 0.
func GetExitCode(err error) int { return GetExitCodeOr(err, 0) }

// GetExitCodeOr returns the exit status code from the first found [ExitCoder]
// in err's error chain. If none is found, it returns the provided value or.
func GetExitCodeOr(err error, or int) int {
	for {
		//goland:noinspection GoTypeAssertionOnErrors
		if e, ok := err.(ExitCoder); ok {
			return e.ExitCode()
		}
		err = Unwrap(err)
		if err == nil {
			break
		}
	}

	return or
}

type exitCodeError struct {
	*embedError
	exitCode int
}

func (e *exitCodeError) SetExitCode(c int) { e.exitCode = c }
func (e *exitCodeError) ExitCode() int     { return e.exitCode }

// GoString prints the error in basic Go syntax.
func (e *exitCodeError) GoString() string {
	return fmt.Sprintf(
		"errors.exitCodeError{exitCode: %d, embedErr: %#v}",
		e.exitCode,
		e.error,
	)
}
