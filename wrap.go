// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"fmt"

	"golang.org/x/xerrors"
)

// Wrap creates a new error that wraps around the causing error, thus extending
// the error chain. It will only create a new error when the provided cause
// error is not nil, otherwise it will return nil.
func Wrap(cause error, text string) error {
	if cause == nil {
		return nil
	}

	ce := toCommonErr(stderrors.New(text), false)
	ce.cause = cause
	ce.Trace(1)
	return ce
}

// Wrapf formats an error message according to a format specifier and provided
// arguments and creates a new error the same way Wrap() does.
func Wrapf(cause error, format string, a ...interface{}) error {
	if cause == nil {
		return nil
	}

	ce := toCommonErr(fmt.Errorf(format, a...), false)
	ce.cause = cause
	ce.Trace(1)
	return ce
}

// An Unwrapper unpacks a wrapped error.
type Unwrapper interface {
	error
	// Unwrap returns the next error in the error chain.
	// If there is no next error, Unwrap returns nil.
	Unwrap() (next error)
}

// Unwrap is an alias of errors.Unwrap. It returns the result of calling the
// Unwrap method on err, if err's type contains an Unwrap method returning
// error. Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
	return stderrors.Unwrap(err)
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

// Opaque is an alias of xerrors.Opaque. It returns an error with the same error
// formatting as err but that does not match err and cannot be unwrapped.
func Opaque(err error) error { return xerrors.Opaque(err) }

// Is is an alias of errors.Is. It reports whether any error in err's chain
// matches target.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
func Is(err, target error) bool {
	return stderrors.Is(Original(err), Original(target))
}

// As is an alias of errors.As. It finds the first error in err's chain that
// matches target, and if so, sets target to that error value and returns true.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error matches target if the error's concrete value is assignable to the value
// pointed to by target, or if the error has a method As(interface{}) bool such that
// As(target) returns true. In the latter case, the As method is responsible for
// setting target.
//
// As will panic if target is not a non-nil pointer to either a type that implements
// error, or to any interface type. As returns false if err is nil.
func As(err error, target interface{}) bool {
	//goland:noinspection GoErrorsAs
	return err != nil && stderrors.As(err, target)
}
