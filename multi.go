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

// Combine returns a MultiError when more than one non-nil errors are provided.
// It returns a single error when only one error is passed, and nil if no
// non-nil errors are provided.
func Combine(errors ...error) error {
	return combine(Filter(errors))
}

func combine(errors []error) error {
	switch len(errors) {
	case 0:
		return nil
	case 1:
		return errors[0]
	}

	return newMultiErr(errors, 2)
}

const panicAppendNilPtr = "errors.Append: dest must not be a nil pointer"

// Append appends multiple non-nil errors to a single multi error dest.
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

type multiErr struct {
	errors []error
	stack  *StackTrace
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

// Errors returns the errors within the multi error.
func (m *multiErr) Errors() []error { return m.errors }

func (m *multiErr) StackTrace() *StackTrace { return m.stack }

func (m *multiErr) Is(target error) bool {
	for _, err := range m.errors {
		if Is(err, target) {
			return true
		}
	}
	return false
}

func (m *multiErr) As(target interface{}) bool {
	for _, err := range m.errors {
		if As(err, target) {
			return true
		}
	}
	return false
}

// Format uses xerrors.FormatError to call the FormatError method of the error
// with a xerrors.Printer configured according to s and v, and writes the
// result to s.
func (m *multiErr) Format(s fmt.State, v rune) { xerrors.FormatError(m, s, v) }

// FormatError prints a summary of the encountered errors to p.
func (m *multiErr) FormatError(p xerrors.Printer) error {
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
