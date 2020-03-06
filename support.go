package errs

import (
	"errors"
)

// ErrorWithKind interfaces provide access to a Kind.
type ErrorWithKind interface {
	error
	Kind() Kind
}

// ErrorWithMessage interfaces provide access to the plain error message without
// returning any of the underlying error messages.
type ErrorWithMessage interface {
	error
	Message() string
}

// ErrorWithStackTrace interfaces provide access to a stack trace.
type ErrorWithStackTrace interface {
	error
	StackTrace() *ST
}

// ErrorWithUnwrap interfaces provide access to an `Unwrap` method which may
// return an underlying error.
type ErrorWithUnwrap interface {
	error
	Unwrap() error
}

// GetKind returns the `Kind` of the error if it implements the `ErrorWithKind`
// interface. If not, it returns `UnknownKind`.
func GetKind(err error) Kind {
	if e, ok := err.(ErrorWithKind); ok {
		return e.Kind()
	}

	return UnknownKind
}

// GetMessage returns the message string of the error if if implements the
// `ErrorWithMessage` interface. If not, it returns an empty string.
func GetMessage(err error) string {
	if e, ok := err.(ErrorWithMessage); ok {
		return e.Message()
	}

	return ""
}

// GetKindMessage returns the message string of the error with its `Kind` as
// prefix. If `Kind` is of `UnknownKind` the prefix is omitted. If message is
// empty, the string value of the `Kind` is returned.
// An empty string is returned when both `Kind` and message are empty.
func GetKindMessage(err error) string {
	kind := GetKind(err)
	if kind == UnknownKind {
		return GetMessage(err)
	}

	msg := GetMessage(err)
	if msg == "" {
		return kind.String()
	}

	return kind.String() + ": " + msg
}

// GetStackTrace returns a pointer to the stack trace of the error if it
// implements the `ErrorWithStackTrace` interface.
func GetStackTrace(err error) *ST {
	if e, ok := err.(ErrorWithStackTrace); ok {
		return e.StackTrace()
	}

	return nil
}

// UnwrapAll returns the complete stack of errors, starting with the supplied
// error and ending with the cause error.
func UnwrapAll(err error) []error {
	res := make([]error, 0, 0)

	for {
		if err == nil {
			break
		}
		res = append(res, err)
		err = errors.Unwrap(err)
	}

	return res
}

// func UnwrapDepth(err error, depth int) (error, bool) {
// 	if depth <= 0 {
// 		panic("UnwrapDepth: unwrapping with a depth lower than 1 is not possible")
// 	}
//
// 	ok := true
// 	for i := 0; i < depth; i++ {
// 		if err == nil {
// 			ok = false
// 			break
// 		}
// 		err = errors.Unwrap(err)
// 	}
//
// 	return err, ok
// }

// UnwrapCause walks through all wrapped errors and returns the first "cause"
// error.
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
