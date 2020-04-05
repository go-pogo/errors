package errs

import (
	"errors"
	"fmt"
)

// Wrap creates a new error with an error cause.
func Wrap(cause error, kind Kind, msg string) error {
	if cause == nil {
		return nil
	}

	return &err{Inner: MakeInner(cause, kind, msg)}
}

// Wrapf creates a new error with an error cause and formatted message.
func Wrapf(cause error, kind Kind, format string, a ...interface{}) error {
	if cause == nil {
		return nil
	}

	return &err{Inner: MakeInner(cause, kind, fmt.Sprintf(format, a...))}
}

// WrapPanic wraps a panicking sequence with the given prefix. It then panics
// again.
func WrapPanic(prefix string) {
	if r := recover(); r != nil {
		panic(fmt.Sprintf("%s: %s", prefix, r))
	}
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
