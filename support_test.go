package errs

import (
	"errors"
	"reflect"
	"testing"

	"github.com/roeldev/go-fail"
)

var emptyErr = reflect.Indirect(reflect.ValueOf(errors.New(""))).Interface()

// func TestPrint(t *testing.T) {
// 	err1 := Err("foo error")
// 	err2 := Wrap(err1)
// 	err3 := Wrap(err2)
//
// 	have := err3.Error()
// 	want := `foo error
//
// Trace:
// support_test.go:13: errors.TestPrint():
// support_test.go:12: errors.TestPrint():
// support_test.go:11: errors.TestPrint():
// >	foo error
// `
//
// 	if have != want {
// 		t.Error(fail.Diff{
// 			Func: "errs.Print",
// 			Msg:  "should create the same error output message",
// 			Have: have,
// 			Want: want,
// 		})
// 	}
// }

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
				})
			}
		})
	}
}

func TestGetStackTrace(t *testing.T) {
	tests := map[string]struct {
		err   error
		stLen uint
	}{
		"nil": {
			err:   nil,
			stLen: 0,
		},
		"primitive": {
			err:   errors.New("foo bar"),
			stLen: 0,
		},
		"error": {
			err:   Err("foo", "bar"),
			stLen: 1,
		},
		"wrapped error": {
			err:   Wrap(Err("baz", "qux")),
			stLen: 2,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var have uint

			st := GetStackTrace(tc.err)
			if st != nil {
				have = st.Len()
			}

			if have != tc.stLen {
				t.Error(fail.Diff{
					Func: "GetStackTrace",
					Msg:  "should return a pointer to the stack trace of the error",
				})
			}
		})
	}
}

func TestUnwrapAll(t *testing.T) {
	tests := map[string]struct {
		err  error
		want []error
		fn   func(e error) []error
	}{
		"nil": {
			err:  Wrap(nil),
			want: []error{},
		},
		"primitive": {
			err:  errors.New("foo bar"),
			want: []error{errors.New("foo bar")},
		},
		"wrapped primitive": {
			err: Wrap(errors.New("bar: baz")),
			fn: func(e error) []error {
				cause := errors.New("bar: baz")
				return []error{
					&wrapErr{st: GetStackTrace(e), err: cause},
					cause,
				}
			},
		},
		"double wrapped primitive": {
			err: Wrap(Wrap(errors.New("qux: xoo"))),
			fn: func(e error) []error {
				cause := errors.New("qux: xoo")
				return []error{
					&wrapErr{st: GetStackTrace(e), err: cause},
					cause,
				}
			},
		},
		"error": {
			err: Err("kind", "err msg"),
			fn: func(e error) []error {
				return []error{&err{
					st:   GetStackTrace(e),
					kind: "kind",
					msg:  "err msg",
				}}
			},
		},
		"wrapped error": {
			err: Wrap(Err("kind", "err msg")),
			fn: func(e error) []error {
				return []error{&err{
					st:   GetStackTrace(e),
					kind: "kind",
					msg:  "err msg",
				}}
			},
		},
		"double wrapped error": {
			err: Wrap(Wrap(Err("kind", "err msg"))),
			fn: func(e error) []error {
				return []error{&err{
					st:   GetStackTrace(e),
					kind: "kind",
					msg:  "err msg",
				}}
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			want := tc.want
			if want == nil {
				want = tc.fn(tc.err)
			}
			if want == nil {
				t.Fatalf("TestUnwrapAll has an invalid test case `%s` that needs to be fixed", name)
			}

			have := UnwrapAll(tc.err)
			if !reflect.DeepEqual(have, want) {
				diff := &fail.Diff{
					Func: "UnwrapAll",
					Msg:  "should unwrap all and return a slice of errors",
					Have: have,
					Want: want,
				}
				t.Error(diff.AllowUnexported(emptyErr, err{}, wrapErr{}, ST{}))
			}
		})
	}
}
