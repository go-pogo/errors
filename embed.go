// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"
	"golang.org/x/xerrors"
)

// An Embedder unwraps an embedded error.
type Embedder interface {
	error
	Unembed() error
}

// Unembed recursively unwraps all Embedder errors and returns the (original)
// error that was wrapped with extra context.
// If err is not an Embedder, Unembed returns err as provided.
func Unembed(err error) error {
	if u, ok := err.(Embedder); ok {
		return Unembed(u.Unembed())
	}
	return err
}

type embedError struct {
	error
	stack *StackTrace
}

// StackTrace returns the StackTrace of this error or, if nil, tries to return
// the first non-nil StackTrace of an embedded error.
func (e *embedError) StackTrace() *StackTrace {
	if e.stack == nil {
		e.stack = GetStackTrace(Unembed(e.error))
	}
	return e.stack
}

// Unembed the underlying error.
func (e *embedError) Unembed() error { return e.error }

// Unwrap the underlying error.
func (e *embedError) Unwrap() error { return e.error }

// Format uses xerrors.FormatError to call the FormatError method of the error
// with a Printer configured according to s and v, and writes the result to s.
func (e *embedError) Format(s fmt.State, v rune) {
	xerrors.FormatError(e, s, v)
}

// FormatError prints the error to the Printer using PrintError and returns the
// next error in the error chain, if any.
func (e *embedError) FormatError(p Printer) error {
	PrintError(p, e)
	return Unwrap(Unembed(e.error))
}

// GoString prints the error in basic Go syntax.
func (e *embedError) GoString() string {
	return fmt.Sprintf("*embedError{error: %#v}", e.error)
}
