package errs

import (
	"testing"

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
