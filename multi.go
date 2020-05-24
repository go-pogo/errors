package errs

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/xerrors"
)

func Filter(errors []error) []error {
	res := make([]error, 0, len(errors))
	for _, err := range errors {
		if err != nil {
			res = append(res, err)
		}
	}
	return res
}

func Multi(errors ...error) error {
	if len(errors) == 1 {
		return TraceSkip(errors[0], 1)
	}

	list := Filter(errors)
	if len(list) == 1 {
		return TraceSkip(list[0], 1)
	}

	return &multiErr{
		errors: errors,
		frames: CaptureFrames(1, 2),
	}
}

// Append
//
// Important note: when using Append with defer, the pointer to the `dest` error
// must be a named return variable. For addition details see
// https://golang.org/ref/spec#Defer_statements.
func Append(dest *error, err error) error {
	if dest == nil {
		panic("errs.Append: dest must not be a nil pointer")
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

func (m *multiErr) As(target interface{}) bool {
	for _, err := range m.errors {
		if errors.As(err, &target) {
			return true
		}
	}
	return false
}

func (m *multiErr) Is(target error) bool {
	for _, err := range m.errors {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func (m *multiErr) Frames() *Frames { return &m.frames }

func (m *multiErr) Errors() []error { return m.errors }

func (m *multiErr) Format(s fmt.State, v rune) { xerrors.FormatError(m, s, v) }

func (m *multiErr) FormatError(p xerrors.Printer) error {
	p.Print(m.Error())
	if p.Detail() {
		for _, err := range m.errors {
			p.Printf("%+v", err)
		}
	}

	return nil
}

func (m *multiErr) Error() string {
	var buf strings.Builder
	_, _ = fmt.Fprint(&buf, "multiple errors occurred:\n")
	for _, e := range m.errors {
		_, _ = fmt.Fprintf(&buf, "* %v\n", e)
	}
	return buf.String()
}
