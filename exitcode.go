// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

// ExitCodeGetter interfaces provide access to an exit code.
type ExitCodeGetter interface {
	error
	ExitCode() int
}

// WithExitCode adds an exit status code to the error which may indicate a
// fatal error. The exit code can be supplied to os.Exit to terminate the
// program immediately.
func WithExitCode(parent error, exitCode int) ExitCodeGetter {
	if parent == nil {
		return nil
	}

	if e, ok := parent.(exitCodeGetterSetter); ok {
		e.setExitCode(exitCode)
		return e
	}
	if _, ok := parent.(OriginalGetter); ok {
		ce := upgrade(parent)
		ce.setExitCode(exitCode)
		return ce
	}

	return &exitCodeErr{
		error:    parent,
		exitCode: exitCode,
	}
}

// GetExitCode returns an exit status code if the error implements the
// ExitCodeGetter interface. If not, it returns 0.
func GetExitCode(err error) int { return GetExitCodeOr(err, 0) }

// GetExitCodeOr returns an exit status code if the error implements the
// ExitCodeGetter interface. If not, it returns the provided value or.
func GetExitCodeOr(err error, or int) int {
	if e, ok := err.(ExitCodeGetter); ok {
		return e.ExitCode()
	}
	return or
}

type exitCodeGetterSetter interface {
	ExitCodeGetter
	setExitCode(c int)
}

type exitCodeErr struct {
	error
	exitCode int
}

func (ce *commonErr) setExitCode(c int)  { ce.exitCode = c }
func (e *exitCodeErr) setExitCode(c int) { e.exitCode = c }

func (e *exitCodeErr) Original() error { return e.error }
func (e *exitCodeErr) ExitCode() int   { return e.exitCode }
func (e *exitCodeErr) Error() string   { return errMsg(e.error.Error(), UnknownKind, e.exitCode) }
