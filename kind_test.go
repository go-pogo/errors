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

func TestKind(t *testing.T) {
	kind := Kind("some msg")
	assert.Equal(t, kind.String(), kind.Error())
}

func TestKindf(t *testing.T) {
	tests := map[string][]interface{}{
		"no args": nil,
		"some %s": {"string"},
		"%s %s":   {"foo", "bar"},
	}
	for f, a := range tests {
		t.Run(f, func(t *testing.T) {
			assert.Equal(t, Kind(fmt.Sprintf(f, a...)), Kindf(f, a...))
		})
	}
}

func TestKind_Is(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		kind := Kind("foobar")
		tests := map[string]error{
			"Kind":  Kind("foobar"),
			"*Kind": &kind,
		}
		for a, err := range tests {
			for b, target := range tests {
				t.Run(a+"/"+b, func(t *testing.T) {
					assert.ErrorIs(t, err, target)
					assert.ErrorIs(t, &kindError{kind: "foobar"}, target)
				})
			}
		}
	})

	t.Run("false", func(t *testing.T) {
		targets := map[string]error{
			"stderror":             stderrors.New("some err"),
			"different msg string": Kind("blabla"),
		}
		for name, target := range targets {
			t.Run(name, func(t *testing.T) {
				assert.NotErrorIs(t, Kind("some err"), target)
				assert.NotErrorIs(t, &kindError{kind: "some err"}, target)
			})
		}
	})
}

func TestKind_As(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		var dest Kind
		assert.True(t, Kind("hi there").As(&dest))
		assert.Exactly(t, Kind("hi there"), dest)

		dest = ""
		assert.True(t, (&kindError{kind: "hi there"}).As(&dest))
		assert.Exactly(t, Kind("hi there"), dest)

	})
	t.Run("false", func(t *testing.T) {
		var dest Kind
		assert.False(t, Kind("hi there").As(dest))
		assert.Exactly(t, Kind(""), dest)

		assert.False(t, (&kindError{kind: "hi there"}).As(dest))
		assert.Exactly(t, Kind(""), dest)
	})
}

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
			assert.ErrorIs(t, have, tc.err)
			assert.ErrorIs(t, have, tc.kind)
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
