// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"
	"reflect"

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

const panicUseWithStackInstead = "errors.New: use errors.WithStack instead to wrap an error with an errors.StackTracer and xerrors.Formatter"

// New creates a new error which implements the StackTracer, Wrapper and
// xerrors.Formatter interfaces. Argument msg can be either a string or Msg.
//
//    err := errors.New("my error message")
//    err := errors.New(errors.Msg("my error message"))
//
// New records a stack trace at the point it was called. Each call returns a
// distinct error value even if msg is identical. It will return nil if msg is
// nil. Use WithStack to wrap an existing error with a StackTracer and
// xerrors.Formatter.
func New(msg interface{}) error {
	if msg == nil {
		return nil
	}

	switch v := msg.(type) {
	case *commonError:
		return v

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
		panic(UnsupportedTypeError{
			Func: "errors.New",
			Type: reflect.TypeOf(v).String(),
		})
	}
}

// Newf formats an error message according to a format specifier and provided
// arguments with fmt.Errorf, and creates a new error similar to New.
//
//    err := errors.Newf("my error %s", "message")
//    err := errors.Newf("my error: %w", cause)
func Newf(format string, args ...interface{}) error {
	if len(args) == 0 {
		return newCommonErr(Msg(format), true)
	}
	return withPossibleCause(newCommonErr(fmt.Errorf(format, args...), true))
}

func (m Msg) Is(target error) bool {
	switch t := target.(type) {
	case Msg:
		return m == t
	case *Msg:
		return m == *t
	default:
		return false
	}
}

func (m Msg) As(target interface{}) bool {
	if t, ok := target.(*Msg); ok {
		*t = m
		return true
	}
	return false
}

func (m Msg) String() string { return string(m) }

func (m Msg) Error() string { return string(m) }

func (m Msg) GoString() string { return `errors.Msg("` + string(m) + `")` }

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

// Unwrap returns the next error in the error chain. It returns nil if there
// is not a next error.
func (ce *commonError) Unwrap() error { return ce.cause }

func (ce *commonError) Is(target error) bool {
	if m, ok := ce.error.(Msg); ok {
		return m.Is(target)
	}
	return false
}

func (ce *commonError) As(target interface{}) bool {
	if t, ok := target.(*commonError); ok {
		*t = *ce
		return true
	}
	return false
}

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
	if ce.cause == nil {
		return fmt.Sprintf("*commonError{error: %#v}", ce.error)
	}

	return fmt.Sprintf(
		"*commonError{error: %#v, cause: %#v}",
		ce.error,
		ce.cause,
	)
}

type UnsupportedTypeError struct {
	Func, Type string
}

func (ut *UnsupportedTypeError) Error() string {
	return ut.Func + ": unsupported type `" + ut.Type + "`"
}
