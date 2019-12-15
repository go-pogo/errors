package errs

import (
	"fmt"
)

// Unknown indicates an error is not created with a distinct kind.
const Unknown Kind = ""

// Kind describes the kind/type of error that has occurred, such as "auth error", "unmarshal error", etc.
type Kind string

func (k Kind) String() string {
	return string(k)
}

// err is a general error message with kind and stack trace information.
type err struct {
	st   *ST   // stack trace of functions that returned the error
	err  error // cause of this error, if any
	kind Kind
	msg  string // message of error that occurred
}

func (err err) StackTrace() *ST { return err.st }
func (err err) Unwrap() error   { return err.err }
func (err err) Kind() Kind      { return err.kind }
func (err err) Message() string { return err.msg }
func (err err) Error() string   { return Print(err) }

// wrapErr is a wrapper for "primitive" errors that do not have stack trace information. It does
// not contain an error message by itself and always displays the message of the underlying
// wrapped error.
type wrapErr struct {
	st  *ST   // stack trace of functions that returned the error
	err error // "primitive" error which contains the real error message
}

func (err wrapErr) StackTrace() *ST { return err.st }
func (err wrapErr) Unwrap() error   { return err.err }
func (err wrapErr) Message() string { return err.err.Error() }
func (err wrapErr) Error() string   { return Print(err) }

// Err creates an error from a message.
func Err(kind Kind, msg string) *err {
	var err err
	return prepError(&err, nil, kind, msg)
}

// Errf creates an error according to a format specifier.
func Errf(kind Kind, msg string, args ...interface{}) *err {
	var err err
	return prepError(&err, nil, kind, fmt.Sprintf(msg, args...))
}

// Wrap wraps an existing error with a new error containing the provided message.
func Wrapf(cause error, kind Kind, msg string, args ...interface{}) error {
	if cause == nil {
		return nil
	}

	var err err
	return prepError(&err, cause, kind, fmt.Sprintf(msg, args...))
}

// Wrap wraps an existing error with information about the stack frame its called with. Errors that
// implement the ErrorWithStackTrace interface add the frame to the existing stack trace. Other
// "simple" errors are wrapped with a WrapError struct.
func Wrap(cause error) error {
	if cause == nil {
		return nil
	}

	if err, ok := cause.(ErrorWithStackTrace); ok {
		err.StackTrace().Capture(1)
		return cause
	}

	var err wrapErr
	return prepWrapError(&err, cause)
}

func prepError(err *err, cause error, kind Kind, msg string) *err {
	err.st = NewStackTrace()
	err.st.Capture(2)

	err.err = cause
	err.kind = kind
	err.msg = msg

	return err
}

func prepWrapError(err *wrapErr, cause error) *wrapErr {
	err.st = NewStackTrace()
	err.st.Capture(2)
	err.err = cause

	return err
}
