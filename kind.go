// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"

	"golang.org/x/xerrors"
)

// UnknownKind is the default Kind for errors that are created without a
// distinct Kind.
const UnknownKind Kind = ""

// Kind describes the kind/type of error that has occurred. For example "auth
// error", "unmarshal error", etc. Errors can be of the same Kind but still
// contain different underlying causes.
// It is recommended to define each Kind as a constant.
type Kind string

func (k Kind) Is(target error) bool {
	switch t := target.(type) {
	case Kind:
		return k == t
	case *Kind:
		return k == *t
	}
	return false
}

func (k Kind) As(target interface{}) bool {
	if t, ok := target.(*Kind); ok {
		*t = k
		return true
	}
	return false
}

// String returns the string representation of Kind.
func (k Kind) String() string { return string(k) }

func (k Kind) Error() string { return string(k) }

func (k Kind) GoString() string { return `errors.Kind("` + string(k) + `")` }

// WithKind wraps an error with Kind, therefore extending the error chain.
func WithKind(err error, kind Kind) error {
	if err == nil || kind == UnknownKind {
		return err
	}

	e := &kindError{
		kind:  kind,
		cause: err,
	}
	if traceStack {
		e.stack = newStackTrace(1)
		skipStackTrace(err, e.stack.Len())
	}
	return e
}

// GetKind returns the Kind of the error if it is added with WithKind. If not,
// it returns UnknownKind.
func GetKind(err error) Kind { return GetKindOr(err, UnknownKind) }

// GetKindOr returns the Kind of the error if it is added with WithKind. If not,
// it returns the provided Kind or.
func GetKindOr(err error, or Kind) Kind {
	err = Unembed(err)
	if e, ok := err.(*kindError); ok {
		return e.kind
	}
	return or
}

type kindError struct {
	kind  Kind
	cause error
	stack *StackTrace
}

func (e *kindError) StackTrace() *StackTrace { return e.stack }

func (e *kindError) Unwrap() error { return e.cause }

func (e *kindError) Is(target error) bool {
	return e.kind.Is(target)
}

func (e *kindError) As(target interface{}) bool {
	return e.kind.As(target)
}

// Format uses xerrors.FormatError to call the FormatError method of the error
// with a xerrors.Printer configured according to s and v, and writes the
// result to s.
func (e *kindError) Format(s fmt.State, v rune) {
	xerrors.FormatError(e, s, v)
}

// FormatError prints the error to the xerrors.Printer using PrintError and
// returns the next error in the error chain.
func (e *kindError) FormatError(p xerrors.Printer) error {
	PrintError(p, e)
	if _, ok := e.cause.(Msg); ok {
		return nil
	}
	if !p.Detail() {
		// prevent showing cause's error message twice
		return Unwrap(Unembed(e.cause))
	}
	return e.cause
}

func (e *kindError) Error() string { return e.kind.String() + ": " + e.cause.Error() }

// GoString prints the error in basic Go syntax.
func (e *kindError) GoString() string {
	return fmt.Sprintf(
		"*kindError{kind: %s, cause: %#v}",
		e.kind.String(),
		e.cause,
	)
}
