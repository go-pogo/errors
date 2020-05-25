package errs

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"
)

const (
	defaultFramesCapacity = 5
	invalidFrameSuffix    = " 0]}"
)

type Frames []xerrors.Frame

func CaptureFrames(n uint, skip uint) Frames {
	fr := make(Frames, 0, n+defaultFramesCapacity)

	var i uint
	for ; i < n; i++ {
		if !fr.Capture(skip + i) {
			break
		}
	}
	return fr
}

func (fr Frames) Len() uint { return uint(len(fr)) }

func isValidFrame(f xerrors.Frame) bool {
	s, x := fmt.Sprintf("%+v", f), invalidFrameSuffix
	return s[len(s)-len(x):] != x
}

func (fr *Frames) Capture(skip uint) (ok bool) {
	f := xerrors.Caller(int(skip) + 1)
	if ok = isValidFrame(f); ok {
		*fr = append(*fr, f)
	}
	return ok
}

func (fr Frames) Format(p xerrors.Printer) {
	for i := len(fr) - 1; i >= 0; i-- {
		fr[i].Format(p)
	}
}

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
