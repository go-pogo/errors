// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package errors implements functions to manipulate errors, record stack frames
// and apply basic formatting to errors.
// It is inspired by the golang.org/x/xerrors package and is design to be a
// drop in replacement for it as well as the standard library's errors package.
// The package contains additional functions, interfaces and structs for working
// with goroutines, multiple errors and custom error types.
//
// The New and Newf functions create errors whose content is a text message and
// who can trace stack frames. Wrap and Wrapf create errors by wrapping an
// existing error with a similar error like New and Newf.
//
// The Unwrap, Is and As functions work on errors that may wrap other errors.
// An error wraps another error if its type has the method
//
//	Unwrap() error
//
// If e.Unwrap() returns a non-nil error w, then we say that e wraps w.
//
// Unwrap unpacks wrapped errors. If its argument's type has an
// Unwrap method, it calls the method once. Otherwise, it returns nil.
//
// A simple way to create wrapped errors is to call Wrap or Wrapf. Another
// options i to create an error with Newf and apply the %w verb to the error
// argument:
//
//	errors.Unwrap(errors.Newf("... %w ...", ..., err, ...))
//
// returns err.
//
// Is, As, Opaque are backwards compatible with the standard library's error
// package and act the same.
package errors
