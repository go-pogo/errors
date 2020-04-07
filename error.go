package errs

import (
	"fmt"
	"strings"
)

const (
	// UnknownKind is used for errors that are created without a distinct `Kind`.
	UnknownKind Kind = ""
	// UnknownError is an error message that is returned when an error has no
	// message and is of `UnknownKind`
	UnknownError string = "unknown error"
)

// New creates a new error.
func New(kind Kind, msg string) error {
	err := &err{MakeInner(nil, kind, msg)}
	err.frames.Capture(1)
	return err
}

// Newf formats an error message according to a format specifier and provided
// arguments and creates a new error the same way `New()` does.
func Newf(kind Kind, format string, a ...interface{}) error {
	err := &err{MakeInner(nil, kind, fmt.Sprintf(format, a...))}
	err.frames.Capture(1)
	return err
}

// Wrap creates a new error that wraps around the causing error, thus extending
// the error chain. In contrast to `New()`, it will only create a new error
// when the cause error is not `nil`.
func Wrap(cause error, kind Kind, msg string) error {
	if cause == nil {
		return nil
	}

	err := &err{MakeInner(cause, kind, msg)}
	err.frames.Capture(1)
	return err
}

// Wrapf formats an error message according to a format specifier and provided
// arguments and creates a new error the same way `Wrap()` does.
func Wrapf(cause error, kind Kind, format string, a ...interface{}) error {
	if cause == nil {
		return nil
	}

	err := &err{MakeInner(cause, kind, fmt.Sprintf(format, a...))}
	err.frames.Capture(1)
	return err
}

type err struct{ Inner }

func (e err) Format(s fmt.State, v rune) { FormatError(e, s, v) }

// Error returns the message of the error with its `Kind` as prefix. If `Kind`
// is of `UnknownKind` the prefix is omitted. If message is empty, the string
// value of the kind is returned. When both kind and message are empty
// "unknown error" will be returned.
func (e err) Error() string {
	if e.kind == "" && e.msg == "" {
		return UnknownError
	}
	if e.kind == "" {
		return e.msg
	}
	if e.msg == "" {
		return e.kind.String()
	}

	return e.kind.String() + ": " + e.msg
}

// Inner is by itself not an error and is designed to be embedded in (custom)
// errors.
type Inner struct {
	frames *Frames // slice of stack trace frames
	cause  error   // cause of this error, if any
	kind   Kind    // specific kind of error
	msg    string  // message of error that occurred
}

func MakeInner(cause error, kind Kind, msg ...string) Inner {
	return Inner{
		frames: new(Frames),
		cause:  cause,
		kind:   kind,
		msg:    strings.Join(msg, " "),
	}
}

// Frames returns a slice of captured `xerrors.Frame` types linked to this error.
func (e Inner) Frames() *Frames { return e.frames }

// Unwrap returns the next error in the error chain. It returns `nil` if there
// is no next error.
func (e Inner) Unwrap() error { return e.cause }

// Kind returns the `Kind` of the error. It returns `UnknownKind` when no `Kind`
// is set.
func (e Inner) Kind() Kind { return e.kind }

// Kind describes the kind/type of error that has occurred. For example "auth
// error", "unmarshal error", etc. This way errors can be of the same `Kind`
// but still contain different error messages or additional fields.
// It is recommended to define each `Kind` as a constant.
type Kind string

// String returns the string representation of `Kind`.
func (k Kind) String() string { return string(k) }
