package errs

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/roeldev/go-fail"
)

func panicOnSomething() {
	panic("panic!")
}

func TestWrapPanic(t *testing.T) {
	defer func() {
		have := recover()
		want := "wrapped: panic!"

		if have != want {
			t.Error(fail.Diff{
				Func: "WrapPanic",
				Msg:  "should wrap the panic with a prefix",
				Have: have,
				Want: want,
			})
		}
	}()

	defer WrapPanic("wrapped")
	panicOnSomething()
}

func TestWrap(t *testing.T) {
	want := New("", "cause")
	// var stLen uint = 0
	// if st := GetStackTrace(want); st != nil {
	// 	stLen = st.Len()
	// }

	have := Trace(want)
	if !reflect.DeepEqual(have, want) {
		t.Error(fail.Diff{
			Func: "Trace",
			Msg:  "should return the same error used as input",
			Have: have,
			Want: want,
		})
	}

	// st := GetStackTrace(have)
	// if st != nil && st.Len() <= stLen {
	// 	t.Error(fail.Msg{
	// 		Func: "Trace",
	// 		Msg:  "should capture an extra stack trace frame",
	// 	})
	// }
}

func TestWrap__nil(t *testing.T) {
	if Trace(nil) != nil {
		t.Error(fail.Msg{
			Func: "Trace",
			Msg:  "should return nil on nil input",
		})
	}
}

func TestWrap__primitive(t *testing.T) {
	cause := errors.New("cause")
	wrapped, ok := Trace(cause).(*traceErr)
	if !ok {
		t.Fatal(fail.Diff{
			Func: "Trace",
			Msg:  "should wrap a primitive error with a wrapErr type",
			Have: reflect.TypeOf(wrapped).String(),
			Want: reflect.TypeOf(&traceErr{}).String(),
		})
	}

	// if wrapped.Trace().Len() < 1 {
	// 	t.Error(fail.Msg{
	// 		Func: "Trace",
	// 		Msg:  "should at least capture a stack trace frame",
	// 	})
	// }z

	if have := wrapped.Unwrap(); have != cause {
		t.Error(fail.Diff{
			Func: "wrapErr.Unwrap",
			Msg:  "should return the same error instance used as input",
			Have: have,
			Want: cause,
		})
	}

	// if have := wrapped.Message(); have != cause.Error() {
	// 	t.Error(fail.Diff{
	// 		Func: "wrapErr.Message",
	// 		Msg:  "should return the same error message as the error used as input",
	// 		Have: have,
	// 		Want: cause.Error(),
	// 	})
	// }
}

var cmpOpts []cmp.Option

func init() {
	errStr := errors.New("")
	// error.errorString
	errorString := reflect.Indirect(reflect.ValueOf(errStr)).Interface()
	// fmt.wrapError
	fmtWrapError := reflect.Indirect(reflect.ValueOf(fmt.Errorf("%w", errStr))).Interface()

	cmpOpts = []cmp.Option{
		cmp.AllowUnexported(Inner{}, traceErr{}, errorString, fmtWrapError),
	}
}

func TestUnwrapAll(t *testing.T) {
	tests := map[string]struct {
		err    error
		wantFn func(e error) []error
	}{
		"nil": {
			err: Trace(nil),
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
		// "wrapped primitive": {
		// 	err: Trace(errors.New("bar: baz")),
		// 	wantFn: func(e error) []error {
		// 		cause := errors.New("bar: baz")
		// 		return []error{
		// 			&wrapErr{st: GetStackTrace(e), err: cause},
		// 			cause,
		// 		}
		// 	},
		// },
		// "double wrapped primitive": {
		// 	err: Trace(Trace(errors.New("qux: xoo"))),
		// 	wantFn: func(e error) []error {
		// 		cause := errors.New("qux: xoo")
		// 		return []error{
		// 			&wrapErr{st: GetStackTrace(e), err: cause},
		// 			cause,
		// 		}
		// 	},
		// },
		"primitive wrap": {
			err: fmt.Errorf("cause: %w", errors.New("foo bar")),
			wantFn: func(e error) []error {
				return []error{e, errors.New("foo bar")}
			},
		},
		// "wrapped primitive wrap": {
		// 	err: Trace(fmt.Errorf("cause: %w", errors.New("foo bar"))),
		// 	wantFn: func(e error) []error {
		// 		cause := errors.Unwrap(e)
		// 		return []error{
		// 			&wrapErr{st: GetStackTrace(e), err: cause},
		// 			cause,
		// 			errors.New("foo bar"),
		// 		}
		// 	},
		// },
		// "error": {
		// 	err: New(Kind("kind"), "err msg"),
		// 	wantFn: func(e error) []error {
		// 		return []error{
		// 			&err{Inner: Inner{
		// 				st:   GetStackTrace(e),
		// 				kind: "kind",
		// 				msg:  "err msg",
		// 			}},
		// 		}
		// 	},
		// },
		// "wrapped error": {
		// 	err: Trace(New(Kind("kind"), "err msg")),
		// 	wantFn: func(e error) []error {
		// 		return []error{
		// 			&err{Inner: Inner{
		// 				st:   GetStackTrace(e),
		// 				kind: "kind",
		// 				msg:  "err msg",
		// 			}},
		// 		}
		// 	},
		// },
		// "double wrapped error": {
		// 	err: Trace(Trace(New("kind", "err msg"))),
		// 	wantFn: func(e error) []error {
		// 		return []error{
		// 			&err{Inner: Inner{
		// 				st:   GetStackTrace(e),
		// 				kind: "kind",
		// 				msg:  "err msg",
		// 			}},
		// 		}
		// 	},
		// },
		// "wrapped error error": {
		// 	err: Wrap(New("baz", "qux"), "foo kind", "bar msg"),
		// 	wantFn: func(e error) []error {
		// 		wrErr := &err{
		// 			Inner: Inner{
		// 				st:   GetStackTrace(errors.Unwrap(e)),
		// 				kind: "baz",
		// 				msg:  "qux",
		// 			},
		// 		}
		//
		// 		return []error{
		// 			&err{Inner: Inner{
		// 				st:   GetStackTrace(e),
		// 				err:  wrErr,
		// 				kind: "foo kind",
		// 				msg:  "bar msg",
		// 			}},
		// 			wrErr,
		// 		}
		// 	},
		// },
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
			err: Trace(errors.New("foo bar")),
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
			err: Trace(fmt.Errorf("cause: %w", errors.New("baz"))),
			wantFn: func(e error) error {
				return errors.New("baz")
			},
		},
		// "error": {
		// 	err: New("qux", "xoo"),
		// 	wantFn: func(e error) error {
		// 		return &err{
		// 			Inner: Inner{
		// 				st:   GetStackTrace(e),
		// 				kind: "qux",
		// 				msg:  "xoo",
		// 			},
		// 		}
		// 	},
		// },
		// "wrapped error": {
		// 	err: Trace(New("qux", "xoo")),
		// 	wantFn: func(e error) error {
		// 		return &err{
		// 			Inner: Inner{
		// 				st:   GetStackTrace(e),
		// 				kind: "qux",
		// 				msg:  "xoo",
		// 			},
		// 		}
		// 	},
		// },
		// "double wrapped error": {
		// 	err: Trace(Trace(New("qux", "xoo"))),
		// 	wantFn: func(e error) error {
		// 		return &err{
		// 			Inner: Inner{
		// 				st:   GetStackTrace(e),
		// 				kind: "qux",
		// 				msg:  "xoo",
		// 			},
		// 		}
		// 	},
		// },
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
