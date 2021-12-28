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
		if err, ok := r.(error); ok {
			Append(dest, WithFormatter(err))
		} else {
			Append(dest, newCommonErr(&panicErr{v: r}, false))
		}
		if st := GetStackTrace(*dest); st != nil {
			st.Skip = 1
		}
	}
}

type panicErr struct{ v interface{} }

func (e *panicErr) Error() string {
	switch v := e.v.(type) {
	case error:
		return fmt.Sprintf("%+v", e.v)
	case string:
		return v
	default:
		return fmt.Sprintf("%#v", e.v)
	}
}
