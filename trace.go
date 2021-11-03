// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

// StackTracer interfaces provide access to a stack of traced Frames.
type StackTracer interface {
	error

	// StackFrames returns a slice of captured xerrors.Frame types associated
	// with the error.
	StackFrames() *Frames
	// Trace captures a xerrors.Frame that describes a frame on the caller's
	// stack. The argument skipFrames is the number of frames to skip over.
	Trace(skipFrames uint)
}

// Trace adds stack trace context to the error by calling StackTracer.Trace
// on the error. If the error is not a StackTracer it is wrapped with an
// UpgradedError that implements this interface.
func Trace(err error) error { return TraceSkip(err, 1) }

// TraceSkip adds stack trace context to the error just like Trace. Unlike
// Trace it passes the skipFrames argument to StackTracer.Trace.
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

// GetStackFrames returns a *Frames if err is a StackTracer or nil otherwise.
func GetStackFrames(err error) *Frames {
	if e, ok := err.(StackTracer); ok {
		return e.StackFrames()
	}
	return nil
}

type tracer struct {
	frames Frames
}

func (t *tracer) StackFrames() *Frames { return &t.frames }

func (t *tracer) Trace(skipFrames uint) { t.frames.capture(skipFrames + 1) }
