// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

// StackTracer interfaces provide access to a stack of traced Frames.
type StackTracer interface {
	error
	StackFrames() *Frames
	Trace(skipFrames uint)
}

// Trace adds stack trace context to the error by calling StackTracer.Trace on
// the error. If the error is not a StackTracer it is wrapped with a Proxy that
// implements this interface.
func Trace(err error) error { return TraceSkip(err, 1) }

// TraceSkip adds stack trace context to the error just like Trace. Unlike Trace
// it passes the skipFrames argument to StackTracer.Trace.
func TraceSkip(err error, skipFrames uint) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(StackTracer); ok {
		e.Trace(skipFrames + 1)
		return e
	}

	ce := toCommonErr(err, true)
	ce.Trace(skipFrames + 1)
	return ce
}

func GetStackFrames(err error) *Frames {
	if e, ok := err.(StackTracer); ok {
		return e.StackFrames()
	}
	return nil
}

type tracer struct {
	frames Frames
}

// StackFrames returns a slice of captured xerrors.Frame types linked to this
// error.
func (e *tracer) StackFrames() *Frames { return &e.frames }

func (e *tracer) Trace(skipFrames uint) { e.frames.capture(skipFrames + 1) }
