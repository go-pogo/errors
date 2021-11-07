// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !notrace
// +build !notrace

package errors

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"
)

var captureFrames = true

func enableCaptureFrames()  { captureFrames = true }
func disableCaptureFrames() { captureFrames = false }

// Trace adds stack trace context to the error by calling StackTracer.Trace
// on the error. If the error is not a StackTracer it is wrapped with an
// OriginalGetter that implements this interface.
func Trace(err error) error { return TraceSkip(err, 1) }

// TraceSkip adds stack trace context to the error just like Trace. Unlike
// Trace it passes the skipFrames argument, after increasing it by 1, to
// StackTracer.Trace.
func TraceSkip(err error, skipFrames uint) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(StackTracer); ok {
		e.Trace(skipFrames + 1)
		return e
	}

	ce := upgrade(err)
	ce.Trace(skipFrames + 1)
	return ce
}

func (t *tracer) Trace(skipFrames uint) {
	if !captureFrames {
		return
	}

	skip := int(skipFrames)
	frames := make([]xerrors.Frame, 0, 10)

	for {
		skip += 1
		f := xerrors.Caller(skip)
		if !isValidFrame(f) {
			break
		}

		frames = append(frames, f)
	}

	if n := len(frames); n > 1 {
		t.frames = append(t.frames, frames[:n-1]...)
	}
}

const invalidFrameSuffix = " 0]}"

func isValidFrame(f xerrors.Frame) bool {
	s, x := fmt.Sprintf("%+v", f), invalidFrameSuffix
	return s[len(s)-len(x):] != x
}

// GetStackFrames returns a *Frames if err is a StackTracer or nil otherwise.
func GetStackFrames(err error) *Frames {
	if e, ok := err.(StackTracer); ok {
		return e.StackFrames()
	}
	return nil
}

func (t *tracer) StackFrames() *Frames { return &t.frames }

// String formats the captured frames and returns its string representation.
func (fr *Frames) String() string {
	var p framesPrinter
	fr.Format(&p)
	return p.b.String()
}

// framesPrinter is a xerrors.Printer that is used to print the string
// representation of Frames.
type framesPrinter struct{ b strings.Builder }

func (p *framesPrinter) Print(args ...interface{}) {
	_, _ = fmt.Fprint(&p.b, args...)
}

func (p *framesPrinter) Printf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(&p.b, format, args...)
}

func (p *framesPrinter) Detail() bool { return true }
