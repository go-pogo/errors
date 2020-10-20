// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"

	"github.com/go-pogo/errors/internal"
)

type Frames []xerrors.Frame

// capture captures a xerrors.Frame that describes a frame on the caller's
// stack. The argument skip is the number of frames to skip over.
// capture(0) returns the frame for the caller of capture.
// It returns a bool false when the captured frame contains a nil pointer.
func (fr *Frames) capture(skip uint) (ok bool) {
	if !internal.CaptureFrames() {
		return true
	}

	f := xerrors.Caller(int(skip) + 1)
	if ok = isValidFrame(f); ok {
		*fr = append(*fr, f)
	}
	return ok
}

// Format formats the captured frames using a xerrors.Printer.
func (fr Frames) Format(p xerrors.Printer) {
	for i := len(fr) - 1; i >= 0; i-- {
		fr[i].Format(p)
	}
}

// String formats the captured frames and returns its string representation.
func (fr *Frames) String() string {
	var p framesPrinter
	fr.Format(&p)
	return p.b.String()
}

const invalidFrameSuffix = " 0]}"

func isValidFrame(f xerrors.Frame) bool {
	s, x := fmt.Sprintf("%+v", f), invalidFrameSuffix
	return s[len(s)-len(x):] != x
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
