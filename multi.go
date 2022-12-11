// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"
)

type MultiError interface {
	error
	Unwrap() []error
	// Deprecated: Use Unwrap instead.
	Errors() []error
}

// Filter returns a slice of errors without nil values in between them. It
// returns the slice with the length of the amount of non-nil errors but keeps
// its original capacity.
func Filter(errors []error) []error {
	n := 0
	for i, err := range errors {
		if err == nil {
			continue
		}
		if i != n {
			errors[i] = nil
			errors[n] = err
		}
		n++
	}
	return errors[:n]
}

// Join returns a MultiError when more than one non-nil errors are provided.
// It returns a single error when only one error is passed, and nil if no
// non-nil errors are provided.
func Join(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}
	return combine(Filter(errs))
}

// Combine returns a MultiError when more than one non-nil errors are provided.
// // It returns a single error when only one error is passed, and nil if no
// // non-nil errors are provided.
//
// Deprecated: Use Join instead.
func Combine(errs ...error) error { return Join(errs...) }

func combine(errs []error) error {
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	}

	return newMultiErr(errs, 2)
}

const (
	panicAppendNilPtr     = "errors.Append: dest must not be a nil pointer"
	panicAppendFuncNilPtr = "errors.AppendFunc: dest must not be a nil pointer"
	panicAppendFuncNilFn  = "errors.AppendFunc: fn must not be nil"
)

// Append appends multiple non-nil errors to a single multi error dest.
// When the value of dest is nil and errs only contains a single error, its
// value is set to the value of dest.
//
// Important: when using Append with defer, the pointer to the dest error
// must be a named return variable. For additional details see
// https://golang.org/ref/spec#Defer_statements.
func Append(dest *error, errs ...error) {
	if dest == nil {
		panic(panicAppendNilPtr)
	}

	for _, err := range errs {
		if err == nil {
			continue
		}

		switch d := (*dest).(type) {
		case nil:
			*dest = err

		case *multiErr:
			if traceStack {
				skipStackTrace(err, d.stack.Len())
			}
			d.errors = append(d.errors, err)

		default:
			*dest = newMultiErr([]error{*dest, err}, 1)
		}
	}
}

// AppendFunc appends the non-nil error result of fn to dest using Append.
func AppendFunc(dest *error, fn func() error) {
	if dest == nil {
		panic(panicAppendFuncNilPtr)
	}
	if fn == nil {
		panic(panicAppendFuncNilFn)
	}
	Append(dest, fn())
}

type multiErr struct {
	stack  *StackTrace
	errors []error
}

func newMultiErr(errors []error, skipFrames uint) *multiErr {
	m := &multiErr{errors: errors}
	if !traceStack {
		return m
	}

	m.stack = newStackTrace(skipFrames + 1)
	skip := m.stack.Len()
	for _, err := range m.errors {
		skipStackTrace(err, skip)
	}
	return m
}

// Unwrap returns the errors within the multi error.
func (m *multiErr) Unwrap() []error { return m.errors }

// Errors returns the errors within the multi error.
//
// Deprecated: Use Unwrap instead.
func (m *multiErr) Errors() []error { return m.errors }

func (m *multiErr) StackTrace() *StackTrace { return m.stack }

// Format uses xerrors.FormatError to call the FormatError method of the error
// with a Printer configured according to s and v, and writes the result to s.
func (m *multiErr) Format(s fmt.State, v rune) { xerrors.FormatError(m, s, v) }

// FormatError prints a summary of the encountered errors to p.
func (m *multiErr) FormatError(p Printer) error {
	p.Print(m.Error())
	if p.Detail() {
		m.stack.Format(p)

		l := len(m.errors)
		for i, err := range m.errors {
			p.Printf("\n[%d/%d] %+v\n", i+1, l, err)
		}
	}

	return nil
}

func (m *multiErr) Error() string {
	var buf strings.Builder
	buf.WriteString("multiple errors occurred:")

	l := len(m.errors)
	for i, e := range m.errors {
		_, _ = fmt.Fprintf(&buf, "\n[%d/%d] %s", i+1, l, e.Error())
		if i < l-1 {
			buf.WriteRune(';')
		}
	}
	return buf.String()
}
