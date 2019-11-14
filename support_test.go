package errs

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/roeldev/go-fail"
)

var cmpOpts []cmp.Option

func init() {
	errStr := errors.New("")
	// error.errorString
	errorString := reflect.Indirect(reflect.ValueOf(errStr)).Interface()
	// fmt.wrapError
	fmtWrapError := reflect.Indirect(reflect.ValueOf(fmt.Errorf("%w", errStr))).Interface()

	cmpOpts = []cmp.Option{
		cmp.AllowUnexported(err{}, wrapErr{}, ST{}, errorString, fmtWrapError),
	}
}

func TestGetKind(t *testing.T) {
	kind := func(str string) *Kind {
		k := Kind(str)
		return &k
	}

	tests := map[string]struct {
		err  error
		want *Kind
	}{
		"nil": {
			err:  nil,
			want: nil,
		},
		"primitive": {
			err:  errors.New("foo bar"),
			want: nil,
		},
		"error": {
			err:  Err("foo", "bar"),
			want: kind("foo"),
		},
		"wrapped error": {
			err:  Wrap(Err("baz", "qux")),
			want: kind("baz"),
		},
		"wrapped error error": {
			err:  Wrapf(Err("baz", "qux"), "foo kind", "bar msg"),
			want: kind("foo kind"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			have := GetKind(tc.err)
			if !reflect.DeepEqual(have, tc.want) {
				t.Error(fail.Diff{
					Func: "GetKind",
					Msg:  "should return a pointer to the Kind of the error, or nil",
					Have: have,
					Want: tc.want,
					Opts: cmpOpts,
				})
			}
		})
	}
}

func TestGetMessage(t *testing.T) {
	tests := map[string]struct {
		err  error
		want string
	}{
		"nil": {
			err:  nil,
			want: "",
		},
		"primitive": {
			err:  errors.New("foo bar"),
			want: "",
		},
		"error": {
			err:  Err("foo", "bar"),
			want: "bar",
		},
		"wrapped error": {
			err:  Wrap(Err("baz", "qux")),
			want: "qux",
		},
		"wrapped error error": {
			err:  Wrapf(Err("baz", "qux"), "foo kind", "bar msg"),
			want: "bar msg",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			have := GetMessage(tc.err)
			if !reflect.DeepEqual(have, tc.want) {
				t.Error(fail.Diff{
					Func: "GetMessage",
					Msg:  "should return a single error message a string, or empty when not available",
					Have: have,
					Want: tc.want,
					Opts: cmpOpts,
				})
			}
		})
	}
}

func TestGetStackTrace(t *testing.T) {
	tests := map[string]struct {
		err     error
		wantNil bool
		wantLen uint
	}{
		"nil": {
			err:     nil,
			wantNil: true,
		},
		"primitive": {
			err:     errors.New("foo bar"),
			wantNil: true,
		},
		"error": {
			err:     Err("foo", "bar"),
			wantLen: 1,
		},
		"wrapped error": {
			err:     Wrap(Err("baz", "qux")),
			wantLen: 2,
		},
		"wrapped error error": {
			err:     Wrapf(Err("baz", "qux"), "foo kind", "bar msg"),
			wantLen: 1,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			st := GetStackTrace(tc.err)
			if st == nil {
				if !tc.wantNil {
					t.Error(fail.Msg{
						Func: "GetStackTrace",
						Msg:  "should return `nil` when error does not have a stack trace",
					})
				}
			} else {
				if st.Len() != tc.wantLen {
					t.Error(fail.Msg{
						Func: "GetStackTrace",
						Msg:  "should return a pointer to the stack trace with the given length of the error",
					})
				}
			}
		})
	}
}

