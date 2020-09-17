package errors

import (
	stderrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetKind(t *testing.T) {
	tests := map[string]struct {
		err  error
		want Kind
	}{
		"with nil": {
			err:  nil,
			want: UnknownKind,
		},
		"with primitive error": {
			err:  stderrors.New("foo bar"),
			want: UnknownKind,
		},
		"with error": {
			err:  New("foo", "bar"),
			want: Kind("foo"),
		},
		"with wrapped error": {
			err:  Trace(New("baz", "qux")),
			want: Kind("baz"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Exactly(t, tc.want, GetKind(tc.err))
		})
	}
}
