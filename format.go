// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"

	"golang.org/x/xerrors"
)

// WithFormatter wraps the error with a xerrors.Formatter that is capable of
// basic error formatting. It returns the provided error as is if it already is
// a xerrors.Formatter, or nil when err is nil.
func WithFormatter(err error) xerrors.Formatter {
	if err == nil {
		return nil
	}

	if f, ok := err.(xerrors.Formatter); ok {
		return f
	}

	return &embedError{error: err}
}

// FormatError calls the FormatError method of err with a xerrors.Printer
// configured according to state and verb, and writes the result to state.
// It will wrap err If err is not a xerrors.Formatter it will wrap err, so it
// is capable of basic error formatting using WithFormatter.
func FormatError(err error, state fmt.State, verb rune) {
	if err == nil {
		return
	}

	f, ok := err.(xerrors.Formatter)
	if !ok {
		f = &embedError{error: err}
	}

	xerrors.FormatError(f, state, verb)
}

// PrintError prints the error err with the provided xerrors.Printer and
// additionally formats and prints the error's stack frames.
func PrintError(printer xerrors.Printer, err error) {
	if err == nil {
		return
	}

	printer.Print(err.Error())
	if !printer.Detail() {
		return
	}
	if stack := GetStackTrace(err); stack != nil {
		stack.Format(printer)
	}
}
