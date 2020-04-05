package errs

import (
	"fmt"
	"reflect"

	"golang.org/x/xerrors"
)

func MakeInner(cause error, kind Kind, msg string) Inner {
	return Inner{
		frames: CaptureFrames(1, 3),
		err:    cause,
		kind:   kind,
		msg:    msg,
	}
}

func MakeInnerWith(args ...interface{}) Inner {
	inner := Inner{
		frames: CaptureFrames(1, 3),
	}

	for _, arg := range args {
		switch arg := arg.(type) {
		case Kind:
			inner.kind = arg
		case *Kind:
			inner.kind = *arg
		case error:
			inner.err = arg
		case string:
			inner.msg = arg
		default:
			panic(fmt.Sprintf(
				"errs.MakeInnerWith: invalid argument of type `%s` provided",
				reflect.TypeOf(arg).String(),
			))
		}
	}

	return inner
}

// Inner is by itself not an error and is designed to be embedded in (custom)
// errors. This adds its methods to the (custom) error.
type Inner struct {
	frames Frames
	err    error  // cause of this error, if any
	kind   Kind   // specific kind of error
	msg    string // message of error that occurred
}

func (inr Inner) Frames() *Frames { return &inr.frames }

// Unwrap returns the next error in the error chain. It returns `nil` if there
// is no next error.
func (inr Inner) Unwrap() error { return inr.err }

// Kind returns the `Kind` of the error. It returns `UnknownKind` of no `Kind`
// is set.
func (inr Inner) Kind() Kind { return inr.kind }

// Message returns the raw error message, without stack trace or any underlying
// errors. It returns an empty string when no message is set.
func (inr Inner) Message() string { return inr.msg }

func (inr Inner) FormatError(p xerrors.Printer) (next error) {
	p.Print(inr.msg)
	inr.frames.Format(p)
	return inr.err
}
