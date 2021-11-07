package errors

import (
	stderrors "errors"

	"golang.org/x/xerrors"
)

// Unwrap is an alias of errors.Unwrap. It returns the result of calling the
// Unwrap method on err, if err's type contains an Unwrap method returning
// error. Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
	return stderrors.Unwrap(err)
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
