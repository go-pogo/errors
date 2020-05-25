package errs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetKind(t *testing.T) {
	tests := map[string]struct {
		err  error
		want Kind
	}{
		"nil": {
			err:  nil,
			want: UnknownKind,
		},
		"primitive": {
			err:  errors.New("foo bar"),
			want: UnknownKind,
		},
		"error": {
			err:  New("foo", "bar"),
			want: Kind("foo"),
		},
		"wrapped error": {
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
