package errs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func panicOnSomething() {
	panic("panic!")
}

func TestWrapPanic(t *testing.T) {
	t.Run("without panic", func(t *testing.T) {
		defer func() {
			assert.Nil(t, recover())
		}()

		defer WrapPanic("wrapped")
	})

	t.Run("with panic", func(t *testing.T) {
		defer func() {
			assert.Equal(t, "wrapped: panic!", recover())
		}()

		defer WrapPanic("wrapped")
		panicOnSomething()
	})
}

func TestMust(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		defer func() {
			assert.Nil(t, recover())
		}()

		var err error
		Must(true, err)
	})

	t.Run("panic on error", func(t *testing.T) {
		errStr := "foo error"
		defer func() {
			assert.Contains(t, recover(), errStr)
		}()

		Must(false, New(UnknownKind, errStr))
	})
}
