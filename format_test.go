// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-pogo/errors/internal"
)

func TestFormatError(t *testing.T) {
	tests := map[string]struct {
		setup      func() error
		traceLines []int
	}{
		"error": {
			setup: func() error {
				return New("some err")
			},
			traceLines: []int{21},
		},
		"traced primitive": {
			setup: func() error {
				return Trace(stderrors.New("primitive"))
			},
			traceLines: []int{27},
		},
		"traced error": {
			setup: func() error {
				err := New("another err")
				return Trace(err)
			},
			traceLines: []int{33, 34},
		},
		"multi error": {
			setup: func() error {
				err1 := New("err1")
				err2 := New("err2")
				return Combine(err1, err2)
			},
			traceLines: []int{40, 41, 42},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.setup()
			str := fmt.Sprintf("%+v", err)

			for _, line := range tc.traceLines {
				assert.Contains(t, str, "format_test.go:"+strconv.Itoa(line))
			}
		})
	}
	t.Run("", func(t *testing.T) {
		internal.DisableCaptureFrames()
		defer internal.EnableCaptureFrames()

		rootCause := stderrors.New("root cause")
		assert.Equal(t,
			fmt.Sprintf("%+v", WithFormatter(rootCause)),
			fmt.Sprintf("%+v", formatErrFixture{error: rootCause}),
		)
	})
}

type formatErrFixture struct{ error }

func (f *formatErrFixture) Format(s fmt.State, v rune) { FormatError(f, s, v) }
