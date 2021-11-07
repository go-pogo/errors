// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"fmt"

	"golang.org/x/xerrors"
)

// New is an alias of errors.New. It returns an error that formats as the given
// text. Each call to New returns a distinct error value even if the text is
// identical.
func New(text string) error {
	return newErr(stderrors.New(text), 1)
}

// Newf formats an error message according to a format specifier and provided
// arguments and creates a new error the same way New does. It serves as an
// alternative to fmt.Errorf.
func Newf(format string, a ...interface{}) error {
	return newErr(fmt.Errorf(format, a...), 1)
}

// An OriginalGetter is capable of returning its original error.
type OriginalGetter interface {
	error
	// Original returns the original error that resides in the OriginalGetter.
	Original() (original error)
}

// Original returns the Original error if err is an OriginalGetter. Otherwise, it
// will return the given error err.
func Original(err error) error {
	if p, ok := err.(OriginalGetter); ok {
		return p.Original()
	} else {
		return err
	}
}

type Frames []xerrors.Frame

// StackTracer interfaces provide access to a stack of traced Frames.
type StackTracer interface {
	error

	// StackFrames returns a slice of captured xerrors.Frame types associated
	// with the error.
	StackFrames() *Frames
	// Trace captures a xerrors.Frame that describes a frame on the caller's
	// stack. The argument skipFrames is the number of frames to skip over.
	Trace(skipFrames uint)
}

type tracer struct {
	frames Frames
}

type commonErr struct {
	error
	tracer

	// upgrade indicates whether this commonErr is the original error (= false)
	// or if the error in the error property is the original error (= true)
	upgrade  bool
	cause    error // cause of this error, if any
	kind     Kind
	exitCode int
}

func newErr(parent error, trace uint) *commonErr {
	ce := &commonErr{
		error:    parent,
		kind:     GetKind(parent),
		exitCode: GetExitCode(parent),
	}
	if trace > 0 {
		ce.Trace(trace + 1)
	}
	return ce
}

// upgrade upgrades the parent error by wrapping it with a commonErr.
func upgrade(parent error) *commonErr {
	if e, ok := parent.(*commonErr); ok {
		return e
	}

	return &commonErr{
		error:    Original(parent),
		upgrade:  true,
		kind:     GetKind(parent),
		exitCode: GetExitCode(parent),
	}
}

func withCause(ce *commonErr, cause error) *commonErr {
	if fr := GetStackFrames(cause); fr != nil {
		*fr = []xerrors.Frame(*fr)[:len(ce.frames)]
	}

	ce.cause = cause
	return ce
}

// Original returns the original error before it was upgraded. This is never the
// case for errors that were created with New, Newf, Wrap of Wrapf.
func (ce *commonErr) Original() error {
	if ce.upgrade {
		return ce.error
	}
	return ce
}

func (ce *commonErr) Kind() Kind { return ce.kind }

func (ce *commonErr) ExitCode() int { return ce.exitCode }

// Format uses xerrors.FormatError to call the FormatError method of the error
// with a xerrors.Printer configured according to s and v, and writes the
// result to s.
func (ce *commonErr) Format(s fmt.State, v rune) { xerrors.FormatError(ce, s, v) }

// FormatError prints the error to the xerrors.Printer using PrintError and
// returns the next error in the error chain, if any.
func (ce *commonErr) FormatError(p xerrors.Printer) error {
	PrintError(p, ce)
	return ce.Unwrap()
}

// todo: implement correct as method
func (ce *commonErr) As(target interface{}) bool {
	return As(ce.error, target)
}

// Unwrap returns the next error in the error chain. It returns nil if there
// is not a next error.
func (ce *commonErr) Unwrap() error {
	if ce.cause != nil {
		return ce.cause
	}
	if ce.upgrade {
		return Unwrap(ce.error)
	}
	return nil
}

func (ce *commonErr) Error() string {
	return errMsg(ce.error.Error(), ce.Kind(), ce.exitCode)
}

// GoString prints a basic error syntax.
func (ce *commonErr) GoString() string {
	return goString(ce, ce.error)
}
