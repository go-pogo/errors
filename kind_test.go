package errs

import (
	"errors"
	"testing"

	"github.com/roeldev/go-fail"
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
			have := GetKind(tc.err)
			if have != tc.want {
				t.Error(fail.Diff{
					Func: "GetKind",
					Msg:  "should return the Kind of the error, or UnknownKind",
					Have: have,
					Want: tc.want,
				})
			}
		})
	}
}
