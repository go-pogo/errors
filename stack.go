// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"
)

// StackTracer interfaces provide access to a stack of traced StackTrace.
type StackTracer interface {
	error

	// StackTrace returns a stack of traces frames.
	StackTrace() *StackTrace
}

const panicUseNewInstead = "errors.WithStack: use errors.New instead to create an error from an errors.Msg"

// WithStack gets a stack trace at the point WithStack was called and adds it to
// the error. If err is nil, WithStack returns nil.
func WithStack(err error) StackTracer {
	if err == nil {
		return nil
	}

	switch err := err.(type) {
	case StackTracer:
		return err

	case Msg, *Msg:
		panic(panicUseNewInstead)

	default:
		e := &embedError{error: err}
		if traceStack {
			e.stack = newStackTrace(1)
			if u := Unwrap(err); u != nil {
				skipStackTrace(u, e.stack.Len())
			}
		}
		return e
	}
}

// GetStackTrace returns a *StackTrace if err is a StackTracer or nil otherwise.
func GetStackTrace(err error) *StackTrace {
	if e, ok := err.(StackTracer); ok {
		return e.StackTrace()
	}
	return nil
}

type StackTrace struct {
	frames []xerrors.Frame

	// Skip n frames when formatting with Format, so overlapping frames from
	// previous errors are not printed.
	Skip uint
}

func newStackTrace(skipFrames uint) *StackTrace {
	st := &StackTrace{frames: make([]xerrors.Frame, 0, 6)}
	Callers(skipFrames+1, &st.frames)
	return st
}

func skipStackTrace(err error, skip uint) {
	if skip == 0 {
		return
	}

	st := GetStackTrace(err)
	if st == nil || st.Len() < skip {
		return
	}

	st.Skip = skip - 1
}

const panicCallersNilPtr = "errors.Callers: dest must be a pointer to a []xerrors.Frame"

// Callers fills the stack *StackTrace with xerrors.Frame's from the point Callers
// is called, skipping the first skipFrames frames.
func Callers(skipFrames uint, dest *[]xerrors.Frame) int {
	if dest == nil {
		panic(panicCallersNilPtr)
	}

	skip := int(skipFrames)

	var n int
	for {
		skip += 1
		f := xerrors.Caller(skip)
		if !isValidFrame(f) {
			break
		}

		*dest = append(*dest, f)
		n++
	}
	if n > 1 {
		n--
		*dest = (*dest)[:n]
	}

	return n
}

const (
	invalidFrameSuffix    = " 0]}"
	invalidFrameSuffixLen = 4
)

func isValidFrame(f xerrors.Frame) bool {
	s := fmt.Sprintf("%+v", f)
	return s[len(s)-invalidFrameSuffixLen:] != invalidFrameSuffix
}

func (st *StackTrace) Frames() []xerrors.Frame { return st.frames }

func (st *StackTrace) Len() uint {
	if nil == st {
		return 0
	}
	return uint(len(st.frames))
}

// Format formats the slice of xerrors.Frame using a xerrors.Printer.
func (st *StackTrace) Format(printer xerrors.Printer) {
	if !printer.Detail() {
		return
	}
	for i := len(st.frames) - 1 - int(st.Skip); i >= 0; i-- {
		st.frames[i].Format(printer)
	}
}

func (st *StackTrace) String() string {
	var p framesPrinter
	for i := len(st.frames) - 1; i >= 0; i-- {
		st.frames[i].Format(&p)
	}
	return p.b.String()
}

// framesPrinter is a xerrors.Printer that is used to print the string
// representation of StackTrace.
type framesPrinter struct{ b strings.Builder }

func (p *framesPrinter) Print(args ...interface{}) {
	_, _ = fmt.Fprint(&p.b, args...)
}

func (p *framesPrinter) Printf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(&p.b, format, args...)
}

func (p *framesPrinter) Detail() bool { return true }
