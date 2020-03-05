package errs

import (
	"errors"
	"reflect"
	"testing"

	"github.com/roeldev/go-fail"
)

const kindFixture Kind = "foo bar baz"

type errWithKindFixture struct{}

func (err errWithKindFixture) Kind() Kind    { return kindFixture }
func (err errWithKindFixture) Error() string { return "" }

func TestNew(t *testing.T) {
	tests := map[string]struct {
		fnName    string
		err       error
		wantCause error
		wantKind  Kind
		wantMsg   string
	}{
		"new with message": {
			err:       New("foo message"),
			wantCause: nil,
			wantKind:  UnknownKind,
			wantMsg:   "foo message",
		},
		"new with kind and message": {
			err:       New(Kind("foo error"), "bar message"),
			wantCause: nil,
			wantKind:  Kind("foo error"),
			wantMsg:   "bar message",
		},
		"new with cause, kind and message": {
			err:       New(errors.New("underlying err"), Kind("qux error"), "caused by"),
			wantCause: errors.New("underlying err"),
			wantKind:  Kind("qux error"),
			wantMsg:   "caused by",
		},
		"wrap": {
			err:       Wrap(errors.New("cause")),
			wantCause: errors.New("cause"),
			wantKind:  UnknownKind,
			wantMsg:   "cause",
		},
		"wrap nil": {
			fnName:    "Wrap",
			err:       Wrap(nil),
			wantCause: nil,
		},
		"wrap cause with kind": {
			fnName:    "Wrap",
			err:       Wrap(errWithKindFixture{}),
			wantCause: errWithKindFixture{},
			wantKind:  kindFixture,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.fnName == "" {
				tc.fnName = name
			}

			if err, ok := tc.err.(ErrorWithStackTrace); ok {
				st := err.StackTrace()
				if st == nil || st.Len() == 0 {
					t.Error(fail.Msg{
						Func: tc.fnName,
						Msg:  "should create an error with a stack trace",
					})
				}
			}

			if err, ok := tc.err.(ErrorWithUnwrap); ok {
				have := err.Unwrap()
				if !reflect.DeepEqual(have, tc.wantCause) {
					t.Error(fail.Diff{
						Func: tc.fnName,
						Msg:  "should create an error with given cause error",
						Have: have,
						Want: tc.wantCause,
					})
				}
			}

			if err, ok := tc.err.(ErrorWithKind); ok {
				have := err.Kind()
				if have != tc.wantKind {
					t.Error(fail.Diff{
						Func: tc.fnName,
						Msg:  "should create an error with given kind",
						Have: have,
						Want: tc.wantKind,
					})
				}
			}

			if err, ok := tc.err.(ErrorWithMessage); ok {
				have := err.Message()
				if have != tc.wantMsg {
					t.Error(fail.Diff{
						Func: tc.fnName,
						Msg:  "should create an error with given message",
						Have: have,
						Want: tc.wantMsg,
					})
				}
			}
		})
	}
}

func TestNew_print(t *testing.T) {
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

func TestErrorPrint(t *testing.T) {
	tests := map[string]error{
		"err": New("foo", "bar baz"),
		// "wrapErr": Wrap(errors.New("foo: bar")),
	}

	for name, err := range tests {
		t.Run(name, func(t *testing.T) {
			have := err.Error()
			want := Print(err)
			if have != want {
				t.Error(fail.Diff{
					Func: name + ".Error",
					Msg:  "should use the Print() util function to create the error message",
					Have: have,
					Want: want,
				})
			}
		})
	}
}
