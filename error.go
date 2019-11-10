package errs

import (
	"errors"
	"fmt"
	"strings"
)

// Unknown indicates an error is not created with a distinct kind.
const Unknown Kind = ""

// Kind describes the kind/type of error that has occurred, such as "auth error", "unmarshal error", etc.
type Kind string

func (k Kind) String() string {
	return string(k)
}

// ErrorWithKind interfaces provide access to a Kind.
type ErrorWithKind interface {
	Kind() Kind
}

// ErrorWithMessage interfaces provide access to the error message without underlying error messages.
type ErrorWithMessage interface {
	Message() string
}

// ErrorWithStackTrace interfaces provide access to a stack trace.
type ErrorWithStackTrace interface {
	StackTrace() *ST
}

// ErrorWithUnwrap interfaces provide access to an Unwrap method which may return an underlying error.
type ErrorWithUnwrap interface {
	Unwrap() error
}

// Error is the error with kind and stack trace information.
type Error struct {
	st   *ST   // stack trace of functions that returned the error
	err  error // cause of this error, if any
	kind Kind
	msg  string // message of error that occurred
}

func (err Error) StackTrace() *ST { return err.st }
func (err Error) Unwrap() error   { return err.err }
func (err Error) Kind() Kind      { return err.kind }
func (err Error) Message() string { return err.msg }
func (err Error) Error() string   { return Print(err) }

// WrapError is a wrapper for "primitive" errors that do not have stack trace information. It does
// not contain an error message by itself and always displays the message of the underlying
// wrapped error.
type WrapError struct {
	st  *ST   // stack trace of functions that returned the error
	err error // "primitive" error which contains the real error message
}

func (err WrapError) StackTrace() *ST { return err.st }
func (err WrapError) Unwrap() error   { return err.err }
func (err WrapError) Message() string { return err.err.Error() }
func (err WrapError) Error() string   { return Print(err) }

func (err WrapError) Kind() Kind {
	if err.err != nil {
		if kErr, ok := err.err.(ErrorWithKind); ok {
			return kErr.Kind()
		}
	}

	return Unknown
}

// prepare error struct, set field values etc.
func prepError(err *Error, cause error, msg string) *Error {
	err.st = NewStackTrace()
	err.st.Capture(2)
	err.err = cause
	err.kind = Unknown
	err.msg = msg

	return err
}

func prepWrapError(err *WrapError, cause error) *WrapError {
	err.st = NewStackTrace()
	err.st.Capture(2)
	err.err = cause

	return err
}

func New(args ...interface{}) {

}

// Err creates an error from a message.
func Err(msg string) *Error {
	var err Error
	return prepError(&err, nil, msg)
}

// Errf creates an error according to a format specifier.
func Errf(msg string, args ...interface{}) *Error {
	var err Error
	return prepError(&err, nil, fmt.Sprintf(msg, args...))
}

// Wrap wraps an existing error with information about the stack frame its called with. Errors that
// implement the ErrorWithStackTrace interface add the frame to the existing stack trace. Other
// "simple" errors are wrapped with a WrapError struct.
func Wrap(cause error, skip ...uint) error {
	if cause == nil {
		return nil
	}

	if wrErr, ok := cause.(ErrorWithStackTrace); ok {
		wrErr.StackTrace().Capture(1)
		return cause
	}

	var err WrapError
	return prepWrapError(&err, cause)
}

// Wrap wraps an existing error with a new error containing the provided message.
func Wrapf(cause error, msg string, args ...interface{}) *Error {
	if cause == nil {
		return nil
	}

	var err Error
	return prepError(&err, cause, fmt.Sprintf(msg, args...))
}

// UnwrapAll returns the complete stack of errors starting with the supplied error.
func UnwrapAll(err error) []error {
	stack := make([]error, 0, 0)

	for {
		if err == nil {
			break
		}
		stack = append(stack, err)
		err = errors.Unwrap(err)
	}

	return stack
}

// UnwrapCause unwraps all errors and returns the first error that started it all.
func UnwrapCause(err error) error {
	for {
		wErr := errors.Unwrap(err)
		if wErr == nil {
			break
		} else {
			err = wErr
		}
	}
	return err
}

// Print returns the complete error stack as a readable formatted string.
func Print(err error) string {
	errorSb := &strings.Builder{}
	traceSb := &strings.Builder{}
	traceSb.WriteString("\n\nTrace:\n")

	for {
		stErr, ok := err.(ErrorWithStackTrace)
		if !ok {
			errorSb.WriteString(err.Error())
			break
		}

		for _, frame := range stErr.StackTrace().frames {
			if frame.IsEmpty() {
				continue
			}

			traceSb.WriteString(frame.String() + ":\n")
		}

		msgErr, ok := err.(ErrorWithMessage)
		if !ok {
			traceSb.WriteString(err.Error())
			break
		}

		msg := msgErr.Message()
		errorSb.WriteString(msg)
		traceSb.WriteString(">\t" + msg + "\n")

		err = errors.Unwrap(err)
		if err == nil {
			break
		}

		errorSb.WriteString(",\n")
		traceSb.WriteRune('\n')
	}

	return errorSb.String() + traceSb.String()
}
