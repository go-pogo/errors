// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"fmt"

	"golang.org/x/xerrors"
)

// A Wrapper provides context around another error, which can be retrieved with
// Unwrap.
type Wrapper interface {
	error
	xerrors.Wrapper
}

// Wrap creates a new error that wraps around the causing error, thus extending
// the error chain. It will only create a new error when the provided cause
// error is not nil, otherwise it will return nil.
// Wrap also records the stack trace at the point it was called.
func Wrap(cause error, text string) error {
	if cause == nil {
		return nil
	}
	return withCause(newCommonErr(stderrors.New(text), true), cause)
}

// Wrapf formats an error message according to a format specifier and provided
// arguments and creates a new error the same way Wrap does.
func Wrapf(cause error, format string, a ...interface{}) error {
	if cause == nil {
		return nil
	}
	return withCause(newCommonErr(fmt.Errorf(format, a...), true), cause)
}

// Opaque is an alias of xerrors.Opaque. It returns an error with the same error
// formatting as err but that does not match err and cannot be unwrapped.
func Opaque(err error) error { return xerrors.Opaque(err) }

// Unwrap is an alias of errors.Unwrap. It returns the result of calling the
// Unwrap method on err, if err's type contains an Unwrap method returning
// error. Otherwise, Unwrap returns nil.
func Unwrap(err error) error { return stderrors.Unwrap(err) }

// UnwrapAll returns the complete chain of errors, starting with the supplied
// error and ending with the root cause error.
func UnwrapAll(err error) []error {
	res := make([]error, 0, 6)
	for {
		if err == nil {
			break
		}
		res = append(res, err)
		err = Unwrap(err)
	}
	return res
}

// Cause walks through all wrapped errors and returns the last error in the
// chain.
func Cause(err error) error {
	for {
		unwrapped := Unwrap(err)
		if unwrapped == nil {
			break
		}

		err = unwrapped
	}
	return err
}
