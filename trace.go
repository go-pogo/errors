package errs

import (
	"fmt"

	"golang.org/x/xerrors"
)

func CaptureFrames(capture int, skip int) Frames {
	fr := make(Frames, 0, capture)
	for i := 0; i < capture; i++ {
		fr.Capture(skip + i)
	}
	return fr
}

type Frames []xerrors.Frame

func (fr *Frames) Capture(skip int) {
	*fr = append(*fr, xerrors.Caller(skip+1))
}

func (fr Frames) Format(p xerrors.Printer) {
	for _, frame := range fr {
		frame.Format(p)
	}
}

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

	return &traceErr{
		frames: CaptureFrames(1, 1),
		err:    err,
	}
}

// traceErr is a wrapper for "primitive" errors that do not have stack trace
// information. It does not contain an error message by itself and always
// displays the message of the underlying wrapped error.
type traceErr struct {
	frames Frames
	err    error
}

func (err traceErr) Frames() *Frames { return &err.frames }

// Unwrap returns the original underlying "primitive" error that resides inside.
func (err traceErr) Unwrap() error { return err.err }

func (err traceErr) Error() string { return fmt.Sprint(err.err) }
