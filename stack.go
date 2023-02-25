// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"
	"runtime"
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

	switch v := err.(type) {
	case StackTracer:
		return v

	case Msg, *Msg:
		panic(panicUseNewInstead)

	default:
		ee := &embedError{error: v}
		if traceStack {
			ee.stack = newStackTrace(1)
			if u := Unwrap(v); u != nil {
				skipStackTrace(u, ee.stack.Len())
			}
		}
		return ee
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
	frames []uintptr

	// Skip n frames when formatting with Format, so overlapping frames from
	// previous errors are not printed.
	Skip uint
}

const framesCap = 16

func newStackTrace(skipFrames uint) *StackTrace {
	st := &StackTrace{frames: make([]uintptr, 0, framesCap)}
	callers(int(skipFrames)+1, &st.frames)
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
func Callers(skipFrames uint, dest *[]uintptr) int {
	if dest == nil {
		panic(panicCallersNilPtr)
	}
	return callers(int(skipFrames)+1, dest)
}

func callers(skip int, dest *[]uintptr) int {
	skip += 2

	var count int
	var pc [framesCap]uintptr
	for {
		n := runtime.Callers(skip+count, pc[:])
		if n == 0 {
			break
		}

		*dest = append(*dest, pc[:n]...)
		count += n
		if n < framesCap {
			break
		}
	}

	count -= 2
	*dest = (*dest)[:count]
	return count
}

func (st *StackTrace) Frames() []uintptr { return st.frames }

func (st *StackTrace) CallersFrames() *runtime.Frames {
	return runtime.CallersFrames(st.frames)
}

func (st *StackTrace) Len() uint {
	if nil == st {
		return 0
	}
	return uint(len(st.frames))
}

// Format formats the slice of xerrors.Frame using a xerrors.Printer.
func (st *StackTrace) Format(printer xerrors.Printer) {
	if printer.Detail() {
		st.printFrames(printer, st.Skip)
	}
}

func (st *StackTrace) String() string {
	var p framesPrinter
	st.printFrames(&p, 0)
	return p.b.String()
}

func (st *StackTrace) printFrames(p Printer, skip uint) {
	var callers []uintptr
	if skip == 0 {
		callers = st.frames
	} else {
		callers = st.frames[:len(st.frames)-int(skip)]
	}

	cf := runtime.CallersFrames(reverse(callers))
	for {
		fr, more := cf.Next()
		p.Printf("%s\n    %s:%d\n", fr.Function, fr.File, fr.Line)
		if !more {
			break
		}
	}
}

func reverse(slice []uintptr) []uintptr {
	n := len(slice)
	for i := n/2 - 1; i >= 0; i-- {
		opp := n - 1 - i
		slice[i], slice[opp] = slice[opp], slice[i]
	}
	return slice
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
