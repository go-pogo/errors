package errs

import (
	"errors"
	"strings"
)

// Print returns the complete error stack as a human-readable formatted string.
func Print(err error) string {
	errorSb := &strings.Builder{}
	traceSb := &strings.Builder{}
	traceSb.WriteString("\n\nTrace:\n")

	for {
		stErr, ok := err.(ErrorWithStackTrace)
		if !ok {
			errorSb.WriteString(err.Error())
			break
		}

		for _, frame := range stErr.StackTrace().frames {
			if frame.IsEmpty() {
				continue
			}

			traceSb.WriteString(frame.String() + ":\n")
		}

		msgErr, ok := err.(ErrorWithMessage)
		if !ok {
			traceSb.WriteString(err.Error())
			break
		}

		if kindErr, ok := err.(ErrorWithKind); ok {
			kind := kindErr.Kind()
			if kind != UnknownKind {
				errorSb.WriteString(kind.String() + ": ")
			}
		}

		msg := msgErr.Message()
		errorSb.WriteString(msg)
		traceSb.WriteString("> " + msg + "\n")

		if wrapErr, ok := err.(wrapErr); ok {
			err = errors.Unwrap(wrapErr.Unwrap())
		} else {
			err = errors.Unwrap(err)
		}

		if err == nil {
			break
		}

		errorSb.WriteString(",\n")
		traceSb.WriteRune('\n')
	}

	return errorSb.String() + traceSb.String()
}
