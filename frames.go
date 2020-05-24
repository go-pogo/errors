package errs

import (
	"golang.org/x/xerrors"
)

const defaultFramesCapacity = 5

type Frames []xerrors.Frame

func CaptureFrames(n uint, skip uint) Frames {
	fr := make(Frames, 0, n+defaultFramesCapacity)

	var i uint
	for ; i < n; i++ {
		fr.Capture(skip + i)
	}
	return fr
}

func (fr Frames) Len() uint { return uint(len(fr)) }

func (fr *Frames) Capture(skip uint) {
	*fr = append(*fr, xerrors.Caller(int(skip)+1))
}

func (fr Frames) Format(p xerrors.Printer) {
	for i := len(fr) - 1; i >= 0; i-- {
		fr[i].Format(p)
	}
}

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
