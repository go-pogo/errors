// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package errors contains additional functions, interfaces and structs for
recording stack frames, applying basic formatting, working with goroutines,
multiple errors and custom error types.

It is inspired by package golang.org/x/xerrors and is designed to be a drop-in
replacement for it, as well as the standard library's errors package.

The New and Newf functions create errors whose content is a text message and
whom can trace stack frames. Wrap and Wrapf create errors by wrapping an
existing error with a similar error like New and Newf.

Stack tracing

Every error can track stack trace information. Just wrap it with errors.Trace
and an additional stack frame is captured and stored within the error.

	err = errors.Trace(err)

Printing the error results in:
	some error: something happened:
		main.doSomething
			.../errors/examples/2_trace/main.go:17
		main.someAction
			.../errors/examples/2_trace/main.go:12

Formatting

Wrap an existing error with errors.WithFormatter to upgrade the error to
include basic formatting. Formatting is done using xerrors.FormatError and
thus the same verbs are supported.

    mt.Printf("%+v", errors.WithFormatter(err))

Catching panics

A convenient function is available to catch panics and store them as an error.

	var err error
	defer errors.CatchPanic(&err)

Backwards compatibility

Unwrap, Is, As are backwards compatible with the standard library's errors
package and act the same.
*/
package errors
