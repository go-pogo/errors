package errs

import (
	"errors"
	"strings"
	"testing"

	"github.com/roeldev/go-fail"
)

func Test__equal_errors(t *testing.T) {
	cause := New("", "some cause")
	tests := map[string]map[string][2]error{

		"New/Newf": {
			"empty": {
				New("", ""),
				Newf("", ""),
			},
			"plain": {
				New("", "foobar"),
				Newf("", "foobar"),
			},
			"formatted": {
				New("", "foo: bar"),
				Newf("", "foo: %s", "bar"),
			},
		},
		"Wrap/Wrapf": {
			"empty": {
				Wrap(cause, "", ""),
				Wrapf(cause, "", ""),
			},
			"plain": {
				Wrap(cause, "", "foobar"),
				Wrapf(cause, "", "foobar"),
			},
			"formatted": {
				Wrap(cause, "", "foo: bar"),
				Wrapf(cause, "", "foo: %s", "bar"),
			},
		},
	}

	for fn, ts := range tests {
		fn = strings.Replace(fn, "/", "&", 1)
		for name, tc := range ts {
			t.Run(fn+"__"+name, func(t *testing.T) {
				if have, want := GetKind(tc[1]), GetKind(tc[0]); have != want {
					t.Error(fail.Diff{
						Func: fn,
						Msg:  "both functions should create the same error kind",
						Have: have,
						Want: want,
					})
				}
				if have, want := GetMessage(tc[1]), GetMessage(tc[0]); have != want {
					t.Error(fail.Diff{
						Func: fn,
						Msg:  "both functions should create the same message",
						Have: have,
						Want: want,
					})
				}
				if have, want := errors.Unwrap(tc[1]), errors.Unwrap(tc[0]); have != want {
					t.Error(fail.Diff{
						Func: fn,
						Msg:  "both functions should have the same underlying cause",
						Have: have,
						Want: want,
					})
				}
			})
		}
	}
}

// func Test(t *testing.T) {
// 	tests := map[string]struct {
// 		fn    string
// 		have  error
// 		inner Inner
// 	}{
// 		"with message": {
// 			fn:   "New",
// 			have: New("", "foo message"),
// 			inner: Inner{
// 				kind: UnknownKind,
// 				msg:  "foo message",
// 			},
// 		},
// 		"with kind and message": {
// 			fn:   "New",
// 			have: New("foo error", "bar message"),
// 			inner: Inner{
// 				kind: Kind("foo error"),
// 				msg:  "bar message",
// 			},
// 		},
// 		"with kind and formatted message": {
// 			fn:   "Newf",
// 			have: Newf("foo error", "%s message", "some"),
// 			inner: Inner{
// 				kind: Kind("foo error"),
// 				msg:  "some message",
// 			},
// 		},
// 		"with cause, kind and message": {
// 			fn:   "Wrap",
// 			have: Wrap(errors.New("underlying err"), "qux error", "caused by"),
// 			inner: Inner{
// 				err:  errors.New("underlying err"),
// 				kind: Kind("qux error"),
// 				msg:  "caused by",
// 			},
// 		},
// 	}
//
// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			newWrap := New(tc.input...)
// 			tcWrap, ok := newWrap.(*err)
// 			if !ok {
// 				t.Fatal(fail.Diff{
// 					Func: "New",
// 					Msg:  "should return a type of `*err`",
// 					Have: reflect.TypeOf(newWrap).String(),
// 					Want: reflect.TypeOf(&err{}).String(),
// 				})
// 			}
//
// 			if tcWrap.StackWrap().Len() == 0 {
// 				t.Error(fail.Msg{
// 					Func: "New",
// 					Msg:  "should create an error with a stack trace",
// 				})
// 			}
//
// 			if have := tcWrap.Unwrap(); !reflect.DeepEqual(have, tc.inner.err) {
// 				t.Error(fail.Diff{
// 					Func: "New",
// 					Msg:  "should create an error with given cause error",
// 					Have: have,
// 					Want: tc.inner.err,
// 				})
// 			}
//
// 			if have := tcWrap.Kind(); have != tc.inner.kind {
// 				t.Error(fail.Diff{
// 					Func: "New",
// 					Msg:  "should create an error with given kind",
// 					Have: have,
// 					Want: tc.inner.kind,
// 				})
// 			}
//
// 			if have := tcWrap.Message(); have != tc.inner.msg {
// 				t.Error(fail.Diff{
// 					Func: "New",
// 					Msg:  "should create an error with given message",
// 					Have: have,
// 					Want: tc.inner.msg,
// 				})
// 			}
// 		})
// 	}
// }

// func TestNew__panic(t *testing.T) {
// 	defer func() {
// 		if r := recover(); r == nil {
// 			t.Error(fail.Msg{
// 				Func: "New",
// 				Msg:  "should panic when receiving no arguments",
// 			})
// 		}
// 	}()
//
// 	_ = New()
// }
