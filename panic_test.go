package errs

import (
	"strings"
	"testing"

	"github.com/roeldev/go-fail"
)

func panicOnSomething() {
	panic("panic!")
}

func TestWrapPanic(t *testing.T) {
	defer func() {
		have := recover().(string)
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

func TestMust_nil_error(t *testing.T) {
	defer func() {
		have := recover()
		if have != nil {
			t.Error(fail.Diff{
				Func: "Must",
				Msg:  "must not panic on nil error",
				Have: have,
				Want: nil,
			})
		}
	}()

	var err error
	Must(true, err)
}

func TestMust_panic_on_error(t *testing.T) {
	defer func() {
		have := recover().(string)
		want := "errs.Must: foo error"

		if !strings.HasPrefix(have, want) {
			t.Error(fail.Diff{
				Func: "Must",
				Msg:  "must panic and include the causing error",
				Have: have,
				Want: want,
			})
		}
	}()

	Must(false, New(UnknownKind, "foo error"))
}
