// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"
	"strings"

	"github.com/go-pogo/errors/internal"
	"golang.org/x/xerrors"
)

// MultiError is an error which unwraps into multiple underlying errors.
type MultiError interface {
	error
	Unwrap() []error
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

// Join returns a [MultiError] when more than one non-nil errors are provided.
// It returns a single error when only one error is passed, and nil if no
// non-nil errors are provided.
func Join(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	errs = Filter(errs)
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	}

	return newMultiErr(errs, 2)
}

// Append creates a [MultiError] from two non-nil errors. If left is already a
// multi error created via this package, the other error is appended to it.
// If either of the errors is nil, the other error is returned.
func Append(left, right error) error {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}

	//goland:noinspection GoTypeAssertionOnErrors
	if m, ok := left.(*multiErr); ok {
		m.append(right)
		return m
	}
	return newMultiErr([]error{left, right}, 1)
}

const (
	panicAppendIntoNilPtr = "errors.AppendInto: dest must not be a nil pointer"
	panicAppendFuncNilPtr = "errors.AppendFunc: dest must not be a nil pointer"
	panicAppendFuncNilFn  = "errors.AppendFunc: fn must not be nil"
)

// AppendInto appends multiple non-nil errors to a single multi error dest.
// When the value of dest is nil and errs only contains a single error, its
// value is set to the value of dest.
//
// Important: when using [AppendInto] with defer, the pointer to the dest error
// must be a named return variable. For additional details see
// https://golang.org/ref/spec#Defer_statements.
func AppendInto(dest *error, errs ...error) (errored bool) {
	if dest == nil {
		panic(panicAppendIntoNilPtr)
	}

	var multi *multiErr
	for _, err := range errs {
		if err == nil {
			continue
		}

		if multi != nil {
			multi.append(err)
			continue
		}

		//goland:noinspection GoTypeAssertionOnErrors
		if *dest == nil {
			*dest = err
		} else if m, ok := (*dest).(*multiErr); ok {
			multi = m
			multi.append(err)
		} else {
			multi = newMultiErr([]error{*dest, err}, 1)
			*dest = multi
		}
	}
	return multi != nil
}

// AppendFunc appends the non-nil error result of fn to dest using
// [AppendInto].
func AppendFunc(dest *error, fn func() error) {
	if dest == nil {
		panic(panicAppendFuncNilPtr)
	}
	if fn == nil {
		panic(panicAppendFuncNilFn)
	}
	AppendInto(dest, fn())
}

// multiErr is an error which unwraps into multiple underlying errors.
type multiErr struct {
	stack *StackTrace
	msg   string
	errs  []error
}

func newMultiErr(errs []error, skipFrames uint) *multiErr {
	m := &multiErr{errs: errs}
	if !internal.TraceStack {
		return m
	}

	m.stack = newStackTrace(skipFrames + 1)
	skip := m.stack.Len()
	for _, err := range m.errs {
		skipStackTrace(err, skip)
	}
	return m
}

func (m *multiErr) append(err error) {
	if internal.TraceStack {
		skipStackTrace(err, m.stack.Len())
	}
	m.errs = append(m.errs, err)
}

func (m *multiErr) StackTrace() *StackTrace { return m.stack }

// Unwrap returns the errors within the [multiErr].
func (m *multiErr) Unwrap() []error { return m.errs }

// Errors returns the errors within the multi error.
//
// Deprecated: Use [Unwrap] instead.
func (m *multiErr) Errors() []error { return m.errs }

// Format uses [xerrors.FormatError] to call the [FormatError] method of the
// error with a [Printer] configured according to s and v, and writes the
// result to s.
func (m *multiErr) Format(s fmt.State, v rune) { xerrors.FormatError(m, s, v) }

// FormatError prints a summary of the encountered errors to p.
func (m *multiErr) FormatError(p Printer) error {
	p.Print(m.Error())
	if !p.Detail() {
		return nil
	}

	m.stack.Format(p)
	p.Print("\n")

	l := len(m.errs)
	for i, err := range m.errs {
		p.Printf("[%d/%d] %+v\n", i+1, l, err)
		//goland:noinspection GoTypeAssertionOnErrors
		if _, ok := err.(StackTracer); ok {
			p.Print("\n")
		}
	}
	return nil
}

func (m *multiErr) Error() string {
	if m.msg != "" {
		return m.msg
	}

	var buf strings.Builder
	buf.WriteString("multiple errors occurred:")

	l := len(m.errs)
	for i, e := range m.errs {
		_, _ = fmt.Fprintf(&buf, "\n[%d/%d] %s", i+1, l, e.Error())
		if i < l-1 {
			buf.WriteRune(';')
		}
	}
	return buf.String()
}
