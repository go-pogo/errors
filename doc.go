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

# StackTrace tracing

Every error can track stack trace information. Just wrap it with errors.WithStack
and a complete stack trace is captured.

	err = errors.WithStack(err)

Printing the error results in a trace similar to:

	some error: something happened:
	    main.main
	        /go-pogo/errors/.examples/3_with_kind/main.go:24
	    main.doSomething
	        /go-pogo/errors/.examples/3_with_kind/main.go:20
	    main.someAction
	        /go-pogo/errors/.examples/3_with_kind/main.go:16

# Formatting

Wrap an existing error with errors.WithFormatter to upgrade the error to
include basic formatting. Formatting is done using xerrors.FormatError and
thus the same verbs are supported.

	fmt.Printf("%+v", errors.WithFormatter(err))

# Catching panics

A convenient function is available to catch panics and store them as an error.

	var err error
	defer errors.CatchPanic(&err)

# Backwards compatibility

Unwrap, Is, As are backwards compatible with the standard library's errors
package and act the same.
*/
package errors
