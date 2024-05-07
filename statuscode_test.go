// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestWithStatusCode(t *testing.T) {
	for name, wantErr := range provideErrors(true) {
		t.Run(name, func(t *testing.T) {
			haveErr := WithStatusCode(wantErr, http.StatusOK)
			assert.Exactly(t, http.StatusOK, GetStatusCode(haveErr))
			assert.ErrorIs(t, haveErr, wantErr)

			// update existing exitcode
			t.Run("update", func(t *testing.T) {
				haveErr2 := WithStatusCode(haveErr, http.StatusContinue)
				assert.Exactly(t, http.StatusContinue, GetStatusCode(haveErr2))
				assert.Same(t, haveErr, haveErr2)
			})
		})
	}

	t.Run("nil", func(t *testing.T) {
		assert.Exactly(t, nil, WithStatusCode(nil, http.StatusNotFound))
	})
}

func TestGetStatusCodeOr(t *testing.T) {
	tests := map[string]struct {
		err  error
		or   int
		want int
	}{
		"nil": {
			err:  nil,
			or:   http.StatusOK,
			want: http.StatusOK,
		},
		"std error": {
			err:  stderrors.New("std err"),
			or:   http.StatusContinue,
			want: http.StatusContinue,
		},
		"std error with exit code": {
			err:  WithStatusCode(stderrors.New("std err"), http.StatusAccepted),
			want: http.StatusAccepted,
		},
		"error": {
			err:  New("some error without exit code"),
			or:   http.StatusNotFound,
			want: http.StatusNotFound,
		},
		"error with exit code": {
			err:  WithStatusCode(New("bar"), http.StatusFound),
			want: http.StatusFound,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Exactly(t, tc.want, GetStatusCodeOr(tc.err, tc.or))
		})
	}
}
