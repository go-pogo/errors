package errs

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"
)

// Filter filters the provided errors and only returns the non-nil errors.
func Filter(errors ...error) []error {
	l := len(errors)
	if l == 0 {
		return errors
	}

	res := make([]error, 0, l)
	for _, err := range errors {
		if err != nil {
			res = append(res, err)
		}
	}
	return res
}

// Combine returns a multi error when there are more than one non-nil errors
// provided. If only one non-nil error is provided, it will act as if
// `TraceSkip()` is called. It returns nil when all provided errors are nil.
func Combine(errors ...error) error {
	errors = Filter(errors...)
	switch len(errors) {
	case 0:
		return nil
	case 1:
		return TraceSkip(errors[0], 1)
	default:
		return &multiErr{
			errors: errors,
			frames: CaptureFrames(1, 2),
		}
	}
}

const panicAppendNilPtr = "errs.Append: dest must not be a nil pointer"

// Append appends multiple non-nil errors to a single multi error `dest`.
//
// Important note: when using Append with defer, the pointer to the `dest` error
// must be a named return variable. For addition details see
// https://golang.org/ref/spec#Defer_statements.
func Append(dest *error, err error) error {
	if dest == nil {
		panic(panicAppendNilPtr)
	}
	if err == nil {
		return *dest
	}

	switch d := (*dest).(type) {
	case nil:
		*dest = err

	case *multiErr:
		d.errors = append(d.errors, err)

	default:
		*dest = &multiErr{
			errors: []error{*dest, err},
			frames: CaptureFrames(1, 2),
		}
	}

	return *dest
}

type multiErr struct {
	errors []error
	frames Frames
}

// Frames returns a slice of captured `xerrors.Frame` types linked to this multi
// error.
func (m *multiErr) Frames() *Frames { return &m.frames }

// Errors returns the errors within the multi error.
func (m *multiErr) Errors() []error { return m.errors }

// Format prints the errors using `xerrors.FormatError()`.
// See the `golang.org/x/xerrors` package for additional information.
func (m *multiErr) Format(s fmt.State, v rune) { xerrors.FormatError(m, s, v) }

// FormatError prints the error using `xerrors.FormatError()` and a formatter
// that implements the `xerrors.Formatter` interface.
// See the `golang.org/x/xerrors` package for additional information.
func (m *multiErr) FormatError(p xerrors.Printer) error {
	p.Print(m.Error())
	if p.Detail() {
		m.frames.Format(p)

		l := len(m.errors)
		for i, err := range m.errors {
			p.Printf("\n[%d/%d] %+v\n", i+1, l, err)
		}
	}

	return nil
}

func (m *multiErr) Error() string {
	var buf strings.Builder
	buf.WriteString("multiple errors occurred:")

	l := len(m.errors)
	for i, e := range m.errors {
		fmt.Fprintf(&buf, "\n[%d/%d] %s", i+1, l, e.Error())
		if i < l-1 {
			buf.WriteRune(';')
		}
	}
	return buf.String()
}
