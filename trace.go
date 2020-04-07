package errs

import (
	"errors"
	"fmt"
)

// Trace wraps an existing error with information about the stack frame its
// called with. Errors that implement the `ErrorWithStackTrace` interface add
// the frame to the existing stack trace. Other "simple" errors are wrapped
// in a `traceErr` type.
func Trace(err error) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(ErrorWithFrames); ok {
		e.Frames().Capture(1)
		return err
	}

	frames := CaptureFrames(1, 2)
	return &traceErr{
		error:  err,
		frames: &frames,
	}
}

// traceErr is a wrapper for "primitive" errors that do not have stack trace
// information. It does not contain an error message by itself and always
// displays the message of the underlying wrapped error.
type traceErr struct {
	error
	frames *Frames
}

func (t traceErr) Frames() *Frames { return t.frames }

func (t traceErr) Unwrap() error { return errors.Unwrap(t.error) }

func (t traceErr) Format(s fmt.State, v rune) { FormatError(t, s, v) }
