package errors

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

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
				assert.Equal(t, tc.want, tc.err1.Error())
				assert.Equal(t, tc.want, tc.err2.Error())
			})
		}
	}
}

func TestWrap(t *testing.T) {
	t.Run("with nil cause", func(t *testing.T) {
		assert.Nil(t, Wrap(nil, UnknownKind, "foobar"))
	})
}

func TestWrapf(t *testing.T) {
	t.Run("with nil cause", func(t *testing.T) {
		assert.Nil(t, Wrapf(nil, UnknownKind, "%s", "foobar"))
	})
}
