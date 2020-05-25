package errs

import (
	"errors"
	"fmt"

	"golang.org/x/xerrors"
)

// FormatError prints the error using `xerrors.FormatError()` and a formatter
// that implements the `xerrors.Formatter` interface.
// See the `golang.org/x/xerrors` package for additional information.
func FormatError(err error, s fmt.State, v rune) {
	f, ok := err.(xerrors.Formatter)
	if !ok {
		f = errorFormatter{err}
	}

	xerrors.FormatError(f, s, v)
}

type errorFormatter struct{ error }

func (f errorFormatter) FormatError(p xerrors.Printer) error {
	p.Print(f.error.Error())
	if p.Detail() {
		frames := GetFrames(f.error)
		if frames != nil {
			frames.Format(p)
		}
	}

	return errors.Unwrap(f.error)
}
