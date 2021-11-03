// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/xerrors"
)

// WithFormatter wraps the error with an UpgradedError that is capable of basic
// error formatting, but only if it is not already wrapped.
func WithFormatter(parent error) xerrors.Formatter {
	if parent == nil {
		return nil
	}

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

	xerrors.FormatError(f, state, verb)
}

// PrintError prints the error err with the provided xerrors.Printer and
// additionally formats and prints the error's stack frames.
func PrintError(printer xerrors.Printer, err error) {
	printer.Print(err.Error())
	if !printer.Detail() {
		return
	}
	if frames := GetStackFrames(err); frames != nil {
		frames.Format(printer)
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

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(strings.Builder)
	},
}

func releaseBuf(buf *strings.Builder) {
	buf.Reset()
	bufPool.Put(buf)
}

func errMsg(msg string, kind Kind, code int) string {
	hasKind, hasCode := kind != UnknownKind, code != 0
	if !hasKind && !hasCode {
		return msg
	}

	buf := bufPool.Get().(*strings.Builder)
	defer releaseBuf(buf)

	if hasKind {
		if msg == "" {
			msg = kind.String()
		} else {
			buf.WriteString(kind.String())
			buf.WriteRune(':')
			buf.WriteRune(' ')
		}
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

const pkgImportPath = "github.com/go-pogo/errors"

func goString(err, parent error) string {
	typ := reflect.TypeOf(err)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	buf := bufPool.Get().(*strings.Builder)
	defer releaseBuf(buf)

	_, _ = fmt.Fprintf(buf, "&\"%s\".%s", pkgImportPath, typ.Name())

	if parent != nil {
		_, _ = fmt.Fprintf(buf, "{error:%#v}", parent)
	} else {
		buf.WriteString("{}")
	}

	return buf.String()
}
