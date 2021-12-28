// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithKind(t *testing.T) {
	tests := map[string]struct {
		err  error
		kind Kind
	}{
		"std error": {
			err:  stderrors.New("root cause error"),
			kind: "foo",
		},
		"wrapped std error with kind": {
			err:  WithKind(stderrors.New("root cause error"), "foobar"),
			kind: "qux xoo",
		},
		"std wrapped error": {
			err:  fmt.Errorf("bar: %w", stderrors.New("absolute horror")),
			kind: "baz",
		},
		"error": {
			err:  New("just some err"),
			kind: "whoops",
		},
		"wrapped error": {
			err:  Wrap(New("just some err"), "some reason"),
			kind: "some error",
		},
		"wrapped error with kind": {
			err:  WithKind(New("just some err"), "who did this"),
			kind: "whoops",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			have := WithKind(tc.err, tc.kind)
			assertErrorIs(t, have, tc.err)
			assertErrorIs(t, have, tc.kind)
			assert.Exactly(t, tc.kind, GetKind(have))
		})
	}

	t.Run("nil", func(t *testing.T) {
		assert.Nil(t, WithKind(nil, "nope"))
	})
}

func TestGetKind(t *testing.T) {
	const (
		foo Kind = "foo"
		bar Kind = "bar"
		qux Kind = "qux"
	)

	tests := map[string]struct {
		err    error
		want   Kind
		orWant Kind
	}{
		"nil": {
			err:    nil,
			orWant: foo,
		},
		"std error": {
			err:    stderrors.New("std err"),
			orWant: qux,
		},
		"std error with kind": {
			err:  WithKind(stderrors.New("std err"), foo),
			want: foo,
		},
		"error": {
			err:    New("some error without kind"),
			orWant: bar,
		},
		"error with kind": {
			err:  WithKind(New("bar"), bar),
			want: bar,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Exactly(t, tc.want, GetKind(tc.err))
			assert.Exactly(t, tc.want, GetKindOr(tc.err, UnknownKind))

			if tc.orWant != UnknownKind {
				assert.Exactly(t, tc.orWant, GetKindOr(tc.err, tc.orWant))
			}
		})
	}
}
