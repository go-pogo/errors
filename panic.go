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

// MustPanicFormat is the template string used by the `Must()` function to
// format its panic message.
var MustPanicFormat = "errors.Must: %+v"

// Must panics when any of the given args is a non-nil error.
// Its message is the error message of the first encountered error.
func Must(args ...interface{}) {
	for _, arg := range args {
		if err, ok := arg.(error); ok && err != nil {
			panic(fmt.Sprintf(MustPanicFormat, err))
		}
	}
}