func TestUnwrapAll(t *testing.T) {
	tests := map[string]struct {
		err    error
		wantFn func(e error) []error
	}{
		"nil": {
			err: Wrap(nil),
			wantFn: func(e error) []error {
				return make([]error, 0, 0)
			},
		},
		"primitive error": {
			err: errors.New("foo bar"),
			wantFn: func(e error) []error {
				return []error{errors.New("foo bar")}
			},
		},
		"wrapped primitive": {
			err: Wrap(errors.New("bar: baz")),
			wantFn: func(e error) []error {
				cause := errors.New("bar: baz")
				return []error{
					&wrapErr{st: GetStackTrace(e), err: cause},
					cause,
				}
			},
		},
		"double wrapped primitive": {
			err: Wrap(Wrap(errors.New("qux: xoo"))),
			wantFn: func(e error) []error {
				cause := errors.New("qux: xoo")
				return []error{
					&wrapErr{st: GetStackTrace(e), err: cause},
					cause,
				}
			},
		},
		"primitive wrap": {
			err: fmt.Errorf("cause: %w", errors.New("foo bar")),
			wantFn: func(e error) []error {
				return []error{e, errors.New("foo bar")}
			},
		},
		"wrapped primitive wrap": {
			err: Wrap(fmt.Errorf("cause: %w", errors.New("foo bar"))),
			wantFn: func(e error) []error {
				cause := errors.Unwrap(e)
				return []error{
					&wrapErr{st: GetStackTrace(e), err: cause},
					cause, // fmt.wrapError
					errors.New("foo bar"),
				}
			},
		},
		"error": {
			err: Err("kind", "err msg"),
			wantFn: func(e error) []error {
				return []error{&err{
					st:   GetStackTrace(e),
					kind: "kind",
					msg:  "err msg",
				}}
			},
		},
		"wrapped error": {
			err: Wrap(Err("kind", "err msg")),
			wantFn: func(e error) []error {
				return []error{&err{
					st:   GetStackTrace(e),
					kind: "kind",
					msg:  "err msg",
				}}
			},
		},
		"double wrapped error": {
			err: Wrap(Wrap(Err("kind", "err msg"))),
			wantFn: func(e error) []error {
				return []error{&err{
					st:   GetStackTrace(e),
					kind: "kind",
					msg:  "err msg",
				}}
			},
		},
		"wrapped error error": {
			err: Wrapf(Err("baz", "qux"), "foo kind", "bar msg"),
			wantFn: func(e error) []error {
				wrErr := &err{
					st:   GetStackTrace(errors.Unwrap(e)),
					kind: "baz",
					msg:  "qux",
				}

				return []error{
					&err{
						st:   GetStackTrace(e),
						err:  wrErr,
						kind: "foo kind",
						msg:  "bar msg",
					},
					wrErr,
				}
			},
		},
	}

	for label, tc := range tests {
		t.Run(label, func(t *testing.T) {
			have := UnwrapAll(tc.err)
			want := tc.wantFn(tc.err)

			if !reflect.DeepEqual(have, want) {
				t.Error(fail.Diff{
					Func: "UnwrapAll",
					Msg:  "should unwrap all and return a slice of errors",
					Have: have,
					Want: want,
					Opts: cmpOpts,
				})
			}
		})
	}
}

func TestUnwrapCause(t *testing.T) {
	tests := map[string]struct {
		err    error
		wantFn func(e error) error
	}{
		"primitive error": {
			err: errors.New("foo bar"),
			wantFn: func(e error) error {
				return errors.New("foo bar")
			},
		},
		"wrapped primitive error": {
			err: Wrap(errors.New("foo bar")),
			wantFn: func(e error) error {
				return errors.New("foo bar")
			},
		},
		"primitive wrap": {
			err: fmt.Errorf("%w", errors.New("foo bar")),
			wantFn: func(e error) error {
				return errors.New("foo bar")
			},
		},
		"wrapped primitive wrap": {
			err: Wrap(fmt.Errorf("cause: %w", errors.New("baz"))),
			wantFn: func(e error) error {
				return errors.New("baz")
			},
		},
		"error": {
			err: Err("qux", "xoo"),
			wantFn: func(e error) error {
				return &err{
					st:   GetStackTrace(e),
					kind: "qux",
					msg:  "xoo",
				}
			},
		},
		"wrapped error": {
			err: Wrap(Err("qux", "xoo")),
			wantFn: func(e error) error {
				return &err{
					st:   GetStackTrace(e),
					kind: "qux",
					msg:  "xoo",
				}
			},
		},
		"double wrapped error": {
			err: Wrap(Wrap(Err("qux", "xoo"))),
			wantFn: func(e error) error {
				return &err{
					st:   GetStackTrace(e),
					kind: "qux",
					msg:  "xoo",
				}
			},
		},
	}

	for label, tc := range tests {
		t.Run(label, func(t *testing.T) {
			have := UnwrapCause(tc.err)
			want := tc.wantFn(tc.err)

			if !reflect.DeepEqual(have, want) {
				t.Error(fail.Diff{
					Func: "UnwrapCause",
					Have: have,
					Want: want,
					Opts: cmpOpts,
				})
			}
		})
	}
}
