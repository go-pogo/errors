package errors

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"

	"github.com/go-pogo/errors/internal"
)

const invalidFrameSuffix = " 0]}"

var DefaultFramesCapacity uint = 5

type Frames []xerrors.Frame

// CaptureFrames captures `n` frames starting from `skip`.
func CaptureFrames(n uint, skip uint) Frames {
	fr := make(Frames, 0, n+DefaultFramesCapacity)

	var i uint
	for ; i < n; i++ {
		if !fr.Capture(skip + i) {
			break
		}
	}
	return fr
}

// Len returns the amount of captures frames as uint.
// Use `len()` when the value is needed as an int.
func (fr Frames) Len() uint { return uint(len(fr)) }

func isValidFrame(f xerrors.Frame) bool {
	s, x := fmt.Sprintf("%+v", f), invalidFrameSuffix
	return s[len(s)-len(x):] != x
}

// Capture captures a `xerrors.Frame` that describes a frame on the caller's
// stack. The argument skip is the number of frames to skip over.
// Capture(0) returns the frame for the caller of Capture.
// It returns a bool false when the captured frame contains a nil pointer.
func (fr *Frames) Capture(skip uint) (ok bool) {
	if !internal.CaptureFrames() {
		return true
	}

	f := xerrors.Caller(int(skip) + 1)
	if ok = isValidFrame(f); ok {
		*fr = append(*fr, f)
	}
	return ok
}

// Format formats the captured frames using xerror's format functionality.
func (fr Frames) Format(p xerrors.Printer) {
	for i := len(fr) - 1; i >= 0; i-- {
		fr[i].Format(p)
	}
}

// String formats the captured frames and returns its string representation.
func (fr Frames) String() string {
	var p framesPrinter
	fr.Format(&p)
	return p.b.String()
}

// framesPrinter uses xerrors format functionality to print a list of frames.
type framesPrinter struct{ b strings.Builder }

func (p *framesPrinter) Print(args ...interface{}) {
	_, _ = fmt.Fprint(&p.b, args...)
}

func (p *framesPrinter) Printf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(&p.b, format, args...)
}

func (p *framesPrinter) Detail() bool { return true }

// ErrorWithFrames interfaces provide access to a stack of frames.
type ErrorWithFrames interface {
	error
	Frames() *Frames
}

func GetFrames(err error) *Frames {
	if e, ok := err.(ErrorWithFrames); ok {
		return e.Frames()
	}
	return nil
}
