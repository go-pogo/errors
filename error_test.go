package errs

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/roeldev/go-fail"
	"golang.org/x/xerrors"
)

var testErrCmpOpts cmp.Options

func init() {
	testErrCmpOpts = cmp.Options{
		cmp.AllowUnexported(err{}, Inner{}, xerrors.Frame{}),
	}
}

func TestErr_Error(t *testing.T) {
	cause := xerrors.New("cause of error")
	tests := map[string]map[string]struct {
		err1 error
		err2 error
		want string
	}{
		"New/Newf": {
			"empty": {
				err1: New("", ""),
				err2: Newf("", ""),
				want: UnknownError,
			},
			"kind only": {
				err1: New("foo", ""),
				err2: Newf("foo", ""),
				want: "foo",
			},
			"message only": {
				err1: New(UnknownKind, "some `foo` happened"),
				err2: Newf(UnknownKind, "some `%s` happened", "foo"),
				want: "some `foo` happened",
			},
			"kind and message": {
				err1: New("foo error", "unexpected `bar`"),
				err2: Newf("foo error", "unexpected `%s`", "bar"),
				want: "foo error: unexpected `bar`",
			},
		},
		"Wrap/Wrapf": {
			"empty": {
				err1: Wrap(cause, "", ""),
				err2: Wrapf(cause, "", ""),
				want: UnknownError,
			},
			"kind only": {
				err1: Wrap(cause, "foo", ""),
				err2: Wrapf(cause, "foo", ""),
				want: "foo",
			},
			"message only": {
				err1: Wrap(cause, UnknownKind, "some `foo` happened"),
				err2: Wrapf(cause, UnknownKind, "some `%s` happened", "foo"),
				want: "some `foo` happened",
			},
			"kind and message": {
				err1: Wrap(cause, "foo error", "unexpected `bar`"),
				err2: Wrapf(cause, "foo error", "unexpected `%s`", "bar"),
				want: "foo error: unexpected `bar`",
			},
		},
	}

	for fn, ts := range tests {
		fn = strings.Replace(fn, "/", "&", 1)
		for name, tc := range ts {
			t.Run(fn+"__"+name, func(t *testing.T) {
				if have := tc.err1.Error(); have != tc.want {
					t.Error(fail.Diff{
						Func: fn,
						Have: have,
						Want: tc.want,
					})
				}
				if have := tc.err2.Error(); have != tc.want {
					t.Error(fail.Diff{
						Func: fn,
						Have: have,
						Want: tc.want,
					})
				}
			})
		}
	}
}

func TestNilWrap(t *testing.T) {
	tests := map[string]error{
		"Wrap":  Wrap(nil, UnknownKind, "foobar"),
		"Wrapf": Wrapf(nil, UnknownKind, "%s", "foobar"),
	}

	for name, have := range tests {
		t.Run(name, func(t *testing.T) {
			if have != nil {
				t.Error(fail.Diff{
					Func: name,
					Have: have,
					Want: nil,
				})
			}
		})
	}
}
