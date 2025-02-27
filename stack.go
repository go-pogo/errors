// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/go-pogo/errors/internal"
	"golang.org/x/xerrors"
)

// StackTracer interfaces provides access to a [StackTrace].
type StackTracer interface {
	error

	// StackTrace returns a stack of trace frames.
	StackTrace() *StackTrace
}

const panicUseNewInstead = "errors.WithStack: use errors.New instead to create an error from an errors.Msg"

// WithStack gets a stack trace at the point WithStack was called and adds it
// to the error. If err is nil, WithStack returns nil.
func WithStack(err error) StackTracer {
	if err == nil {
		return nil
	}

	//goland:noinspection GoTypeAssertionOnErrors
	switch v := err.(type) {
	case StackTracer:
		return v

	case Msg, *Msg:
		panic(panicUseNewInstead)

	default:
		ee := &embedError{error: v}
		if internal.TraceStack {
			ee.stack = newStackTrace(1)
			if u := Unwrap(v); u != nil {
				skipStackTrace(u, ee.stack.Len())
			}
		}
		return ee
	}
}

// GetStackTrace returns a [StackTrace] if err is a [StackTracer] or nil
// otherwise. It will always return nil when the "notrace" build tag is set.
func GetStackTrace(err error) *StackTrace {
	if !internal.TraceStack {
		return nil
	}

	//goland:noinspection GoTypeAssertionOnErrors
	if e, ok := err.(StackTracer); ok {
		return e.StackTrace()
	}
	return nil
}

type StackTrace struct {
	frames   []uintptr
	reversed bool

	// Skip n frames when formatting with [Format], so overlapping frames from
	// previous errors are not printed.
	Skip uint
}

const framesDepth = 16

func newStackTrace(skipFrames uint) *StackTrace {
	st := &StackTrace{frames: make([]uintptr, 0, framesDepth)}

	skip := int(skipFrames) + 2
	var pc [framesDepth]uintptr
	for {
		n := runtime.Callers(skip, pc[:])
		if n == 0 {
			break
		}

		st.frames = append(st.frames, pc[:n]...)
		skip += n

		if n < framesDepth {
			break
		}
	}

	st.frames = st.frames[:len(st.frames)-2]
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

func (st *StackTrace) reverseFrames() {
	if st.reversed {
		return
	}

	n := len(st.frames)
	for i := n/2 - 1; i >= 0; i-- {
		opp := n - 1 - i
		st.frames[i], st.frames[opp] = st.frames[opp], st.frames[i]
	}
}

type Frame uintptr

// PC is the program counter for the location in this frame.
func (fr Frame) PC() uintptr { return uintptr(fr) }

// Func returns a [runtime.Func] describing the function that contains the
// given program counter address, or else nil.
func (fr Frame) Func() *runtime.Func { return runtime.FuncForPC(fr.PC()) }

// FileLine returns the file name and line number of the source code
// corresponding to the program counter [PC].
func (fr Frame) FileLine() (file string, line int) {
	if f := fr.Func(); f == nil {
		return "", 0
	} else {
		return f.FileLine(fr.PC())
	}
}

// Frames returns a slice of [Frame]. Use [StackTrace.CallersFrames] instead if
// you want to access the whole stack trace of frames.
func (st *StackTrace) Frames() []Frame {
	st.reverseFrames()
	frames := make([]Frame, len(st.frames))
	for i, pc := range st.frames {
		frames[i] = Frame(pc)
	}
	return frames
}

// CallersFrames returns a [runtime.Frames] by calling [runtime.CallersFrames]
// with the captured stack trace frames as callers argument.
func (st *StackTrace) CallersFrames() *runtime.Frames {
	st.reverseFrames()
	return runtime.CallersFrames(st.frames)
}

// Len returns the amount of captures frames.
func (st *StackTrace) Len() uint {
	if nil == st {
		return 0
	}
	return uint(len(st.frames))
}

// Format formats the slice of [xerrors.Frame] using a [xerrors.Printer]. It
// will skip n frames according to [StackTrace.Skip], when printing so no
// overlapping frames with underlying errors are displayed.
func (st *StackTrace) Format(printer xerrors.Printer) {
	if printer.Detail() {
		st.printFrames(printer, st.Skip)
	}
}

// String returns a formatted string of the complete stack trace.
func (st *StackTrace) String() string {
	var b strings.Builder
	st.printFrames(&framesPrinter{&b}, 0)
	return b.String()
}

func (st *StackTrace) printFrames(p Printer, skip uint) {
	st.reverseFrames()
	PrintFrames(p, runtime.CallersFrames(st.frames[skip:]))
}

// PrintFrames prints a complete stack of [runtime.Frames] using [Printer] p.
func PrintFrames(p Printer, cf *runtime.Frames) {
	for {
		f, more := cf.Next()
		p.Printf("%s\n    %s:%d\n", f.Function, f.File, f.Line)
		if !more {
			break
		}
	}
}

// framesPrinter is a [xerrors.Printer] that is used to print the string
// representation of [StackTrace].
type framesPrinter struct{ b io.Writer }

func (p *framesPrinter) Print(args ...interface{}) {
	_, _ = fmt.Fprint(p.b, args...)
}

func (p *framesPrinter) Printf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(p.b, format, args...)
}

func (*framesPrinter) Detail() bool { return true }
