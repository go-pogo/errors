package errs

import (
	"errors"
	"fmt"

	"golang.org/x/xerrors"
)

// FormatError prints the error using `xerrors.FormatError()` and a formatter
// that implements the `xerrors.Formatter` interface. See the
// `golang.org/x/xerrors` package for additional information.
func FormatError(err error, s fmt.State, v rune) {
	xerrors.FormatError(formatter{err}, s, v)
}

type formatter struct{ error }

func (f formatter) FormatError(p xerrors.Printer) error {
	p.Print(f.error.Error())
	if p.Detail() {
		frames := GetFrames(f.error)
		if frames != nil {
			frames.Format(p)
		}
	}

	return errors.Unwrap(f.error)
}
