package errs

import (
	"golang.org/x/xerrors"
)

func CaptureFrames(n int, skip int) Frames {
	fr := make(Frames, 0, n)
	for i := 0; i < n; i++ {
		fr.Capture(skip + i)
	}
	return fr
}

type Frames []xerrors.Frame

func (fr *Frames) Capture(skip int) {
	*fr = append(*fr, xerrors.Caller(skip+1))
}

func (fr Frames) Format(p xerrors.Printer) {
	for i := len(fr) - 1; i >= 0; i-- {
		fr[i].Format(p)
	}
}
