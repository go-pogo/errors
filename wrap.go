// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"fmt"
)

// An Unwrapper unpacks a wrapped error.
type Unwrapper interface {
	error
	// Unwrap returns the next error in the error chain.
	// If there is no next error, Unwrap returns nil.
	Unwrap() (next error)
}

// Wrap creates a new error that wraps around the causing error, thus extending
// the error chain. It will only create a new error when the provided cause
// error is not nil, otherwise it will return nil.
func Wrap(cause error, text string) error {
	if cause == nil {
		return nil
	}
	return withCause(newErr(stderrors.New(text), 1), cause)
}

// Wrapf formats an error message according to a format specifier and provided
// arguments and creates a new error the same way Wrap() does.
func Wrapf(cause error, format string, a ...interface{}) error {
	if cause == nil {
		return nil
	}
	return withCause(newErr(fmt.Errorf(format, a...), 1), cause)
}

// UnwrapAll returns the complete chain of errors, starting with the supplied
// error and ending with the (upgraded) root cause error.
func UnwrapAll(err error) []error {
	var res []error
	for {
		if err == nil {
			break
		}
		res = append(res, err)
		err = Unwrap(err)
	}
	return res
}

// RootCause walks through all wrapped errors and returns the last (upgraded)
// error in the chain, which is the root cause error.
// To get the original non-upgraded root cause error use
//
//   Original(RootCause(err))
//
func RootCause(err error) error {
	for {
		unwrapped := Unwrap(err)
		if unwrapped == nil {
			break
		}

		err = unwrapped
	}
	return err
}
