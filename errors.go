package errs

import (
	"fmt"

	"golang.org/x/xerrors"
)

// New creates a new error.
func New(kind Kind, msg string) error {
	return &err{
		Inner: MakeInner(nil, kind, msg),
	}
}

// Newf creates a new error with formatted message.
func Newf(kind Kind, format string, a ...interface{}) error {
	return &err{
		Inner: MakeInner(nil, kind, fmt.Sprintf(format, a...)),
	}
}

// err is a general error with type `Inner` embedded in it.
type err struct{ Inner }

// Format calls `xerrors.FormatError()` which formats the error according to
// s and v. See `xerrors.FormatError()` for additional information.
func (err err) Format(s fmt.State, v rune) { xerrors.FormatError(err, s, v) }

// Error prints the error in a human readable form.
func (err err) Error() string { return fmt.Sprint(err) }

// UnknownKind is used for errors that are created without a distinct `Kind`.
const UnknownKind Kind = ""

// Kind describes the kind/type of error that has occurred. For example "auth
// error", "unmarshal error", etc. This way errors can be of the same `Kind`
// but still contain different error messages or additional fields.
// It is recommended to define each `Kind` as a constant.
type Kind string

// String returns the string representation of `Kind`.
func (k Kind) String() string { return string(k) }

// ErrorWithKind interfaces provide access to a `Kind`.
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

// ErrorWithFrames interfaces provide access to a stack of frames.
type ErrorWithFrames interface {
	error
	Frames() *Frames
}

// ErrorWithUnwrap interfaces provide access to an `Unwrap` method which may
// return an underlying error.
type ErrorWithUnwrap interface {
	error
	Unwrap() error
}
