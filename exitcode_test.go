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
	tests := map[string]struct {
		err          error
		code1, code2 int
	}{
		"std error": {
			err:   stderrors.New("some error"),
			code1: 1,
			code2: 2,
		},
		"error": {
			err:   New("absolute horror"),
			code1: 99,
			code2: 0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			have := WithExitCode(tc.err, tc.code1)
			assert.Exactly(t, tc.code1, GetExitCode(have))
			assertErrorIs(t, have, tc.err)

			// update existing exitcode
			t.Run("update", func(t *testing.T) {
				have2 := WithExitCode(have, tc.code2)
				assert.Exactly(t, tc.code2, GetExitCode(have2))
				assert.Same(t, have, have2)
			})
		})
	}

	t.Run("nil", func(t *testing.T) {
		assert.Exactly(t, nil, WithExitCode(nil, 666))
	})
}

func TestGetExitCodeOr(t *testing.T) {
	tests := map[string]struct {
		err  error
		or   int
		want int
	}{
		"nil": {
			err:  nil,
			or:   12,
			want: 12,
		},
		"std error": {
			err:  stderrors.New("std err"),
			or:   23,
			want: 23,
		},
		"std error with exit code": {
			err:  WithExitCode(stderrors.New("std err"), 12),
			want: 12,
		},
		"error": {
			err:  New("some error without exit code"),
			or:   99,
			want: 99,
		},
		"error with exit code": {
			err:  WithExitCode(New("bar"), 34),
			want: 34,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Exactly(t, tc.want, GetExitCodeOr(tc.err, tc.or))
		})
	}
}
