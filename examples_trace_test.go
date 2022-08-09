// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !notrace
// +build !notrace

package errors

import (
	"encoding/json"
	stdfmt "fmt"
	"path/filepath"
	"runtime"
	"strings"
)

func ExamplePrintWithDetails() {
	doSomething := func() error {
		return New("something happened")
	}

	err := doSomething()
	var fmt = new(ignoreThisLineUseFmtPackageInActualCode)
	fmt.Printf("%+v\n", err)
	// Output:
	// something happened:
	//     github.com/go-pogo/errors.ExamplePrintWithDetails
	//         /path/to/errors/examples_trace_test.go:23
	//     github.com/go-pogo/errors.ExamplePrintWithDetails.func1
	//         /path/to/errors/examples_trace_test.go:20
}

func ExampleWithStack() {
	type Result struct{}

	doSomething := func() (*Result, error) {
		dest := new(Result)
		err := json.Unmarshal([]byte("invalid"), &dest) // this wil result in an error
		return dest, WithStack(err)
	}

	_, err := doSomething()
	var fmt = new(ignoreThisLineUseFmtPackageInActualCode)
	fmt.Printf("%+v\n", err)
	// Output:
	// invalid character 'i' looking for beginning of value:
	//     github.com/go-pogo/errors.ExampleWithStack
	//         /path/to/errors/examples_trace_test.go:43
	//     github.com/go-pogo/errors.ExampleWithStack.func1
	//         /path/to/errors/examples_trace_test.go:40
}

// ignoreThisLineUseFmtPackageInActualCode limits the frame trace stack to the
// last two entries and normalizes the file paths within these frame.
// It is only used to make sure the above examples stay as close to real code
// as possible and to keep the output the same on different systems.
type ignoreThisLineUseFmtPackageInActualCode struct{}

func (fmt *ignoreThisLineUseFmtPackageInActualCode) Printf(format string, err error) {
	st := GetStackTrace(err)
	if st == nil {
		panic("only use this method when a stack trace is present!")
	}

	// skip test runtime related stack traces
	st.Skip = st.Len() - 2

	// get actual path of this file
	_, f, _, _ := runtime.Caller(0)
	f = strings.ReplaceAll(filepath.Dir(f), "\\", "/")

	output := stdfmt.Sprintf(format, err)

	// normalize file paths of stack trace entries
	output = strings.ReplaceAll(output, f, "/path/to/errors")
	stdfmt.Println(output)
}
