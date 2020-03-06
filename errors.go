package errs

// UnknownKind is used for errors that are created without a distinct `Kind`.
const UnknownKind Kind = ""

// Kind describes the kind/type of error that has occurred. For example "auth
// error", "unmarshal error", etc. This way errors can be of the same `Kind` but
// still contain different error messages.
// It is recommended to define each `Kind` as a constant.
type Kind string

// String returns the string representation of `Kind`.
func (k Kind) String() string { return string(k) }

// New creates an error from the provided arguments.
func New(args ...interface{}) error {
	if len(args) < 1 {
		panic("errs.New: no arguments provided")
	}

	defer WrapPanic("errs.New")
	return &err{Inner: MakeInnerWith(args...)}
}

// Wrap wraps an existing error with information about the stack frame its
// called with. Errors that implement the `ErrorWithStackTrace` interface add
// the frame to the existing stack trace. Other "simple" errors are wrapped
// in a `wrapErr` type.
func Wrap(cause error) error {
	if cause == nil {
		return nil
	}

	if err, ok := cause.(ErrorWithStackTrace); ok {
		err.StackTrace().Capture(1)
		return cause
	}

	return &wrapErr{
		st:  NewStackTraceCapture(1),
		err: cause,
	}
}

// err is a general error. The type `Inner` is embedded in it.
type err struct{ Inner }

// Error returns the human-readable error report using the `Print` function.
func (err err) Error() string { return Print(err) }

// wrapErr is a wrapper for "primitive" errors that do not have stack trace
// information. It does not contain an error message by itself and always
// displays the message of the underlying wrapped error.
type wrapErr struct {
	st  *ST
	err error
}

func (err wrapErr) StackTrace() *ST { return err.st }

// Unwrap returns the original underlying "primitive" error that resides inside.
func (err wrapErr) Unwrap() error { return err.err }

// Message returns the original error of the underlying error.
func (err wrapErr) Message() string { return err.err.Error() }

// Error returns the human-readable error report using the `Print` function.
func (err wrapErr) Error() string { return Print(err) }
