// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

type ExitCoder interface {
	error
	ExitCode() int
}

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

func GetExitCode(err error) int {
	if e, ok := err.(ExitCoder); ok {
		return e.ExitCode()
	}

	return 0
}

type exitCodeErr struct {
	error
	exitCode int
}

func (e *exitCodeErr) Original() error { return e.error }
func (e *exitCodeErr) ExitCode() int   { return e.exitCode }
func (e *exitCodeErr) Error() string   { return errMsg(e.error.Error(), UnknownKind, e.exitCode) }
