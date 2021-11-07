// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithExitCode(t *testing.T) {
	t.Run("std error", func(t *testing.T) {
		rootCause := stderrors.New("root cause error")
		have := WithExitCode(rootCause, 23)

		t.Run("add", func(t *testing.T) {
			want := &exitCodeErr{
				error:    rootCause,
				exitCode: 23,
			}
			assertErrorIs(t, have, rootCause)
			assert.Exactly(t, want, have)
			assert.Exactly(t, want.exitCode, GetExitCode(have))
		})
		t.Run("overwrite", func(t *testing.T) {
			have = WithExitCode(have, 45)
			want := &exitCodeErr{
				error:    rootCause,
				exitCode: 45,
			}
			assertErrorIs(t, have, rootCause)
			assert.Exactly(t, want, have)
			assert.Exactly(t, want.exitCode, GetExitCode(have))
		})
	})

	t.Run("common error", func(t *testing.T) {
		rootCause := New("root cause error")
		have := WithExitCode(rootCause, 23)

		t.Run("set", func(t *testing.T) {
			want := upgrade(Original(rootCause))
			want.exitCode = 23

			assertErrorIs(t, have, rootCause)
			assert.Exactly(t, want, have)
			assert.Exactly(t, want.exitCode, GetExitCode(have))
		})
		t.Run("overwrite", func(t *testing.T) {
			have = WithExitCode(have, 45)
			want := upgrade(Original(rootCause))
			want.exitCode = 45

			assertErrorIs(t, have, rootCause)
			assert.Exactly(t, want, have)
			assert.Exactly(t, want.exitCode, GetExitCode(have))
		})
	})

	t.Run("nil", func(t *testing.T) {
		assert.Exactly(t, nil, WithExitCode(nil, 666))
	})
}

func TestGetExitCode(t *testing.T) {
	tests := map[string]struct {
		err    error
		want   int
		orWant map[int]int
	}{
		"with nil": {
			err:    nil,
			orWant: map[int]int{1: 1, 2: 2},
		},
		"std error": {
			err:    stderrors.New("std err"),
			orWant: map[int]int{1: 1, 2: 2},
		},
		"std error with kind": {
			err:    WithExitCode(stderrors.New("std err"), 12),
			want:   12,
			orWant: map[int]int{1: 12, 2: 12},
		},
		"common error": {
			err:    New("some error without kind"),
			orWant: map[int]int{1: 0, 2: 0},
		},
		"common error with kind": {
			err:    WithExitCode(New("bar"), 34),
			want:   34,
			orWant: map[int]int{1: 34, 2: 34},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Exactly(t, tc.want, GetExitCode(tc.err))
			assert.Exactly(t, tc.want, GetExitCodeOr(tc.err, 0))

			for or, want := range tc.orWant {
				t.Run("", func(t *testing.T) {
					assert.Exactly(t, want, GetExitCodeOr(tc.err, or))
				})
			}
		})
	}
}
