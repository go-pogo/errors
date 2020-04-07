package errs

import (
	"errors"
)

// ErrorWithKind interfaces provide access to a `Kind`.
type ErrorWithKind interface {
	error
	Kind() Kind
}

// ErrorWithFrames interfaces provide access to a stack of frames.
type ErrorWithFrames interface {
	error
	Frames() *Frames
}

// ErrorWithUnwrap interfaces provide access to an underlying error further down
// the error chain, if any.
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

func GetFrames(err error) *Frames {
	if e, ok := err.(ErrorWithFrames); ok {
		return e.Frames()
	}
	return nil
}

// UnwrapAll returns the complete chain of errors, starting with the supplied
// error and ending with the error that started the chain.
func UnwrapAll(err error) []error {
	var res []error
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
		}

		err = wErr
	}
	return err
}
