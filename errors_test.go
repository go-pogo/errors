package errs

import (
	"errors"
	"reflect"
	"testing"

	"github.com/roeldev/go-fail"
)

func TestNew(t *testing.T) {
	tests := map[string]struct {
		input []interface{}
		inner Inner
	}{
		"with message": {
			input: []interface{}{"foo message"},
			inner: Inner{
				kind: UnknownKind,
				msg:  "foo message",
			},
		},
		"with kind and message": {
			input: []interface{}{Kind("foo error"), "bar message"},
			inner: Inner{
				kind: Kind("foo error"),
				msg:  "bar message",
			},
		},
		"with cause, kind and message": {
			input: []interface{}{errors.New("underlying err"), Kind("qux error"), "caused by"},
			inner: Inner{
				err:  errors.New("underlying err"),
				kind: Kind("qux error"),
				msg:  "caused by",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			newErr := New(tc.input...)
			tcErr, ok := newErr.(*err)
			if !ok {
				t.Fatal(fail.Diff{
					Func: "New",
					Msg:  "should return a type of `*err`",
					Have: reflect.TypeOf(newErr).String(),
					Want: reflect.TypeOf(&err{}).String(),
				})
			}

			if tcErr.StackTrace().Len() == 0 {
				t.Error(fail.Msg{
					Func: "New",
					Msg:  "should create an error with a stack trace",
				})
			}

			if have := tcErr.Unwrap(); !reflect.DeepEqual(have, tc.inner.err) {
				t.Error(fail.Diff{
					Func: "New",
					Msg:  "should create an error with given cause error",
					Have: have,
					Want: tc.inner.err,
				})
			}

			if have := tcErr.Kind(); have != tc.inner.kind {
				t.Error(fail.Diff{
					Func: "New",
					Msg:  "should create an error with given kind",
					Have: have,
					Want: tc.inner.kind,
				})
			}

			if have := tcErr.Message(); have != tc.inner.msg {
				t.Error(fail.Diff{
					Func: "New",
					Msg:  "should create an error with given message",
					Have: have,
					Want: tc.inner.msg,
				})
			}
		})
	}
}

func TestNew__panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error(fail.Msg{
				Func: "New",
				Msg:  "should panic when receiving no arguments",
			})
		}
	}()

	_ = New()
}

func TestWrap(t *testing.T) {
	want := New("cause")
	var stLen uint = 0
	if st := GetStackTrace(want); st != nil {
		stLen = st.Len()
	}

	have := Wrap(want)
	if have != want {
		t.Error(fail.Diff{
			Func: "Wrap",
			Msg:  "should return the same error used as input",
			Have: have,
			Want: want,
		})
	}

	st := GetStackTrace(have)
	if st != nil && st.Len() <= stLen {
		t.Error(fail.Msg{
			Func: "Wrap",
			Msg:  "should capture an extra stack trace frame",
		})
	}
}

func TestWrap__nil(t *testing.T) {
	if Wrap(nil) != nil {
		t.Error(fail.Msg{
			Func: "Wrap",
			Msg:  "should return nil on nil input",
		})
	}
}

func TestWrap__primitive(t *testing.T) {
	cause := errors.New("cause")
	wrapped, ok := Wrap(cause).(*wrapErr)
	if !ok {
		t.Fatal(fail.Diff{
			Func: "Wrap",
			Msg:  "should wrap a primitive error with a wrapErr type",
			Have: reflect.TypeOf(wrapped).String(),
			Want: reflect.TypeOf(&wrapErr{}).String(),
		})
	}

	if wrapped.StackTrace().Len() < 1 {
		t.Error(fail.Msg{
			Func: "Wrap",
			Msg:  "should at least capture one stack trace frame",
		})
	}

	if have := wrapped.Unwrap(); have != cause {
		t.Error(fail.Diff{
			Func: "wrapErr.Unwrap",
			Msg:  "should return the same error instance used as input",
			Have: have,
			Want: cause,
		})
	}

	if have := wrapped.Message(); have != cause.Error() {
		t.Error(fail.Diff{
			Func: "wrapErr.Message",
			Msg:  "should return the same error message as the error used as input",
			Have: have,
			Want: cause.Error(),
		})
	}
}
