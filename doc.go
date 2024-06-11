// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package errors contains additional functions, interfaces and structs for
recording stack frames, applying basic formatting, working with goroutines,
multiple errors and custom error types.

It is inspired by package golang.org/x/xerrors and is designed to be a drop-in
replacement for it, as well as the standard library's errors package.

The errors.New and errors.Errorf functions create errors whose content is a
text message and whom can trace stack frames. errors.Wrap and errors.Wrapf
create errors by wrapping an existing error with a similar error like
errors.New and errors.Errorf.

# Msg

Instead of defining error messages as global variables, it is possible to define
them as constants using errors.Msg.

	const ErrSomethingWentWrong errors.Msg = "something went wrong"

# Formatting

Wrap an existing error with errors.WithFormatter to upgrade the error to
include basic formatting.
Formatting is done using xerrors.FormatError and thus the same verbs are
supported. Any error created with this package implements the fmt.Formatter
and xerrors.Formatter interfaces.

	fmt.Printf("%+v", errors.WithFormatter(err))

# Stack tracing

Every error can track stack trace information. Just wrap it with
errors.WithStack and a complete stack trace is captured.

	err = errors.WithStack(err)

An errors.StackTrace can be retrieved using errors.GetStackTrace.
Printing the error results in a trace similar to:

	invalid character 'i' looking for beginning of value:
		github.com/go-pogo/errors.ExampleWithStack
			/path/to/errors/examples_trace_test.go:43
		github.com/go-pogo/errors.ExampleWithStack.func1
			/path/to/errors/examples_trace_test.go:40

# Disable stack tracing

Stack tracing comes with a performance cost. For production environments this
cost can be undesirable. To disable stack tracing, compile your Go program with
the "notrace" tag.

	go build -tags=notrace

# Catching panics

A convenient function is available to catch panics and store them as an error.

	var err error
	defer errors.CatchPanic(&err)

# Backwards compatibility

Unwrap, Is, As are backwards compatible with the standard library's errors
package and act the same.
*/
package errors
