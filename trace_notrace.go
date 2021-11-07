// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build notrace
// +build notrace

package errors

var captureFrames = false

func enableCaptureFrames()  {}
func disableCaptureFrames() {}

func Trace(err error) error { return err }

func TraceSkip(err error, skipFrames uint) error { return err }

func (t *tracer) Trace(skipFrames uint) {}

func GetStackFrames(err error) *Frames { return nil }

func (t *tracer) StackFrames() *Frames { return nil }

func (fr *Frames) String() string { return "" }
