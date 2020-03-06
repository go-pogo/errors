package errs

import (
	"errors"
	"fmt"
	"strings"
)

// Print returns the complete error stack as a human-readable formatted string.
func Print(err error) string {
	p := printer{
		esb:  new(strings.Builder),
		tsb:  new(strings.Builder),
		eSep: "\n",
		tSep: ",\n",
	}

	p.tsb.WriteString("\n\nTrace:\n")

	for {
		stErr, ok := err.(ErrorWithStackTrace)
		if !ok {
			p.WritePlainError(err)
			break
		}

		p.WriteErrorWithStackTrace(stErr)

		if wrapErr, ok := err.(wrapErr); ok {
			err = errors.Unwrap(wrapErr.Unwrap())
		} else {
			err = errors.Unwrap(err)
		}

		if err == nil {
			break
		}
	}

	return p.String()
}

type printer struct {
	esb  *strings.Builder // error messages
	tsb  *strings.Builder // stack trace
	eSep string           // separator for error messages
	tSep string           // separator for stack traces
}

func (p printer) WritePlainError(err error) {
	p.esb.WriteString(err.Error())
}

func (p printer) WriteErrorWithStackTrace(err ErrorWithStackTrace) {
	msg := GetKindMessage(err)
	p.esb.WriteString(msg + p.eSep)

	st := err.StackTrace()
	if st == nil {
		return
	}

	for _, frame := range st.frames {
		if frame.IsEmpty() {
			continue
		}

		p.WriteFrame(frame)
	}

	fmt.Fprintf(p.tsb, "> %s\n%s", msg, p.tSep)
}

func (p printer) WriteFrame(f Frame) {
	fmt.Fprintf(p.tsb, "%s:%d: %s()\n", f.Path, f.Line, f.Func)
}

func (p printer) String() string {
	return strings.TrimSuffix(p.esb.String(), p.eSep) +
		strings.TrimSuffix(p.tsb.String(), p.tSep)
}
