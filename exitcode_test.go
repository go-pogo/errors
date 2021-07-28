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
			want := toCommonErr(Original(rootCause), true)
			want.exitCode = 23

			assertErrorIs(t, have, rootCause)
			assert.Exactly(t, want, have)
			assert.Exactly(t, want.exitCode, GetExitCode(have))
		})
		t.Run("overwrite", func(t *testing.T) {
			have = WithExitCode(have, 45)
			want := toCommonErr(Original(rootCause), true)
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
		err  error
		want int
	}{
		"with nil": {
			err:  nil,
			want: 0,
		},
		"std error": {
			err:  stderrors.New("std err"),
			want: 0,
		},
		"std error with kind": {
			err:  WithExitCode(stderrors.New("std err"), 12),
			want: 12,
		},
		"common error": {
			err:  New("some error without kind"),
			want: 0,
		},
		"common error with kind": {
			err:  WithExitCode(New("bar"), 34),
			want: 34,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Exactly(t, tc.want, GetExitCode(tc.err))
		})
	}
}
