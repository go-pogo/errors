// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/xerrors"
)

// WithFormatter wraps the error with an UpgradedError that is capable of basic
// error formatting, but only if it is not already wrapped.
func WithFormatter(parent error) xerrors.Formatter {
	switch e := parent.(type) {
	case *formatterErr:
		return e

	case UpgradedError:
		return toCommonErr(parent, true)
	}

	return &formatterErr{error: parent}
}

// FormatError calls the FormatError method of err with an xerrors.Printer
// configured according to state and verb, and writes the result to state.
// If err is not an xerrors.Formatter it will wrap the error with an
// UpgradedError that is capable of basic error formatting using WithFormatter.
func FormatError(err error, state fmt.State, verb rune) {
	f, ok := err.(xerrors.Formatter)
	if !ok {
		f = &formatterErr{err}
	}

	xerrors.FormatError(f, s, verb)
}

func PrintError(printer xerrors.Printer, err error) {
	printer.Print(err.Error())
	if printer.Detail() {
		frames := GetStackFrames(err)
		if frames != nil {
			frames.Format(printer)
		}
	}
}

type formatterErr struct{ error }

func (e *formatterErr) Original() error { return e.error }

// Format formats the error using FormatError.
func (e *formatterErr) Format(s fmt.State, v rune) { FormatError(e, s, v) }

// FormatError prints the error to the xerrors.Printer using PrintError and
// returns the next error in the error chain, if any.
func (e *formatterErr) FormatError(p xerrors.Printer) error {
	PrintError(p, e)
	return Unwrap(e.error)
}

// GoString prints a basic error syntax.
func (e *formatterErr) GoString() string {
	return goString(e, e.error)
}

const fullPkgName = "github.com/go-pogo/errors"

func goString(err, parent error) string {
	typ := reflect.TypeOf(err)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	var buf strings.Builder
	_, _ = fmt.Fprintf(&buf, "&\"%s\".%s", fullPkgName, typ.Name())

	if parent != nil {
		_, _ = fmt.Fprintf(&buf, "{error:%#v}", parent)
	} else {
		buf.WriteString("{}")
	}

	return buf.String()
}
