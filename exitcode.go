// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

// ExitCoder interfaces provide access to an exit code.
type ExitCoder interface {
	error
	ExitCode() int
}

// WithExitCode adds an exit status code to the error which may indicate a
// fatal error. The exit code can be supplied to os.Exit to terminate the
// program immediately.
func WithExitCode(parent error, exitCode int) ExitCoder {
	if parent == nil {
		return nil
	}

	switch e := parent.(type) {
	case *exitCodeErr:
		e.exitCode = exitCode
		return e

	case UpgradedError:
		ce := toCommonErr(parent, true)
		ce.exitCode = exitCode
		return ce
	}

	return &exitCodeErr{
		error:    parent,
		exitCode: exitCode,
	}
}

// GetExitCode returns an exit status code if the error implements the
// ExitCoder interface. If not, it returns 0.
func GetExitCode(err error) int { return GetExitCodeOr(err, 0) }

// GetExitCodeOr returns an exit status code if the error implements the
// ExitCoder interface. If not, it returns the provided value or.
func GetExitCodeOr(err error, or int) int {
	if e, ok := err.(ExitCoder); ok {
		return e.ExitCode()
	}
	return or
}

type exitCodeErr struct {
	error
	exitCode int
}

func (e *exitCodeErr) Original() error { return e.error }
func (e *exitCodeErr) ExitCode() int   { return e.exitCode }
func (e *exitCodeErr) Error() string   { return errMsg(e.error.Error(), UnknownKind, e.exitCode) }
