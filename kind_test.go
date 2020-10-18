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
			want := toCommonErr(Original(rootCause), true)
			want.kind = kind1

			assertErrorIs(t, have, rootCause)
			assert.Exactly(t, want, have)
			assert.Exactly(t, want.kind, GetKind(have))
		})
		t.Run("overwrite", func(t *testing.T) {
			have = WithKind(have, kind2)
			want := toCommonErr(Original(rootCause), true)
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
	tests := map[string]struct {
		err  error
		want Kind
	}{
		"with nil": {
			err:  nil,
			want: UnknownKind,
		},
		"std error": {
			err:  stderrors.New("std err"),
			want: UnknownKind,
		},
		"std error with kind": {
			err:  WithKind(stderrors.New("std err"), "xoo"),
			want: Kind("xoo"),
		},
		"common error": {
			err:  New("some error without kind"),
			want: UnknownKind,
		},
		"common error with kind": {
			err:  WithKind(New("bar"), "foo"),
			want: Kind("foo"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Exactly(t, tc.want, GetKind(tc.err))
		})
	}
}
