package errs

import (
	"errors"
)

// ErrorWithUnwrap interfaces provide access to an underlying error further down
// the error chain, if any.
type ErrorWithUnwrap interface {
	error
	Unwrap() error
}

// UnwrapAll returns the complete chain of errors, starting with the supplied
// error and ending with the error that started the chain.
func UnwrapAll(err error) []error {
	var res []error
	for {
		if err == nil {
			break
		}
		if t, ok := err.(*traceErr); ok {
			// skip traceErrs, they only contain stack trace frames and not an
			// error message of its own
			err = t.error
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
