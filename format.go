// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"
	"golang.org/x/xerrors"
)

// A Formatter formats error messages and prints them to a Printer.
// It is an alias of xerrors.Formatter.
type Formatter = xerrors.Formatter

// A Printer prints a formatted error. It is an alias of xerrors.Printer.
type Printer = xerrors.Printer

// WithFormatter wraps the error with a Formatter that is capable of basic
// error formatting. It returns the provided error as is if it already is a
// Formatter, or nil when err is nil.
func WithFormatter(err error) Formatter {
	if err == nil {
		return nil
	}

	if f, ok := err.(Formatter); ok {
		return f
	}

	return &embedError{error: err}
}

// FormatError calls the FormatError method of err with a Printer configured
// according to state and verb, and writes the result to state. It will wrap
// err If err is not a Formatter it will wrap err, so it is capable of basic
// error formatting using WithFormatter.
func FormatError(err error, state fmt.State, verb rune) {
	if err == nil {
		return
	}

	f, ok := err.(Formatter)
	if !ok {
		f = &embedError{error: err}
	}

	xerrors.FormatError(f, state, verb)
}

// PrintError prints the error err with the provided Printer and formats and
// prints the error's stack frames.
func PrintError(p Printer, err error) {
	if err == nil {
		return
	}

	p.Print(err.Error())
	if !p.Detail() {
		return
	}
	if stack := GetStackTrace(err); stack != nil {
		stack.Format(p)
	}
}
