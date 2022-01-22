// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"
)

// WrapPanic wraps a panicking sequence with the given prefix.
// It then panics again.
func WrapPanic(prefix string) {
	if r := recover(); r != nil {
		panic(fmt.Sprintf("%s: %s", prefix, r))
	}
}

// Must panics when any of the given args is a non-nil error.
// Its message is the error message of the first encountered error.
func Must(args ...interface{}) {
	for _, arg := range args {
		if err, ok := arg.(error); ok && err != nil {
			panic(fmt.Sprintf("errors.Must: %+v", err))
		}
	}
}

// CatchPanic recovers from a panic and wraps it in an error. It then calls
// Append with the provided dest *error and wrapped panic.
// Use CatchPanic directly with defer. It is not possible to use CatchPanic
// inside a deferred function, like:
//      defer func(){ CatchPanic(&err }()
func CatchPanic(dest *error) {
	if r := recover(); r != nil {
		Append(dest, newCommonErr(&panicError{v: r}, false))
		if st := GetStackTrace(*dest); st != nil {
			st.Skip = 1
		}
	}
}

type panicError struct{ v interface{} }

func (p *panicError) Unwrap() error {
	if e, ok := p.v.(error); ok {
		return e
	}
	return nil
}

func (p *panicError) Error() string {
	switch v := p.v.(type) {
	case error:
		return fmt.Sprintf("panic: %+v", p.v)
	case string:
		return "panic: " + v
	default:
		return fmt.Sprintf("panic: %v", p.v)
	}
}

func (p *panicError) GoString() string {
	return fmt.Sprintf("*panicError{v: %#v}", p.v)
}
