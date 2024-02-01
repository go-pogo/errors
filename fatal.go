// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"
	"os"
)

// FatalOnErr prints the error to stderr and exits the program with an exit
// code that is not 0. When err is an ExitCoder its exit code is used.
func FatalOnErr(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\nFatal error: %+v\n", err)
		os.Exit(GetExitCodeOr(err, 1))
	}
}

// PanicOnErr panics when err is not nil.
func PanicOnErr(err error) {
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
}
