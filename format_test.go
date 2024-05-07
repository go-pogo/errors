// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"bytes"
	stderrors "errors"
	"fmt"
	"github.com/go-pogo/errors/internal"
	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
	"strings"
	"testing"
	"time"
)

func TestWithFormatter(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Exactly(t, nil, WithFormatter(nil))
	})

	t.Run("std error", func(t *testing.T) {
		err := stderrors.New("some err")
		have := WithFormatter(err)
		assert.Implements(t, (*xerrors.Formatter)(nil), have)
		assert.Same(t, err, Unembed(have))
	})

	tests := map[string]error{
		"error":          New("whoops"),
		"multiErr":       newMultiErr([]error{New("hi"), stderrors.New("there")}, 0),
		"with exit code": WithExitCode(New("cause"), 1),
		"with formatter": WithFormatter(New("cause")),
		"with kind":      WithKind(Msg("my bad"), "failure"),
		// "with stack":     WithStack(fmt.Errorf("failure: %w", stderrors.New("my bad"))),
		"with time": WithTime(New("cause"), time.Now()),
	}

	for name, have := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Implements(t, (*xerrors.Formatter)(nil), have)
		})
	}
}

func TestFormatError(t *testing.T) {
	internal.DisableTraceStack()
	defer internal.EnableTraceStack()

	tests := map[string]struct {
		err          error
		want1, want2 string
	}{
		"nil": {},
		"std error": {
			err:   stderrors.New("some err"),
			want1: "some err",
			want2: "some err",
		},
		"error": {
			err:   New("oh noes"),
			want1: "oh noes",
			want2: "oh noes",
		},
		"wrap wrap": {
			err:   Wrap(Wrap(New("some err"), "failure"), "whoops"),
			want1: "whoops: failure: some err",
			want2: "whoops:\n  - failure:\n  - some err",
		},
		"kind": {
			err:   WithKind(New("some err"), "failure"),
			want1: "failure: some err",
			want2: "failure: some err:\n  - some err",
		},
		"kind Msg": {
			err:   WithKind(Msg("some err"), "failure"),
			want1: "failure: some err",
			want2: "failure: some err",
		},
		"wrap kind": {
			err:   Wrap(WithKind(New("some err"), "failure"), "whoops"),
			want1: "whoops: failure: some err",
			want2: "whoops:\n  - failure: some err:\n  - some err",
		},
		"kind wrap": {
			err:   WithKind(Wrap(New("some err"), "failure"), "whoops"),
			want1: "whoops: failure: some err",
			want2: "whoops: failure:\n  - failure:\n  - some err",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// without details
			t.Run("plain", func(t *testing.T) {
				var state fmtStateHelper
				FormatError(tc.err, &state, 'v')
				assert.Exactly(t, tc.want1, state.String())

				if tc.err != nil {
					assert.Exactly(t, tc.want1, fmt.Sprintf("%v", tc.err))
				}
			})
			// with details
			t.Run("details", func(t *testing.T) {
				state := fmtStateHelper{flags: "+"}
				FormatError(tc.err, &state, 'v')
				assert.Exactly(t, tc.want2, state.String())

				if tc.err != nil {
					assert.Exactly(t, tc.want2, fmt.Sprintf("%+v", tc.err))
				}
			})
		})
	}
}

func TestPrintError(t *testing.T) {
	tests := map[string]struct {
		err  error
		want string
	}{
		"nil": {},
		"std error": {
			err:  stderrors.New("some err"),
			want: "some err",
		},
		"error": {
			err:  New("oh noes"),
			want: "oh noes",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var printer fmtStateHelper
			PrintError(&printer, tc.err)
			assert.Exactly(t, tc.want, printer.String())

			printer.Reset()
			printer.flags = "+"
			PrintError(&printer, tc.err)
			have := printer.String()
			assert.Truef(t, strings.HasPrefix(have, tc.want), "should have prefix `%s`", tc.want)
		})
	}

}

type fmtStateHelper struct {
	flags string
	bytes.Buffer
}

func (ts *fmtStateHelper) Width() (int, bool) { return 0, false }

func (ts *fmtStateHelper) Precision() (int, bool) { return 0, false }

func (ts *fmtStateHelper) Flag(f int) bool {
	return strings.ContainsRune(ts.flags, rune(f))
}

func (ts *fmtStateHelper) Detail() bool { return ts.Flag('+') }

func (ts *fmtStateHelper) Print(args ...interface{}) {
	_, _ = fmt.Fprint(ts, args...)
}

func (ts *fmtStateHelper) Printf(f string, args ...interface{}) {
	_, _ = fmt.Fprintf(ts, f, args...)
}
