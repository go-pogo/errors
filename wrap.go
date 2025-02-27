// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"fmt"
	"reflect"

	"golang.org/x/xerrors"
)

// A Wrapper provides context around another error, which can be retrieved with
// [Unwrap].
type Wrapper interface {
	error
	xerrors.Wrapper
}

// Wrap creates a new error, which implements the [StackTracer], [Wrapper] and
// [Formatter] interfaces, that wraps around the causing error. Argument msg
// can be either a string or [Msg].
//
//	err = errors.Wrap(err, "my error message")
//	err = errors.Wrap(err, errors.Msg("my error message"))
//
// Wrap records a stack trace at the point it was called. Each call returns a
// distinct error value even if cause and msg are identical.
// Wrap will return nil when cause is nil, and it will return the provided
// cause when msg is nil.
func Wrap(cause error, msg interface{}) error {
	if cause == nil || msg == nil {
		return cause
	}

	var parent error
	switch v := msg.(type) {

	case string:
		parent = Msg(v)
	case *string:
		parent = Msg(*v)

	case Msg:
		parent = v
	case *Msg:
		parent = *v

	default:
		panic(unsupportedType("errors.Wrap", reflect.TypeOf(v).String()))
	}

	return withCause(newCommonErr(parent, true, 1), cause)
}

// Wrapf formats an error message according to a format specifier and provided
// arguments with [fmt.Errorf], and creates a new error similar to [Wrap].
//
//	err = errors.Wrapf(err, "my error %s", "message")
func Wrapf(cause error, format string, args ...interface{}) error {
	if cause == nil {
		return nil
	}
	return withCause(newCommonErr(fmt.Errorf(format, args...), true, 1), cause)
}

// Opaque is an alias of [xerrors.Opaque]. It returns an error with the same
// error formatting as err but that does not match err and cannot be unwrapped.
func Opaque(err error) error { return xerrors.Opaque(err) }

// Unwrap is an alias of [errors.Unwrap]. It returns the result of calling the
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

// IsCause indicates if error is the root cause of a possible error chain.
func IsCause(err error) bool { return Unwrap(Unembed(err)) == nil }

// Is reports whether any error in err's chain matches target. It is fully
// compatible with both [errors.Is] and [xerrors.Is].
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
func Is(err, target error) bool { return stderrors.Is(err, target) }

// As is an alias of [errors.As]. It finds the first error in err's chain that
// matches target, and if so, sets target to that error value and returns true.
func As(err error, target interface{}) bool {
	//goland:noinspection GoErrorsAs
	return err != nil && stderrors.As(err, target)
}
