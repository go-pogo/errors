package errs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func panicOnSomething() {
	panic("panic!")
}

func TestWrapPanic(t *testing.T) {
	defer func() {
		assert.Equal(t, "wrapped: panic!", recover())
	}()

	defer WrapPanic("wrapped")
	panicOnSomething()
}

func TestMust_nil_error(t *testing.T) {
	defer func() {
		assert.Nil(t, recover())
	}()

	var err error
	Must(true, err)
}

func TestMust_panic_on_error(t *testing.T) {
	errStr := "foo error"
	defer func() {
		assert.Contains(t, recover(), errStr)
	}()

	Must(false, New(UnknownKind, errStr))
}
