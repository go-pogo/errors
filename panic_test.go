// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

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

		Must(false, New(errStr))
	})
}

func TestCatchPanic(t *testing.T) {
	disableCaptureFrames()
	defer enableCaptureFrames()

	t.Run("panic string", func(t *testing.T) {
		var want error
		defer func() {
			assert.Equal(t, newErr(&panicErr{"paniek!"}, 0), want)
		}()
		defer CatchPanic(&want)
		panic("paniek!")
	})

	t.Run("panic error", func(t *testing.T) {
		var have, want error
		defer func() {
			assert.Same(t, want, have)
			assert.Equal(t, "panic error", have.Error())
		}()
		defer CatchPanic(&want)

		have = New("panic error")
		panic(have)
	})
}
