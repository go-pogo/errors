package errs

import (
	"errors"
	"strings"
)

// ErrorWithKind interfaces provide access to a Kind.
type ErrorWithKind interface {
	Kind() Kind
}

// ErrorWithMessage interfaces provide access to the error message without underlying error messages.
type ErrorWithMessage interface {
	Message() string
}

// ErrorWithStackTrace interfaces provide access to a stack trace.
type ErrorWithStackTrace interface {
	StackTrace() *ST
}

// ErrorWithUnwrap interfaces provide access to an Unwrap method which may return an underlying error.
type ErrorWithUnwrap interface {
	Unwrap() error
}

// GetKind returns a pointer to the Kind of the error if it implements the ErrorWithKind interface.
func GetKind(err error) *Kind {
	if e, ok := err.(ErrorWithKind); ok {
		k := e.Kind()
		return &k
	}

	return nil
}

// GetMessage returns the message string of the error if if implements the ErrorWithMessage
// interface. If not, it returns an empty string.
func GetMessage(err error) string {
	if e, ok := err.(ErrorWithMessage); ok {
		return e.Message()
	}

	return ""
}

// GetStackTrace returns a pointer to the stack trace of the error if it implements the
// ErrorWithStackTrace interface.
func GetStackTrace(err error) *ST {
	if e, ok := err.(ErrorWithStackTrace); ok {
		return e.StackTrace()
	}

	return nil
}

// UnwrapAll returns the complete stack of errors starting with the supplied error.
func UnwrapAll(err error) []error {
	res := make([]error, 0, 0)

	for {
		if err == nil {
			break
		}
		res = append(res, err)
		err = errors.Unwrap(err)
	}

	return res
}

// UnwrapCause walks through all wrapped errors and returns the first "cause" error.
func UnwrapCause(err error) error {
	for {
		wErr := errors.Unwrap(err)
		if wErr == nil {
			break
		} else {
			err = wErr
		}
	}
	return err
}

// Print returns the complete error stack as a readable formatted string.
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

		msg := msgErr.Message()
		errorSb.WriteString(msg)
		traceSb.WriteString(">\t" + msg + "\n")

		err = errors.Unwrap(err)
		if err == nil {
			break
		}

		errorSb.WriteString(",\n")
		traceSb.WriteRune('\n')
	}

	return errorSb.String() + traceSb.String()
}
