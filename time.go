// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"
	"time"
)

// Timer interfaces provide access to a [time.Time] indicating when the error
// occurred.
type Timer interface {
	error
	Time() time.Time
}

type TimerSetter interface {
	Timer
	SetTime(time.Time)
}

// WithTime adds time information to the error. It does so by wrapping the
// error with a [Timer], or update/set the time when error implements
// [TimerSetter]. It will return nil when the provided error is nil.
func WithTime(err error, when time.Time) Timer {
	if err == nil {
		return nil
	}
	//goland:noinspection GoTypeAssertionOnErrors
	if e, ok := err.(TimerSetter); ok {
		e.SetTime(when)
		return e
	}

	return &dateTimeError{
		embedError: &embedError{error: err},
		time:       when,
	}
}

// GetTime returns the [time.Time] of the last found [Timer] in err's error
// chain. If none is found, it returns the provided value or.
func GetTime(err error) (time.Time, bool) {
	var dt time.Time
	var has bool

	for {
		//goland:noinspection GoTypeAssertionOnErrors
		if e, ok := err.(Timer); ok {
			dt = e.Time()
			has = true
		}
		err = Unwrap(err)
		if err == nil {
			break
		}
	}

	return dt, has
}

type dateTimeError struct {
	*embedError
	time time.Time
}

func (e *dateTimeError) SetTime(t time.Time) { e.time = t }
func (e *dateTimeError) Time() time.Time     { return e.time }

// GoString prints the error in basic Go syntax.
func (e *dateTimeError) GoString() string {
	return fmt.Sprintf(
		"errors.dateTimeError{time: %s, embedErr: %#v}",
		e.time.String(),
		e.error,
	)
}
