// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"
)

// StatusCoder interfaces provide access to a (http) status code.
type StatusCoder interface {
	error
	StatusCode() int
}

type StatusCoderSetter interface {
	StatusCoder
	SetStatusCode(int)
}

// WithStatusCode adds a (http) status code to the error which can be retrieved
// using GetStatusCode and may be set to a http.ResponseWriter.
func WithStatusCode(err error, statusCode int) StatusCoder {
	if err == nil {
		return nil
	}

	if e, ok := err.(StatusCoderSetter); ok {
		e.SetStatusCode(statusCode)
		return e
	}

	return &statusCodeError{
		embedError: &embedError{error: err},
		statusCode: statusCode,
	}
}

// GetStatusCode returns the status code if the error implements the StatusCoder
// interface. If not, it returns 0.
func GetStatusCode(err error) int { return GetStatusCodeOr(err, 0) }

// GetStatusCodeOr returns the status code from the first found StatusCoder
// in err's error chain. If none is found, it returns the provided value or.
func GetStatusCodeOr(err error, or int) int {
	for {
		if e, ok := err.(StatusCoder); ok {
			return e.StatusCode()
		}
		err = Unwrap(err)
		if err == nil {
			break
		}
	}

	return or
}

type statusCodeError struct {
	*embedError
	statusCode int
}

func (e *statusCodeError) SetStatusCode(c int) { e.statusCode = c }
func (e *statusCodeError) StatusCode() int     { return e.statusCode }

// GoString prints the error in basic Go syntax.
func (e *statusCodeError) GoString() string {
	return fmt.Sprintf(
		"*statusCodeError{statusCode: %d, embedErr: %#v}",
		e.statusCode,
		e.error,
	)
}
