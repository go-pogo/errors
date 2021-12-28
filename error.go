// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"

	"golang.org/x/xerrors"
)

// Msg is a string alias which can also be used as a basic error. This is
// particularly useful for defining constants of known errors in your library
// or application.
//
//    const ErrMyErrorMessage errors.Msg = "my error message"
//    const ErrAnotherError   errors.Msg = "just another error"
//
// A new error can be constructed from any Msg with New and is considered to be
// equal when comparing with Is.
//
//    err := errors.New(ErrMyErrorMessage)
//    errors.Is(err, ErrMyErrorMessage) // true
type Msg string

// String returns the string representation of Kind.
func (p Msg) String() string { return string(p) }

func (p Msg) Error() string { return string(p) }

const panicUseWithStackInstead = "errors.New: use errors.WithStack instead to wrap an error with an errors.StackTracer and xerrors.Formatter"

// New creates a new error which implements the StackTracer, Wrapper and
// xerrors.Formatter interfaces. Argument msg can be either a string or Msg.
//
//    err := errors.New("my error message")
//    err := errors.New(errors.Msg("my error message"))
//
// Each call to New returns a distinct error value even if msg is identical.
// Use WithStack to wrap an existing error with a StackTracer and
// xerrors.Formatter.
func New(msg interface{}) error {
	if msg == nil {
		return nil
	}

	switch v := msg.(type) {
	case string:
		return newCommonErr(Msg(v), true)
	case *string:
		return newCommonErr(Msg(*v), true)

	case Msg:
		return newCommonErr(v, true)
	case *Msg:
		return newCommonErr(*v, true)

	case error:
		panic(panicUseWithStackInstead)

	default:
		panic(fmt.Sprintf("errors.New: unsupported type `%T`", v))
	}
}

// Newf formats an error message according to a format specifier and provided
// arguments with fmt.Errorf, and creates a new error similar to New.
//
//    err := errors.Newf("my error %s", "message")
//    err := errors.Newf("my error: %w", causingErr)
func Newf(format string, args ...interface{}) error {
	return withPossibleCause(newCommonErr(fmt.Errorf(format, args...), true))
}

type commonError struct {
	error
	cause error
	stack *StackTrace
}

func newCommonErr(parent error, trace bool) *commonError {
	ce := &commonError{error: parent}
	if traceStack && trace {
		ce.stack = newStackTrace(2)
	}
	return ce
}

func withCause(ce *commonError, cause error) *commonError {
	ce.cause = cause
	if traceStack && ce.stack != nil {
		skipStackTrace(cause, ce.stack.Len())
	}
	return ce
}

func withPossibleCause(ce *commonError) *commonError {
	if w, ok := ce.error.(xerrors.Wrapper); ok {
		if cause := w.Unwrap(); cause != nil {
			return withCause(ce, cause)
		}
	}
	return ce
}

func (ce *commonError) StackTrace() *StackTrace { return ce.stack }

func (ce *commonError) Unwrap() error { return ce.cause }

// Format uses xerrors.FormatError to call the FormatError method of the error
// with a xerrors.Printer configured according to s and v, and writes the
// result to s.
func (ce *commonError) Format(s fmt.State, v rune) {
	xerrors.FormatError(ce, s, v)
}

// FormatError prints the error to the xerrors.Printer using PrintError and
// returns the next error in the error chain, if any.
func (ce *commonError) FormatError(p xerrors.Printer) error {
	PrintError(p, ce)
	return ce.cause
}

// GoString prints the error in basic Go syntax.
func (ce *commonError) GoString() string {
	return fmt.Sprintf(
		"*commonError{error: %#v, cause: %#v}",
		ce.error,
		ce.cause,
	)
}
