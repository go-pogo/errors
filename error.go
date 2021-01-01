// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/xerrors"
)

// New is an alias of errors.New. It returns an error that formats as the given
// text. Each call to New returns a distinct error value even if the text is
// identical.
func New(text string) error {
	err := toCommonErr(stderrors.New(text), false)
	err.Trace(1)
	return err
}

// Newf formats an error message according to a format specifier and provided
// arguments and creates a new error the same way New does. It serves as an
// alternative to fmt.Errorf.
func Newf(format string, a ...interface{}) error {
	err := toCommonErr(fmt.Errorf(format, a...), false)
	err.Trace(1)
	return err
}

// Upgrade upgrades the given standard error by wrapping it with a Proxy that
// can record stack frames and has basic error formatting.
// The original parent error can always be retrieved by calling Original on the
// result of Upgrade. Thus
//
//   Original(Upgrade(err)) == err
//
// equals true.
func Upgrade(parent error) error {
	return toCommonErr(parent, true)
}

// toCommonErr upgrades the parent error by wrapping it with a commonErr.
func toCommonErr(parent error, upgrade bool) *commonErr {
	if e, ok := parent.(*commonErr); ok {
		return e
	}

	ce := &commonErr{
		error:   Original(parent),
		upgrade: upgrade,
	}

	switch e := parent.(type) {
	case *kindErr:
		ce.kind = e.kind
	}

	return ce
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

// Original returns the original error before it was upgraded. This is never the
// case for errors that were created with New, Newf, Wrap of Wrapf.
func (ce *commonErr) Original() error {
	if ce.upgrade {
		return ce.error
	}
	return ce
}

func (ce *commonErr) Kind() Kind {
	if ce.kind != UnknownKind {
		return ce.kind
	}
	if e, ok := ce.error.(Kinder); ok {
		return e.Kind()
	}

	return UnknownKind
}

func (ce *commonErr) ExitCode() int { return ce.exitCode }

// Format formats the error using FormatError.
func (ce *commonErr) Format(s fmt.State, v rune) { FormatError(ce, s, v) }

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
// is no next error.
func (ce *commonErr) Unwrap() error {
	if ce.cause != nil {
		return ce.cause
	}
	return Unwrap(ce.error)
}

func (ce *commonErr) Error() string {
	return errMsg(ce.error.Error(), ce.Kind(), ce.exitCode)
}

// GoString prints a basic error syntax.
func (ce *commonErr) GoString() string {
	return goString(ce, ce.error)
}

func errMsg(msg string, kind Kind, code int) string {
	hasKind, hasCode := kind != UnknownKind, code != 0
	if !hasKind && !hasCode {
		return msg
	}

	var buf strings.Builder
	if hasKind {
		buf.WriteString(kind.String())
		buf.WriteRune(':')
		buf.WriteRune(' ')
	}
	if hasCode {
		buf.WriteRune('[')
		buf.WriteString(strconv.Itoa(code))
		buf.WriteRune(']')
		buf.WriteRune(' ')
	}

	buf.WriteString(msg)
	return buf.String()
}
