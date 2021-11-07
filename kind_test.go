// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	stderrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithKind(t *testing.T) {
	kind1 := Kind("foobar")
	kind2 := Kind("updated err")

	t.Run("std error", func(t *testing.T) {
		rootCause := stderrors.New("root cause error")
		have := WithKind(rootCause, kind1)

		t.Run("add", func(t *testing.T) {
			want := &kindErr{
				error: rootCause,
				kind:  kind1,
			}
			assertErrorIs(t, have, rootCause)
			assert.Exactly(t, want, have)
			assert.Exactly(t, want.kind, GetKind(have))
		})
		t.Run("overwrite", func(t *testing.T) {
			have = WithKind(have, kind2)
			want := &kindErr{
				error: rootCause,
				kind:  kind2,
			}
			assertErrorIs(t, have, rootCause)
			assert.Exactly(t, want, have)
			assert.Exactly(t, want.kind, GetKind(have))
		})
	})

	t.Run("common error", func(t *testing.T) {
		rootCause := New("root cause error")
		have := WithKind(rootCause, kind1)

		t.Run("set", func(t *testing.T) {
			want := upgrade(Original(rootCause))
			want.kind = kind1

			assertErrorIs(t, have, rootCause)
			assert.Exactly(t, want, have)
			assert.Exactly(t, want.kind, GetKind(have))
		})
		t.Run("overwrite", func(t *testing.T) {
			have = WithKind(have, kind2)
			want := upgrade(Original(rootCause))
			want.kind = kind2

			assertErrorIs(t, have, rootCause)
			assert.Exactly(t, want, have)
			assert.Exactly(t, want.kind, GetKind(have))
		})
	})

	t.Run("nil", func(t *testing.T) {
		assert.Exactly(t, nil, WithKind(nil, "some kind"))
	})
}

func TestGetKind(t *testing.T) {
	const (
		foo Kind = "foo"
		bar Kind = "bar"
		baz Kind = "baz"
		xoo Kind = "xoo"
	)

	tests := map[string]struct {
		err    error
		want   Kind
		orWant map[Kind]Kind
	}{
		"with nil": {
			err:    nil,
			want:   UnknownKind,
			orWant: map[Kind]Kind{foo: foo, bar: bar},
		},
		"std error": {
			err:    stderrors.New("std err"),
			want:   UnknownKind,
			orWant: map[Kind]Kind{foo: foo, bar: bar},
		},
		"std error with kind": {
			err:    WithKind(stderrors.New("std err"), xoo),
			want:   xoo,
			orWant: map[Kind]Kind{foo: xoo, bar: xoo},
		},
		"common error": {
			err:    New("some error without kind"),
			want:   UnknownKind,
			orWant: map[Kind]Kind{foo: UnknownKind, bar: UnknownKind},
		},
		"common error with kind": {
			err:    WithKind(New("bar"), baz),
			want:   baz,
			orWant: map[Kind]Kind{foo: baz, bar: baz},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Exactly(t, tc.want, GetKind(tc.err))
			assert.Exactly(t, tc.want, GetKindOr(tc.err, UnknownKind))

			for or, want := range tc.orWant {
				t.Run("", func(t *testing.T) {
					assert.Exactly(t, want, GetKindOr(tc.err, or))
				})
			}
		})
	}
}
